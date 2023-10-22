<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [分布式锁](#%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81)
  - [难点](#%E9%9A%BE%E7%82%B9)
    - [死锁](#%E6%AD%BB%E9%94%81)
      - [四大必要 条件分别为：](#%E5%9B%9B%E5%A4%A7%E5%BF%85%E8%A6%81-%E6%9D%A1%E4%BB%B6%E5%88%86%E5%88%AB%E4%B8%BA)
    - [惊群效应](#%E6%83%8A%E7%BE%A4%E6%95%88%E5%BA%94)
    - [脑裂](#%E8%84%91%E8%A3%82)
  - [Consul分布式锁源码分析](#consul%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [锁结构体](#%E9%94%81%E7%BB%93%E6%9E%84%E4%BD%93)
    - [锁选项](#%E9%94%81%E9%80%89%E9%A1%B9)
    - [lock,unlock](#lockunlock)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 分布式锁

## 难点

### 死锁
定义：当两个以上的运算单元，双方都在等待对方停止运行，以获取系统资源，但是没有一方提前退出时，就称为死锁

#### 四大必要 条件分别为：

- 互斥条件
- 不可抢占条件
- 占用并申请条件
- 循环等待条件

解决方式： 在分布式环境中，我们使用 session + TTL 来打破 循环等待条件

当一个客户端尝试操作一把分布式锁的时候，我们必须校验其 session 是否为锁的拥有者，则无法进行操作。

当一个客户端已经持有一把分布式锁后发生了 掉线 ，在超出了 TTL 时间后无法连接上，则回收其锁的拥有权。

### 惊群效应
举一个例子说明惊群效应：老农养了 N 只鸽子，而老农在进行喂食的时候，每次抛出食物的时候，都会引起所有鸽子的注意，纷纷来抢夺食物，这就是 惊群效应 。

zookeeper的解决方式：

    当客户端来请求该锁的时候， ZooKeeper 会生成一个 ${lock-name}-${i} 的临时目录，此后每次请求该锁的时候，就会生成 ${lock-name}-${i+1} 的目录，
    如果此时在 ${lock-name} 中拥有最小的 i 的客户端会获得该锁，而该客户端使用结束之后，就会删除掉自己的临时目录，并通知后续节点进行锁的获取
    
    没错，这个 i 是 ZooKeeper 解决惊群效应的利器，它被称为 顺序节点 。

### 脑裂
脑裂是集群环境中肯定会遇到的问题，其出现的主要原因为 网络波动

举最常见的 双机热备 的场景，节点A 和 节点B 是相同功能的两个服务，它们彼此通过心跳保持联系，并作为对方的备份。
但如果此时 A 与 B 的网络连接被中断了， A 会尝试占用 B 的资源，而 B 会尝试占用 A 的资源，这就是 脑裂 问题。

当集群中出现 脑裂 的时候，往往会出现多个 master 的情况，这样数据的一致性会无法得到保障，从而导致整个服务无法正常运行。

而在 ZooKeeper 的分布式锁场景中，如果 客户端A 已经得到了锁，但是却因为网络波动原因断开了与 ZooKeeper 的连接，
那么下一个顺序节点 客户端B 就会获得锁，但是因为 客户端A 此时依然还持有该锁 ，从而发生了 脑裂 问题。

解决脑裂问题有两种方式，
1. 可以将集群中的服务作为 P2P 节点，避免 Leader 与 Salve 的切换，
2. 另一种方案更简单一些，那就是当客户端与 ZooKeeper 非正常断开连接的时候， ZooKeeper 应该尝试向客户端发起 多次重试 机制，并在一段时间后依然无法连接上，再让下一个顺序客户端获取锁


## Consul分布式锁源码分析
路径： /github.com/hashicorp/consul/api@v1.3.0/lock.go

### 锁结构体
```go
// Lock 分布式锁数据结构
type Lock struct {
	c    *Client  // 提供访问consul的API客户端
	opts *LockOptions // 分布式锁的可选项

	isHeld       bool // 该锁当前是否已经被持有
	sessionRenew chan struct{} // 通知锁持有者需要更新session
	lockSession  string // 锁持有者的session
	l            sync.Mutex // 锁变量的互斥锁
}
```
### 锁选项
```go
// LockOptions is used to parameterize the Lock behavior.
type LockOptions struct {
	Key              string        // 锁的 Key，必填项，且必须有 KV 的写权限
	Value            []byte        // 锁的内容，以下皆为选填项
	Session          string        // 锁的session，用于判断锁是否被创建
	SessionOpts      *SessionEntry // 自定义创建session条目，用于创建session，避免惊群
	SessionName      string        // 自定义锁的session名称，默认为 "Consul API Lock"
	SessionTTL       string        // 自定义锁的TTL时间，默认为 "15s"
	MonitorRetries   int           // 自定义监控的重试次数，避免脑裂问题
	MonitorRetryTime time.Duration // 自定义监控的重试时长，避免脑裂问题
	LockWaitTime     time.Duration // 自定义锁的等待时长，避免死锁问题
	LockTryOnce      bool          // 是否只重试一次，默认为false，则为无限重试
}
```
从 LockOptions 中带有 session / TTL / monitor / wait 等字眼的成员变量可以看出，consul 已经考虑到解决我们上一节提到的三个难点

初始化
```go
func (c *Client) LockOpts(opts *LockOptions) (*Lock, error) {
	if opts.Key == "" {
		return nil, fmt.Errorf("missing key")
	}
	if opts.SessionName == "" {
		opts.SessionName = DefaultLockSessionName
	}
    // 15s 的 SessionTTL 用于解决死锁、脑裂问题
	if opts.SessionTTL == "" {
		opts.SessionTTL = DefaultLockSessionTTL
	} else {
		if _, err := time.ParseDuration(opts.SessionTTL); err != nil {
			return nil, fmt.Errorf("invalid SessionTTL: %v", err)
		}
	}
    // 2s 的 MonitorRetryTime 是一个长期运行的协程用于监听当前锁持有者，用于解决脑裂问题。
	if opts.MonitorRetryTime == 0 {
		opts.MonitorRetryTime = DefaultMonitorRetryTime
	}

	//15s 的 LockWaitTime 用于设置尝试获取锁的超时时间，用于解决死锁问题。
	if opts.LockWaitTime == 0 {
		opts.LockWaitTime = DefaultLockWaitTime
	}
	l := &Lock{
		c:    c,
		opts: opts,
	}
	return l, nil
}
```

### lock,unlock

lock
```go
// Lock尝试获取一个可用的锁，可以通过一个非空的 stopCh 来提前终止获取
// 如果返回的锁发生异常，则返回一个被关闭的 chan struct ，应用程序必须要处理该情况
func (l *Lock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
    // 先锁定本地互斥锁
	l.l.Lock()
	defer l.l.Unlock()

	// 确认本地已经获取到分布式锁了   
	if l.isHeld {
		return nil, ErrLockHeld
	}

	//  检查是否需要创建session
	l.lockSession = l.opts.Session
	if l.lockSession == "" {
		s, err := l.createSession()
		if err != nil {
			return nil, fmt.Errorf("failed to create session: %v", err)
		}

		l.sessionRenew = make(chan struct{})
		l.lockSession = s
		session := l.c.Session()
		go session.RenewPeriodic(l.opts.SessionTTL, s, nil, l.sessionRenew)

		// 如果我们无法锁定该分布式锁，清除本地session
		defer func() {
			if !l.isHeld {
				close(l.sessionRenew)
				l.sessionRenew = nil
			}
		}()
	}

	//  准备向consul KV发送查询锁操作的参数
	kv := l.c.KV()
	qOpts := &QueryOptions{
		WaitTime: l.opts.LockWaitTime,
	}

	start := time.Now()
	attempts := 0
WAIT:
    // 判断是否需要退出锁争夺的循环
	select {
	case <-stopCh:
		return nil, nil
	default:
	}

    // 处理只重试一次的逻辑
	// 配置该锁只重试一次且已经重试至少一次了
	if l.opts.LockTryOnce && attempts > 0 {
        // 获取当前时间偏移量
		elapsed := time.Since(start)
		if elapsed > l.opts.LockWaitTime {
            // 当超过设置中的剩余等待时间
			return nil, nil
		}

		// 重设剩余等待时间
		qOpts.WaitTime = l.opts.LockWaitTime - elapsed
	}
    // 已尝试次数自增1
	attempts++

	// 阻塞查询该存在的分布式锁，直至无法获取成功
	pair, meta, err := kv.Get(l.opts.Key, qOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to read lock: %v", err)
	}
	if pair != nil && pair.Flags != LockFlagValue {
		return nil, ErrLockConflict
	}
	locked := false
	if pair != nil && pair.Session == l.lockSession {
		goto HELD
	}
	if pair != nil && pair.Session != "" {
		qOpts.WaitIndex = meta.LastIndex
		goto WAIT
	}

	// Try to acquire the lock
	pair = l.lockEntry(l.lockSession)
	locked, _, err = kv.Acquire(pair, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %v", err)
	}

	// Handle the case of not getting the lock
	if !locked {
		// Determine why the lock failed
		qOpts.WaitIndex = 0
		pair, meta, err = kv.Get(l.opts.Key, qOpts)
		if pair != nil && pair.Session != "" {
			//If the session is not null, this means that a wait can safely happen
			//using a long poll
			qOpts.WaitIndex = meta.LastIndex
			goto WAIT
		} else {
			// If the session is empty and the lock failed to acquire, then it means
			// a lock-delay is in effect and a timed wait must be used
			select {
			case <-time.After(DefaultLockRetryTime):
				goto WAIT
			case <-stopCh:
				return nil, nil
			}
		}
	}

HELD:
	// Watch to ensure we maintain leadership
	leaderCh := make(chan struct{})
	go l.monitorLock(l.lockSession, leaderCh)

	// Set that we own the lock
	l.isHeld = true

	// Locked! All done
	return leaderCh, nil
}
```

unlock
```go
// Unlock 尝试释放 consul 分布式锁，如果发生异常则返回 error
func (l *Lock) Unlock() error {
    // 在释放锁之前必须先把 Lock 结构锁住
	l.l.Lock()
	defer l.l.Unlock()

    // 确认我们依然持有该锁
	if !l.isHeld {
		return ErrLockNotHeld
	}

	// 提前先将锁的持有权释放
	l.isHeld = false

	// 清除刷新 session 通道
	if l.sessionRenew != nil {
		defer func() {
			close(l.sessionRenew)
			l.sessionRenew = nil
		}()
	}

	// 获取当前 session 持有的锁信息
	lockEnt := l.lockEntry(l.lockSession)
	l.lockSession = ""

	// 将持有的锁尝试释放
	kv := l.c.KV()
	_, _, err := kv.Release(lockEnt, nil)
	if err != nil {
		return fmt.Errorf("failed to release lock: %v", err)
	}
	return nil
}
```



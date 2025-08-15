<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Zookeeper](#zookeeper)
  - [ZooKeeper 基础知识基本分为三大模块](#zookeeper-%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86%E5%9F%BA%E6%9C%AC%E5%88%86%E4%B8%BA%E4%B8%89%E5%A4%A7%E6%A8%A1%E5%9D%97)
    - [Zookeeper的数据模型](#zookeeper%E7%9A%84%E6%95%B0%E6%8D%AE%E6%A8%A1%E5%9E%8B)
      - [1 PERSISTENT（持久节点)](#1-persistent%E6%8C%81%E4%B9%85%E8%8A%82%E7%82%B9)
      - [2 EPHEMERAL](#2-ephemeral)
      - [3 PERSISTENT_SEQUENTIAL](#3-persistent_sequential)
      - [4 EPHEMERAL_SEQUENTIAL](#4-ephemeral_sequential)
    - [watch](#watch)
      - [客户端](#%E5%AE%A2%E6%88%B7%E7%AB%AF)
      - [服务端](#%E6%9C%8D%E5%8A%A1%E7%AB%AF)
    - [ACL](#acl)
      - [权限模式：Scheme](#%E6%9D%83%E9%99%90%E6%A8%A1%E5%BC%8Fscheme)
      - [授权对象（ID）](#%E6%8E%88%E6%9D%83%E5%AF%B9%E8%B1%A1id)
      - [权限信息（Permission）](#%E6%9D%83%E9%99%90%E4%BF%A1%E6%81%AFpermission)
  - [数据存储底层实现](#%E6%95%B0%E6%8D%AE%E5%AD%98%E5%82%A8%E5%BA%95%E5%B1%82%E5%AE%9E%E7%8E%B0)
    - [数据日志](#%E6%95%B0%E6%8D%AE%E6%97%A5%E5%BF%97)
    - [快照日志](#%E5%BF%AB%E7%85%A7%E6%97%A5%E5%BF%97)
  - [客户端基本使用](#%E5%AE%A2%E6%88%B7%E7%AB%AF%E5%9F%BA%E6%9C%AC%E4%BD%BF%E7%94%A8)
  - [github.com/go-zookeeper/zk](#githubcomgo-zookeeperzk)
    - [获取数据](#%E8%8E%B7%E5%8F%96%E6%95%B0%E6%8D%AE)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Zookeeper
ZooKeeper 是一个分布式的，开放源码的分布式应用程序协调服务，是 Google 的 Chubby 一个开源的实现，是 Hadoop 和 Hbase 的重要组件。
它是一个为分布式应用提供一致性服务的软件，提供的功能包括：配置维护、域名服务、分布式同步、组服务等.

- 配置中心:发布者将数据发布到ZK节点上，供订阅者动态获取数据，实现配置信息的集中式管理和动态更新。例如全局的配置信息、服务式服务框架的服务地址列表等就非常适合使用
- 负载均衡: 消息中间件中发布者和订阅者的负载均衡linkedin开源的Kafka和阿里开源的metaq（RocketMQ的前身）都是通过zookeeper来做到生产者、消费者的负载均衡。
- 分布式通知/协调




## ZooKeeper 基础知识基本分为三大模块

- 数据模型
- ACL 权限控制
- Watch 监控

### Zookeeper的数据模型
![](.zookeeper_images/zookeeper_data_structure.png)
Zookeeper数据模型的结构与Unix文件系统很类似，整体上可以看作是一颗树，每一个节点称做一个ZNode。
每一个 Znode 默认能够存储1MB的数据，每个ZNode都可以通过其路径唯一标识。


![](.zookeeper_images/zookeeper_node.png)
zookeeper 中数据基本单元叫节点，节点之下可包含子节点，最后以树级方式程现。每个节点拥有唯一的路径path。客户端基于PATH上传节点数据，zookeeper 收到后会实时通知对该路径进行监听的客户端。

zookeeper 中节点叫znode存储结构上跟文件系统类似，以树级结构进行存储。不同之外在于znode没有目录的概念，不能执行类似cd之类的命令。znode结构包含如下：

- path:唯一路径
- childNode：子节点
- stat:状态属性
- type:节点类型

节点类型:
临时节点（ephemeral）、持久节点（persistent）、顺序节点（sequence）。节点类型在创建时确定，之后不可修改。

#### 1 PERSISTENT（持久节点)

持久节点除非手动删除，否则节点一直存在于 Zookeeper 上

#### 2 EPHEMERAL
临时节点临时节点的生命周期与客户端会话绑定，一旦客户端会话失效（客户端与zookeeper 连接断开不一定会话失效），那么这个客户端创建的所有临时节点都会被移除


#### 3 PERSISTENT_SEQUENTIAL
持久顺序节点基本特性同持久节点，只是增加了顺序属性，节点名后边会追加一个由父节点维护的自增整型数字。

#### 4 EPHEMERAL_SEQUENTIAL
临时顺序节点基本特性同临时节点，增加了顺序属性，节点名后边会追加一个由父节点维护的自增整型数字。


```shell
# 节点状态信息
[zk: localhost:2181(CONNECTED) 9] stat /china
cZxid = 0x2
ctime = Sun Jun 08 11:15:06 CST 2025
mZxid = 0x2
mtime = Sun Jun 08 11:15:06 CST 2025
pZxid = 0x2
cversion = 0
dataVersion = 0
aclVersion = 0
ephemeralOwner = 0x0
dataLength = 3
numChildren = 0
```
字段解释
* czxid The zxid of the change that caused this znode to be created.
* mzxid The zxid of the change that last modified this znode.
* pzxid The zxid of the change that last modified children of this znode.
* ctime The time in milliseconds from epoch when this znode was created.
* mtime The time in milliseconds from epoch when this znode was last modified.
* version The number of changes to the data of this znode.
* cversion The number of changes to the children of this znode.
* aversion The number of changes to the ACL of this znode.
* ephemeralOwner The session id of the owner of this znode if the znode is an ephemeral node. If it is not an ephemeral node, it will be zero.
* dataLength The length of the data field of this znode.
* numChildren The number of children of this znode.

### watch


#### 客户端

```go
// 添加 watcher
func (c *Conn) addWatcher(path string, watchType watchType) <-chan Event {
	c.watchersLock.Lock()
	defer c.watchersLock.Unlock()

	ch := make(chan Event, 1)
	wpt := watchPathType{path, watchType}
	c.watchers[wpt] = append(c.watchers[wpt], ch)
	return ch
}

```


#### 服务端

Zookeeper 服务端处理 Watch 事件基本有 2 个过程：

1. 解析收到的请求是否带有 Watch 注册事件
2. 将对应的 Watch 事件存储到 WatchManager

```java
// https://github.com/apache/zookeeper/blob/c26634f34490bb0ea7a09cc51e05ede3b4e320ee/zookeeper-server/src/main/java/org/apache/zookeeper/server/FinalRequestProcessor.java

public class FinalRequestProcessor implements RequestProcessor {

    private Record handleGetDataRequest(Record request, ServerCnxn cnxn, List<Id> authInfo) throws KeeperException, IOException {
        GetDataRequest getDataRequest = (GetDataRequest) request;
        String path = getDataRequest.getPath();
        DataNode n = zks.getZKDatabase().getNode(path);
        if (n == null) {
            throw new KeeperException.NoNodeException();
        }
        zks.checkACL(cnxn, zks.getZKDatabase().aclForNode(n), ZooDefs.Perms.READ, authInfo, path, null);
        Stat stat = new Stat();
        byte[] b = zks.getZKDatabase().getData(path, stat, getDataRequest.getWatch() ? cnxn : null);
        return new GetDataResponse(b, stat);
    }
    
}
```

注册 watch 
```java
public class WatchManager implements IWatchManager {
    @Override
    public synchronized boolean addWatch(String path, Watcher watcher, WatcherMode watcherMode) {
        if (isDeadWatcher(watcher)) {
            LOG.debug("Ignoring addWatch with closed cnxn");
            return false;
        }

        Set<Watcher> list = watchTable.get(path);
        if (list == null) {
            // don't waste memory if there are few watches on a node
            // rehash when the 4th entry is added, doubling size thereafter
            // seems like a good compromise
            list = new HashSet<>(4);
            watchTable.put(path, list);
        }
        list.add(watcher);

        Map<String, WatchStats> paths = watch2Paths.get(watcher);
        if (paths == null) {
            // cnxns typically have many watches, so use default cap here
            paths = new HashMap<>();
            watch2Paths.put(watcher, paths);
        }

        WatchStats stats = paths.getOrDefault(path, WatchStats.NONE);
        WatchStats newStats = stats.addMode(watcherMode);

        if (newStats != stats) {
            paths.put(path, newStats);
            if (watcherMode.isRecursive()) {
                ++recursiveWatchQty;
            }
            return true;
        }

        return false;
    }

}
```

设置数据后触发事件
```java
public class DataTree {
    public Stat setData(String path, byte[] data, int version, long zxid, long time) throws NoNodeException {
        Stat s = new Stat();
        DataNode n = nodes.get(path);
        if (n == null) {
            throw new NoNodeException();
        }
        List<ACL> acl;
        byte[] lastData;
        synchronized (n) {
            acl = getACL(n);
            lastData = n.data;
            nodes.preChange(path, n);
            n.data = data;
            n.stat.setMtime(time);
            n.stat.setMzxid(zxid);
            n.stat.setVersion(version);
            n.copyStat(s);
            nodes.postChange(path, n);
        }

        // first do a quota check if the path is in a quota subtree.
        String lastPrefix = getMaxPrefixWithQuota(path);
        long bytesDiff = (data == null ? 0 : data.length) - (lastData == null ? 0 : lastData.length);
        // now update if the path is in a quota subtree.
        long dataBytes = data == null ? 0 : data.length;
        if (lastPrefix != null) {
            updateQuotaStat(lastPrefix, bytesDiff, 0);
        }
        nodeDataSize.addAndGet(getNodeSize(path, data) - getNodeSize(path, lastData));

        updateWriteStat(path, dataBytes);
        dataWatches.triggerWatch(path, EventType.NodeDataChanged, zxid, acl);
        return s;
    }

}
```

### ACL 

一个 ACL 权限设置通常可以分为 3 部分，分别是：权限模式（Scheme）、授权对象（ID）、权限信息（Permission）。最终组成一条例如“scheme:id:permission”格式的 ACL 请求信息。


#### 权限模式：Scheme
权限模式就是用来设置 ZooKeeper 服务器进行权限验证的方式。ZooKeeper 的权限验证方式大体分为两种类型，一种是范围验证，另外一种是口令验证。

所谓的范围验证就是说 ZooKeeper 可以针对一个 IP 或者一段 IP 地址授予某种权限。比如我们可以让一个 IP 地址为“ip：192.168.0.11”的机器对服务器上的某个数据节点具有写入的权限。
或者也可以通过“ip:192.168.0.11/22”给一段 IP 地址的机器赋权。

另一种权限模式就是口令验证，也可以理解为用户名密码的方式.
在 ZooKeeper 中这种验证方式是 Digest 认证，我们知道通过网络传输相对来说并不安全，所以“绝不通过明文在网络发送密码”也是程序设计中很重要的原则之一，而 Digest 这种认证方式首先在客户端传送“username:password”这种形式的权限表示符后，ZooKeeper 服务端会对密码 部分使用 SHA-1 和 BASE64 算法进行加密，以保证安全性。
```shell
 echo -n user-1:password-1 | openssl dgst -binary -sha1 | openssl base64
```

最后一种授权模式是 world 模式，其实这种授权模式对应于系统中的所有用户，本质上起不到任何作用。


#### 授权对象（ID）

#### 权限信息（Permission）

在 ZooKeeper 中已经定义好的权限有 5 种：

* 数据节点（create）创建权限，授予权限的对象可以在数据节点下创建子节点；
* 数据节点（write）更新权限，授予权限的对象可以更新该数据节点；
* 数据节点（read）读取权限，授予权限的对象可以读取该节点的内容以及子节点的信息；
* 数据节点（delete）删除权限，授予权限的对象可以删除该数据节点的子节点；
* 数据节点（admin）管理者权限，授予权限的对象可以对该数据节点体进行 ACL 权限设置。


需要注意的一点是，每个节点都有维护自身的 ACL 权限数据，即使是该节点的子节点也是有自己的 ACL 权限而不是直接继承其父节点的权限。


## 数据存储底层实现

ZooKeeper 中的数据存储，可以分为两种类型：数据日志文件和快照文件，

### 数据日志

在 ZooKeeper 服务运行的过程中，数据日志是用来记录 ZooKeeper 服务运行状态的数据文件.

### 快照日志  
快照日志文件主要用来存储 ZooKeeper 服务中的事务性操作日志，并通过数据快照文件实现集群之间服务器的数据同步功能。

存储到本地磁盘中的数据快照文件，是经过 ZooKeeper 序列化后的二进制格式文件，通常我们无法直接查看，但如果想要查看，也可以通过 ZooKeeper 自带的 SnapshotFormatter 类来实现。


## 客户端基本使用
https://zookeeper.apache.org/doc/r3.9.3/zookeeperCLI.html

```shell
# 部署命令
$ mkdir zookeeper_data
$ docker run -d -e TZ="Asia/Shanghai" -p 2181:2181 -v $PWD/zookeeper_data:/data --name zookeeper --restart always zookeeper:3.8.4
```

```shell
# 连接 zk 服务器
$ docker exec -it zookeeper zkCli.sh -server localhost:2181

# 查看子节点-ls
ls /brokers


# create [-s] [-e] path data   
# 其中 -s 为有序节点， -e 临时节点
# 创建持久节点:创建一个名称为 china 的 znode，其值为 999
create /china 999

# 创建持久顺序节点:在/china 节点下创建了顺序子节点 beijing、 shanghai、 guangzhou，它们的数据内容分别为 bj、 sh、 gz
create -s /china/beijing bj
create -s /china/shanghai sh
create -s /china/guangzhou gz

# 创建临时节点
create -e /china/aaa A

# 创建临时顺序节点
create -e  -s /china/bbb B

# 获取节点信息 get
## 获取持久节点数据
get /china

# 更新节点数据内容-set

# 删除节点-delete
```



## github.com/go-zookeeper/zk


连接
```go
// github.com/go-zookeeper/zk@v1.0.4/conn.go
func Connect(servers []string, sessionTimeout time.Duration, options ...connOption) (*Conn, <-chan Event, error) {
	// 校验连接地址
	if len(servers) == 0 {
		return nil, nil, errors.New("zk: server list must not be empty")
	}

	// 格式优化,缺少端口 2181 会自动补充
	srvs := FormatServers(servers)

	// Randomize the order of the servers to avoid creating hotspots
	stringShuffle(srvs)

	ec := make(chan Event, eventChanSize)
	conn := &Conn{
		dialer:         net.DialTimeout,
		hostProvider:   NewDNSHostProvider(),
		conn:           nil,
		state:          StateDisconnected,
		eventChan:      ec,
		shouldQuit:     make(chan struct{}),
		connectTimeout: 1 * time.Second,
		sendChan:       make(chan *request, sendChanSize),
		requests:       make(map[int32]*request),
		watchers:       make(map[watchPathType][]chan Event),
		passwd:         emptyPassword,
		logger:         DefaultLogger,
		logInfo:        true, // default is true for backwards compatability
		buf:            make([]byte, bufferSize),
		resendZkAuthFn: resendZkAuth,
	}

	// Set provided options.
	for _, option := range options {
		option(conn)
	}

	if err := conn.hostProvider.Init(srvs); err != nil {
		return nil, nil, err
	}

	conn.setTimeouts(int32(sessionTimeout / time.Millisecond))
	// TODO: This context should be passed in by the caller to be the connection lifecycle context.
	ctx := context.Background()

	go func() {
		// 连接操作
		conn.loop(ctx)
		conn.flushRequests(ErrClosing)
		conn.invalidateWatches(ErrClosing)
		close(conn.eventChan)
	}()
	return conn, ec, nil
}
```


```go
func (c *Conn) loop(ctx context.Context) {
	for {
		// 创建连接
		if err := c.connect(); err != nil {
			// c.Close() was called
			return
		}

		// 认证
		err := c.authenticate()
		switch {
		case err == ErrSessionExpired:
			c.logger.Printf("authentication failed: %s", err)
			c.invalidateWatches(err)
		case err != nil && c.conn != nil:
			c.logger.Printf("authentication failed: %s", err)
			c.conn.Close()
		case err == nil:
			if c.logInfo {
				c.logger.Printf("authenticated: id=%d, timeout=%d", c.SessionID(), c.sessionTimeoutMs)
			}
			c.hostProvider.Connected()        // mark success
			c.closeChan = make(chan struct{}) // channel to tell send loop stop

			var wg sync.WaitGroup

			// 发送
			wg.Add(1)
			go func() {
				defer c.conn.Close() // causes recv loop to EOF/exit
				defer wg.Done()

				if err := c.resendZkAuthFn(ctx, c); err != nil {
					c.logger.Printf("error in resending auth creds: %v", err)
					return
				}

				if err := c.sendLoop(); err != nil || c.logInfo {
					c.logger.Printf("send loop terminated: %v", err)
				}
			}()

			// 接收
			wg.Add(1)
			go func() {
				defer close(c.closeChan) // tell send loop to exit
				defer wg.Done()

				var err error
				if c.debugCloseRecvLoop {
					err = errors.New("DEBUG: close recv loop")
				} else {
					err = c.recvLoop(c.conn)
				}
				if err != io.EOF || c.logInfo {
					c.logger.Printf("recv loop terminated: %v", err)
				}
				if err == nil {
					panic("zk: recvLoop should never return nil error")
				}
			}()

			c.sendSetWatches()
			// 等待读写关闭
			wg.Wait()
		}

		c.setState(StateDisconnected)

		select {
		case <-c.shouldQuit:
			c.flushRequests(ErrClosing)
			return
		default:
		}

		if err != ErrSessionExpired {
			err = ErrConnectionClosed
		}
		c.flushRequests(err)

		if c.reconnectLatch != nil {
			select {
			case <-c.shouldQuit:
				return
			case <-c.reconnectLatch:
			}
		}
	}
}

```

### 获取数据

```go
// 获取 znode
func (c *Conn) Get(path string) ([]byte, *Stat, error) {
	if err := validatePath(path, false); err != nil {
		return nil, nil, err
	}

	res := &getDataResponse{}
	_, err := c.request(opGetData, &getDataRequest{Path: path, Watch: false}, res, nil)
	if err == ErrConnectionClosed {
		return nil, nil, err
	}
	return res.Data, &res.Stat, err
}

// 获取 Znode 并设置 watch 
func (c *Conn) GetW(path string) ([]byte, *Stat, <-chan Event, error) {
	if err := validatePath(path, false); err != nil {
		return nil, nil, nil, err
	}

	var ech <-chan Event
	res := &getDataResponse{}
	_, err := c.request(opGetData, &getDataRequest{Path: path, Watch: true}, res, func(req *request, res *responseHeader, err error) {
		if err == nil {
			ech = c.addWatcher(path, watchTypeData)
		}
	})
	if err != nil {
		return nil, nil, nil, err
	}
	return res.Data, &res.Stat, ech, err
}
```


## 参考
- https://zookeeper.apache.org/doc/current/zookeeperOver.html
- [ZAB协议概述与选主流程详解](https://github.com/h2pl/JavaTutorial/blob/master/docs/distributed/practice/%E6%90%9E%E6%87%82%E5%88%86%E5%B8%83%E5%BC%8F%E6%8A%80%E6%9C%AF%EF%BC%9AZAB%E5%8D%8F%E8%AE%AE%E6%A6%82%E8%BF%B0%E4%B8%8E%E9%80%89%E4%B8%BB%E6%B5%81%E7%A8%8B%E8%AF%A6%E8%A7%A3.md)
- [zookeeper 全解](https://blog.csdn.net/General_zy/article/details/129233373)
- [Zookeeper基础篇1-Zookeeper安装和客户端使用](https://juejin.cn/post/7098311052831653919)
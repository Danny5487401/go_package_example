<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [raft 协议在consul中应用](#raft-%E5%8D%8F%E8%AE%AE%E5%9C%A8consul%E4%B8%AD%E5%BA%94%E7%94%A8)
  - [consul 源码分析(v1.9.11)](#consul-%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90v1911)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# raft 协议在consul中应用
hashicorp/raft 是raft协议的一种golang实现，由hashicorp公司实现并开源，已经在consul等软件中使用。
它封装了raft协议的leader选举、log同步等底层实现，基于它能够相对比较容易的构建强一致性的分布式系统.

## consul 源码分析(v1.9.11) 


代理agent接口:服务端和客户端需要实现的接口
```go
// github.com/hashicorp/consul/agent/agent.go

// delegate 定义了代理的公共接口
type delegate interface {
    Encrypted() bool // 数据是否进行了加密
    GetLANCoordinate() (lib.CoordinateSet, error) // 获取局域网内为 Coordinate 角色的服务端
    Leave() error  // 离开集群
    LANMembers() []serf.Member // 局域网内的所有成员
    LANMemberAllSegments() ([]serf.Member, error) // 局域网内所有段区的成员
    LANSegmentMembers(segment string) ([]serf.Member, error) // 局域网内某个段区的所有成员
    LocalMember() serf.Member // 本机成员
    JoinLAN(addrs []string) (n int, err error) // 加入一个或多个段区
    RemoveFailedNode(node string) error // 尝试移除某个异常的节点
    RPC(method string, args interface{}, reply interface{}) error // 远程调用
    SnapshotRPC(args *structs.SnapshotRequest, in io.Reader, out io.Writer, replyFn structs.SnapshotReqlyFn) error // 发起快照存档的远程调用
    Shutdown() error // 关闭代理
    Stats() map[string]map[string]string // 用于获取应用当前状态
}

// Agent 代理
type Agent struct {
    // 代理的运行时配置，支持 hot reload
    config *config.RuntimeConfig
    
    // ... 日志相关
    
    // 内存中收集到的应用、主机等状态信息
    MemSink *metrics.InmemSink
    
    // 代理的公共接口，而配置项则决定代理的角色
    delegate delegate
    
    // 本地策略执行的管理者
    acls *aclManager
        
    // 保证策略执行者的权限，可实时更新，覆盖本地配置文件
    tokens *token.Store
    
    // 存储本地节点、应用、心跳的状态，用于反熵
    State *local.State
    
    // 负责维持本地与远端的状态同步
    sync *ae.StateSyncer
    
    // ...各种心跳包(Monitor/HTTP/TCP/TTL/Docker 等)
    
    // 用于接受其他节点发送过来的事件
    eventCh chan serf.UserEvent
    
    // 用环形队列存储接受到的所有事件，用 index 指向下一个插入的节点
    // 使用读写锁保证数据安全，当一个事件被插入时，会通知 group 中所有的订阅者
    eventBuf    []*UserEvent
    eventIndex  int
    eventLock   sync.RWMutex
    eventNotify NotifyGroup
    
    // 重启，并返回通知是否重启成功
    reloadCh chan chan error
    
    // 关闭代理前的操作
    shutdown     bool
    shutdownCh   chan struct{}
    shutdownLock sync.Mutex
    
    // 添加到局域网成功的回调函数
    joinLANNotifier notifier
    
    // 返回重试加入局域网失败的错误信息
    retryJoinCh chan error
    
    // 并发安全的存储当前所有节点的唯一名称，用于 RPC传输
    endpoints     map[string]string
    endpointsLock sync.RWMutex
    
    // ...为代理提供 DNS/HTTP 的API

    // 追踪当前代理正在运行的所有监控
    watchPlans []*watch.Plan
}

// 决定当前启动的是 server 还是 client 的关键代码在于
func (a *Agent) Start() error {
    // ...
    
    if c.ServerMode {
        server, err := consul.NewServerLogger(consulCfg, a.logger, a.tokens)
        // error handler
        a.delegate = server // 主要差别在这里
    } else {
        client, err := consul.NewClientLogger(consulCfg, a.logger, a.tokens)
        // error handler
        a.delegate = client // 主要差别在这里
    }
    
    a.sync.ClusterSize = func() int { return len(a.delegate.LANMembers()) }
    
    // ...
}
```


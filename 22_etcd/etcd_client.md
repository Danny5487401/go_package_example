# etcd-clientV3源码分析


## rpc定义的proto文件
ETCD核心模块：
- KV:创建、更新、获取和删除键值对。
- Watch:，监视键的更改。
- Lease:实现键值对过期，客户端用来续租、保持心跳。
- Cluster
- Maintenance
- Auth

```protobuf
// /Users/python/go/pkg/mod/go.etcd.io/etcd/api/v3@v3.5.2/etcdserverpb/rpc.proto

syntax = "proto3";
package etcdserverpb;

import "gogoproto/gogo.proto";
import "etcd/api/mvccpb/kv.proto";
import "etcd/api/authpb/auth.proto";

// for grpc-gateway
import "google/api/annotations.proto";

option (gogoproto.marshaler_all) = true;
option (gogoproto.unmarshaler_all) = true;

service KV {
  // Range gets the keys in the range from the key-value store.
  rpc Range(RangeRequest) returns (RangeResponse) {
      option (google.api.http) = {
        post: "/v3/kv/range"
        body: "*"
    };
  }

  // Put puts the given key into the key-value store.
  // A put request increments the revision of the key-value store
  // and generates one event in the event history.
  rpc Put(PutRequest) returns (PutResponse) {
      option (google.api.http) = {
        post: "/v3/kv/put"
        body: "*"
    };
  }

  // DeleteRange deletes the given range from the key-value store.
  // A delete request increments the revision of the key-value store
  // and generates a delete event in the event history for every deleted key.
  rpc DeleteRange(DeleteRangeRequest) returns (DeleteRangeResponse) {
      option (google.api.http) = {
        post: "/v3/kv/deleterange"
        body: "*"
    };
  }

  // Txn processes multiple requests in a single transaction.
  // A txn request increments the revision of the key-value store
  // and generates events with the same revision for every completed request.
  // It is not allowed to modify the same key several times within one txn.
  rpc Txn(TxnRequest) returns (TxnResponse) {
      option (google.api.http) = {
        post: "/v3/kv/txn"
        body: "*"
    };
  }

  // Compact compacts the event history in the etcd key-value store. The key-value
  // store should be periodically compacted or the event history will continue to grow
  // indefinitely.
  rpc Compact(CompactionRequest) returns (CompactionResponse) {
      option (google.api.http) = {
        post: "/v3/kv/compaction"
        body: "*"
    };
  }
}

```

proto 文件是定义服务端和客户端通信接口的标准。包括
- 客户端该传什么样的参数

- 服务端该返回什么参数

- 客户端该怎么调用

- 是阻塞还是非阻塞

- 是同步还是异步

## 连接配置
```go
type Config struct {
	// ETCD服务器地址，注意需要提供ETCD集群所有节点的ip
	Endpoints []string `json:"endpoints"`

	// 设置了此间隔时间，每 AutoSyncInterval 时间ETCD客户端都会
	// 自动向ETCD服务端请求最新的ETCD集群的所有节点列表
	// 默认为0，即不请求
	AutoSyncInterval time.Duration `json:"auto-sync-interval"`

	// 建立底层的GRPC连接的超时时间
	DialTimeout time.Duration `json:"dial-timeout"`

	// 这个配置和下面的 DialKeepAliveTimeout
	// 都是用来打开GRPC提供的 KeepAlive
	// 功能，作用主要是保持底层TCP连接的有效性，
	// 及时发现连接断开的异常。
	// 默认不打开 keepalive
	DialKeepAliveTime time.Duration `json:"dial-keep-alive-time"`

	// 客户端发送 keepalive 的 ping 后，等待服务端的 ping ack 包的时长
	// 超过此时长会报 `translation is closed`
	DialKeepAliveTimeout time.Duration `json:"dial-keep-alive-timeout"`
		


	// 最大可发送字节数，默认为2MB
	// 也就是说，我们ETCD的一条KV记录最大不能超过2MB，
	// 如果要设置超过2MB的KV值，
	// 只修改这个配置也是无效的，因为ETCD服务端那边的限制也是2MB。
        // 需要先修改ETCD服务端启动参数：`--max-request-bytes`，再修改此值。
	MaxCallSendMsgSize int

        // 最大可接收的字节数，默认为`Int.MaxInt32`
        // 一般不需要改动
	MaxCallRecvMsgSize int

	// HTTPS证书配置
	TLS *tls.Config
	
	// 上下文，一般用于取消操作
	ctx.Context

	// 设置此值，会拒绝连接到低版本的ETCD
	// 什么是低版本呢？
	// 写死了，小于v3.2的版本都是低版本。
	RejectOldCluster bool `json:"reject-old-cluster"`

	// GRPC 的连接配置，具体可参考GRPC文档
	DialOptions []grpc.DialOption

	// zap包的Logger配置 
	// ETCD用的日志包就是zap
	Logger *zap.Logger
	LogConfig *zap.Config

    // 也是 keepalive 中的设置，
    // true则表示无论有没有活跃的GRPC连接，都执行ping
    // false的话，没有活跃的连接也就不会发送ping。
    PermitWithoutStream bool `json:"permit-without-stream"`
	

}

```

### 响应头
etcd API 的所有响应都有一个附加的响应标头，其中包括响应的集群元数据：
```protobuf
type ResponseHeader struct {
	// cluster_id is the ID of the cluster which sent the response.
    //产生响应的集群的 ID。
	ClusterId uint64 `protobuf:"varint,1,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
	// member_id is the ID of the member which sent the response.
    //产生响应的成员的 ID。
    //应用服务可以通过 Cluster_ID 和 Member_ID 字段来确保，当前与之通信的正是预期的那个集群或者成 
    // 员。
	MemberId uint64 `protobuf:"varint,2,opt,name=member_id,json=memberId,proto3" json:"member_id,omitempty"`
	// revision is the key-value store revision when the request was applied.
	// For watch progress responses, the header.revision indicates progress. All future events
	// recieved in this stream are guaranteed to have a higher revision number than the
	// header.revision number.
    //产生响应时键值存储的修订版本号。
    //应用服务可以使用修订号字段来获得当前键值存储库最新的修订号。应用程序指定历史修订版以进行查询，如果希望在请求时知道最新修订版，此功能特别有用。
	Revision int64 `protobuf:"varint,3,opt,name=revision,proto3" json:"revision,omitempty"`
	// raft_term is the raft term when the request was applied.
    //产生响应时，成员的 Raft 称谓。
    //应用服务可以使用 Raft_Term 来检测集群何时完成一个新的 leader 选举。
	RaftTerm             uint64   `protobuf:"varint,4,opt,name=raft_term,json=raftTerm,proto3" json:"raft_term,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
```

## 初始化
客户端结构体
```go
 
// Client provides and manages an etcd v3 client session.
type Client struct {
	Cluster      // 向集群里增加 etcd 服务端节点之类，属于管理员操作。
	KV           //我们主要使用的功能，即操作 K-V。
	Lease        //租约相关操作，比如申请一个 TTL=10 秒的租约。
	Watcher      //观察订阅，从而监听最新的数据变化。
	Auth         //管理 etcd 的用户和权限，属于管理员操作。
	Maintenance  //维护 etcd，比如主动迁移 etcd 的 leader 节点，属于管理员操作。
 
	conn *grpc.ClientConn
 
	cfg      Config
	creds    grpccredentials.TransportCredentials
	resolver *resolver.EtcdManualResolver
	mu       *sync.RWMutex
 
	ctx    context.Context
	cancel context.CancelFunc
 
	// Username is a user name for authentication.
	Username string
	// Password is a password for authentication.
	Password        string
	authTokenBundle credentials.Bundle
 
	callOpts []grpc.CallOption
 
	lgMu *sync.RWMutex
	lg   *zap.Logger
}
```

```go
func newClient(cfg *Config) (*Client, error) {
    // ... 
	// 1. 初始化了一个client的实例
	ctx, cancel := context.WithCancel(baseCtx)
	client := &Client{
		conn:     nil,
		cfg:      *cfg,
		creds:    creds,
		ctx:      ctx,
		cancel:   cancel,
		mu:       new(sync.RWMutex),
		callOpts: defaultCallOpts,
		lgMu:     new(sync.RWMutex),
	}
	// ...
    // 2. 创建解析器
    // resolver（解析器）其实是grpc中的概念，比如：DNS解析器，域名转化为真实的ip；服务注册中心，也是一种把服务名转化为真实ip的解析服务。
	client.resolver = resolver.New(cfg.Endpoints...)

	if len(cfg.Endpoints) < 1 {
		client.cancel()
		return nil, fmt.Errorf("at least one Endpoint is required in client config")
	}
	
	// 3. dialWithBalancer() 建立了到ETCD的服务端链接
	conn, err := client.dialWithBalancer()
	if err != nil {
		client.cancel()
		client.resolver.Close()
		// TODO: Error like `fmt.Errorf(dialing [%s] failed: %v, strings.Join(cfg.Endpoints, ";"), err)` would help with debugging a lot.
		return nil, err
	}
	client.conn = conn

	// 4. 是做一些功能接口的初始化
	client.Cluster = NewCluster(client)
	client.KV = NewKV(client)
	client.Lease = NewLease(client)
	client.Watcher = NewWatcher(client)
	client.Auth = NewAuth(client)
	client.Maintenance = NewMaintenance(client)

	//get token with established connection
	ctx, cancel = client.ctx, func() {}
	if client.cfg.DialTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, client.cfg.DialTimeout)
	}
	// 这个token是我们开启了ETCD的账号密码功能后，通过账号密码获取到了token，然后才能访问ETCD提供的GRPC接口
	err = client.getToken(ctx)
	if err != nil {
		client.Close()
		cancel()
		//TODO: Consider fmt.Errorf("communicating with [%s] failed: %v", strings.Join(cfg.Endpoints, ";"), err)
		return nil, err
	}
	cancel()

	if cfg.RejectOldCluster {
		if err := client.checkVersion(); err != nil {
			client.Close()
			return nil, err
		}
	}

	go client.autoSync()
	return client, nil
}
```

连接grpc
```go
func (c *Client) dial(creds grpccredentials.TransportCredentials, dopts ...grpc.DialOption) (*grpc.ClientConn, error) {
	//  首先，ETCD通过这行代码，向GRPC框架加入了一些自己的
	// 配置，比如：KeepAlive特性（配置里提到的配置项）、
	// TLS证书配置、还有最重要的重试策略
	
	opts, err := c.dialSetupOpts(creds, dopts...)
	if err != nil {
		return nil, fmt.Errorf("failed to configure dialer: %v", err)
	}
	if c.Username != "" && c.Password != "" {
		// 绑定token用
		c.authTokenBundle = credentials.NewBundle(credentials.Config{})
		opts = append(opts, grpc.WithPerRPCCredentials(c.authTokenBundle.PerRPCCredentials()))
	}

	opts = append(opts, c.cfg.DialOptions...)

	dctx := c.ctx
	if c.cfg.DialTimeout > 0 {
		var cancel context.CancelFunc
		dctx, cancel = context.WithTimeout(c.ctx, c.cfg.DialTimeout)
		defer cancel() // TODO: Is this right for cases where grpc.WithBlock() is not set on the dial options?
	}
	target := fmt.Sprintf("%s://%p/%s", resolver.Schema, c, authority(c.Endpoints()[0]))
	
	//  最终调用grpc.DialContext()建立连接
	conn, err := grpc.DialContext(dctx, target, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

```
设置连接选项
```go
// dialSetupOpts gives the dial opts prior to any authentication.
func (c *Client) dialSetupOpts(creds grpccredentials.TransportCredentials, dopts ...grpc.DialOption) (opts []grpc.DialOption, err error) {
	if c.cfg.DialKeepAliveTime > 0 {
		params := keepalive.ClientParameters{
			Time:                c.cfg.DialKeepAliveTime,
			Timeout:             c.cfg.DialKeepAliveTimeout,
			PermitWithoutStream: c.cfg.PermitWithoutStream,
		}
		opts = append(opts, grpc.WithKeepaliveParams(params))
	}
	opts = append(opts, dopts...)

	if creds != nil {
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// Interceptor retry and backoff.
	// TODO: Replace all of clientv3/retry.go with RetryPolicy:
	// https://github.com/grpc/grpc-proto/blob/cdd9ed5c3d3f87aef62f373b93361cf7bddc620d/grpc/service_config/service_config.proto#L130
	rrBackoff := withBackoff(c.roundRobinQuorumBackoff(defaultBackoffWaitBetween, defaultBackoffJitterFraction))
	opts = append(opts,
		// Disable stream retry by default since go-grpc-middleware/retry does not support client streams.
		// Streams that are safe to retry are enabled individually.
		// Client端的Stream重试不被支持
		grpc.WithStreamInterceptor(c.streamClientInterceptor(withMax(0), rrBackoff)),
		grpc.WithUnaryInterceptor(c.unaryClientInterceptor(withMax(defaultUnaryMaxRetries), rrBackoff)),
	)

	return opts, nil
}

```

## KV 查询 Get
```go
Get(ctx context.Context, key string, opts ...OpOption) (*GetResponse, error)
```
OpOption 为可选的函数传参，
- 传参为WithRange(end)时，Get 将返回 [key，end) 范围内的键；
- 传参为 WithFromKey() 时，Get 返回大于或等于 key 的键；
- 当通过 rev> 0 传递 WithRev(rev) 时，Get 查询给定修订版本的键；如果压缩了所查找的修订版本，则返回请求失败，并显示 ErrCompacted。
- 传递 WithLimit(limit) 时，返回的 key 数量受 limit 限制；传参为 WithSort 时，将对键进行排序。
# etcd-clientV3源码分析


## rpc定义的proto文件
ETCD核心模块：KV,Watch,Lease,Cluster,Maintenance,Auth

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

service Watch {
  // Watch watches for events happening or that have happened. Both input and output
  // are streams; the input stream is for creating and canceling watchers and the output
  // stream sends events. One watch RPC can watch on multiple key ranges, streaming events
  // for several watches at once. The entire event history can be watched starting from the
  // last compaction revision.
  rpc Watch(stream WatchRequest) returns (stream WatchResponse) {
      option (google.api.http) = {
        post: "/v3/watch"
        body: "*"
    };
  }
}
```

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

## 初始化
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
#Client源码分析
因为gRPC没有提供服务注册，服务发现的功能，所以需要开发者自己编写服务发现的逻辑：也就是Resolver——解析器。

在得到了解析的结果，也就是一连串的IP地址之后，需要对其中的IP进行选择，也就是Balancer

连接对象
```go
type ClientConn struct {
	ctx    context.Context
	cancel context.CancelFunc

	target       string
	parsedTarget resolver.Target
	authority    string
	dopts        dialOptions
	csMgr        *connectivityStateManager

	balancerBuildOpts balancer.BuildOptions
	blockingpicker    *pickerWrapper

	safeConfigSelector iresolver.SafeConfigSelector

	mu              sync.RWMutex
	resolverWrapper *ccResolverWrapper
	sc              *ServiceConfig
	conns           map[*addrConn]struct{}
	// Keepalive parameter can be updated if a GoAway is received.
	mkp             keepalive.ClientParameters
	curBalancerName string
	balancerWrapper *ccBalancerWrapper
	retryThrottler  atomic.Value

	firstResolveEvent *grpcsync.Event

	channelzID int64 // channelz unique identification number
	czData     *channelzData

	lceMu               sync.Mutex // protects lastConnectionError
	lastConnectionError error
}
```

##一. grpc.Dial 方法实际上是对于 grpc.DialContext 的封装，区别在于 ctx 是直接传入 context.Background。

首先要做的就是调用Dial或DialContext函数来初始化一个clientConn对象，而resolver是这个连接对象的一个重要的成员，
所以我们首先看一看clientConn对象创建过程中，resolver是怎么设置进去的。

客户端启动时，一定会调用grpc的Dial或DialContext函数来创建连接，而这两个函数都需要传入一个名为target的参数，target，就是连接的目标，也就是server了，
接下来，我们就看一看，DialContext函数里是如何处理这个target的.

首先，创建了一个clientConn对象，并把target赋给了对象中的target：

```go
func DialContext(ctx context.Context, target string, opts ...DialOption) (conn *ClientConn, err error) {

	// 1.创建ClientConn结构体
	cc := &ClientConn{
		target:            target, //将target连接对象赋给了对象中的target
		//...
	}
	
	// 2.解析target
	cc.parsedTarget = grpcutil.ParseTarget(cc.target, cc.dopts.copts.Dialer != nil)
	
	// 3.根据解析的target找到合适的resolverBuilder
	resolverBuilder := cc.getResolver(cc.parsedTarget.Scheme)
	
	// 4.创建Resolver
	rWrapper, err := newCCResolverWrapper(cc, resolverBuilder)
	
	// 5.完事
	return cc, nil

```

也就是在根据解析的结果，包括scheme和endpoint这两个参数，获取一个resolver的builder
```go
func (cc *ClientConn) getResolver(scheme string) resolver.Builder {
    // 先查看是否在配置中存在resolver
	for _, rb := range cc.dopts.resolvers {
		if scheme == rb.Scheme() {
			return rb
		}
	}
    // 如果配置中没有相应的resolver，再从注册的resolver中寻找
	return resolver.Get(scheme)
}
```

Get函数是通过m这个map，去查找有没有scheme对应的resolver的builder，那么m这个map是什么时候插入的值呢？这个在resolver的Register函数里
```go
func Register(b Builder) {
	m[b.Scheme()] = b
}
```

那么谁会去调用这个Register函数，向map中写入resolver呢 ？

    有两个人会去调，首先，grpc实现了一个默认的解析器，也就是"passthrough"，这个看名字就理解了，就是透传，所谓透传就是，什么都不做，那么什么时候需要透传呢？
    当你调用DialContext的时候，如果传入的target本身就是一个ip+port，这个时候，自然就不需要再解析什么了。
    那么"passthrough"对应的这个默认的解析器是什么时候注册到m这个map中的呢？这个调用在passthrough包的init函数里
```go

func init() {
	resolver.Register(&passthroughBuilder{})
}
```

ResolverWrapper的创建
```go
func newCCResolverWrapper(cc *ClientConn, rb resolver.Builder) (*ccResolverWrapper, error) {
	ccr := &ccResolverWrapper{
		cc:   cc,
		done: grpcsync.NewEvent(),
	}
  
  // 根据传入的Builder，创建resolver，并放入wrapper中
  ccr.resolver, err = rb.Build(cc.parsedTarget, ccr, rbo)
 	return ccr, nil
}
```

为了解耦Resolver和Balancer，我们希望能够有一个中间的部分，接收到Resolver解析到的地址，然后对它们进行负载均衡。
因此，在接下来的代码阅读过程中，我们可以带着这个问题：Resolver和Balancer的通信过程是什么样的


在创建Resolver的时候，我们需要在Build方法里面初始化Resolver的各种状态。并且，因为Build方法中有一个target的参数，我们会在创建Resolver的时候，需要对这个target进行解析。

也就是说，创建Resolver的时候，会进行第一次的域名解析。并且，这个解析过程，是由开发者自己设计的。

到了这里我们会自然而然的接着考虑，解析之后的结果应该保存为什么样的数据结构，又应该怎么去将这个结果传递下去呢？

我们拿最简单的passthroughResolver来举例
```go
func (*passthroughBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &passthroughResolver{
		target: target,
		cc:     cc,
	}
  // 创建Resolver的时候，进行第一次的解析
	r.start()
	return r, nil
}

// 对于passthroughResolver来说，正如他的名字，直接将参数作为结果返回
func (r *passthroughResolver) start() {
	r.cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: r.target.Endpoint}}})
```
我们可以看到，对于一个Resolver，需要将解析出的地址，传入resolver.State中，然后调用r.cc.UpdateState方法。

那么这个r.cc.UpdateState又是什么呢？

他就是我们上面提到的ccResolverWrapper。

这个时候逻辑就很清晰了，gRPC的ClientConn通过调用ccResolverWrapper来进行域名解析，而具体的解析过程则由开发者自己决定。在解析完毕后，将解析的结果返回给ccResolverWrapper


balancer的选择
我们因此也可以进行推测：在ccResolverWrapper中，会将解析出的结果以某种形式传递给Balancer
```go
func (ccr *ccResolverWrapper) UpdateState(s resolver.State) error {
    //...
	// 将Resolver解析的最新状态保存下来
	ccr.curState = s

	//...
	
	// 对状态进行更新
	if err := ccr.cc.updateResolverState(ccr.curState, nil); err == balancer.ErrBadResolverState {
        return balancer.ErrBadResolverState
    }
    return nil
}
```
总结

    其主要功能是创建与给定目标的客户端连接，其承担了以下职责

    1.初始化 ClientConn
    2.初始化（基于进程 LB）负载均衡配置
    3.初始化 channelz
    4.初始化重试规则和客户端一元/流式拦截器
    5.初始化协议栈上的基础信息
    6.相关 context 的超时控制
    7.初始化并解析地址信息
    8.创建与服务端之间的连接

我们可以有几个核心方法一直在等待/处理信号
```go
func (ac *addrConn) connect()
func (ac *addrConn) resetTransport()
func (ac *addrConn) createTransport(addr resolver.Address, copts transport.ConnectOptions, connectDeadline time.Time)
func (ac *addrConn) getReadyTransport()
```

##二. 实例化
```go
type GreeterClient interface {
	// Sends a greeting
	//  rpc SayHello (HelloRequest) returns (HelloReply) {}
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type greeterClient struct {
	cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}
```
##三. 调用

底层http2连接对应的是一个grpc的stream，而stream的创建有两种方式

	1.一种就是我们主动去创建一个stream池，这样当有请求需要发送时，我们可以直接使用我们创建好的stream，

	2.除了我们自己创建，我们使用protoc为我们生成的客户端接口里，也会为我们实现stream的创建，也就是说这个完全是可以不用我们自己费心的

```go
// Invoke sends the RPC request on the wire and returns after response is
// received.  This is typically called by generated code.
//
// All errors returned by Invoke are compatible with the status package.
func (cc *ClientConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...CallOption) error {
	// allow interceptor to see all applicable call options, which means those
	// configured as defaults from dial option as well as per-call options
	opts = combine(cc.dopts.callOptions, opts)

	if cc.dopts.unaryInt != nil {
		return cc.dopts.unaryInt(ctx, method, args, reply, cc, invoke, opts...)
	}
	return invoke(ctx, method, args, reply, cc, opts...)
}
```
在没有设置拦截器的情况下，会直接调invoke

```go
func invoke(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, opts ...CallOption) error {

	cs, err := newClientStream(ctx, unaryStreamDesc, cc, method, opts...)
	if err != nil {
		return err
	}
	if err := cs.SendMsg(req); err != nil {
		return err
	}
	return cs.RecvMsg(reply)
}
/*
newClientStream：获取传输层 Transport 并组合封装到 ClientStream 中返回，在这块会涉及负载均衡、超时控制、 Encoding、 Stream 的动作，与服务端基本一致的行为。
cs.SendMsg：发送 RPC 请求出去，但其并不承担等待响应的功能。
cs.RecvMsg：阻塞等待接受到的 RPC 方法响应结果。
*/
```
##四。关闭链接
```go

func (cc *ClientConn) Close() error {
    defer cc.cancel()
    ...
    cc.csMgr.updateState(connectivity.Shutdown)
    ...
    cc.blockingpicker.close()
    if rWrapper != nil {
        rWrapper.close()
    }
    if bWrapper != nil {
        bWrapper.close()
    }

    for ac := range conns {
        ac.tearDown(ErrClientConnClosing)
    }
    if channelz.IsOn() {
        ...
        channelz.AddTraceEvent(cc.channelzID, ted)
        channelz.RemoveEntry(cc.channelzID)
    }
    return nil
}


```
该方法会取消 ClientConn 上下文，同时关闭所有底层传输。涉及如下

	* Context Cancel
	* 清空并关闭客户端连接
	* 清空并关闭解析器连接
	* 清空并关闭负载均衡连接
	* 添加跟踪引用
	* 移除当前通道信息
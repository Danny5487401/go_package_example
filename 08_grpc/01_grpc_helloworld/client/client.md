#client源码分析

一,grpc.Dial 方法实际上是对于 grpc.DialContext 的封装，区别在于 ctx 是直接传入 context.Background。

```go
func DialContext(ctx context.Context, target string, opts ...DialOption) (conn *ClientConn, err error) {
    cc := &ClientConn{
        target:            target,
        csMgr:             &connectivityStateManager{},
        conns:             make(map[*addrConn]struct{}),
        dopts:             defaultDialOptions(),
        blockingpicker:    newPickerWrapper(),
        czData:            new(channelzData),
        firstResolveEvent: grpcsync.NewEvent(),
    }
    ...
    chainUnaryClientInterceptors(cc)
    chainStreamClientInterceptors(cc)

    ...
}
```

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

二，实例化
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
三。调用
	实际上调用的还是 grpc.invoke 方法。

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
四。关闭链接
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

	*Context Cancel
	*清空并关闭客户端连接
	*清空并关闭解析器连接
	*清空并关闭负载均衡连接
	*添加跟踪引用
	*移除当前通道信息
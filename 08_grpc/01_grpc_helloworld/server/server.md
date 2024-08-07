<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [grpc源码分析](#grpc%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [一. grpc.NewServer()分析](#%E4%B8%80-grpcnewserver%E5%88%86%E6%9E%90)
  - [二. 注册](#%E4%BA%8C-%E6%B3%A8%E5%86%8C)
  - [三. s.Serve(lis)](#%E4%B8%89-sservelis)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# grpc源码分析
```go
// UnimplementedGreeterServer can be embedded to have forward compatible implementations.
type UnimplementedGreeterServer struct {
}

func (*UnimplementedGreeterServer) SayHello(context.Context, *HelloRequest) (*HelloReply, error) {
    return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
```

这里pb.UnimplementedGreeterServer被嵌入了server结构，所以即使没有实现SayHello方法，编译也能通过。

但是，我们通常要强制server在编译期就必须实现对应的方法，所以生产中建议不嵌入。





## 一. grpc.NewServer()分析
```go
func NewServer(opt ...ServerOption) *Server {
    opts := defaultServerOptions
    for _, o := range opt {
        o(&opts)
    }
    s := &Server{
        lis:    make(map[net.Listener]bool), // 监听地址列表
        opts:   opts,  //服务选项，这块包含 Credentials、Interceptor 以及一些基础配置
        conns:  make(map[io.Closer]bool),  //客户端连接句柄列表
        m:      make(map[string]*service),  //服务信息映射
        quit:   make(chan struct{}),  //退出信号
        done:   make(chan struct{}),  //完成信号
        czData: new(channelzData),  //用于存储 ClientConn，addrConn 和 Server 的channelz 相关数据。
    }
    s.cv = sync.NewCond(&s.mu)  //当优雅退出时，会等待这个信号量，直到所有 RPC 请求都处理并断开才会继续处理
    ...

    return s
}
```

1. 入参为选项参数options
2. 自带一组defaultServerOptions，最大发送size、最大接收size、连接超时、发送缓冲、接收缓冲
3. s.cv = sync.NewCond(&s.mu) 条件锁，用于关闭连接
4. 全局参数 EnableTraciing ，会调用golang.org/x/net/trace 这个包


## 二. 注册
```go
func RegisterGreeterServer(s *grpc.Server, srv GreeterServer) {
    s.RegisterService(&_Greeter_serviceDesc, srv)
}

//Greeter_serviceDesc解释

var _Greeter_serviceDesc = grpc.ServiceDesc{
    ServiceName: "Greeter",  //服务名称
    HandlerType: (*GreeterServer)(nil),  //服务接口，用于检查用户提供的实现是否满足接口要求
    Methods: []grpc.MethodDesc{
        //一元方法集，注意结构内的 Handler 方法，其对应最终的 RPC 处理方法，在执行 RPC 方法的阶段会使用
        {
            MethodName: "SayHello",
            Handler:    _Greeter_SayHello_Handler,
        },
    },
    Streams:  []grpc.StreamDesc{},  //流式方法集
    Metadata: "01_grpc_helloworld/proto/helloworld.proto",  //元数据，是一个描述数据属性的东西
}
```

## 三. s.Serve(lis)

1. listener 放到内部的map中
2. for循环，进行tcp连接，这一部分和http源码中的ListenAndServe极其类似
```go
// /Users/python/go/pkg/mod/google.golang.org/grpc@v1.47.0/server.go
func (s *Server) Serve(lis net.Listener) error {
	......
	s.serve = true
	......
	for {
		rawConn, err := lis.Accept()
		if err != nil {
			......
		}
		......
		s.serveWG.Add(1)
		go func() {
			s.handleRawConn(lis.Addr().String(), rawConn)
			s.serveWG.Done()
		}()

}
```
3. 在协程中进行handleRawConn
```go
func (s *Server) handleRawConn(lisAddr string, rawConn net.Conn) {
	// 如果 gRPC server 已经关闭, 那么同样关闭这个 TCP 连接
	if s.quit.HasFired() {
		rawConn.Close()
		return
	}
	// 为这个 TCP 连接设置 deadline
	rawConn.SetDeadline(time.Now().Add(s.opts.connectionTimeout))

	// Finish handshaking (HTTP2)
	// RPC 连接阶段, server 和 client 之间进行 HTTP2 的握手
	st := s.newHTTP2Transport(rawConn)
	rawConn.SetDeadline(time.Time{})
	if st == nil {
		return
	}

	if !s.addConn(lisAddr, st) {
		return
	}
	// RPC 交互阶段, 在新的 goroutine 中处理来自 client 的数据
	go func() {
		s.serveStreams(st)
		s.removeConn(lisAddr, st)
	}()
}
```
4. 将tcp连接封装对应的creds认证信息

5. 新建newHTTP2Transport传输层连接
```go
func (s *Server) newHTTP2Transport(c net.Conn) transport.ServerTransport {
	// 组装 ServerConfig
	config := &transport.ServerConfig{
		MaxStreams:            s.opts.maxConcurrentStreams,
		ConnectionTimeout:     s.opts.connectionTimeout,
		......
	}
	
	// 根据 config 的配置信息, 和 client 进行 HTTP2 的握手
	st, err := transport.NewServerTransport(c, config)
	if err != nil {
		......
	}

	return st
}
```
transport.NewServerTransport 是一个很长的函数, 我们下面一点点来看.



6. 在协程中进行serveStreams，而http1这里为阻塞的
7. 函数HandleStreams中参数为2个函数，前者为处理请求，后者用于trace
8. 进入handleStream，前半段被拆为service，后者为method，通过map查找
9. method在processUnaryRPC处理，stream在processStreamingRPC处理，这两块内部就比较复杂了，涉及到具体的算法，以后有时间细读



package main

import (
	"go_grpc_example/08_grpc/01_grpc_helloworld/proto"

	"google.golang.org/grpc"

	"context"
	"net"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{
		Message: "hello," + request.Name,
	}, nil
}

func main() {
	//1。初始化
	g := grpc.NewServer()

	// 2.注册服务 service放在 m map[string]*service 中
	proto.RegisterGreeterServer(g, &Server{})

	// 3.gRPC的应用层是基于HTTP2的，所以这里不出意外，监听的是tcp端口
	lis, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	_ = g.Serve(lis)
}

/*
grpc源码分析
	// UnimplementedGreeterServer can be embedded to have forward compatible implementations.
	type UnimplementedGreeterServer struct {
	}

	func (*UnimplementedGreeterServer) SayHello(context.Context, *HelloRequest) (*HelloReply, error) {
		return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
	}
	这里pb.UnimplementedGreeterServer被嵌入了server结构，所以即使没有实现SayHello方法，编译也能通过。

	但是，我们通常要强制server在编译期就必须实现对应的方法，所以生产中建议不嵌入。

	一。grpc.NewServer()分析中
	1。入参为选项参数options
	2。自带一组defaultServerOptions，最大发送size、最大接收size、连接超时、发送缓冲、接收缓冲
	3，s.cv = sync.NewCond(&s.mu) 条件锁，用于关闭连接
	4。全局参数 EnableTraciing ，会调用golang.org/x/net/trace 这个包

	二。s.Serve(lis)
	1.listener 放到内部的map中
	2.for循环，进行tcp连接，这一部分和http源码中的ListenAndServe极其类似
	3.在协程中进行handleRawConn
	4.将tcp连接封装对应的creds认证信息
	5.新建newHTTP2Transport传输层连接
	6.在协程中进行serveStreams，而http1这里为阻塞的
	7.函数HandleStreams中参数为2个函数，前者为处理请求，后者用于trace
	8.进入handleStream，前半段被拆为service，后者为method，通过map查找
	9.method在processUnaryRPC处理，stream在processStreamingRPC处理，这两块内部就比较复杂了，涉及到具体的算法，以后有时间细读

*/

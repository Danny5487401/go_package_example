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

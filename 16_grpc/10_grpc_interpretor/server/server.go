package main

import (
	"context"
	"fmt"
	"go_test_project/16_grpc/10_grpc_interpretor/proto"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{
		Message: "hello," + request.Name,
	}, nil
}

func main() {
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 继续处理请求
		fmt.Println("接收到新请求,开始时间")
		res, err := handler(ctx, req)
		fmt.Println("请求处理完成，结束时间")
		return res, err
	}
	opt := grpc.UnaryInterceptor(interceptor)
	g := grpc.NewServer(opt)

	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		panic("faild to listen:" + err.Error())
	}
	_ = g.Serve(lis)
}

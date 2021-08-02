package main

import (
	"context"
	"fmt"
	"go_grpc_example/08_grpc/06_grpc_interpretor/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			fmt.Println("get metadata errors")
			return resp, status.Errorf(codes.Unauthenticated, "无token认证信息")
		}
		/*
			timestamp [2021-03-20 15:10:53]
			:authority [127.0.0.1:9000]
			content-type [application/grpc]
		*/
		var (
			appid  string
			appkey string
		)
		if val1, ok := md["appid"]; ok {
			appid = val1[0]
		}
		if val1, ok := md["appkey"]; ok {
			appkey = val1[0]
		}

		if appid != "123456" || appkey != "i am a key" {
			return resp, status.Errorf(codes.Unauthenticated, "认证信息错误")
		}

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

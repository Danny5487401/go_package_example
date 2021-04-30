package main

import (
	"context"
	"net"

	"go_test_project/13_grpc_error_test/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc"
)

type Server struct {

}

func (s *Server)SayHello(ctx context.Context,request *proto.HelloRequest)  (*proto.HelloReply,error){
	return nil,status.Errorf(codes.AlreadyExists,"账号已存在%s",request.Name)
	//return &proto.HelloReply{
	//	Message:"hello," + request.Name,
	//},nil
}

func main()  {
	g := grpc.NewServer()

	proto.RegisterGreeterServer(g,&Server{})
	lis,err := net.Listen("tcp","127.0.0.1:9000")
	if err != nil{
		panic("failed to listen:" + err.Error())
	}
	_ = g.Serve(lis)
}

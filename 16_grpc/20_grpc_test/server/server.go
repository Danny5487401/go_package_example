package main

import (
	"golang.org/x/net/context"
	"net"

	"go_grpc_exam"
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
	g := grpc.NewServer()

	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	_ = g.Serve(lis)
}

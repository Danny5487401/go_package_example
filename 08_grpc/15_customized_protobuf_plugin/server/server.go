package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Danny5487401/go_package_example/08_grpc/15_customized_protobuf_plugin/helloworld_protobuf"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello1(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v,age:%v", in.GetName(), in.GetAge())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
func (s *server) SayHello2(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v,age:%v", in.GetName(), in.GetAge())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//1。初始化
	s := grpc.NewServer()
	// 2.注册服务 service放在 m map[string]*service 中
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

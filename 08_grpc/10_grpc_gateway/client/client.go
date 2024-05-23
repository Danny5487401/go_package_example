package main

import (
	"context"
	pb "github.com/Danny5487401/go_package_example/08_grpc/10_grpc_gateway/proto_without_buf/helloworld"
	"google.golang.org/grpc"
	"log"
)

const (
	address     = "localhost:50051"
	defaultName = "danny"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: defaultName})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}

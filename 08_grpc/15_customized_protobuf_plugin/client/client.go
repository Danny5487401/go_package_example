package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/Danny5487401/go_package_example/08_grpc/15_customized_protobuf_plugin/helloworld_protobuf"
	"google.golang.org/grpc"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	opts := make([]grpc.DialOption, 0)
	opts = append(opts, grpc.WithInsecure())                               // 不安全链接
	opts = append(opts, grpc.WithBlock(), grpc.WithTimeout(time.Second*2)) // 保证连接上,连接不上会一直阻塞，加个时间限制
	conn, err := grpc.DialContext(context.Background(), *addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello2(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

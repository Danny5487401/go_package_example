package main

import (
	"context"
	"fmt"
	"time"

	"go_grpc_example/14_metadata_test/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := proto.NewGreeterClient(conn)

	md := metadata.Pairs("timestamp", time.Now().Format("2006-01-02 15:04:05"),
		"name", "Joy")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	r, err := c.SayHello(ctx, &proto.HelloRequest{
		Name: "danny",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}

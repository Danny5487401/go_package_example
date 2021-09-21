package main

import (
	"fmt"
	"go_grpc_example/08_grpc/01_grpc_helloworld/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	//1.新建一个conn连接，
	conn, err := grpc.Dial("127.0.0.1:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := proto.NewGreeterClient(conn)

	r, err := c.SayHello(context.Background(), &proto.HelloRequest{
		Name: "danny",
	}) // 核心调用的是 Invoke 方法，具体实现要看grpc.ClientConn中
	// grpc.ClientConn中实现了Invoke方法，在call.go文件中，详情都在invoke

	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}

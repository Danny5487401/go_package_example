package main

import (
	"context"
	"fmt"
	"go_grpc_example/08_grpc/13_grpc_error_test/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	// 超时机制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err = c.SayHello(ctx, &proto.HelloRequest{
		Name: "danny",
	})

	//_,err = c.SayHello(context.Background(),&proto.HelloRequest{
	//	Name: "danny",
	//})
	//if err != nil{
	//	panic(err)
	//}
	if err != nil {
		sta, ok := status.FromError(err)
		if !ok {
			panic("解析error失败")
		}
		fmt.Println(sta.Message())
		fmt.Println(sta.Code())

	}
	//fmt.Println(r.Message)
}

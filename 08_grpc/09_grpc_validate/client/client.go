package main

import (
	"context"
	"fmt"
	"github.com/Danny5487401/go_package_example/08_grpc/09_grpc_validate/proto"
	"google.golang.org/grpc"
)

type customCredential struct{}

func main() {
	var opts []grpc.DialOption

	//opts = append(opts, grpc.WithUnaryInterceptor(interceptor))
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	//rsp, _ := c.Search(context.Background(), &empty.Empty{})
	rsp, err := c.SayHello(context.Background(), &proto.Person{
		Id:     1000,
		Email:  "540021730@qq.com",
		Mobile: "18621815637",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

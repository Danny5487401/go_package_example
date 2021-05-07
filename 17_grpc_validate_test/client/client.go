package  main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"go_test_project/17_grpc_validate_test/proto"
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
		Id: 1000,
		Email: "540021730@qq.com",
		Mobile: "18621815637",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}
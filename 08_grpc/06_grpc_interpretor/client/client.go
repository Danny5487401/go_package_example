package main

import (
	"context"
	"fmt"
	"go_grpc_example/08_grpc/06_grpc_interpretor/proto"
	"time"

	"google.golang.org/grpc"
)

func main() {
	interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		fmt.Printf("耗时%s\n", time.Since(start))
		return err
	}
	//opt := grpc.WithUnaryInterceptor(interceptor)
	//conn,err := grpc.Dial("127.0.0.1:9000",grpc.WithInsecure(),opt)  //grpc.WithInsecure(),opt两个相同

	// 方法二
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(interceptor))
	conn, err := grpc.Dial("127.0.0.1:9000", opts...)

	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := proto.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{
		Name: "danny",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}

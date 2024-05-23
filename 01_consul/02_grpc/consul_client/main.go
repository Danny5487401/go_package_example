package main

import (
	"context"
	"fmt"
	pb "github.com/Danny5487401/go_package_example/08_grpc/01_grpc_helloworld/proto"
	"log"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important  consul实现了下面的两个接口
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.DialContext(
		context.Background(),
		// whoami： 是名称
		"consul://tencent.danny.games:8500/helloworld",
		grpc.WithInsecure(),
		//grpc.WithBlock(),

		//grpc.WithBalancerName() 已经弃用的方法
		//关于serverConfig https://github.com/grpc/grpc/blob/master/doc/service_config.md
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //json格式
	)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "danny"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	fmt.Sprintln(r.GetMessage())

}

package main

import (
	"context"
	"fmt"
	"go_package_example/08_grpc/05_grpc_load_balance/proto"
	"log"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important  consul实现了下面的两个接口
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(
		// whoami： 是名称
		"consul://192.168.16.111:8500/user-srv?wait=14s&tag=srv",
		grpc.WithInsecure(),

		//grpc.WithBalancerName() 已经弃用的方法
		//关于serverConfig https://github.com/grpc/grpc/blob/master/doc/service_config.md
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// 模拟十次同时快速请求，保证在round_robin在一次
	for i := 0; i < 10; i++ {
		userSrvClient := proto.NewUserClient(conn)
		rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
			Pn:    2,
			PSize: 2,
		})
		if err != nil {
			panic(err)
		}
		for index, data := range rsp.Data {
			str := fmt.Sprintf("索引是%v，对应数据是%v", index, data)
			fmt.Println(str)
		}
	}
}

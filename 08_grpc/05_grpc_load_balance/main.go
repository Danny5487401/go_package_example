package main

import (
	"context"
	"fmt"
	"go_grpc_example/08_grpc/05_grpc_load_balance/proto"
	"log"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important  consul实现了下面的两个接口
	"google.golang.org/grpc"
)

/*
gRPC已提供了简单的负载均衡策略（如：Round Robin），我们只需实现它提供的Builder和Resolver接口，就能完成gRPC客户端负载均衡。

	Builder接口：创建一个resolver（本文称之服务发现），用于监视名称解析更新。
type Builder interface {
	Build(target Target, cc ClientConn, opts BuildOption) (Resolver, error)//为给定目标创建一个新的resolver，当调用grpc.Dial()时执行
	Scheme() string //返回此resolver支持的方案
}
	Resolver接口：监视指定目标的更新，包括地址更新和服务配置更新。
type Resolver interface {
	ResolveNow(ResolveNowOption) // 被 gRPC 调用，以尝试再次解析目标名称。只用于提示，可忽略该方法。
	Close()
}


*/

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

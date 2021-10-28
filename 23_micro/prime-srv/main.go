package main

//
//import (
//	"github.com/micro/go-micro"
//	"github.com/micro/go-micro/util/log"
//
//	"go_grpc_example/23_micro/prime-srv/handler"
//	"go_grpc_example/23_micro/proto/prime"
//)
//
///*
//service端服务
//*/
//func main() {
//	// 声明服务
//	srv := micro.NewService(
//		micro.Name("Micro-prime"))
//	//初始化
//	srv.Init(
//		// 定义钩子函数
//		micro.BeforeStart(func() error {
//			log.Log("[srv]启动之前")
//		}),
//		micro.AfterStart(func() error {
//			log.Log("[srv]启动之后")
//		}))
//	//挂载接口
//	_ = prime.RegisterPrimeHandler(srv.Server(), handler.Handler())
//	//启动
//	if err := srv.Run(); err != nil {
//		panic(err)
//	}
//}

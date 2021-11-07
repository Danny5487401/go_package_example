package main

//
///*
//go-micro v1.18.0
//	注册中心：默认etcd，不是msdn
//*/
//
//import (
//	"encoding/json"
//	"github.com/micro/go-micro"
//	"github.com/micro/go-micro/broker"
//	"github.com/micro/go-micro/client"
//	"github.com/micro/go-micro/server"
//	"github.com/micro/go-micro/util/log"
//
//	logProto "go_grpc_example/23_micro/proto/log"
//	"go_grpc_example/23_micro/proto/sum"
//	"go_grpc_example/23_micro/sum-srv/handler"
//
//	"context"
//)
//
//func main() {
//	// 声明服务
//	srv := micro.NewService(
//		// 框架api方式传参数
//		micro.Name("Micro-sum"))
//	//初始化
//	srv.Init(
//		// 装饰器模式 wrapper
//		micro.WrapHandler(reqLogger(srv.Client())),
//
//		// 定义钩子函数
//		micro.BeforeStart(func() error {
//			log.Log("[srv]启动之前")
//		}),
//		micro.AfterStart(func() error {
//			log.Log("[srv]启动之后")
//		}))
//	//挂载接口
//	_ = sum.RegisterSumHandler(srv.Server(), handler.Handler())
//	//启动
//	if err := srv.Run(); err != nil {
//		panic(err)
//	}
//}
//
//// 日志wrapper--broker异步消息
//func reqLogger(cli client.Client) server.HandlerWrapper {
//	// 推送器
//	pub := micro.NewPublisher("topic.log", cli)
//	// 初始化动作
//	return func(handlerFunc server.HandlerFunc) server.HandlerFunc {
//		// 中间状态
//		return func(ctx context.Context, req server.Request, rsp interface{}) error {
//			log.Info("请求准备发送日志")
//			evt := logProto.LogEvt{
//				Msg: "hello1",
//			}
//			body, _ := json.Marshal(evt)
//			_ = pub.Publish(ctx, broker.Message{
//				Header: map[string]string{
//					"serviceName": "sum",
//				},
//				Body: body,
//			})
//
//			// 发送日志之后才会调用求和方法
//			return handlerFunc(ctx, req, rsp)
//
//		}
//	}
//}

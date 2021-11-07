package main

//
//import (
//	"context"
//	"github.com/micro/go-micro"
//	"github.com/micro/go-micro/util/log"
//	"github.com/micro/go-plugins/broker/rabbitmq"
//	_ "github.com/micro/go-plugins/broker/rabbitmq"
//
//	proto "go_grpc_example/23_micro/proto/log"
//)
//
//type Sub struct {
//}
//
//func (s *Sub) Process(ctx context.Context, evt *proto.LogEvt) (err error) {
//	log.Logf("收到日志%s"evt.Msg)
//	return nil
//
//}
//
//func main() {
//	srv := micro.NewService(
//		micro.Name("micro-log"),
//		// 代码声明定义插件  http 广播broker是自带的,--broker=http
//		// 插件库 https://github.com/micro/go-plugins
//		micro.Broker(rabbitmq.NewBroker),
//
//	)
//	srv.Init()
//	// 监听
//	_ = micro.RegisterSubscriber("topic.log", srv.Server(), &Sub{})
//
//	_ = srv.Run()
//}

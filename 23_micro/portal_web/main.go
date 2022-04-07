package main

//
//import (
//	"context"
//	"github.com/micro/go-micro/web"
//	"go_package_example/23_micro/proto/sum"
//	"net/http"
//	"strconv"
//)
//
//var (
//	srvClient sum.SumService
//)
//
//func main() {
//	service := web.NewService(
//		web.Name("web端服务"),
//		web.Address(":9000"),
//		web.StaticDir("html"),
//	)
//	// 初始化
//	_ = service.Init()
//	srvClient = sum.NewSumService("Micro-sum", service.Options().Service.Client()) // 名字与srv注册相同
//	// 挂载方法
//	service.Handle("/learning/sum", Sum)
//	// 启动
//	if err := service.Run(); err != nil {
//		panic(err)
//	}
//}
//func Sum(w http.ResponseWriter, r *http.Request) {
//	inputString := r.URL.Query().Get("input")
//	input, _ := strconv.ParseInt(inputString, 10, 10)
//	req := &sum.SumRequest{
//		Input: input,
//	}
//	// 发送请求
//	rsp, err := srvClient.GetSum(context.Background(), req)
//	if err != nil {
//	}
//	_, _ = w.Write([]byte(strconv.Itoa(rsp.Output)))
//
//}

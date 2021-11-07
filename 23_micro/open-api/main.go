package main

//
//import (
//	"context"
//	"encoding/json"
//	"github.com/micro/go-micro"
//
//	api "github.com/micro/go-micro/api/proto"
//)
//
///*
//	聚合服务
//*/
//
//type Open struct {
//}
//
//func (o Open) Fetch(ctx context.Context, req *api.Request, rsp *api.Response) error {
//	ret, _ := json.Marshal(map[string]interface{}{
//		"sum":   1,
//		"prime": 2,
//	})
//
//	rsp.Body = string(ret)
//}
//func main() {
//	// 声明服务
//	srv := micro.NewService(
//		micro.Name("Micro-openapi"))
//	//聚合服务
//
//	//启动
//	if err := srv.Run(); err != nil {
//		panic(err)
//	}
//}

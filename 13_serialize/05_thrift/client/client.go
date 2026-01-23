package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Danny5487401/go_package_example/13_serialize/05_thrift/gen-go/thrift/example"
	"github.com/apache/thrift/lib/go/thrift"
)

func main() {

	transport, err := thrift.NewTHttpClient("http://localhost:9090")
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// 创建客户端协议
	protocolFactory := thrift.NewTJSONProtocolFactory()
	client := example.NewExampleServiceClientFactory(transport, protocolFactory)

	// 打开连接
	if err := transport.Open(); err != nil {
		log.Panicf("Error opening transport: %v", err)
	}
	defer transport.Close()

	// 调用远程服务
	if err := client.SayHello(context.Background(), "world"); err != nil {
		log.Panicf("Error saying hello: %v", err)
	}
	fmt.Println("Done!")
}

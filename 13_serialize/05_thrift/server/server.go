package main

import (
	"fmt"
	"net/http"

	"github.com/Danny5487401/go_package_example/13_serialize/05_thrift/gen-go/thrift/example"
	"github.com/apache/thrift/lib/go/thrift"
)

func main() {
	// 初始化 processor
	handler := &ExampleServiceImpl{}
	processor := example.NewExampleServiceProcessor(handler)

	// 创建 Thrift HTTP 处理函数
	protocolFactory := thrift.NewTJSONProtocolFactory()
	thriftHandler := thrift.NewThriftHandlerFunc(processor, protocolFactory, protocolFactory)

	// 创建 HTTP 服务器
	http.Handle("/", http.HandlerFunc(thriftHandler))
	fmt.Println("Starting the server on :9090...")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}
}

package main

import (
	"github.com/apache/thrift/lib/go/thrift"
)

func main() {
	handler := &ExampleServiceImpl{}
	processor := example.NewExampleServiceProcessor(handler)

	// 创建服务器传输对象
	transportFactory := thrift.NewTTransportFactory()
	confN := &thrift.TConfiguration{}
	protocolFactory := thrift.NewTBinaryProtocolFactoryConf(confN)
	serverTransport, err := thrift.NewTServerSocket(":9090")
	if err != nil {
		panic(err)
	}

	// 创建简单的单线程服务器
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	println("Starting the server...")
	server.Serve()
}

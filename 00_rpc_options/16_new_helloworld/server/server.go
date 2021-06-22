package main

import (
	"go_test_project/00_rpc_options/16_new_helloworld/handler"
	"go_test_project/00_rpc_options/16_new_helloworld/server_proxy"
	"net"
	"net/rpc"
)

//type HelloService struct {
//
//}
//
//func (s *HelloService) Hello(request string, reply *string) error{
//	*reply = "hello, " + request
//	return nil
//}

func main() {
	// 1.实例话一个server
	listener, _ := net.Listen("tcp", ":1234")
	//2. 注册处理逻辑 handler
	_ = server_proxy.RegisterHelloService(&handler.HelloService{})
	//_ = rpc.RegisterName(handler.HelloServiceName, &HelloService{})
	// 3.启动服务
	for {
		conn, _ := listener.Accept() // 当一个新连接进来的时候，创建套接字
		go rpc.ServeConn(conn)
	}

}

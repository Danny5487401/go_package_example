package main

import (
	"github.com/Danny5487401/go_package_example/00_rpc_options/02_new_helloworld_withStub/handler"
	"github.com/Danny5487401/go_package_example/00_rpc_options/02_new_helloworld_withStub/server_proxy"
	"net"
	"net/rpc"
)

/*
手动实现rpc的stub
*/
// 抽离实现逻辑
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
	// 修改前我们不关注名字
	//_ = rpc.RegisterName(handler.HelloServiceName, &HelloService{})
	// 修改后
	_ = server_proxy.RegisterHelloService(&handler.HelloService{})

	// 3.启动服务
	for {
		conn, _ := listener.Accept() // 当一个新连接进来的时候，创建套接字
		go rpc.ServeConn(conn)
	}

}

package main

import (
	"net"
	"net/rpc"
)

type HelloService struct {

}

func (s *HelloService) Hello(request string, reply *string) error{
	*reply = "hello, " + request
	return nil
}

func main(){
	// 1.实例话一个server
	listener,_ := net.Listen("tcp",":1234")
	//2. 注册处理逻辑 handler
	_ = rpc.RegisterName("HelloService", &HelloService{})
	// 3.启动服务
	conn, _ := listener.Accept() // 当一个新连接进来的时候，创建套接字
	rpc.ServeConn(conn)
}

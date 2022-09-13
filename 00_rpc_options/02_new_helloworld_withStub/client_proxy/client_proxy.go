package client_proxy

import (
	"go_package_example/00_rpc_options/02_new_helloworld_withStub/handler"
	"net/rpc"
)

type HelloServiceStub struct {
	*rpc.Client
}

// 初始化
func NewHelloServiceClient(protocol, address string) HelloServiceStub {
	conn, err := rpc.Dial(protocol, address)
	if err != nil {
		panic("connect error")
	}
	return HelloServiceStub{conn}
}

func (c *HelloServiceStub) Hello(request string, reply *string) error {
	err := c.Call(handler.HelloServiceName+".Hello", request, reply)
	if err != nil {
		return err
	}
	return nil
}

package client_proxy

import (
	"go_test_project/00_rpc_options/16_new_helloworld/handler"
	"net/rpc"
)

type HelloServiceStub struct {
	*rpc.Client
}

// 初始化
func NewHelloServiceClient(protol, address string) HelloServiceStub {
	conn, err := rpc.Dial(protol, address)
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

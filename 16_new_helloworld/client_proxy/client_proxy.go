package client_proxy

import (
	"net/rpc"
	"go_test_project/new_helloworld/handler"
)

type HelloServiceStub struct {
	*rpc.Client
}

// 初始化
func NewHelloServiceClient(protol,address string) HelloServiceStub  {
	conn,err := rpc.Dial(protol,address)
	if err != nil{
		panic("connect error")
	}
	return  HelloServiceStub{conn}
}

func (c *HelloServiceStub)Hello(request string,reply *string) error  {
	err :=c.Call(handler.HelloServiceName+".Hello",request,reply)
	if err != nil{
		return err
	}
	return nil
}
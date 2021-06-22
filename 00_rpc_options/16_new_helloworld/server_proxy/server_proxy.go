package server_proxy

import (
	"go_test_project/00_rpc_options/16_new_helloworld/handler"
	"net/rpc"
)

type HelloServicer interface {
	Hello(request string, reply *string) error
}

func RegisterHelloService(srv HelloServicer) error {
	return rpc.RegisterName(handler.HelloServiceName, srv)
}

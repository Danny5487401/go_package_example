package server_proxy

import (
	"github.com/Danny5487401/go_package_example/00_rpc_options/02_new_helloworld_withStub/handler"
	"net/rpc"
)

type HelloServer interface {
	Hello(request string, reply *string) error
}

func RegisterHelloService(srv HelloServer) error {
	return rpc.RegisterName(handler.HelloServiceName, srv)
}

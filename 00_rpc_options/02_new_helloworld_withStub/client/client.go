package main

import (
	"fmt"
	"github.com/Danny5487401/go_package_example/00_rpc_options/02_new_helloworld_withStub/client_proxy"
)

func main() {
	//修改前
	//// 1.建立连接
	//client,_ := rpc.Dial("tcp","localhost:1234")
	//
	//var reply string
	//err := client.Call(handler.HelloServiceName+".Hello","danny",&reply)
	//if err != nil{
	//	panic("调用失败")
	//}
	//fmt.Println(reply)

	// 修改后
	client := client_proxy.NewHelloServiceClient("tcp", "localhost:1234")
	var reply string
	err := client.Hello("danny", &reply)
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
}

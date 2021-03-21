package main

import (
	"fmt"
	"go_test_project/new_helloworld/client_proxy"
)

func main(){
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
	client := client_proxy.NewHelloServiceClient("tcp","localhost:1234")
	var reply string
	err := client.Hello("danny",&reply)
	if err != nil{
		panic(err)
	}
	fmt.Println(reply)
}

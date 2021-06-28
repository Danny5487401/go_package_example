package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main(){
	// 1.建立连接
	conn,_ := net.Dial("tcp","localhost:1234")

	var reply string // 申请空间
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	err := client.Call("HelloService.Hello","danny",&reply)
	if err != nil{
		panic("调用失败")
	}
	fmt.Println(reply)

}

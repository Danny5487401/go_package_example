package main

import (
	"fmt"
	"net/rpc"
)

func main(){
	// 1.建立连接
	client,_ := rpc.Dial("tcp","localhost:1234")

	var reply *string = new(string) // 申请空间
	err := client.Call("HelloService.Hello","danny",reply)
	if err != nil{
		panic("调用失败")
	}
	fmt.Println(*reply)

}

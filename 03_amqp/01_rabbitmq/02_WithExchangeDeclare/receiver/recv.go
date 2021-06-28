package main

import (
	"fmt"
	"go_grpc_example/03_amqp/01_rabbitmq/02_WithExchangeDeclare/rmq"
	"time"
)

func callback(d rmq.MSG) {
	fmt.Println("Ok")
	fmt.Println(string(d.Body))
}
func dlxCallback(d rmq.MSG) {
	fmt.Println("Dlx备用机")
	fmt.Println(string(d.Body))
}

func main() {
	if err := rmq.Init("03_amqp/01_rabbitmq/02_WithExchangeDeclare/rmq.json"); err != nil {
		fmt.Println(err)
	}

	if err := rmq.Pop("myPoper", callback); err != nil {
		fmt.Println(err)
	}
	// 注意选择备用机
	if err := rmq.Pop("dlxPoper", dlxCallback); err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Duration(1000) * time.Second)

	if err := rmq.Fini(); err != nil {
		fmt.Println(err)
	}

}

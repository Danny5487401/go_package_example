package main

import (
	"fmt"
	"github.com/Danny5487401/go_package_example/03_amqp/01_rabbitmq/02_WithExchangeDeclare/rmq"
)

func main() {
	if err := rmq.Init("03_amqp/01_rabbitmq/02_WithExchangeDeclare/rmq.json"); err != nil {
		fmt.Println(err)
	}

	if err := rmq.Push("myPusher", "myQueue", []byte("Hello rabbitmq5!")); err != nil {
		fmt.Println(err)
	}

	if err := rmq.Fini(); err != nil {
		fmt.Println(err)
	}
}

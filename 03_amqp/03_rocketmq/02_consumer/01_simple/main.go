package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	srvs := []string{"http://rocketmq-nameserver.rocketmq:9876"}
	topic := "rocketmq_topic"
	c, err := rocketmq.NewPushConsumer(
		// 消费组名称
		consumer.WithGroupName("go_testGroup"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver(srvs)),
	)
	if err != nil {
		log.Fatalln(err)
	}
	err = c.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("subscribe callback: %v \n", msgs[i])
		}

		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Minute)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown Consumer error: %s", err.Error())
	}
}

package main

// 发送普通消息

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {

	srvs := []string{"http://rocketmq-nameserver.rocketmq:9876"}
	// 生成生产者对象
	p, err := rocketmq.NewProducer(
		// 这里使用域名解析
		producer.WithNsResolver(primitive.NewPassthroughResolver(srvs)),
		producer.WithRetry(2),
	)
	if err != nil {
		log.Fatalln(err)
	}
	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	topic := "rocketmq_topic"

	for i := 0; i < 2; i++ {
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i)),
		}
		res, err := p.SendSync(context.Background(), msg)

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}

package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// 不能不是const 常量
var (
	topic  = "user-active"
	broker = "tencent.danny.games"
)

func main() {
	// 生成生产者客户端
	// bootstrap.servers是metadata.broker.list的别名,定义中间人
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker,
		// 默认消息体大小1000000
		"message.max.bytes": "1000000"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("创建生产者 %v\n", p)
	defer p.Close()

	doneChan := make(chan bool)
	go func() {
		// 通知主程序
		defer close(doneChan)
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("发送失败: %v\n", m.TopicPartition.Error)
				} else {
					// 打印发送成功后的消息
					fmt.Printf("发送消息到topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
				return

			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	value := "Hello Go!"
	// 发送消息
	p.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: []byte(value)}
	// wait for delivery report goroutine to finish
	_ = <-doneChan
}

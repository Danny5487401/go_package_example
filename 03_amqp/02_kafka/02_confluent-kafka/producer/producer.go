package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

/*
confluent-kafka-go
	使用了c库
生产者幂等性：
	if enable.idempotence is set). Requires broker version >= 0.11.0 要求版本大于0.11.0,要求acks=all
*/
func main() {
	// 生产者客户端
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "tencent.danny.games",
		"enable.idempotence": "true",
		"acks":               "all"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "myTopic"
	for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client", "欢迎来到kafka-go"} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

/*\
源码分析
	// Message represents a Kafka message
	type Message struct {
		TopicPartition TopicPartition
		Value          []byte
		Key            []byte
		Timestamp      time.Time
		TimestampType  TimestampType
		Opaque         interface{}
		Headers        []Header
	}
	消息字段,Value中默认带了id字段，没写为0.
*/

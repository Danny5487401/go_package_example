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
	bootstrapServers := "172.16.7.30:30001"
	userName := "user1"
	password := "qwl2pTlW6e"
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.username":     userName,
		"sasl.password":     password,
		"acks":              "all"},
	)

	if err != nil {
		fmt.Println("Failed to create producer:", err)
		return
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		// 必须消费events，默认100 万
		// You will need to read from the Events() channel to know if messages were successfully sent, and free any per-message state the client keeps for the application.
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
	topic := "test1"
	for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client", "欢迎来到kafka-go"} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

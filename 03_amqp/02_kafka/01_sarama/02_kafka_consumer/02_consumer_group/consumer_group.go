package main

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

/*

版本要求：panic: kafka: invalid configuration (consumer groups require Version to be >= V0_10_2_0)

*/

var Addr = []string{"tencent.danny.games:9092", "tencent.danny.games:9093", "tencent.danny.games:9094"}
var Topics = []string{"kafka_test1", "kafka_test2"}

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	config := sarama.NewConfig() // 默认 c.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup(Addr, "my-group", config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		topics := Topics
		handler := exampleConsumerGroupHandler{}

		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}

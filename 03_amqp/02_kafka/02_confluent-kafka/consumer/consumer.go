package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

/*
消息重复和消费幂等
	消息队列Kafka版消费的语义是at least once， 也就是至少投递一次，保证消息不丢失，但是无法保证消息不重复。在出现网络问题、客户端重启时均有可能造成少量重复消息，此时应用消费端如果对消息重复比较敏感（例如订单交易类），则应该做消息幂等。
	以数据库类应用为例，常用做法是：
		1.发送消息时，传入key作为唯一流水号ID。
		2.消费消息时，判断key是否已经消费过，如果已经消费过了，则忽略，如果没消费过，则消费一次
消费阻塞以及堆积
	消费端最常见的问题就是消费堆积，最常造成堆积的原因是：
		1.消费速度跟不上生产速度，此时应该提高消费速度，详情请参见提高消费速度。
		2.消费端产生了阻塞。
		3.消费端拿到消息后，执行消费逻辑，通常会执行一些远程调用，如果这个时候同步等待结果，则有可能造成一直等待，消费进程无法向前推进。

	消费端应该竭力避免堵塞消费线程，如果存在等待调用结果的情况，建议设置等待的超时时间，超时后作为消费失败进行处理

分区个数
	分区个数主要影响的是消费者的并发数量。
	对于同一个Consumer Group内的消费者来说，一个分区最多只能被一个消费者消费。因此，消费实例的个数不要大于分区的数量，否则会有消费实例分配不到任何分区而处于空跑状态。
	一般来说，不建议分区数小于12，否则可能影响消费发送性能；也不建议超过100个，否则易引发消费端Rebalance。
	控制台的默认分区个数是12，可以满足绝大部分场景的需求。您可以根据业务使用量进行增加。

消息队列Kafka版订阅者在订阅消息时的基本流程是：
	Poll数据。
	执行消费逻辑。
	再次poll数据
*/

func main() {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "tencent.danny.games",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"myTopic", "^aRegex.*[Tt]opic"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}

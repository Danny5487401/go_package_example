package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

/*
	消费者；Consumer or Consumer-Group API.
*/
var Addr = []string{"tencent.danny.games:9092", "tencent.danny.games:9093", "tencent.danny.games:9094"}

const Topic = "danny_kafka_log"

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	client, err := sarama.NewClient(Addr, config)
	// 关闭防止泄漏
	defer client.Close()
	if err != nil {
		panic(err)
	}
	consumer, err := sarama.NewConsumerFromClient(client)

	defer consumer.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("starting consumer success")
	partitionList, err := consumer.Partitions(Topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println("分区列表", partitionList)
	for _, partitionId := range partitionList {
		// retrieve partitionConsumer for every partitionId
		// 最新的位置开始读
		//partitionConsumer, err := consumer.ConsumePartition("redis_log", partitionId, sarama.OffsetNewest)
		partitionConsumer, err := consumer.ConsumePartition(Topic, partitionId, sarama.OffsetOldest)
		if err != nil {
			panic(err)
		}
		fmt.Println(partitionConsumer)
		go func(pc *sarama.PartitionConsumer) {
			defer (*pc).Close()
			// block
			for message := range (*pc).Messages() {
				// 返回的数据
				value := string(message.Value)
				log.Printf("Partitionid: %d; offset:%d, value: %s\n", message.Partition, message.Offset, value)
			}

		}(&partitionConsumer)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	// 保证一直
	select {
	case <-signals:

	}
}

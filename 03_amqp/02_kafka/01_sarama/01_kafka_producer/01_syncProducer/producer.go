package main

import (
	"fmt"
	"github.com/IBM/sarama"
)

// Addr 集群地址
var Addr = []string{"tencent.danny.games:9092", "tencent.danny.games:9093", "tencent.danny.games:9094"}

const Topic = "danny_kafka_log"

func main() {
	// 默认配置
	config := sarama.NewConfig()

	// 生产者配置
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follower确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个分区--随机
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel 返回

	// 构建一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = Topic
	msg.Value = sarama.StringEncoder("开始发送kafka消息")

	// 连接kafka
	producer, err := sarama.NewSyncProducer(Addr, config)
	if err != nil {
		fmt.Printf("producer closed,err : %s", err.Error())
		return
	}
	fmt.Println("连接kafka成功")
	defer producer.Close()

	// 发送信息
	pid, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed ,err: ", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)

}

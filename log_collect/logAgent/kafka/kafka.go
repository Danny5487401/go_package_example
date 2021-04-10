package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

// 专门往 kafka写日志的文件

var (
	producer sarama.SyncProducer  //声明全局的生产者
)

func Init(address []string) (err error ){
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要leader和follower确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个分区--随机
	config.Producer.Return.Successes = true //成功交付的消息将在success channel 返回


	// 连接kafka
	producer, err = sarama.NewSyncProducer(address,config)
	if err != nil{
		fmt.Printf("producer closed,err : %s",err.Error())
		return
	}
	fmt.Println("连接kafka成功")
	return
}

func SendToKafka(topic,kafkaMsg string)  {
	// 构建一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(kafkaMsg)

	// 发送到kafka
	pid,offset,err := producer.SendMessage(msg)
	if err != nil{
		fmt.Println("send msg failed ,err: ",err)
		return
	}
	fmt.Println("发送成功")
	fmt.Printf("pid:%v offset:%v\n",pid,offset)
}

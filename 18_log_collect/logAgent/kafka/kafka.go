package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

// 专门往 kafka写日志的文件

type LogData struct {
	topic string
	data  string
}

var (
	producer    sarama.SyncProducer //声明全局的生产者
	LogDataChan chan *LogData       // 缓冲
)

func Init(address []string, chanMaxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follower确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个分区--随机
	config.Producer.Return.Successes = true                   //成功交付的消息将在success channel 返回

	// 连接kafka
	producer, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		fmt.Printf("producer closed,err : %s", err.Error())
		return
	}
	fmt.Println("连接kafka成功")

	// 初始化 logDataChan,动态配置缓冲区
	LogDataChan = make(chan *LogData, chanMaxSize)

	// 初始化后等待数据发送到kafka：从通道中取数据发送到kafka
	go SendDataToKafka()
	return
}

// 真正发送到kafka的函数
func SendDataToKafka() {
	for {
		select {
		case msgDetail := <-LogDataChan:
			// 构建一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = msgDetail.topic
			msg.Value = sarama.StringEncoder(msgDetail.data)

			// 发送到kafka
			pid, offset, err := producer.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed ,err: ", err)
				return
			}
			fmt.Println("数据发送成功")
			fmt.Printf("pid:%v offset:%v\n", pid, offset)
		default:
			time.Sleep(time.Microsecond * 50)
		}

	}
}

// 把日志数据发送到内部的chan中
func SendToChan(topic, data string) {
	msg := &LogData{
		topic: topic,
		data:  data,
	}
	LogDataChan <- msg
}

// ---后期不需要的函数---
func SendToKafka(topic, kafkaMsg string) {
	// 构建一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(kafkaMsg)

	// 发送到kafka
	pid, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed ,err: ", err)
		return
	}
	fmt.Println("发送成功")
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}

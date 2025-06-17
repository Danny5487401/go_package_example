package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
)

// Addr 集群地址
var Addr = []string{"my-kafka.kafka.svc.cluster.local:9092"}

const kafkaTopic = "test1"

func main() {
	// 默认配置
	config := sarama.NewConfig()

	// 生产者配置
	config.Producer.RequiredAcks = sarama.WaitForAll // 默认WaitForLocal, WaitForAll 要求发送完数据需要leader和follower确认
	config.Producer.Return.Successes = true          // 成功交付的消息将在 success channel 返回
	config.Producer.Timeout = 5 * time.Second

	config.Net.SASL.Enable = true
	config.Net.SASL.User = "user1"
	config.Net.SASL.Password = "Pnw5pgUQUp"
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.ClientID = "producer-1"

	logger, _ := zap.NewProduction()

	sarama.Logger = &zapLogger{logger: logger}

	// 测试日志
	sarama.Logger.Printf("Sarama log with zap: %s", "test")

	// 构建一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = kafkaTopic
	msg.Value = sarama.StringEncoder("开始发送kafka消息")

	// 连接kafka
	producer, err := sarama.NewSyncProducer(Addr, config)
	if err != nil {
		fmt.Printf("producer closed,err: %s", err.Error())
		return
	}
	fmt.Println("连接kafka成功")
	defer producer.Close()

	// 发送信息
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed ,err: ", err)
		return
	}
	fmt.Printf("partId:%v offset:%v\n", partition, offset)

}

type zapLogger struct {
	logger *zap.Logger
}

func (z *zapLogger) Print(v ...interface{}) {
	z.logger.Sugar().Info(v...)
}

func (z *zapLogger) Printf(format string, v ...interface{}) {
	z.logger.Sugar().Infof(format, v...)
}

func (z *zapLogger) Println(v ...interface{}) {
	z.logger.Sugar().Info(v...)
}

// 	这么会不知道调用者是谁,需要 skip
//	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.Lshortfile|log.Lmsgprefix) // 重定向日志,默认io.Discard

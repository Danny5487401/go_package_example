package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/iimos/saramaprom"
	"go.uber.org/zap"
)

// Addr 集群地址
var Addr = []string{"my-kafka.kafka.svc.cluster.local:9092"}

const kafkaTopic = "test1"

func main() {
	// 默认配置
	config := sarama.NewConfig()
	clientId := "producer-1"
	// 暴露自身指标
	err := saramaprom.ExportMetrics(context.Background(), config.MetricRegistry, saramaprom.Options{
		Label:     clientId,
		Namespace: "sarama-test1",
	})

	// 生产者配置
	config.Producer.RequiredAcks = sarama.WaitForAll // 默认WaitForLocal, WaitForAll 要求发送完数据需要leader和follower确认
	config.Producer.Return.Successes = true          // 成功交付的消息将在 success channel 返回
	config.Producer.Timeout = 5 * time.Second

	config.Net.SASL.Enable = true
	config.Net.SASL.User = "user1"
	config.Net.SASL.Password = "Pnw5pgUQUp"
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.ClientID = clientId

	logger, _ := zap.NewProduction()

	// sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.Lshortfile|log.Lmsgprefix) // 重定向日志,默认io.Discard, 这么写会不知道调用者是谁,需要 skip
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

	http.Handle("/metrics", promhttp.Handler())
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		http.ListenAndServe("127.0.0.1:8866", nil) //多个进程不可监听同一个端口
	}()
	time.Sleep(500 * time.Second)
	wg.Done()
	wg.Wait()

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

/*
{"level":"info","ts":1750143311.16097,"caller":"01_syncProducer/producer.go:71","msg":"Sarama log with zap: test"}
{"level":"info","ts":1750143311.161108,"caller":"01_syncProducer/producer.go:75","msg":"Initializing new client"}
{"level":"info","ts":1750143311.16121,"caller":"01_syncProducer/producer.go:71","msg":"client/metadata fetching metadata for all topics from broker my-kafka.kafka.svc.cluster.local:9092\n"}
{"level":"info","ts":1750143311.355187,"caller":"01_syncProducer/producer.go:71","msg":"Connected to broker at my-kafka.kafka.svc.cluster.local:9092 (unregistered)\n"}
{"level":"info","ts":1750143311.375333,"caller":"01_syncProducer/producer.go:71","msg":"client/brokers registered new broker #0 at my-kafka-controller-0.my-kafka-controller-headless.kafka.svc.cluster.local:9092"}
{"level":"info","ts":1750143311.375369,"caller":"01_syncProducer/producer.go:71","msg":"client/brokers registered new broker #1 at my-kafka-controller-1.my-kafka-controller-headless.kafka.svc.cluster.local:9092"}
{"level":"info","ts":1750143311.375378,"caller":"01_syncProducer/producer.go:71","msg":"client/brokers registered new broker #2 at my-kafka-controller-2.my-kafka-controller-headless.kafka.svc.cluster.local:9092"}
{"level":"info","ts":1750143311.375398,"caller":"01_syncProducer/producer.go:75","msg":"Successfully initialized new client"}
{"level":"info","ts":1750143311.375896,"caller":"01_syncProducer/producer.go:71","msg":"producer/broker/1 starting up\n"}
{"level":"info","ts":1750143311.3759391,"caller":"01_syncProducer/producer.go:71","msg":"producer/broker/1 state change to [open] on test1/0\n"}
连接kafka成功
{"level":"info","ts":1750143313.3082922,"caller":"01_syncProducer/producer.go:71","msg":"Connected to broker at my-kafka-controller-1.my-kafka-controller-headless.kafka.svc.cluster.local:9092 (registered as #1)\n"}
partId:0 offset:1
{"level":"info","ts":1750143313.326383,"caller":"01_syncProducer/producer.go:75","msg":"Producer shutting down."}
{"level":"info","ts":1750143313.326485,"caller":"01_syncProducer/producer.go:75","msg":"Closing Client"}
{"level":"info","ts":1750143313.326545,"caller":"01_syncProducer/producer.go:71","msg":"producer/broker/1 input chan closed\n"}
{"level":"info","ts":1750143313.326572,"caller":"01_syncProducer/producer.go:71","msg":"producer/broker/1 shut down\n"}

*/

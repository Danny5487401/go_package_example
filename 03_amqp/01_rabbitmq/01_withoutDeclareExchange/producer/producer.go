package main

import (
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// 连接 RabbitMQ
	conn, err := amqp.Dial("amqp://admin:admin@ali.danny.games:5672/")
	failOnError(err, "连接失败")
	defer conn.Close()

	// 建立一个 channel ( 一个TCP连接可有有多个channel ）
	ch, err := conn.Channel()
	failOnError(err, "打开通道失败")
	defer ch.Close()

	// 创建一个名字叫 "hello_queue" 的队列
	q, err := ch.QueueDeclare(
		"hello_queue", // 队列名字
		// ：1、不丢失是相对的，如果宕机时有消息没来得及存盘，还是会丢失的。2、存盘影响性能。
		false, // 持久化
		false, // delete when unused
		false, // exclusive
		false, //  阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。
		nil,   // arguments
	)
	failOnError(err, "创建队列失败")

	// 构建一个消息
	body := "Hello World10!"
	msg := amqp.Publishing{
		ContentType: "text/plain", // （内容类型）
		Body:        []byte(body), //消息主体（有效载荷)
	}

	// 构建一个生产者，将消息 放入队列
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		msg)
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

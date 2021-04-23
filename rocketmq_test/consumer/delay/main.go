package main
// 延迟消息

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("go_delay_consume_Group"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"106.14.35.115:9876"})),
	)
	err := c.Subscribe("go_delay_topic", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		for _, msg := range msgs {
			t := time.Now().UnixNano()/int64(time.Millisecond) - msg.BornTimestamp
			fmt.Printf("Receive message[msgId=%s] %d ms later\n,content:%s", msg.MsgId, t,msg.Message.Body)
		}

		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("Shutdown Consumer error: %s", err.Error())
	}
}

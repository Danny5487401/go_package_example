package main

// 事务消息

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type DemoListener struct {
	// 本地事务
	localTrans       *sync.Map
	// 操作的元素
	transactionIndex int32
}

func NewDemoListener() *DemoListener {
	return &DemoListener{
		localTrans: new(sync.Map),
	}
}

// 执行本地事务逻辑
func (dl *DemoListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {

	// 执行原子操作
	nextIndex := atomic.AddInt32(&dl.transactionIndex, 1)
	fmt.Printf("执行本地事务nextIndex: %v for transactionID: %v\n", nextIndex, msg.TransactionId)

	// 共三种状态：CommitMessageState  RollbackMessageState UnknowState
	status := nextIndex % 3
	// 保存状态
	dl.localTrans.Store(msg.TransactionId, primitive.LocalTransactionState(status+1))

	fmt.Printf("dl :%#v",dl)
	// 模拟状态
	return primitive.UnknowState
}

// 执行本地事务回查
func (dl *DemoListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Printf("当前时间%v ,msg transactionID : %v\n", time.Now(), msg.TransactionId)

	// 获取数据
	v, existed := dl.localTrans.Load(msg.TransactionId)
	if !existed {
		// 数据不存在
		fmt.Printf("unknow msg: %v, return Commit", msg)
		return primitive.CommitMessageState
	}
	// 数据存在
	state := v.(primitive.LocalTransactionState)
	switch state {
	case 1:
		fmt.Printf("checkLocalTransaction COMMIT_MESSAGE: %v\n", msg)
		return primitive.CommitMessageState
	case 2:
		fmt.Printf("checkLocalTransaction ROLLBACK_MESSAGE: %v\n", msg)
		return primitive.RollbackMessageState
	case 3:
		fmt.Printf("checkLocalTransaction unknown: %v\n", msg)
		return primitive.UnknowState
	default:
		fmt.Printf("checkLocalTransaction default COMMIT_MESSAGE: %v\n", msg)
		return primitive.CommitMessageState
	}
}

func main() {
	p, _ := rocketmq.NewTransactionProducer(
		NewDemoListener(),
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"106.14.35.115:9876"})),
		producer.WithRetry(1),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s\n", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 2; i++ {
		res, err := p.SendMessageInTransaction(context.Background(),
			primitive.NewMessage("Transaction_go_topic", []byte("Hello RocketMQ again "+strconv.Itoa(i))))

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	time.Sleep(5 * time.Minute)
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}

package main

/*
Zap 跟 logrus 以及目前主流的 Go 语言 log 类似，提倡采用结构化的日志格式，而不是将所有消息放到消息体中，
	简单来讲，日志有两个概念：字段和消息。字段用来结构化输出错误相关的上下文环境，而消息简明扼要的阐述错误本身

比如，用户不存在的错误消息可以这么打印:

log.Error(“User does not exist”, zap.Int(“uid”, uid)
上面 User does not exist 是消息， 而 uid 是字段

日志属于 io 密集型的组件, 规避反射 这种类型操作是贯穿在整个 zap 的逻辑中.
zap 每打印1条日志，至少需要2次内存分配:
1.创建 field 时分配内存。

2. 将组织好的日志格式化成目标 []byte 时分配内存
zap 通过 sync.Pool 提供的对象池，复用了大量可以复用的对象，避开了 gc 这个大麻烦
*/

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"
)

func main() {

	logger, _ := zap.NewProduction() //生产环境
	//logger,_ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	url := "https://www.baidu.com"

	// 方式一 :兼容Printf
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL1",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	// {"level":"info","ts":1625733799.522482,"caller":"cosole/main.go:36","msg":"failed to fetch URL1","url":"https://www.baidu.com","attempt":3,"backoff":1}
	sugar.Infof("Failed to fetch URL2: %s", url)
	// {"level":"info","ts":1625733829.883971,"caller":"cosole/main.go:43","msg":"Failed to fetch URL2: https://www.baidu.com"}

	// 方式二 :无反射机制
	logger.Info("failed to fetch url",
		zap.String("url", url),
		zap.Int("num", 3))
	// 结果键值对方式{"level":"info","ts":1625733829.883981,"caller":"cosole/main.go:46","msg":"failed to fetch url","url":"https://www.baidu.com","num":3}

	// 错误栈帧调用
	errorStacktraceDemo()

}

// 错误栈帧查看
func errorStacktraceDemo() {
	logger, err := zap.NewDevelopment()
	defer logger.Sync()
	if err != nil {
		panic(err)
	}
	logger.Info("errorField", zap.Error(errors.New("demo err")))

	fmt.Println(zap.Stack("default stack").String) //main65行->main52行>proc.go 203行
	fmt.Println("------")
	fmt.Println(zap.StackSkip("skip 2", 2).String) // 跳过两行 proc.go 203行
	logger.Info("stacktrace default", zap.Stack("default stack"))
	logger.Info("stacktrace skip 2", zap.StackSkip("skip 2", 2))
}

package main

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"
)

func main() {

	logger, _ := zap.NewProduction() //生产环境
	//logger,_ := zap.NewDevelopment()  //开发环境
	defer logger.Sync() // flushes buffer, if any
	url := "https://www.baidu.com"

	// 方式一 :兼容Printf
	sugarPrint(logger)

	// 方式二 :无反射机制
	logger.Info("failed to fetch url3",
		zap.String("url", url),
		zap.Int("num", 3))
	// 结果键值对方式{"level":"info","ts":1625733829.883981,"caller":"cosole/consumer.go:46","msg":"failed to fetch url","url":"https://www.baidu.com","num":3}

	// 错误栈帧调用
	errorStacktraceDemo(logger)

}

func sugarPrint(logger *zap.Logger) {
	url := "https://www.google.com"
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL1",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	// {"level":"info","ts":1625733799.522482,"caller":"cosole/consumer.go:36","msg":"failed to fetch URL1","url":"https://www.baidu.com","attempt":3,"backoff":1}

	sugar.Infof("Failed to fetch URL2: %s", url)
	// {"level":"info","ts":1625733829.883971,"caller":"cosole/consumer.go:43","msg":"Failed to fetch URL2: https://www.baidu.com"}

}

// 错误栈帧查看
func errorStacktraceDemo(logger *zap.Logger) {
	// 抛出错误
	logger.Info("抛出错误errorField", zap.Error(errors.New("demo err")))
	// 定义错误
	fmt.Println(strings.Repeat("---", 30))
	fmt.Println(zap.Stack("default stack1").String)

	fmt.Println(strings.Repeat("===", 30))
	fmt.Println(zap.StackSkip("跳过skip前 2个栈", 2).String)

	fmt.Println(strings.Repeat("***", 30))
	logger.Info("栈追踪默认", zap.Stack("默认栈1"))

	fmt.Println(strings.Repeat("###", 30))
	logger.Info("栈追踪 跳过 2层", zap.StackSkip("skip 2", 2))

	fmt.Println(strings.Repeat("$$$", 30))
	logger.Error("栈追踪 默认", zap.Stack("默认栈2"))

}

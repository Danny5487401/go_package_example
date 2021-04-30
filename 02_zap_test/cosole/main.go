package main

import (
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction() //生产环境
	//logger,_ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	url := "https://www.baidu.com"

	// 方式一
	//sugar := logger.Sugar()
	//sugar.Infow("failed to fetch URL",
	//	// Structured context as loosely typed key-value pairs.
	//	"url", url,
	//	"attempt", 3,
	//	"backoff", time.Second,
	//)
	//sugar.Infof("Failed to fetch URL: %s", url)

	// 方式二 无反射机制
	logger.Info("failed to fetch url",
		zap.String("url",url),
		zap.Int("num",3))

}
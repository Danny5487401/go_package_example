package main

import (
	"fmt"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	// 使用全局并发安全
	zap.ReplaceGlobals(logger)

	err := fmt.Errorf("数据有误")
	// 需要看错误栈
	zap.S().Error("错误", err)

}

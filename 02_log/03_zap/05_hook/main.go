package main

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DemoHook 用于生成一个示例Hook函数，Hook函数实现当前日志大于Error级别时打印Message的长度
// 这里使用函数的方法返回一个Hook函数的好处时外层函数可以根据需要传入一些参数供Hook函数使用
func DemoHook() func(zapcore.Entry) error {
	return func(e zapcore.Entry) error {
		if e.Level < zapcore.ErrorLevel {
			return nil
		}

		fmt.Printf("message length:%d\n", len(e.Message))
		return nil
	}
}

func main() {
	logger := zap.NewExample()                         // 生成一个logger
	logger = logger.WithOptions(zap.Hooks(DemoHook())) // 将回调函数注册到logger生成一个新的logger
	logger.Info("123456789")                           // 不打印长度信息
	logger.Error("123456789")                          // 打印长度信息
}

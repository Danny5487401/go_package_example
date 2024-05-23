// Package main is a demo for zapcore for sentry capture error

package main

import (
	"errors"
	"fmt"
	"github.com/Danny5487401/go_package_example/25_sentry/zap_sentry/sentryzapcore"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 生成SentryCore对象并添加到Logger中
func main() {
	dsn := "http://437e75d208004309b74cac7dad719c3f@tencent.danny.games:9000/2"
	// 默认使用IoCore
	logger := zap.NewExample()

	// sentrycore配置
	cfg := sentryzapcore.SentryCoreConfig{
		Level: zap.ErrorLevel,
		Tags: map[string]string{
			"source":    "ticket-system",       // 这里需要获取traceId, requestId信息,每次都带上的信息
			"requestId": "http-requestId-1234", // web可以从context中获取
		},
	}
	// 生成sentry客户端
	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              dsn,
		Debug:            true,
		AttachStacktrace: true,
	})
	if err != nil {
		fmt.Println(err)
	}
	// 生成sentryCore
	sCore := sentryzapcore.NewSentryCore(cfg, sentryClient)
	// 添加sentryCore到默认logger产生新的logger，使用该logger即可自动上报sentry
	logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, sCore) //添加一个自定义的core
	}))

	logger.Info("info log")

	logger.Error("这条error日志sentry自动收集", zap.String("field1", "value1"), zap.Error(errors.New("this ia an error")))
	time.Sleep(2 * time.Second) // sleep避免程序结束太快而导致上报失败
}

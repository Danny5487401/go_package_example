package _1_jaeger_xorm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go" // opentracing协议
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/log/zap"
	zap2 "go.uber.org/zap"
	"io"
	"xorm.io/xorm"
	xormLog "xorm.io/xorm/log"
)

func initJaeger() (closer io.Closer, err error) {
	// 根据配置初始化Tracer 返回Closer
	tracer, closer, err := (&config.Configuration{
		ServiceName: "xormWithTracing",
		Disabled:    false,
		Sampler: &config.SamplerConfig{
			Type: jaeger.SamplerTypeConst,
			// param的值在0到1之间，设置为1则将所有的Operation输出到Reporter
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "ali.danny.games:6831",
		},
	}).NewTracer()
	if err != nil {
		return
	}

	// 设置全局Tracer - 如果不设置将会导致上下文无法生成正确的Span
	opentracing.SetGlobalTracer(tracer)
	return
}

// xorm 1.0.2已经支持Hook钩子函数注入操作上下文
func NewEngineForHook() (engine *xorm.Engine, err error) {
	// XORM创建引擎
	engine, err = xorm.NewEngine("mysql", "xorm:chuanzhi@(ali.danny.games:3306)/xormtest?charset=utf8mb4")
	if err != nil {
		return
	}

	// 使用我们的钩子函数
	engine.AddHook(NewTracingHook())
	return
}

// 请注意，这种方法只适合在Xorm 1.0版本以上和1.0.2版本以下
// 非并发安全，请慎重使用
func NewEngineForLogger() (engine *xorm.Engine, err error) {
	// XORM创建引擎
	engine, err = xorm.NewEngine("mysql", "xorm:chuanzhi@(ali.danny.games:3306)/xormtest?charset=utf8mb4")
	if err != nil {
		return
	}

	// 创建自定义的日志实例
	_l, err := zap2.NewDevelopment()
	if err != nil {
		return
	}

	// 将日志实例设置到XORM的引擎中
	engine.SetLogger(&CustomCtxLogger{
		logger:  zap.NewLogger(_l),
		level:   xormLog.LOG_DEBUG,
		showSQL: true,
	})
	return
}

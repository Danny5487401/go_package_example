package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go/log/zap"
	xormLog "xorm.io/xorm/log"
)

/*
	下面都是实现自定义ContextLogger的部分，这里使用OpenTracing自带的Zap日志（语法糖裁剪版）
	也可以使用Zap、logrus、原生自带的log实现xorm的ContextLogger接口
*/

type CustomCtxLogger struct {
	logger  *zap.Logger
	level   xormLog.LogLevel
	showSQL bool
	span    opentracing.Span
}

// BeforeSQL implements ContextLogger
func (l *CustomCtxLogger) BeforeSQL(ctx xormLog.LogContext) {
	// ----> 重头戏在这里，需要从Context上下文中创建一个新的Span来对SQL执行进行链路监控
	l.span, _ = opentracing.StartSpanFromContext(ctx.Ctx, "XORM SQL Execute")
}

// AfterSQL implements ContextLogger
func (l *CustomCtxLogger) AfterSQL(ctx xormLog.LogContext) {
	// defer结束掉span
	defer l.span.Finish()

	// 原本的SimpleLogger里面会获取一次SessionId
	var sessionPart string
	v := ctx.Ctx.Value("__xorm_session_id")
	if key, ok := v.(string); ok {
		sessionPart = fmt.Sprintf(" [%s]", key)
		l.span.LogFields(tracerLog.String("session_id", sessionPart))
	}

	// 将Ctx中全部的信息写入到Span中
	l.span.LogFields(tracerLog.String("SQL", ctx.SQL))
	l.span.LogFields(tracerLog.Object("args", ctx.Args))
	l.span.SetTag("execute_time", ctx.ExecuteTime)

	if ctx.ExecuteTime > 0 {
		l.logger.Infof("[SQL]%s %s %v - %v", sessionPart, ctx.SQL, ctx.Args, ctx.ExecuteTime)
	} else {
		l.logger.Infof("[SQL]%s %s %v", sessionPart, ctx.SQL, ctx.Args)
	}
}

// Errorf implement ILogger
func (l *CustomCtxLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, v...))
	return
}

// Debugf implement ILogger
func (l *CustomCtxLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
	return
}

// Infof implement ILogger
func (l *CustomCtxLogger) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
	return
}

// Warnf implement ILogger ---> 这里偷懒了，直接用Info代替
func (l *CustomCtxLogger) Warnf(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
	return
}

// Level implement ILogger
func (l *CustomCtxLogger) Level() xormLog.LogLevel {
	return l.level
}

// SetLevel implement ILogger
func (l *CustomCtxLogger) SetLevel(lv xormLog.LogLevel) {
	l.level = lv
	return
}

// ShowSQL implement ILogger
func (l *CustomCtxLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

// IsShowSQL implement ILogger
func (l *CustomCtxLogger) IsShowSQL() bool {
	return l.showSQL
}

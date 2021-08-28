package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"io"
	"net/http"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

/*
jaeger消耗：
	作用于应用层，提取，注入，生成span,序列化成thrift，发送到远程
处理方式：
	选择合理的采集策略，constant全量采集,probabilist随机采集,rate limiting每秒采集数,remote远程
主要概念
	1.提取:找到父亲
	2.注入:为了孩子能找到爸爸----跨进程使用
	3.异步report：finish()完成时

*/

func main() {
	// 初始化Tracer
	_, err := initJaeger()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//defer closer.Close()
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("test1")

	// 上下文，减少注入
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	go mysql(ctx)

	//注入
	h := http.Header{}
	h.Add("header1", "headerValue1")
	tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(h))

	// 正常业务
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	span.SetTag("sum", sum)
	span.Finish()

}

func mysql(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mysql")
	span.SetTag("mysqlKey", "mysqlValue")
	span.Finish()
}

func initJaeger() (closer io.Closer, err error) {
	// 根据配置初始化Tracer 返回Closer
	tracer, closer, err := (&config.Configuration{
		ServiceName: "jaegerTest",
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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	service     = "trace-demo"
	environment = "production"
	id          = 1
)

var (
	tp *tracesdk.TracerProvider
)

func MainBaggageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")

	fmt.Println("index handler")

	tr := tp.Tracer("component-main")
	spanCtx, span := tr.Start(context.Background(), "index-handler")
	defer span.End()

	time.Sleep(time.Second * 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go funA(spanCtx, wg)
	funcBWithBaggage(spanCtx)

	wg.Wait()
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")

	fmt.Println("index handler")

	tr := tp.Tracer("component-main")
	spanCtx, span := tr.Start(context.Background(), "index-handler")
	defer span.End()

	time.Sleep(time.Second * 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	// funcA目的就是来实现往span里面写tags
	go funA(spanCtx, wg)
	// funB 为了实现跨服务的trace
	funB(spanCtx)

	wg.Wait()
}

// tracerProvider is 返回一个openTelemetry TraceProvider，这里用的是jaeger
func tracerProvider(url string) error {
	fmt.Println("init traceProvider")

	// 创建jaeger provider
	// 可以直接连collector也可以连agent
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return nil
}

func funA(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println("do function a")

	// Use the global TracerProvider.
	tr := otel.Tracer("component-main")

	// 如果有调用子方法的，需要用这个spanctx，不然会挂到父span上面
	_, span := tr.Start(ctx, "func-a")

	// 只能有特定数据类型
	span.SetAttributes(attribute.KeyValue{
		Key:   "isGetHere",
		Value: attribute.BoolValue(true),
	})

	span.SetAttributes(attribute.KeyValue{
		Key:   "current time",
		Value: attribute.StringValue(time.Now().Format("2006-01-02 15:04:05")),
	})

	type _LogStruct struct {
		CurrentTime time.Time `json:"current_time"`
		PassByWho   string    `json:"pass_by_who"`
		Name        string    `json:"name"`
	}

	logTest := _LogStruct{
		CurrentTime: time.Now(),
		PassByWho:   "postman",
		Name:        "func-a",
	}

	b, _ := json.Marshal(logTest)

	span.SetAttributes(attribute.Key("这是测试日志的key").String(string(b)))

	time.Sleep(time.Second * 1)

	defer span.End()
}

func funB(ctx context.Context) {

	fmt.Println("do function b")

	tr := otel.Tracer("component-main")

	spanCtx, span := tr.Start(ctx, "func-b")

	fmt.Println("trace:", span.SpanContext().TraceID().String(), ", span: ", span.SpanContext().SpanID())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:8090/service-2", nil)

	// header写入trace-id和span-id
	req.Header.Set("trace-id", span.SpanContext().TraceID().String())
	req.Header.Set("span-id", span.SpanContext().SpanID().String())

	p := otel.GetTextMapPropagator()
	p.Inject(spanCtx, propagation.HeaderCarrier(req.Header))

	// 发送请求
	_, _ = client.Do(req)

	//结束当前请求的span
	defer span.End()
}

func funcBWithBaggage(ctx context.Context) {

	fmt.Println("do function b baggage")

	tr := otel.Tracer("component-main")
	spanCtx, span := tr.Start(ctx, "func-b-with-baggage")

	client := &http.Client{}
	// 请求服务2
	req, _ := http.NewRequest("POST", "http://localhost:8090/service-2-baggage", nil)

	// 使用baggage写入trace id和span id
	p := propagation.Baggage{}

	traceMember, _ := baggage.NewMember("trace-id", span.SpanContext().TraceID().String())
	spanMember, _ := baggage.NewMember("span-id", span.SpanContext().SpanID().String())

	b, _ := baggage.New(traceMember, spanMember)

	ctxBaggage := baggage.ContextWithBaggage(spanCtx, b)

	fmt.Println("trace id : ", span.SpanContext().TraceID().String())

	//req.Header.Set("baggage", "trace-id="+span.SpanContext().TraceID().String())
	p.Inject(ctxBaggage, propagation.HeaderCarrier(req.Header))

	// 发送请求
	_, _ = client.Do(req)

	//结束当前请求的span
	defer span.End()
}

func main() {

	var err error

	err = tracerProvider("http://tencent.danny.games:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	http.HandleFunc("/baggage", MainBaggageHandler)
	http.HandleFunc("/", MainHandler)
	http.ListenAndServe("127.0.0.1:8060", nil)

}

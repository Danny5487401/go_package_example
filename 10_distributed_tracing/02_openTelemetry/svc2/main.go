package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	service     = "trace-svc-demo"
	environment = "production"
	id          = 1
)

var (
	tp *tracesdk.TracerProvider
)

// tracerProvider is 返回一个openTelemetry TraceProvider，这里用的是jaeger
func tracerProvider(url string) error {
	fmt.Println("init traceProvider")

	// 创建jaeger provider
	// 可以直接连collector也可以连agent

	//exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost("127.0.0.1"), jaeger.WithAgentPort("6831")))
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
			attribute.Int("进程ID", os.Getpid()),
		)),
	)

	return nil
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello Danny2")

	carrier := propagation.HeaderCarrier{}
	carrier.Set("trace-id", r.Header.Get("trace-id"))

	var propagator = otel.GetTextMapPropagator()
	pctx := propagator.Extract(r.Context(), carrier)
	tr := tp.Tracer("component-main")

	traceID := r.Header.Get("trace-id")
	spanID := r.Header.Get("span-id")

	fmt.Println("parent trace-id : ", traceID)

	traceid, _ := trace.TraceIDFromHex(traceID)
	spanid, _ := trace.SpanIDFromHex(spanID)

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceid,
		SpanID:     spanid,
		TraceFlags: trace.FlagsSampled, //这个没写，是不会记录的
		TraceState: trace.TraceState{},
		Remote:     true,
	})

	// 不用pctx，不会把spanctx当做parentCtx
	sct := trace.ContextWithRemoteSpanContext(pctx, spanCtx)

	_, span := tr.Start(sct, "func-c")

	sc := span.SpanContext()
	fmt.Println("trace:", sc.TraceID().String(), ", span: ", sc.SpanID())

	defer span.End()

	time.Sleep(time.Second * 2)
}

func MainHandlerWithBaggage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")

	var propagator = propagation.TextMapPropagator(propagation.Baggage{})

	pctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	tr := tp.Tracer("component-main")

	bag := baggage.FromContext(pctx)

	traceid, _ := trace.TraceIDFromHex(bag.Member("trace-id").Value())
	spanid, _ := trace.SpanIDFromHex(bag.Member("span-id").Value())

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceid,
		SpanID:     spanid,
		TraceFlags: trace.FlagsSampled, //这个没写，是不会记录的
		TraceState: trace.TraceState{},
		Remote:     true,
	})

	// 不用pctx，不会把spanctx当做parentCtx
	sct := trace.ContextWithRemoteSpanContext(pctx, spanCtx)

	_, span := tr.Start(sct, "func-c-with-baggage")

	sc := span.SpanContext()

	fmt.Println("trace id:", sc.TraceID().String(), ", span id: ", sc.SpanID())

	span.SetAttributes(attribute.Key("svc-2-baggage").String("service-2-baggage"))
	defer span.End()

	// 必须放在span start之后
	time.Sleep(time.Second * 2)
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

	http.HandleFunc("/service-2", MainHandler)                    // 用来接收通过request header方式的请求
	http.HandleFunc("/service-2-baggage", MainHandlerWithBaggage) //  用来接收通过baggage item方式的请求
	http.ListenAndServe("127.0.0.1:8090", nil)
}

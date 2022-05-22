package _1_jaeger_xorm

import (
	"github.com/opentracing/opentracing-go"

	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

//func Test_initJaeger(t *testing.T) {
//	// 初始化Tracer
//	closer, err := initJaeger()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer closer.Close()
//
//	// 用全局Tracer分裂出一个新的Span
//	span := opentracing.GlobalTracer().StartSpan("test Jaeger")
//	defer span.Finish()
//	span.SetTag("my name is", "Avtion")
//	time.Sleep(3 * time.Second)
//}

// XORM技术文档范例
type User struct {
	Id   int64
	Name string `xorm:"varchar(25) notnull unique 'usr_name' comment('姓名')"`
}

// 新方式进行上下文注入，要求 xorm 1.0.2版本
func TestNewEngineForHook(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	// 初始化XORM引擎
	engine, err := NewEngineForHook()
	if err != nil {
		t.Fatal(err)
	}

	// 初始化Tracer
	closer, err := initJaeger()
	if err != nil {
		t.Fatal(err)
	}
	defer closer.Close()

	// 生成新的Span - 注意将span结束掉，不然无法发送对应的结果
	span := opentracing.StartSpan("xorm sync")
	defer span.Finish()

	// 把生成的Root Span写入到Context上下文，获取一个子Context
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	// 将子上下文传入Session
	session := engine.Context(ctx)

	// Sync2同步表结构
	if err := session.Sync2(&User{}); err != nil {
		t.Fatal(err)
	}

	// 插入一条数据
	if _, err := session.InsertOne(&User{Name: fmt.Sprintf("test-%d", rand.Intn(1<<10))}); err != nil {
		t.Fatal()
	}
}

// 旧方式进行上下文注入
//func TestNewEngine(t *testing.T) {
//	// 初始化XORM引擎
//	engine, err := NewEngineForLogger()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// 初始化Tracer
//	closer, err := initJaeger()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer closer.Close()
//
//	// 生成新的Span - 注意将span结束掉，不然无法发送对应的结果
//	span := opentracing.StartSpan("xorm sync")
//	defer span.Finish()
//
//	// 把生成的Root Span写入到Context上下文，获取一个子Context
//	ctx := opentracing.ContextWithSpan(context.Background(), span)
//
//	// 将子上下文传入Session
//	session := engine.Context(ctx)
//
//	// Sync2同步表结构
//	if err := session.Sync2(&User{}); err != nil {
//		t.Fatal(err)
//	}
//
//	// 插入一条数据
//	if _, err := session.InsertOne(&User{Name: "test"}); err != nil {
//		t.Fatal()
//	}
//}

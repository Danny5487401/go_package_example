package _1_yaeger_xorm

import (
	"context"
	"github.com/opentracing/opentracing-go"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"xorm.io/builder"
	"xorm.io/xorm/contexts"
)

// 使用新定义类型
// context.WithValue方法中注释写到
// 提供的键需要可比性，而且不能是字符串或者任意内建类型，避免不同包之间
// 调用到相同的上下文Key发生碰撞，context的key具体类型通常为struct{}，
// 或者作为外部静态变量（即开头字母为大写）,类型应该是一个指针或者interface类型
type xormHookSpan struct{}

var xormHookSpanKey = &xormHookSpan{}

type TracingHook struct {
	// 注意Hook伴随DB实例的生命周期，所以我们不能在Hook里面寄存span变量
	// 否则就会发生并发问题
	before func(c *contexts.ContextHook) (context.Context, error)
	after  func(c *contexts.ContextHook) error
}

// 让编译器知道这个是xorm的Hook，防止编译器无法检查到异常
var _ contexts.Hook = &TracingHook{}

func (h *TracingHook) BeforeProcess(c *contexts.ContextHook) (context.Context, error) {
	return h.before(c)
}

func (h *TracingHook) AfterProcess(c *contexts.ContextHook) error {
	return h.after(c)
}

func NewTracingHook() *TracingHook {
	return &TracingHook{
		before: before,
		after:  after,
	}
}

func before(c *contexts.ContextHook) (context.Context, error) {
	// 这里一定要注意，不要拿第二个返回值作为上下文进行替换，而是用自己的key
	span, _ := opentracing.StartSpanFromContext(c.Ctx, "xorm sql execute")
	c.Ctx = context.WithValue(c.Ctx, xormHookSpanKey, span)
	return c.Ctx, nil
}

func after(c *contexts.ContextHook) error {
	// 自己实现opentracing的SpanFromContext方法，断言将interface{}转换成opentracing的span
	sp, ok := c.Ctx.Value(xormHookSpanKey).(opentracing.Span)
	if !ok {
		// 没有则说明没有span
		return nil
	}
	defer sp.Finish()

	// 记录我们需要的内容
	if c.Err != nil {
		sp.LogFields(tracerLog.Object("err", c.Err))
	}

	// 使用xorm的builder将查询语句和参数结合，方便后期调试
	sql, _ := builder.ConvertToBoundSQL(c.SQL, c.Args)

	sp.LogFields(tracerLog.String("SQL", sql))
	sp.LogFields(tracerLog.Object("args", c.Args))
	sp.SetTag("execute_time", c.ExecuteTime)

	return nil
}

/*
源码分析：xorm -> hook.go
	钩子函数的接口
	type Hook interface {
		BeforeProcess(c *ContextHook) (context.Context, error)
		AfterProcess(c *ContextHook) error
	}
	type Hooks struct {
		hooks []Hook
	}

	func (h *Hooks) AddHook(hooks ...Hook) {
		h.hooks = append(h.hooks, hooks...)
	}

	// 上下文ContextHook
	type ContextHook struct {
		// 开始时间
		start       time.Time
		// 上下文
		Ctx         context.Context
		// SQL语句
		SQL         string        // log content or SQL
		// SQL参数
		Args        []interface{} // if it's a SQL, it's the arguments
		// 查询结果
		Result      sql.Result
		// 执行时间
		ExecuteTime time.Duration
		// 如果发生错误，会赋值
		Err         error // SQL executed error
	}


	调用逻辑xorm -> db.go

	func (db *DB) beforeProcess(c *contexts.ContextHook) (context.Context, error) {
		if db.NeedLogSQL(c.Ctx) {
			// <-- 重要，这里是将日志上下文转化成值传递
			// 所以不能修改context.Context的内容
			db.Logger.BeforeSQL(log.LogContext(*c))
		}
		// Hook是指针传递，所以可以修改context.Context的内容
		ctx, err := db.hooks.BeforeProcess(c)
		if err != nil {
			return nil, err
		}
		return ctx, nil
	}

	func (db *DB) afterProcess(c *contexts.ContextHook) error {
		// 和beforeProcess同理，日志上下文不能修改context.Context的内容
		// 而hook可以
		err := db.hooks.AfterProcess(c)
		if db.NeedLogSQL(c.Ctx) {
			db.Logger.AfterSQL(log.LogContext(*c))
		}
		return err
	}

	这一段就是实际SQL查询过程中调用日志和Hook的过程，从这里可以非常明显的看到日志模块传入的是值而不是指针，从而导致了我们无法修改日志模块中的上下文实现span的传递，
	只能利用全局日志实例来传递span，这直接出现了并发安全问题
	而Hook的传递使用的是指针传递，将contexts.ContextHook的指针传入钩子函数执行流程，允许我们直接操作Ctx


*/

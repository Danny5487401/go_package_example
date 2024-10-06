<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [slog](#slog)
  - [使用](#%E4%BD%BF%E7%94%A8)
  - [源码](#%E6%BA%90%E7%A0%81)
    - [三大组件](#%E4%B8%89%E5%A4%A7%E7%BB%84%E4%BB%B6)
    - [打印过程](#%E6%89%93%E5%8D%B0%E8%BF%87%E7%A8%8B)
    - [handler](#handler)
  - [升级第三方库](#%E5%8D%87%E7%BA%A7%E7%AC%AC%E4%B8%89%E6%96%B9%E5%BA%93)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# slog



在 Go1.21 中引入的 log/slog 软件包试图弥补原有日志软件包的不足，即日志缺乏结构化和级别特性.


slog 设计之初大量参考了社区中现有日志包方案，相比于 log，主要解决了两个问题：

- log 不支持日志级别。
- log 日志不是结构化的。

## 使用

设置日志级别

```go
func SetLogLoggerLevel(level Level) (oldLevel Level) {
	oldLevel = logLoggerLevel.Level()
	logLoggerLevel.Set(level)
	return
}

```

level 存储结构
```go

type LevelVar struct {
	val atomic.Int64
}

// Set sets v's level to l.
func (v *LevelVar) Set(l Level) {
	v.val.Store(int64(l))
}
```


```go

type Level int


const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)
```
可以发现这几个日志级别并不是连续的，这是 slog 团队经过深思熟虑后的结果。故意这样设计，是为了方便我们增加自定义日志级别。使我们可以在任意两个日志级别之间定义自己的日志级别


## 源码

### 三大组件
slog从逻辑上分为前端(front)和后端(backend)

结构体 Logger 是一个结构体，面向用户侧
```go
// go1.22.2/src/log/slog/logger.go
type Logger struct {
	handler Handler // for structured logging
}

```

slog前端就是slog提供给使用者的API，不过，很遗憾slog依旧像log那样没有抽取出Logger接口，而是定义了一个Logger结构体，这也意味着我们依旧无法在整个Go社区统一前端API；


```go
type Handler interface {
    // 开启的日志级别
	Enabled(context.Context, Level) bool

	// 处理 Record.
	Handle(context.Context, Record) error

	// 扩充并返回新的handler
	WithAttrs(attrs []Attr) Handler
	
	//
	// A Handler should treat WithGroup as starting a Group of Attrs that ends
	// at the end of the log event. That is,
	//
	//     logger.WithGroup("s").LogAttrs(ctx, level, msg, slog.Int("a", 1), slog.Int("b", 2))
	//
	// should behave like
	//
	//     logger.LogAttrs(ctx, level, msg, slog.Group("s", slog.Int("a", 1), slog.Int("b", 2)))
	//
	// If the name is empty, WithGroup returns the receiver.
	WithGroup(name string) Handler
}

```
接口类型的存在，让slog的后端扩展性更强，我们除了可以使用slog提供的两个内置Handler实现：TextHandler和JSONHandler之外，还可以基于第三方log包定义或完全自定义后端Handler的实现。




Record 是一条日志条目，一个 Record 实例就代表了一条日志记录

```go
type Record struct {
	// 当前这条日志记录的时间
	Time time.Time

	// The log message.
	Message string

	// The level of the event.
	Level Level

	// The program counter at the time the record was constructed, as determined
	// by runtime.Callers. If zero, no program counter is available.
	//
	// The only valid use for this value is as an argument to
	// [runtime.CallersFrames]. In particular, it must not be passed to
	// [runtime.FuncForPC].
	PC uintptr

	// Allocation optimization: an inline array sized to hold
	// the majority of log calls (based on examination of open-source
	// code). It holds the start of the list of Attrs.
	front [nAttrsInline]Attr

	// The number of Attrs in front.
	nFront int

	// The list of Attrs except for those in front.
	// Invariants:
	//   - len(back) > 0 iff nFront == len(front)
	//   - Unused array elements are zero. Used to detect mistakes.
	back []Attr
}
```
每条日志记录都由时间、级别、消息等参数和一组键值对组成。

从 slog 架构逻辑上讲，其中 Logger 被称为前端，Handler 被称为后端，而 Record 就是连接二者的桥梁。


### 打印过程

```go
func Info(msg string, args ...any) {
	Default().log(context.Background(), LevelInfo, msg, args...)
}

```


```go


func (l *Logger) log(ctx context.Context, level Level, msg string, args ...any) {
	// 判断 level 级别
	if !l.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	// 内部调用栈是否忽略
	if !internal.IgnorePC {
		var pcs [1]uintptr
		// skip [runtime.Callers, this function, this function's caller]
		runtime.Callers(3, pcs[:])
		pc = pcs[0]
	}
	// 初始化record
	r := NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.Handler().Handle(ctx, r)
}
```


### handler
默认 handler
获取并初始化默认 logger

```go
var defaultLogger atomic.Pointer[Logger]

func init() {
	// defaultHandler的output函数就是log.Output：
	defaultLogger.Store(New(newDefaultHandler(loginternal.DefaultOutput)))
}

// Default returns the default [Logger].
func Default() *Logger { return defaultLogger.Load() }


type defaultHandler struct {
	ch *commonHandler
	// internal.DefaultOutput, except for testing
	output func(pc uintptr, data []byte) error
}

func newDefaultHandler(output func(uintptr, []byte) error) *defaultHandler {
	return &defaultHandler{
		ch:     &commonHandler{json: false},
		output: output,
	}
}

```

loginternal.DefaultOutput 初始化

```go
// /go1.22.2/src/log/log.go

var std = New(os.Stderr, "", LstdFlags)

func init() {
	internal.DefaultOutput = func(pc uintptr, data []byte) error {
		return std.output(pc, 0, func(buf []byte) []byte {
			return append(buf, data...)
		})
	}
}

```

```go
func (h *defaultHandler) Handle(ctx context.Context, r Record) error {
	buf := buffer.New()
	buf.WriteString(r.Level.String())
	buf.WriteByte(' ')
	buf.WriteString(r.Message)
	state := h.ch.newHandleState(buf, true, " ")
	defer state.free()
	state.appendNonBuiltIns(r)
	return h.output(r.PC, *buf)
}
```

jsonhandler 说明

```go
func NewJSONHandler(w io.Writer, opts *HandlerOptions) *JSONHandler {
	if opts == nil {
		opts = &HandlerOptions{}
	}
	return &JSONHandler{
		&commonHandler{
			json: true,
			w:    w,
			opts: *opts,
			mu:   &sync.Mutex{},
		},
	}
}
```

```go
func (h *JSONHandler) Handle(_ context.Context, r Record) error {
	return h.commonHandler.handle(r)
}
```

textHandler 和JsonHandler 共同实现

```go
func (h *commonHandler) handle(r Record) error {
	state := h.newHandleState(buffer.New(), true, "")
	defer state.free()
	if h.json {
		state.buf.WriteByte('{')
	}
	// Built-in attributes. They are not in a group.
	stateGroups := state.groups
	state.groups = nil // So ReplaceAttrs sees no groups instead of the pre groups.
	rep := h.opts.ReplaceAttr
	// time
	if !r.Time.IsZero() {
		key := TimeKey
		val := r.Time.Round(0) // strip monotonic to match Attr behavior
		if rep == nil {
			state.appendKey(key)
			state.appendTime(val)
		} else {
			state.appendAttr(Time(key, val))
		}
	}
	// level
	key := LevelKey
	val := r.Level
	if rep == nil {
		state.appendKey(key)
		state.appendString(val.String())
	} else {
		state.appendAttr(Any(key, val))
	}
	// source
	if h.opts.AddSource {
		state.appendAttr(Any(SourceKey, r.source()))
	}
	key = MessageKey
	msg := r.Message
	if rep == nil {
		state.appendKey(key)
		state.appendString(msg)
	} else {
		state.appendAttr(String(key, msg))
	}
	state.groups = stateGroups // Restore groups passed to ReplaceAttrs.
	state.appendNonBuiltIns(r)
	state.buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(*state.buf)
	return err
}
```

## 升级第三方库
活跃的Golang社区有许多为 slog 定制的升级框架，以下是一些常用框架。

- github.com/samber/slog-multi：处理器链，如流水线、路由器、扇出等。
- github.com/samber/slog-sampling：丢弃重复的日志条目。
- slog-shim: 为 1.21 以下的 Go 版本提供向后兼容的 slog 支持。
- sloggen: 生成各种辅助工具。
- sloglint:可确保代码的一致性


## 参考

- [万字解析 Go 官方结构体化日志包 slog](https://www.cnblogs.com/cheyunhua/p/18269295)
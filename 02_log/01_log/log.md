<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [log](#log)
  - [打印流程](#%E6%89%93%E5%8D%B0%E6%B5%81%E7%A8%8B)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# log
log 包作为 Go 标准库，仅支持日志的基本功能，不支持记录结构化日志、日志切割、Hook 等高级功能，所以更适合小型项目使用，比如一个单文件的脚本。

## 打印流程
```go
// go1.22.2/src/log/log.go

func Printf(format string, v ...any) {
	std.output(0, 2, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
}
```

使用标准 默认值
```go
// 标准错误输出,并且使用 日期+时间 属性格式
var std = New(os.Stderr, "", LstdFlags)


// 初始化
func New(out io.Writer, prefix string, flag int) *Logger {
	l := new(Logger)
	l.SetOutput(out)
	l.SetPrefix(prefix)
	l.SetFlags(flag)
	return l
}
```

`log.New` 函数接收三个参数，分别用来指定：日志输出位置（一个 `io.Writer` 对象）、日志前缀（字符串，每次打印日志都会跟随输出）、日志属性（定义好的常量)



```go
func (l *Logger) output(pc uintptr, calldepth int, appendOutput func([]byte) []byte) error {
	// 判断是否丢弃
	if l.isDiscard.Load() {
		return nil
	}
    // 获取当前时间
	now := time.Now() 

	// Load prefix and flag once so that their value is consistent within
	// this call regardless of any concurrent changes to their value.
	prefix := l.Prefix()
	flag := l.Flags()

	var file string
	var line int
	if flag&(Lshortfile|Llongfile) != 0 { // 通过位运算来判断是否需要获取文件名和行号
		if pc == 0 {
			var ok bool
			_, file, line, ok = runtime.Caller(calldepth)
			if !ok {
				file = "???"
				line = 0
			}
		} else {
			fs := runtime.CallersFrames([]uintptr{pc})
			f, _ := fs.Next()
			file = f.File
			if file == "" {
				file = "???"
			}
			line = f.Line
		}
	}
    // 复用对象 *[]byte
	buf := getBuffer()
	defer putBuffer(buf)
	// 格式化日志头信息（如：日期时间、文件名和行号、前缀）并写入 buf
	formatHeader(buf, now, prefix, flag, file, line)
	// 追加日志内容到 buf
	*buf = appendOutput(*buf)
	if len(*buf) == 0 || (*buf)[len(*buf)-1] != '\n' {
		*buf = append(*buf, '\n')
	}

	// 加锁，保证并发安全
	l.outMu.Lock()
	defer l.outMu.Unlock()
	// 调用 Logger 对象的 out 属性的 Write 方法输出日志
	_, err := l.out.Write(*buf)
	return err
}
```

## 源码分析

结构体
```go
type Logger struct {
	outMu sync.Mutex // 锁，保证并发情况下对其属性操作是原子性的
	out   io.Writer // destination for output

	prefix    atomic.Pointer[string] // 每行日志前缀 (but see Lmsgprefix)
	flag      atomic.Int32           // 属性 properties,用来控制日志输出格式
	isDiscard atomic.Bool // 当 out = io.Discard 是，此值为 true
}

```

属性 flag

```go
const (
	Ldate         = 1 << iota     // 当前时区日期: 2009/01/23
	Ltime                         // 当前时区时间: 01:23:23
	Lmicroseconds                 // 当前时区时间,精确到微秒 microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // 全文件名和行号: /a/b/c/d.go:23
	Lshortfile                    // 当前文件名和行号: d.go:23. overrides Llongfile
	LUTC                          // 使用UTC而不是本地市区
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime // 标准初始值
)

```


## 参考

- [深入探究 Go log 标准库](https://mp.weixin.qq.com/s?__biz=MzkzMjQ1NjkyNw==&mid=2247483722&idx=1&sn=898f0a4b868dea760f30d73513b935a8&chksm=c25a31faf52db8ecc52b395f8a06aeb0859c568eff89b87304a1d71a5a20b64636b93c00b527&scene=21#wechat_redirect)
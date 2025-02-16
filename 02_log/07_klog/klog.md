<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/kubernetes/klog](#githubcomkubernetesklog)
  - [主要](#%E4%B8%BB%E8%A6%81)
  - [初始化](#%E5%88%9D%E5%A7%8B%E5%8C%96)
  - [打印过程](#%E6%89%93%E5%8D%B0%E8%BF%87%E7%A8%8B)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/kubernetes/klog


forked from github.com/golang/glog.

需要开发功能的因素：

- glog 提出了很多“陷阱”，并介绍了容器化环境中的挑战，所有这些都没有得到很好的记录。
- glog 没有提供一种简单的方法来测试日志，这降低了使用它的软件的稳定性
- 长期目标是实现一个日志接口，该接口允许我们添加上下文，更改输出格式等


klog 包是一个提供了 INFO/ERROR/V 日志记录的模块。它提供了函数 Info、Warning、Error、Fatel，以及诸如 Infof 的格式化版本。它还提供了由 -v 和 -vmodule=file=2 标志控制的 V-style 日志。

日志先写到缓冲区中，并且周期性的使用 Flush 写入到文件中。程序在退出之前会调用 Flush 以保证写入所有日志。



## 主要 
```go
type loggingT struct {
	settings

	// flushD holds a flushDaemon that frequently flushes log file buffers.
	// Uses its own mutex.
	flushD *flushDaemon

	// mu protects the remaining elements of this structure and the fields
	// in settingsT which need a mutex lock.
	mu sync.Mutex

	// pcs is used in V to avoid an allocation when computing the caller's PC.
	pcs [1]uintptr
	// vmap is a cache of the V Level for each V() call site, identified by PC.
	// It is wiped whenever the vmodule flag changes state.
	vmap map[uintptr]Level
}
```

## 初始化

解析的变量
```go
func init() {
	commandLine.StringVar(&logging.logDir, "log_dir", "", "If non-empty, write log files in this directory (no effect when -logtostderr=true)")
	commandLine.StringVar(&logging.logFile, "log_file", "", "If non-empty, use this log file (no effect when -logtostderr=true)")
	commandLine.Uint64Var(&logging.logFileMaxSizeMB, "log_file_max_size", 1800,
		"Defines the maximum size a log file can grow to (no effect when -logtostderr=true). Unit is megabytes. "+
			"If the value is 0, the maximum file size is unlimited.")
	commandLine.BoolVar(&logging.toStderr, "logtostderr", true, "log to standard error instead of files")
	commandLine.BoolVar(&logging.alsoToStderr, "alsologtostderr", false, "log to standard error as well as files (no effect when -logtostderr=true)")
	logging.setVState(0, nil, false)
	commandLine.Var(&logging.verbosity, "v", "number for the log level verbosity")
	commandLine.BoolVar(&logging.addDirHeader, "add_dir_header", false, "If true, adds the file directory to the header of the log messages")
	commandLine.BoolVar(&logging.skipHeaders, "skip_headers", false, "If true, avoid header prefixes in the log messages")
	commandLine.BoolVar(&logging.oneOutput, "one_output", false, "If true, only write logs to their native severity level (vs also writing to each lower severity level; no effect when -logtostderr=true)")
	commandLine.BoolVar(&logging.skipLogHeaders, "skip_log_headers", false, "If true, avoid headers when opening log files (no effect when -logtostderr=true)")
	logging.stderrThreshold = severityValue{
		Severity: severity.ErrorLog, // Default stderrThreshold is ERROR.
	}
	commandLine.Var(&logging.stderrThreshold, "stderrthreshold", "logs at or above this threshold go to stderr when writing to files and stderr (no effect when -logtostderr=true or -alsologtostderr=true)")
	commandLine.Var(&logging.vmodule, "vmodule", "comma-separated list of pattern=N settings for file-filtered logging")
	commandLine.Var(&logging.traceLocation, "log_backtrace_at", "when logging hits line file:N, emit a stack trace")

	logging.settings.contextualLoggingEnabled = true
	logging.flushD = newFlushDaemon(logging.lockAndFlushAll, nil)
}

```


## 打印过程

```go
// k8s.io/klog/v2@v2.130.1/klog.go

var logging loggingT

// info 级别打印
func Info(args ...interface{}) {
	// logging.logger 为全局logging loggingT的 logger 字段
	logging.print(severity.InfoLog, logging.logger, logging.filter, args...)
}

```

```go
func (l *loggingT) print(s severity.Severity, logger *logWriter, filter LogFilter, args ...interface{}) {
	l.printDepth(s, logger, filter, 1, args...) // Caller depth 默认是1
}


func (l *loggingT) printDepth(s severity.Severity, logger *logWriter, filter LogFilter, depth int, args ...interface{}) {
	if false {
		_ = fmt.Sprint(args...) //  // cause vet to treat this function like fmt.Print
	}
    // header 写入
	buf, file, line := l.header(s, depth)
	l.printWithInfos(buf, file, line, s, logger, filter, depth+1, args...)
}



func (l *loggingT) printWithInfos(buf *buffer.Buffer, file string, line int, s severity.Severity, logger *logWriter, filter LogFilter, depth int, args ...interface{}) {
	// If a logger is set and doesn't support writing a formatted buffer,
	// we clear the generated header as we rely on the backing
	// logger implementation to print headers.
	if logger != nil && logger.writeKlogBuffer == nil {
		buffer.PutBuffer(buf)
		buf = buffer.GetBuffer()
	}
	if filter != nil { // 过滤字段
		args = filter.Filter(args)
	}
	fmt.Fprint(buf, args...)
	if buf.Len() == 0 || buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	// 输出到文件并且释放 buffer
	l.output(s, logger, buf, depth, file, line, false)
}


// output writes the data to the log files and releases the buffer.
func (l *loggingT) output(s severity.Severity, logger *logWriter, buf *buffer.Buffer, depth int, file string, line int, alsoToStderr bool) {
	var isLocked = true
	l.mu.Lock()
	defer func() {
		if isLocked {
			// Unlock before returning in case that it wasn't done already.
			l.mu.Unlock()
		}
	}()

	if l.traceLocation.isSet() {
		if l.traceLocation.match(file, line) {
			buf.Write(dbg.Stacks(false))
		}
	}
	data := buf.Bytes()
	if logger != nil {
		if logger.writeKlogBuffer != nil {
			logger.writeKlogBuffer(data)
		} else {
			if len(data) > 0 && data[len(data)-1] == '\n' {
				data = data[:len(data)-1]
			}
			// TODO: set 'severity' and caller information as structured log info
			// keysAndValues := []interface{}{"severity", severityName[s], "file", file, "line", line}
			if s == severity.ErrorLog {
				logger.WithCallDepth(depth+3).Error(nil, string(data))
			} else {
				logger.WithCallDepth(depth + 3).Info(string(data))
			}
		}
	} else if l.toStderr {
		os.Stderr.Write(data)
	} else {
		if alsoToStderr || l.alsoToStderr || s >= l.stderrThreshold.get() {
			os.Stderr.Write(data)
		}

		if logging.logFile != "" { // 如果指定了文件
			// Since we are using a single log file, all of the items in l.file array
			// will point to the same file, so just use one of them to write data.
			if l.file[severity.InfoLog] == nil {
				// 传入severity.InfoLog,所以只创建这级别的syncBuffer,如果第一次会写入相关 header 信息
				if err := l.createFiles(severity.InfoLog); err != nil {
					os.Stderr.Write(data) // Make sure the message appears somewhere.
					l.exit(err)
				}
			}
			_, _ = l.file[severity.InfoLog].Write(data)
		} else {
			if l.file[s] == nil {
				if err := l.createFiles(s); err != nil {
					os.Stderr.Write(data) // Make sure the message appears somewhere.
					l.exit(err)
				}
			}

			if l.oneOutput { 
				_, _ = l.file[s].Write(data)
			} else { // 级联写入:高严重等级的日志也将写入到低严重等级的日志文件中。
				
				switch s {
				case severity.FatalLog:
					_, _ = l.file[severity.FatalLog].Write(data)
					fallthrough
				case severity.ErrorLog:
					_, _ = l.file[severity.ErrorLog].Write(data)
					fallthrough
				case severity.WarningLog:
					_, _ = l.file[severity.WarningLog].Write(data)
					fallthrough
				case severity.InfoLog:
					_, _ = l.file[severity.InfoLog].Write(data)
				}
			}
		}
	}
	if s == severity.FatalLog {
		// If we got here via Exit rather than Fatal, print no stacks.
		if atomic.LoadUint32(&fatalNoStacks) > 0 {
			l.mu.Unlock()
			isLocked = false
			timeoutFlush(ExitFlushTimeout)
			OsExit(1)
		}
		// Dump all goroutine stacks before exiting.
		// First, make sure we see the trace for the current goroutine on standard error.
		// If -logtostderr has been specified, the loop below will do that anyway
		// as the first stack in the full dump.
		if !l.toStderr {
			os.Stderr.Write(dbg.Stacks(false))
		}

		// Write the stack trace for all goroutines to the files.
		trace := dbg.Stacks(true)
		logExitFunc = func(error) {} // If we get a write error, we'll still exit below.
		for log := severity.FatalLog; log >= severity.InfoLog; log-- {
			if f := l.file[log]; f != nil { // Can be nil if -logtostderr is set.
				_, _ = f.Write(trace)
			}
		}
		l.mu.Unlock()
		isLocked = false
		timeoutFlush(ExitFlushTimeout)
		OsExit(255) // C++ uses -1, which is silly because it's anded with 255 anyway.
	}
	buffer.PutBuffer(buf)

	if stats := severityStats[s]; stats != nil {
		atomic.AddInt64(&stats.lines, 1)
		atomic.AddInt64(&stats.bytes, int64(len(data)))
	}
}

```

hearder 格式
```go
Log lines have this form:
0 1 2 |3 4 
012345678901234567890123456789|01234567890 
Lmmdd hh:mm:ss.uuuuuu _______ |file:line] msg... 
                      threadid



前 29 字节固定，file 和 line 不定长


	L                表示日志级别，有 "IWEF"（例如，"I" 表示 INFO）
	mm               月份（零填充，例如五月是 "05"）
	dd               日期（零填充）
	hh:mm:ss.uuuuuu  Time in hours, minutes and fractional seconds
	threadid         The space-padded thread ID as returned by GetTID()
	file             The file name
	line             The line number
	msg              The user-supplied message


func (buf *Buffer) FormatHeader(s severity.Severity, file string, line int, now time.Time) {
	// ...
	_, month, day := now.Date()
	hour, minute, second := now.Clock()
	// Lmmdd hh:mm:ss.uuuuuu threadid file:line]
	buf.Tmp[0] = severity.Char[s]
	buf.twoDigits(1, int(month))
	buf.twoDigits(3, day)
	buf.Tmp[5] = ' '
	buf.twoDigits(6, hour)
	buf.Tmp[8] = ':'
	buf.twoDigits(9, minute)
	buf.Tmp[11] = ':'
	buf.twoDigits(12, second)
	buf.Tmp[14] = '.'
	buf.nDigits(6, 15, now.Nanosecond()/1000, '0')
	buf.Tmp[21] = ' '
	buf.nDigits(7, 22, Pid, ' ') // TODO: should be TID
	buf.Tmp[29] = ' '
	buf.Write(buf.Tmp[:30])
	buf.WriteString(file)
	buf.Tmp[0] = ':'
	n := buf.someDigits(1, line)
	buf.Tmp[n+1] = ']'
	buf.Tmp[n+2] = ' '
	buf.Write(buf.Tmp[:n+3])
}

```


## 参考


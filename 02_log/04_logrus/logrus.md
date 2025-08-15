<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/sirupsen/logrus](#githubcomsirupsenlogrus)
  - [特点](#%E7%89%B9%E7%82%B9)
  - [第三方使用-->calico](#%E7%AC%AC%E4%B8%89%E6%96%B9%E4%BD%BF%E7%94%A8--calico)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/sirupsen/logrus

它是一个结构化、插件化的日志记录库。完全兼容 golang 标准库中的日志模块。
它还内置了 2 种日志输出格式 JSONFormatter 和 TextFormatter，来定义输出的日志格式。

## 特点

- 与 Go log 标准库 API 完全兼容，这意味着任何使用 log 标准库的代码都可以将日志库无缝切换到 Logrus。



## 第三方使用-->calico

```go
// https://github.com/projectcalico/calico/blob/6135946c488e3d15741a746593c714b985205fc5/felix/logutils/logutils.go
func ConfigureLogging(configParams *config.Config) {
	// Parse the log levels, defaulting to panic if in doubt.
	logLevelScreen := logutils.SafeParseLogLevel(configParams.LogSeverityScreen)
	logLevelFile := logutils.SafeParseLogLevel(configParams.LogSeverityFile)
	logLevelSyslog := logutils.SafeParseLogLevel(configParams.LogSeveritySys)

	// Work out the most verbose level that is being logged.
	mostVerboseLevel := logLevelScreen
	if logLevelFile > mostVerboseLevel {
		mostVerboseLevel = logLevelFile
	}
	if logLevelSyslog > mostVerboseLevel {
		mostVerboseLevel = logLevelScreen
	}
	// 设置日志级别
	log.SetLevel(mostVerboseLevel)

	// Screen target.
	var dests []*logutils.Destination
	if configParams.LogSeverityScreen != "" {
		dests = append(dests, getScreenDestination(configParams, logLevelScreen))
	}

	// File target.  We record any errors so we can log them out below after finishing set-up
	// of the logger.
	var fileDirErr, fileOpenErr error
	if configParams.LogSeverityFile != "" && configParams.LogFilePath != "" {
		var destination *logutils.Destination
		destination, fileDirErr, fileOpenErr = getFileDestination(configParams, logLevelFile)
		if fileDirErr == nil && fileOpenErr == nil && destination != nil {
			dests = append(dests, destination)
		}
	}

	// Syslog target.  Again, we record the error if we fail to connect to syslog.
	var sysErr error
	if configParams.LogSeveritySys != "" {
		var destination *logutils.Destination
		destination, sysErr = getSyslogDestination(configParams, logLevelSyslog)
		if sysErr == nil && destination != nil {
			dests = append(dests, destination)
		}
	}

	// 设置 hook 
	hook := logutils.NewBackgroundHook(
		logutils.FilterLevels(mostVerboseLevel),
		logLevelSyslog,
		dests,
		counterDroppedLogs,
		logutils.WithDebugFileRegexp(configParams.LogDebugFilenameRegex),
	)
	hook.Start()
	log.AddHook(hook)

	// Disable logrus' default output, which only supports a single destination.  We use the
	// hook above to fan out logs to multiple destinations.
	log.SetOutput(&logutils.NullWriter{})

	// Do any deferred error logging.
	if fileDirErr != nil {
		log.WithError(fileDirErr).WithField("file", configParams.LogFilePath).
			Fatal("Failed to create log file directory.")
	}
	if fileOpenErr != nil {
		log.WithError(fileOpenErr).WithField("file", configParams.LogFilePath).
			Fatal("Failed to open log file.")
	}
	if sysErr != nil {
		// We don't bail out if we can't connect to syslog because our default is to try to
		// connect but it's very common for syslog to be disabled when we're run in a
		// container.
		log.WithError(sysErr).Error(
			"Failed to connect to syslog. To prevent this error, either set config " +
				"parameter LogSeveritySys=none or configure a local syslog service.")
	}
}

```

hook

```go
// https://github.com/projectcalico/calico/blob/92d22390223e8aac663bf113d7ee932300696678/libcalico-go/lib/logutils/logutils.go
func NewBackgroundHook(
	levels []log.Level,
	syslogLevel log.Level,
	destinations []*Destination,
	counter prometheus.Counter,
	opts ...BackgroundHookOpt,
) *BackgroundHook {
	bh := &BackgroundHook{
		destinations: destinations,
		levels:       levels,
		syslogLevel:  syslogLevel,
		counter:      counter,
	}
	for _, opt := range opts {
		opt(bh)
	}
	return bh
}

func (h *BackgroundHook) Levels() []log.Level {
	return h.levels
}

func (h *BackgroundHook) Fire(entry *log.Entry) (err error) {
	if entry.Buffer != nil {
		defer entry.Buffer.Truncate(0)
	}

	if entry.Level >= log.DebugLevel && h.debugFileNameRE != nil {
		// This is a debug log, check if debug logging is enabled for this file.
		fileName, _ := getFileInfo(entry)
		if fileName == FileNameUnknown || !h.debugFileNameRE.MatchString(fileName) {
			return nil
		}
	}

	var serialized []byte
	if serialized, err = entry.Logger.Formatter.Format(entry); err != nil {
		return
	}

	// entry's buffer will be reused after we return but we're about to send the message over
	// a channel so we need to take a copy.
	bufCopy := make([]byte, len(serialized))
	copy(bufCopy, serialized)

	ql := QueuedLog{
		Level:   entry.Level,
		Message: bufCopy,
	}

	if entry.Level <= h.syslogLevel {
		// syslog gets its own log string since our default log string duplicates a lot of
		// syslog metadata.  Only calculate that string if it's needed.
		ql.SyslogMessage = FormatForSyslog(entry)
	}

	var waitGroup *sync.WaitGroup
	if entry.Level <= log.FatalLevel || entry.Data[FieldForceFlush] == true {
		// If the process is about to be killed (or we're asked to do so), flush the log.
		waitGroup = &sync.WaitGroup{}
		ql.WaitGroup = waitGroup
	}

	for _, dest := range h.destinations {
		if ql.Level > dest.Level {
			continue
		}
		if waitGroup != nil {
			// Thread safety: we must call add before we send the wait group over the
			// channel (or the background thread could be scheduled immediately and
			// call Done() before we call Add()).  Since we don't know if the send
			// will succeed that leads to the need to call Done() on the 'default:'
			// branch below to correctly pair Add()/Done() calls.
			waitGroup.Add(1)
		}

		if ok := dest.Send(ql); !ok { // 如果阻塞,返回 false 
			// Background thread isn't keeping up.  Drop the log and count how many
			// we've dropped.
			if waitGroup != nil {
				waitGroup.Done()
			}
			// Increment the number of dropped logs
			dest.counter.Inc()
		}
	}
	if waitGroup != nil {
		waitGroup.Wait()
	}
	return
}

func (h *BackgroundHook) Start() {
	for _, d := range h.destinations {
		go d.LoopWritingLogs()
	}
}
```


## 参考

- [Go 每日一库之 logrus](https://darjun.github.io/2020/02/07/godailylib/logrus/)
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/golang/glog](#githubcomgolangglog)
  - [特点](#%E7%89%B9%E7%82%B9)
  - [Flush Daemon](#flush-daemon)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/golang/glog

Glog是著名google开源C++日志库glog的golang版本，具有轻量级、简单、稳定和高效等特性。 目前被用在大型的容器云开源项目Kubernetes中



##  特点

- 支持四种日志等级INFO < WARING < ERROR < FATAL，不支持DEBUG等级。
- 每个日志等级对应一个日志文件，低等级的日志文件中除了包含该等级的日志，还会包含高等级的日志。
- 日志文件可以根据大小切割，但是不能根据日期切割。
- 日志文件名称格式：program.host.userName.log.log_level.date-time.pid，不可自定义。
- 固定日志输出格式：Lmmdd hh:mm:ss.uuuuuu threadid file:line] msg…，不可自定义。
- 程序开始时必须调用flag.Parse()解析命令行参数，退出时必须调用glog.Flush()确保将缓存区日志输出


Glog的代码主要分两个文件：

- glog.go：主要实现log等级定义、输出以及vlog。
- glog_file.go：主要实现日志文件目录和各等级日志文件的创建


## Flush Daemon


Glog在初始化的时候，会定义一些命令行参数，同时启动flush守护进程。Flush守护进程会间隔30s周期性地flush缓冲区中的log
```go
// github.com/golang/glog@v1.1.0/glog_file.go

func init() {
    // ...
	go sinks.file.flushDaemon()
}

func (s *fileSink) flushDaemon() {
	tick := time.NewTicker(30 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			s.Flush()
		case sev := <-s.flushChan:
			s.flush(sev)
		}
	}
}



// Flush flushes all the logs and attempts to "sync" their data to disk.
func (s *fileSink) Flush() error {
	return s.flush(logsink.Info)
}

// flush flushes all logs of severity threshold or greater.
func (s *fileSink) flush(threshold logsink.Severity) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var firstErr error
	updateErr := func(err error) {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	// Flush from fatal down, in case there's trouble flushing.
	for sev := logsink.Fatal; sev >= threshold; sev-- {
		file := s.file[sev]
		if file != nil {
			updateErr(file.Flush())
			updateErr(file.Sync())
		}
	}

	return firstErr
}

```

## 参考

- https://github.com/google/glog

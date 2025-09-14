package main

import (
	"flag"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()
	glog.Infof("no level message")
	glog.V(3).Info("LEVEL 3 message") // 使用日志级别 3
	glog.V(4).Info("LEVEL 4 message") // 使用日志级别 4
	glog.V(5).Info("LEVEL 5 message") // 使用日志级别 5
	glog.V(8).Info("LEVEL 8 message") // 使用日志级别 8
}

// # 日志级别小于或等于 4 的日志将被打印出来：
// $ go run main.go -v=4 -log_dir=log -alsologtostderr
// I0216 11:53:12.421374   88707 main.go:12] LEVEL 3 message
// I0216 11:53:12.422333   88707 main.go:13] LEVEL 4 message

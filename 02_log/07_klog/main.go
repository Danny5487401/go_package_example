package main

import (
	"flag"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/textlogger"
)

func main() {
	klog.InitFlags(nil) // 注册 flag.CommandLine
	// By default klog writes to stderr. Setting logtostderr to false makes klog
	// write to a log file.
	flag.Set("logtostderr", "false")                  // 日志输出到stderr，不输出到日志文件。false为关闭
	flag.Set("alsologtostderr", "true")               // 日志输出到stderr，不输出到日志文件。false为关闭
	flag.Set("log_file", "02_log/07_klog/myfile.log") // 设置文件路径
	flag.Set("v", "4")                                // log level 级别
	flag.Parse()
	klog.Infof("nice to meet u:%s", "danny") // I0218 15:13:50.214041   30383 main.go:18] nice to meet u :danny
	klog.V(3).Info("klogv3")

	vLevel := 5
	config := textlogger.NewConfig(textlogger.Verbosity(vLevel))
	log := textlogger.NewLogger(config).WithName("MyName").WithValues("user", "you")
	log.Info("hello", "val1", 1, "val2", map[string]int{"k": 1})                    // I0218 15:13:50.214041   30383 main.go:18] nice to meet u :danny
	log.V(3).Info("nice to meet you", "name", "joy")                                // I0218 15:24:53.992436   31066 main.go:24] "nice to meet you" logger="MyName" user="you" name="joy"
	log.Error(nil, "uh oh", "trouble", true, "reasons", []float64{0.1, 0.11, 3.14}) // E0218 15:13:50.221743   30383 main.go:25] "uh oh" logger="MyName" user="you" trouble=true reasons=[0.1,0.11,3.14]
	log.Error(myError{"an error occurred"}, "goodbye", "code", -1)                  // E0218 15:13:50.221762   30383 main.go:26] "goodbye" err="an error occurred" logger="MyName" user="you" code=-1
	klog.Flush()
}

type myError struct {
	str string
}

func (e myError) Error() string {
	return e.str
}

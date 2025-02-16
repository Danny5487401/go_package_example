package main

import (
	"flag"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/textlogger"
)

func main() {
	klog.InitFlags(nil)
	// By default klog writes to stderr. Setting logtostderr to false makes klog
	// write to a log file.
	flag.Set("logtostderr", "false")                  //日志输出到stderr，不输出到日志文件。false为关闭
	flag.Set("alsologtostderr", "true")               //日志输出到stderr，不输出到日志文件。false为关闭
	flag.Set("log_file", "02_log/07_klog/myfile.log") // 设置文件路径
	//flag.Set("v", "3")                                        // log level 级别
	flag.Parse()
	klog.Info("nice to meet you")
	klog.Flush()

	config := textlogger.NewConfig(textlogger.Verbosity(1))
	log := textlogger.NewLogger(config).WithName("MyName").WithValues("user", "you")
	log.Info("hello", "val1", 1, "val2", map[string]int{"k": 1})
	log.V(3).Info("nice to meet you")
	log.Error(nil, "uh oh", "trouble", true, "reasons", []float64{0.1, 0.11, 3.14})
	log.Error(myError{"an error occurred"}, "goodbye", "code", -1)
	klog.Flush()
}

type myError struct {
	str string
}

func (e myError) Error() string {
	return e.str
}

package main

import (
	"flag"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	bar()
	glog.V(3).Info("LEVEL 3 message")
	glog.V(5).Info("LEVEL 5 message")
}

/*
$ go run main.go extra.go -v=3 -log_dir=log -alsologtostderr -vmodule=extra=5 -log_backtrace_at=extra.go:6
I0216 12:20:16.166016   90062 extra.go:6] LEVEL 4: level 4 message in bar.go

goroutine 1 [running]:
github.com/golang/glog.Verbose.Info(...)
/Users/python/go/pkg/mod/github.com/golang/glog@v1.1.0/glog.go:387
main.bar()
/Users/python/Downloads/git_download/go_package_example/02_log/06_glog/02_vmodule_and_trace_location/extra.go:6 +0x74
main.main()
/Users/python/Downloads/git_download/go_package_example/02_log/06_glog/02_vmodule_and_trace_location/main.go:12 +0x80

I0216 12:20:16.167927   90062 main.go:13] LEVEL 3 message
*/

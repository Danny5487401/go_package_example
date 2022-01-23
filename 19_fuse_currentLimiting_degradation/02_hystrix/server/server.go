package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"net/http"
	"time"
)

type Handle struct{}

func (h *Handle) ServeHTTP(r http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	fmt.Println("请求的路径是", path)
	switch path {
	case "/foo":
		h.async(r, request)
	default:
		h.synchronize(r, request)
	}

}

func (h *Handle) async(r http.ResponseWriter, request *http.Request) {
	output := make(chan bool, 1)
	errors := hystrix.Go("myCommand", func() error {
		// talk to other services
		output <- true
		return nil
	}, nil)

	select {
	case out := <-output:
		_ = out
		msg := "异步success"
		r.Write([]byte(msg))
		// success
	case err := <-errors:
		msg := "同步failure" + err.Error()
		// failure
		r.Write([]byte(msg))
	}
}

func (h *Handle) synchronize(r http.ResponseWriter, request *http.Request) {

	msg := "同步success"
	// 同步的方法
	_ = hystrix.Do("myCommand", func() error {
		_, err := http.Get("https://www.baidu.com")
		if err != nil {
			fmt.Printf("请求失败:%v", err)
			return err
		}
		return nil
	}, func(err error) error {
		fmt.Printf("handle  error:%v\n", err)
		msg = "异步error"
		return nil
	})
	r.Write([]byte(msg))
}

func main() {
	// 配置策略-mycommand为全局
	hystrix.ConfigureCommand("myCommand", hystrix.CommandConfig{
		Timeout:                int(3 * time.Second), // 执行 command 的超时时间。
		MaxConcurrentRequests:  10,                   // MaxConcurrentRequests：command 的最大并发量
		SleepWindow:            5000,                 // 当熔断器被打开后，SleepWindow 的时间就是控制过多久后去尝试服务是否可用了。
		RequestVolumeThreshold: 20,                   // 一个统计窗口 10 秒内请求数量。达到这个请求数量后才去判断是否要开启熔断
		ErrorPercentThreshold:  30,                   // 错误百分比，请求数量大于等于 RequestVolumeThreshold 并且错误率到达这个百分比后就会启动熔断
	})

	http.ListenAndServe(":8090", &Handle{})
}

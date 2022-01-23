package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var client *http.Client

func init() {
	tr := &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 1 * time.Second,
	}
	client = &http.Client{Transport: tr}
}

type info struct {
	Data interface{} `json:"data"`
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 11; i++ {
		wg.Add(1)
		go func(int2 int) {
			defer wg.Done()
			//req, err := http.NewRequest("GET", "http://localhost:8090", nil)
			req, err := http.NewRequest("GET", "http://localhost:8090/foo", nil)
			if err != nil {
				fmt.Printf("初始化http客户端处错误:%v", err)
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("初始化http客户端处错误:%v", err)
				return
			}
			defer resp.Body.Close()
			nByte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("读取http数据失败:%v", err)
				return
			}
			fmt.Printf("接收到到值:%v\n", string(nByte))
		}(i)
	}
	wg.Wait()

	fmt.Printf("请求完毕\n")
}

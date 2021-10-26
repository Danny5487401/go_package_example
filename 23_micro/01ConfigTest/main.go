package main

import (
	"fmt"
	"github.com/micro/go-micro/config/encoder/json"
	"github.com/micro/go-micro/config/source"
	"time"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
)

func main() {
	//指定编码器

	// 加载配置文件
	if err := config.Load(file.NewSource(
		file.WithPath("23_micro/01ConfigTest/config/config.json"),
		source.WithEncoder(json.NewEncoder()),
	)); err != nil {
		fmt.Println(err)
		return
	}

	var mysql Mysql

	if err := config.Get("mysql", "database").Scan(&mysql); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mysql.Name, mysql.Address, mysql.Port)
	// 监控变化
	go func() {
		for range time.Tick(time.Second) {
			conf := config.Map()
			fmt.Println(conf)
		}
	}()
	select {}
}

// 定义我们的额数据结构
type Mysql struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

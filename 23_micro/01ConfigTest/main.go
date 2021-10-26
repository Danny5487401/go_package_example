package main

import (
	"fmt"
	"github.com/micro/go-micro/config/encoder/json"
	"github.com/micro/go-micro/config/encoder/yaml"
	"github.com/micro/go-micro/config/source"
	"time"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
)

func main() {

	// 配置来源
	jsonSource := file.NewSource(
		//从文件中读取
		file.WithPath("23_micro/01ConfigTest/config/config.json"),
		//指定json编码器
		source.WithEncoder(json.NewEncoder()))

	yamlSource := file.NewSource( //从文件中读取
		file.WithPath("23_micro/01ConfigTest/config/config.yaml"),
		//指定json编码器
		source.WithEncoder(yaml.NewEncoder()))

	// 后面读取的优先级越高，所以yaml的配置会覆盖json的配置
	if err := config.Load(jsonSource, yamlSource); err != nil {
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
			fmt.Printf("配置是%+v\n", conf)
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

package main

import (
	"fmt"
	"time"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
)

/*
Config 特征
	1 动态加载：根据需要动态加载多个资源文件。 go config 在后台管理并监控配置文件，并自动更新到内存中
	2 资源可插拔： 从任意数量的源中进行选择以加载和合并配置。后台资源源被抽象为内部使用的标准格式，并通过编码器进行解码。源可以是环境变量，标志，文件，etcd，k8s configmap等
	3 可合并的配置：如果指定了多种配置，它们会合并到一个视图中。
	4 监控变化：可以选择是否监控配置的指定值，热重启
	5 安全恢复： 万一配置加载错误或者由于未知原因而清除，可以指定回退值进行回退

*/

func main() {
	// 加载配置文件
	if err := config.Load(file.NewSource(
		file.WithPath("23_micro/config/config.json"),
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

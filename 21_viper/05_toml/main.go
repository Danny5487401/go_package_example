package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
)

func main() {
	config, _ := toml.LoadFile("21_viper/05_toml/test.toml") //加载toml文件
	os := config.Get("os").(string)                          //读取key对应的值.括号为指定数据类型，也可以忽略
	fmt.Println(os)

	fmt.Println("string 描述: ", config.Get("server.desc"))
	fmt.Println("route 路由: ", config.GetArray("server.route"))

}

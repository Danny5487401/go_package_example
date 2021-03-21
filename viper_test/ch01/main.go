package main

import (
	"fmt"
	"github.com/spf13/viper"
)

// 方式二 映射
type ServerConfig struct {
	ServiceName string `mapstructure:"name"`
	Port int `mapstructure:"port"`

}

func main (){
	v := viper.New()
	v.SetConfigFile("viper_test/ch01/config.yaml")  // 注意路径问题  看goland edit configuration:working directory
	if err :=v.ReadInConfig();err!=nil{
		panic(err)
	}
	serverConfig:=ServerConfig{}
	if err := v.Unmarshal(&serverConfig);err!=nil{
		panic(err)
	}
	fmt.Println(serverConfig)
	// 方式一
	fmt.Printf("%v",v.Get("name"))
}

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	remote "github.com/yoyofxteam/nacos-viper-remote"
)

var appName string

func main() {
	configViper := viper.New()

	// 注册 nacos 配置
	remote.SetOptions(&remote.Option{
		Url:         "tencent.danny.games",
		Port:        8848,
		NamespaceId: "public",
		GroupName:   "DEFAULT_GROUP",
		Config:      remote.Config{DataId: "config_dev.yaml"},
		Auth: &remote.Auth{
			User:     "nacos",
			Password: "nacos"},
	})

	// 初始化远程 viper
	remoteViper := viper.New()
	err := remoteViper.AddRemoteProvider("nacos", "tencent.danny.games", "")
	remoteViper.SetConfigType("yaml")
	// 加载远程配置
	err = remoteViper.ReadRemoteConfig() //sync get remote configs to remoteViper instance memory . for example , remoteViper.GetString(key)

	if err == nil {
		configViper = remoteViper
		fmt.Println("监听远程配置")
		provider := remote.NewRemoteProvider("yaml")
		respChan := provider.WatchRemoteConfigOnChannel(configViper)

		go func(rc <-chan bool) {
			for {
				<-rc
				fmt.Printf("监听到配置: %s\n", configViper.GetString("app.age"))
			}
		}(respChan)
	}

	go func() {
		for {
			time.Sleep(time.Second * 30) // delay after each request
			appName = configViper.GetString("app.age")
			fmt.Println("每次拉取配置:" + appName)
		}
	}()

	onExit()

}

// 监听信号
func onExit() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("Program Exit...", s)

		default:
			fmt.Println("other signal", s)
		}
	}
}

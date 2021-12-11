package main

import (
	"fmt"
	"github.com/spf13/viper"
	remote "github.com/yoyofxteam/nacos-viper-remote"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var appName string

func main() {
	config_viper := viper.New()

	remote.SetOptions(&remote.Option{
		Url:         "tencent.danny.games",
		Port:        8848,
		NamespaceId: "public",
		GroupName:   "DEFAULT_GROUP",
		Config:      remote.Config{DataId: "config_dev"},
		Auth: &remote.Auth{User: "nacos",
			Password: "nacos"},
	})

	remoteViper := viper.New()
	err := remoteViper.AddRemoteProvider("nacos", "tencent.danny.games", "")
	remoteViper.SetConfigType("yaml")
	err = remoteViper.ReadRemoteConfig() //sync get remote configs to remoteViper instance memory . for example , remoteViper.GetString(key)

	if err == nil {
		config_viper = remoteViper
		fmt.Println("使用远程配置")
		provider := remote.NewRemoteProvider("yaml")
		respChan := provider.WatchRemoteConfigOnChannel(config_viper)

		go func(rc <-chan bool) {
			for {
				<-rc
				fmt.Printf("监听到配置: %s\n", config_viper.GetString("app.age"))
			}
		}(respChan)
	}

	go func() {
		for {
			time.Sleep(time.Second * 30) // delay after each request
			appName = config_viper.GetString("app.age")
			fmt.Println("每次拉取配置:" + appName)
		}
	}()

	onExit()

}

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

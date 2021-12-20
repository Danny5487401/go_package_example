package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

type MysqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructue:"port"`
}
type ServerConfig struct {
	ServiceName string      `mapstructure:"name"`
	MysqlInfo   MysqlConfig `mapstructure:"mysql"`
}

func GetEnvInfo(env string) bool {
	// 获取环境变量区分环境
	viper.AutomaticEnv()
	return viper.GetBool(env) //必须重启goland才行
}

func main() {
	// 获取环境变量
	debug := GetEnvInfo("ENV_DEBUG")
	//将线上线下文件配置文件隔离
	configFileName := "21_viper/01_read_n_watch_config/config-prod.yaml"
	if debug {
		configFileName = "21_viper/01_read_n_watch_config/config-debug.yaml"
	}

	// 初始化
	v := viper.New()

	// 指定读取配置文件
	v.SetConfigFile(configFileName) // 注意路径问题  看goland edit configuration:working directory

	// 查找并读取配置文件，
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	// 获取配置
	serverConfig := ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println(serverConfig)

	// 监听动态变化
	v.WatchConfig()
	v.OnConfigChange(func(n fsnotify.Event) {
		fmt.Println("configFile change:", n.Name)
		_ = v.ReadInConfig()
		serverConfig := ServerConfig{}
		if err := v.Unmarshal(&serverConfig); err != nil {
			panic(err)
		}
		fmt.Println(serverConfig)
	})

	time.Sleep(time.Second * 300)
}

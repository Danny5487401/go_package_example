package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type MysqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructue:"port"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructue:"port"`
	Db   int    `mapstructure:"db"`
}

type ServerConfig struct {
	ServiceName string      `mapstructure:"name"`
	MysqlInfo   MysqlConfig `mapstructure:"mysql"`
	RedisInfo   RedisConfig `mapstructure:"redis"`
	Author      string      `mapstructure:"author"`
}

func GetEnvInfo(env string) bool {
	// 获取环境变量区分环境
	viper.AutomaticEnv()

	return viper.GetBool(env) //必须重启goland才行
}

func main() {
	// 初始化
	v := viper.New()

	// 默认值
	v.SetDefault("ContentDir", "content")
	v.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 12)

	// flag
	pflag.String("author", "", "--author xxx")

	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)

	// 获取环境变量
	os.Setenv("Name", "Danny")
	//默认是AllowEmptyEnv(false)，这里设置为true
	v.AllowEmptyEnv(true)
	debug := GetEnvInfo("ENV_DEBUG")
	//将线上线下文件配置文件隔离
	configFileName := "21_viper/01_read_n_watch_config/config-prod.yaml"
	if debug {
		configFileName = "21_viper/01_read_n_watch_config/config-debug.yaml"
	}

	// 指定读取配置文件
	v.SetConfigFile(configFileName) // 注意路径问题  看goland edit configuration:working directory

	// 查找并读取配置文件
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
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

	time.Sleep(time.Second * 60)
}

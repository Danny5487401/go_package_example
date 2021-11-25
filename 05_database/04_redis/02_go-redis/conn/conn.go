package conn

import (
	"fmt"
	"go_grpc_example/05_database/04_redis/02_go-redis/hooks"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
)

func GetRedisDB() *redis.Client {
	if RedisDB == nil {
		initRedisEngine()
	}
	return RedisDB
}

type RedisConfig struct {
	Addr         string `yaml:"addr"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	Password     string `yaml:"password"`
	IdleTimeout  int    `yam:"idle_timeout"`
	DialTimeout  int    `yam:"dial_timeout"`
	ReadTimeout  int    `yam:"read_timeout"`
	WriteTimeout int    `yam:"write_timeout"`
}

var RedisDB *redis.Client
var wg sync.WaitGroup

func initRedisEngine() {
	conf := &RedisConfig{
		Addr:         "ali.danny.games:6379",
		Password:     "root",
		DB:           2,
		PoolSize:     10, // 连接池大小
		IdleTimeout:  30, // 客户端关闭空闲连接的时间
		DialTimeout:  1,
		ReadTimeout:  1,
		WriteTimeout: 1,
	}

	//连接服务器
	RedisDB = redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		DB:           conf.DB,
		PoolSize:     conf.PoolSize,
		Password:     conf.Password,
		IdleTimeout:  time.Duration(conf.IdleTimeout) * time.Second,
		DialTimeout:  time.Duration(conf.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
	})
	//心跳
	pong, err := RedisDB.Ping().Result()
	if err != nil {
		fmt.Println("连接失败", err.Error())
		os.Exit(1)
	}
	RedisDB.AddHook(hooks.NewLogHook())
	fmt.Println(pong) // Output: PONG

}

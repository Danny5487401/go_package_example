package main

import (
	"fmt"
	"github.com/imdario/mergo"
	"log"
)

type redisConfig struct {
	Address string
	Port    int
	DBs     []int
}

var defaultConfig = redisConfig{
	Address: "127.0.0.1",
	Port:    6381,
}

func main() {
	var config redisConfig
	config.DBs = []int{2, 3}

	// WithOverrideEmptySlice: 源对象的空切片覆盖目标对象的对应字段
	if err := mergo.Merge(&config, defaultConfig, mergo.WithOverride, mergo.WithOverrideEmptySlice); err != nil {
		log.Fatal(err)
	}

	fmt.Println("redis address: ", config.Address)
	fmt.Println("redis port: ", config.Port)
	fmt.Println("redis dbs: ", config.DBs)
}

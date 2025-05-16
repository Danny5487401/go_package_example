package main

import (
	"fmt"
	"log"

	"github.com/imdario/mergo"
)

type redisConfig struct {
	Address string
	Port    int
	DB      int
}

var defaultConfig = redisConfig{
	Address: "127.0.0.1",
	Port:    6381,
	DB:      1,
}

func main() {
	var config = redisConfig{
		Address: "192.168.0.0.1", // 有值默认不覆盖
	}

	if err := mergo.Merge(&config, defaultConfig, mergo.WithOverride); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("redis info address: %+v \n", config) // redis info address: {Address:127.0.0.1 Port:6381 DB:1}

	var m = map[string]interface{}{
		"address": "192.168.0.0.1", // 有值默认不覆盖, 要求首字母小写
	}
	if err := mergo.Map(&m, defaultConfig); err != nil {
		log.Fatal(err)
	}

	fmt.Println(m) // map[address:192.168.0.0.1 dB:1 port:6381]
}

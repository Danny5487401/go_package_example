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
	var config redisConfig

	if err := mergo.Merge(&config, defaultConfig); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("redis info address: %+v \n", config) // redis info address: {Address:127.0.0.1 Port:6381 DB:1}

	var m = make(map[string]interface{})
	if err := mergo.Map(&m, defaultConfig); err != nil {
		log.Fatal(err)
	}

	fmt.Println(m) // map[address:127.0.0.1 dB:1 port:6381]
}

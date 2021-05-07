package main

// 读写json到redis

import (
	// "github.com/garyburd/redigo/redis" // 旧路径
	"fmt"

	"encoding/json"
	"github.com/gomodule/redigo/redis" // 新路径
)

func main() {
	conn, err := redis.Dial("tcp", "106.14.35.115:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 写入redis
	key := "profile"
	imap := map[string]string{"username": "danny", "phoneNumber": "888"}
	value, _ := json.Marshal(imap)

	n, err := conn.Do("SETNX", key, value)
	if err != nil {
		fmt.Println(err)
	}
	if n == int64(1) {
		fmt.Println("success")
	}

	// 读取redis
	var imapGet map[string]string

	valueGet, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		fmt.Println(err)
	}

	errShal := json.Unmarshal(valueGet, &imapGet)
	if errShal != nil {
		fmt.Println(err)
	}
	fmt.Println(imapGet["username"])
	fmt.Println(imapGet["phoneNumber"])
}

package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
)

//获取redis连接池
func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: time.Duration(24) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "106.14.35.115:6379")
			if err != nil {
				panic(err.Error())
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return err
			}
			return err
		},
	}
}

func setValue() {
	conn := pool.Get()
	defer conn.Close()
	reply, err := conn.Do("SET", "age", "18", "EX", "30")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
	fmt.Println(reply)
}

func main() {
	pool = newPool()
	setValue()
}

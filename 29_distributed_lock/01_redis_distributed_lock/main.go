package main

// 使用redis实现分布式锁

import (
	"context"
	"fmt"

	goredislib "github.com/go-redis/redis/v7"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v7"
)

func main() {

	client := goredislib.NewClient(&goredislib.Options{
		Addr:     "106.14.35.115:6379",
		DB:       2,
		PoolSize: 10, // 连接池大小
		Password: "root",
	})

	pool := goredis.NewPool(client)

	rs := redsync.New(pool)

	mutex := rs.NewMutex("test-redsync")
	ctx := context.Background()

	if err := mutex.LockContext(ctx); err != nil {
		panic(err)
	}
	fmt.Println("加锁成功")
	if _, err := mutex.UnlockContext(ctx); err != nil {
		panic(err)
	}
	fmt.Println("解锁成功")
}

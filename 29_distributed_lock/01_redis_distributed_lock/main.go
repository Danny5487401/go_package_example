package main

// 使用redis实现分布式锁

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/exp/rand"
	"log"
	"sync"
	"time"
)

/*
实现思路
	获取：通过redis的setnx命令去获取锁，如果成功了，定时的更新锁的过期时间。
	释放：使用lua的脚本实现原子性的操作，保证，释放的确实是自己的锁。
*/

func incr(i int) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"106.14.35.115:6379"},
		//Username: "",
		Password: "root",
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("初始化错误%v", err)

	}

	// lua的脚本
	scr := redis.NewScript(`if redis.call('get',KEYS[1]) == ARGV[1] 
									then return redis.call('del',KEYS[1]) 
									else return 0 
								end `)

	sha, err := scr.Load(client.Context(), client).Result()
	if err != nil {
		fmt.Println("加载lua脚本错误", err)
		return
	}

	var lockKey = "counter_lock"
	var value = "abc" + fmt.Sprintf("%d", i)
	ctx := context.Background()
	sta := 0

	delChan := make(chan struct{})

	go func() {
		for {
			if sta == 1 {
				fmt.Printf("第%v台机器开始lock\n", i)
				doSomeThing(ctx)
				sta = 0
				delChan <- struct{}{}
			}
			time.Sleep(time.Second * 2)

		}
	}()

	// lock：每两秒去更新
	tick := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-delChan:
			ret := client.EvalSha(client.Context(), sha, []string{
				lockKey,
			}, value)
			if result, err := ret.Result(); err != nil || result == 0 {
				fmt.Println(fmt.Sprintf("第%d个机器 释放锁失败,%v,%v", i, result, err))
				sta = 0
			} else {
				fmt.Sprintf("第%d个机器 释放锁成功", i)
			}
			continue
		case <-tick.C:
			if sta == 0 {
				resp := client.SetNX(ctx, lockKey, value, time.Second*10)
				lockSuccess, err := resp.Result()
				if err != nil || !lockSuccess {
					continue
				}
				sta = 1
			} else {
				resp := client.SetEX(ctx, lockKey, 1, time.Second*10)
				sr, err := resp.Result()
				if err != nil {
					fmt.Println(i, "更新锁失败: ", err)
					continue
				}
				fmt.Println("更新锁成功", sr, err)
			}
		default:
			time.Sleep(time.Second * 1)

		}

	}

}

// 执行业务
func doSomeThing(ctx context.Context) {
	num := time.Duration(rand.Int63n(10))
	fmt.Printf("将执行%v 秒\n", int64(num))
	time.Sleep(time.Second * num)
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(in int) {
			defer wg.Done()
			incr(in)
		}(i)
	}
	wg.Wait()
}
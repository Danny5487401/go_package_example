package main

import (
	"fmt"
	"time"

	"go.uber.org/ratelimit"
)

func main() {
	rl := ratelimit.New(100, ratelimit.WithoutSlack) // 给定限流器每秒可以通过 100 个请求，也就是平均每个请求间隔 10ms

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}

	// Output:
	// 0 0
	// 1 10ms
	// 2 10ms
	// 3 10ms
	// 4 10ms
	// 5 10ms
	// 6 10ms
	// 7 10ms
	// 8 10ms
	// 9 10ms
}

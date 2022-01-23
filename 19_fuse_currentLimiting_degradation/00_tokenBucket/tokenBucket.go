package main

import (
	"golang.org/x/time/rate"

	"context"
	"fmt"
	"time"
)

func main() {
	// 1。构建限流器对象
	// r代表每秒可以向 Token 桶中产生10 token,b 代表 Token 桶的容量大小100
	// 方式一
	//limiter := rate.NewLimiter(10, 100)
	// 方式二
	//  Every 方法来指定向桶中放置 Token 的间隔
	limit := rate.Every(100 * time.Millisecond) // 表示每 100ms 往桶中放一个 Token。本质上也是一秒钟往桶里放 10 个
	limiter := rate.NewLimiter(limit, 100)

	// 2。使用限流器
	// Limiter 提供了三类方法供程序消费 Token，可以每次消费一个 Token，也可以一次性消费多个 Token。
	//每种方法代表了当 Token 不足时，各自不同的对应手段，可以阻塞等待桶中Token补充，也可以直接返回取Token失败。
	//func (lim *Limiter) Wait(ctx context.Context) (err error)
	//Wait 实际上就是 WaitN(ctx,1)
	//func (lim *Limiter) WaitN(ctx context.Context, n int) (err error)
	// 设置一秒的等待超时时间
	ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
	err := limiter.Wait(ctx)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// func (lim *Limiter) Allow() bool
	//func (lim *Limiter) AllowN(now time.Time, n int) bool
	// AllowN 方法表示，截止到某一时刻，目前桶中数目是否至少为 n 个，满足则返回 true，同时从桶中消费 n 个 token。
	//反之不消费桶中的Token，返回false
	if limiter.AllowN(time.Now(), 2) {
		fmt.Println("event allowed")
	} else {
		fmt.Println("event not allowed")
	}

	// func (lim *Limiter) Reserve() *Reservation
	// func (lim *Limiter) ReserveN(now time.Time, n int) *Reservation
	// 当调用完成后，无论 Token 是否充足，都会返回一个 *Reservation 对象。你可以调用该对象的Delay()方法，该方法返回的参数类型为time.Duration，
	//	反映了需要等待的时间，必须等到等待时间之后，才能进行接下来的工作。如果不想等待，可以调用Cancel()方法，该方法会将 Token 归还
	r := limiter.Reserve()
	//是否愿意等待
	if !r.OK() {
		// Not allowed to act! Did you remember to set lim.burst to be > 0 ?
		return
	}
	//如果愿意等待，将等待时间抛给用户 time.Sleep代表用户需要等待的时间。
	time.Sleep(r.Delay())
	Act() // 一段时间后生成生成新的令牌，开始执行相关逻辑

	// 3.动态调整速率
	//Limiter 支持可以调整速率和桶大小
	//SetLimit(Limit) 改变放入 Token 的速率
	//SetBurst(int) 改变 Token 桶大小

}

//执行业务逻辑
func Act() {

}

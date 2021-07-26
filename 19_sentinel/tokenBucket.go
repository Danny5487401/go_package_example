package main

import (
	"golang.org/x/time/rate"

	"context"
	"fmt"
	"time"
)

/*
背景：
	在高并发业务场景下，保护系统时，常用的"三板斧"有："熔断、降级和限流"
一。限流算法常用的几种实现方式有如下四种：
	1。计数器
	2。滑动窗口
	3。漏桶 Uber公司开源的限流器uber-go/ratelimit是漏桶算法实现的
	4。令牌桶
二。令牌桶：
	介绍：
		令牌桶是反向的"漏桶"，它是以恒定的速度往木桶里加入令牌，木桶满了则不再加入令牌。
		服务收到请求时尝试从木桶中取出一个令牌，如果能够得到令牌则继续执行后续的业务逻辑。如果没有得到令牌，
		直接返回访问频率超限的错误码或页面等，不继续执行后续的业务逻辑。
	特点：
		由于木桶内只要有令牌，请求就可以被处理，所以令牌桶算法可以支持突发流量。同时由于往木桶添加令牌的速度是恒定的，且木桶的容量有上限，
		所以单位时间内处理的请求书也能够得到控制，起到限流的目的。假设加入令牌的速度为 1token/10ms，桶的容量为500，
		在请求比较的少的时候（小于每10毫秒1个请求）时，木桶可以先"攒"一些令牌（最多500个）。当有突发流量时，
		一下把木桶内的令牌取空，也就是有500个在并发执行的业务逻辑，之后要等每10ms补充一个新的令牌才能接收一个新的请求
三。官方参考
	Golang 官方提供的扩展库里就自带了限流算法的实现，即 golang.org/x/time/rate。
	该限流器也是基于 Token Bucket(令牌桶) 实现的
四。源码分析
	time/rate包的Limiter类型对限流器进行了定义，所有限流功能都是通过基于Limiter类型实现的
	type Limiter struct {
		mu     sync.Mutex
		limit  Limit
		burst  int // 令牌桶的大小
		tokens float64
		last time.Time // 上次更新tokens的时间
		lastEvent time.Time // 上次发生限速器事件的时间（通过或者限制都是限速器事件）
	}
	字段解释
	limit：limit字段表示往桶里放Token的速率，它的类型是Limit，是int64的类型别名。
		设置limit时既可以用数字指定每秒向桶中放多少个Token，也可以指定向桶中放Token的时间间隔，
		其实指定了每秒放Token的个数后就能计算出放每个Token的时间间隔了。
	burst: 令牌桶的大小。
	tokens: 桶中的令牌。
	last: 上次往桶中放 Token 的时间。
	lastEvent：上次发生限速器事件的时间（通过或者限制都是限速器事件
五。适用场景
	适合电商抢购或者微博出现热点事件这种场景，因为在限流的同时可以应对一定的突发流量。如果采用漏桶那样的均匀速度处理请求的算法，
	在发生热点时间的时候，会造成大量的用户无法访问，对用户体验的损害比较大。
*/

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
	if !r.OK() {
		// Not allowed to act! Did you remember to set lim.burst to be > 0 ?
		return
	}
}

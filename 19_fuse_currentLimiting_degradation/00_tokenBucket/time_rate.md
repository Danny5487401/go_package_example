<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [官方扩展库 Token Bucket(令牌桶)限流算法 golang.org/x/time/rate](#%E5%AE%98%E6%96%B9%E6%89%A9%E5%B1%95%E5%BA%93-token-bucket%E4%BB%A4%E7%89%8C%E6%A1%B6%E9%99%90%E6%B5%81%E7%AE%97%E6%B3%95-golangorgxtimerate)
  - [适用场景](#%E9%80%82%E7%94%A8%E5%9C%BA%E6%99%AF)
  - [golang.org/x/time/rate源码分析](#golangorgxtimerate%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [常量结构](#%E5%B8%B8%E9%87%8F%E7%BB%93%E6%9E%84)
    - [方法](#%E6%96%B9%E6%B3%95)
    - [Reservation结构体](#reservation%E7%BB%93%E6%9E%84%E4%BD%93)
    - [初始化](#%E5%88%9D%E5%A7%8B%E5%8C%96)
    - [获取](#%E8%8E%B7%E5%8F%96)
    - [CancelAt](#cancelat)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 官方扩展库 Token Bucket(令牌桶)限流算法 golang.org/x/time/rate
该限流器也是基于 Token Bucket(令牌桶) 实现的

## 适用场景
适合电商抢购或者微博出现热点事件这种场景，因为在限流的同时可以应对一定的突发流量。如果采用漏桶那样的均匀速度处理请求的算法，
在发生热点时间的时候，会造成大量的用户无法访问，对用户体验的损害比较大。


## golang.org/x/time/rate源码分析
time/rate包的Limiter类型对限流器进行了定义，所有限流功能都是通过基于Limiter类型实现的
```go
type Limiter struct {
    mu     sync.Mutex 
    limit  Limit  
    burst  int //令牌桶的最大数量， 如果burst为0，则除非limit == Inf，否则不允许处理任何事件
    tokens float64  //令牌桶中可用的令牌数量
    last time.Time // 上次更新tokens的时间
    lastEvent time.Time //lastEvent记录速率受限制(桶中没有令牌)的时间点，该时间点可能是过去的，也可能是将来的(Reservation预定的结束时间点)

}
```

字段解释
- limit：limit字段表示往桶里放Token的速率，它的类型是Limit，是int64的类型别名。
设置limit时既可以用数字指定每秒向桶中放多少个Token，也可以指定向桶中放Token的时间间隔，
其实指定了每秒放Token的个数后就能计算出放每个Token的时间间隔了。
- burst: 令牌桶的大小。
- tokens: 桶中的令牌。
- last: 上次往桶中放 Token 的时间。
- lastEvent：上次发生限速器事件的时间（通过或者限制都是限速器事件)

Limiter是限流器中最核心的结构体，用于限流(控制事件发生的频率)，在初始化后默认是满的，并以每秒r个令牌的速率重新填充直到达到桶的容量(burst)，如果r == Inf表示无限制速率。

注意：由于令牌桶的令牌可以预约，所有令牌桶中的tokens可能为负数。



tokens更新的策略：

1. 成功获取到令牌或成功预约(Reserve)到令牌

2. 预约取消时(Cancel)并且需要还原令牌到令牌桶中时

3. 重新设置限流器的速率时(SetLimit)

4. 重新设置限流器的容量时(SetBurst)

### 常量结构
```go
//定义某个时间的最大频率
//表示每秒的事件数
type Limit float64
 
//Inf表示无速率限制
const Inf = Limit(math.MaxFloat64)
```

### 方法
- Allow： 如果没有令牌，则直接返回false

- Reserve：如果没有令牌，则返回一个reservation预约，

- Wait：如果没有令牌，则等待直到获取一个令牌或者其上下文被取消。

### Reservation结构体
```go
type Reservation struct {
	ok        bool //到截至时间是否可以获取足够的令牌
	lim       *Limiter  //用于指向具体的限流器
	tokens    int   //需要获取的令牌数量
	timeToAct time.Time  //预约的时间
	// This is the Limit at reservation time, it can change later.
	// 在预约的时候，限流器的产生速度，其实是可以通过变量lim来获取限流器产生令牌的速度，那么为什么还要单独的整出这个参数呢？是因为Limiter中的limit是可变的
	limit Limit
}
```
Reservation可以理解成预定令牌的操作，timeToAct是本次预约需要等待到的指定时间点才有足够预约的令牌。


### 初始化
```go
// 初始化Limiter，指定每秒允许处理事件的上限为r，允许令牌桶的最大容量为b
func NewLimiter(r Limit, b int) *Limiter {
	return &Limiter{
		limit: r,
		burst: b,
	}
}

```
```go
func Every(interval time.Duration) Limit {
	if interval <= 0 {
		return Inf
	}
	return 1 / Limit(interval.Seconds())
}
```
Every将事件的最小时间间隔转换为限制


### 获取
```go

func (lim *Limiter) Allow() bool {
	return lim.AllowN(time.Now(), 1)
}
```
从令牌桶中获取一个令牌，成功获取到则返回true
```go
func (lim *Limiter) AllowN(now time.Time, n int) bool {
	return lim.reserveN(now, n, 0).ok
}
```
从令牌桶中获取n个令牌，成功获取到则返回true

```go
func (lim *Limiter) Wait(ctx context.Context) (err error) {
	return lim.WaitN(ctx, 1)
}
```
获取一个令牌，如果没有则等待直到获取令牌或者上下文ctx取消
```go
func (lim *Limiter) WaitN(ctx context.Context, n int) (err error) {
	//同步获取令牌桶的最大容量burst和限流器的速率limit
	lim.mu.Lock()
	burst := lim.burst
	limit := lim.limit
	lim.mu.Unlock()
 
	//如果n大于令牌桶的最大容量，则返回error
	if n > burst && limit != Inf {
		return fmt.Errorf("rate: Wait(n=%d) exceeds limiter's burst %d", n, lim.burst)
	}
	// Check if ctx is already cancelled
	//判断上下文ctx是否已经被取消，如果已经取消则返回error
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
 
	// Determine wait limit
	now := time.Now()
	waitLimit := InfDuration
	if deadline, ok := ctx.Deadline(); ok { //如果可以获取上下文的截至时间，则更新可以等待的时间waitLimit
		waitLimit = deadline.Sub(now)
	}
 
	// Reserve
	//调用reserveN获取Reversation
	r := lim.reserveN(now, n, waitLimit)
	if !r.ok { //没有足够的时间获取令牌，则返回error
		return fmt.Errorf("rate: Wait(n=%d) would exceed context deadline", n)
	}
	// Wait if necessary
	//需要等待的时间
	delay := r.DelayFrom(now)
	if delay == 0 {
		return nil
	}
	t := time.NewTimer(delay)
	defer t.Stop()
	select {
	case <-t.C:
		// We can proceed.
		return nil
	case <-ctx.Done():
		// Context was canceled before we could proceed.  Cancel the
		// reservation, which may permit other events to proceed sooner.
		r.Cancel()
		return ctx.Err()
	}
}
```
WaitN方法获取n个令牌，直到成功获取或者ctx取消

- 如果n大于令牌桶的最大容量则返回error

- 如果上下文被取消或者等待的时间大于上下文的截至时间，则返回error

- 如果速率限制为Inf则不会限流

无论是Wait、Allow或者Reserve其实都会调用advance和reserveN方法，所以这两个方法是整个限流器rate实现的核心。

```go
func (lim *Limiter) reserveN(now time.Time, n int, maxFutureReserve time.Duration) Reservation {
	lim.mu.Lock()
 
	//如果没有限流则直接返回
	if lim.limit == Inf {
		lim.mu.Unlock()
		return Reservation{
			ok:        true,  //桶中有足够的令牌
			lim:       lim,
			tokens:    n,
			timeToAct: now,
		}
	}
 
	//更新令牌桶的状态，tokens为目前可用的令牌数量
	now, last, tokens := lim.advance(now)
 
	// Calculate the remaining number of tokens resulting from the request.
	//可用的令牌数tokens减去需要获取的令牌数(n)
	tokens -= float64(n)
 
	// Calculate the wait duration
	//如果tokens小于0，则说明桶中没有足够的令牌，计算出产生这些缺数的令牌需要多久(waitDuration)
	//计算出产生出缺数的令牌(即-tokens)需要多长时间
	var waitDuration time.Duration
	if tokens < 0 {
		waitDuration = lim.limit.durationFromTokens(-tokens)
	}
 
	// Decide result
	//如果n小于等于令牌桶的容量，并且可以等待到足够的令牌(即 waitDuration <= maxFutureReserve),则ok为true。表示可以获取到足够的令牌
	ok := n <= lim.burst && waitDuration <= maxFutureReserve
 
	// Prepare reservation
	r := Reservation{
		ok:    ok,
		lim:   lim,
		limit: lim.limit,
	}
	if ok {
		r.tokens = n // 需要的令牌数
		r.timeToAct = now.Add(waitDuration)  //计算获取到足够令牌的结束时间点
	}
 
	// Update state
	if ok {
		lim.last = now  //更新tokens的时间
		lim.tokens = tokens  //更新令牌桶目前可用的令牌数tokens
		lim.lastEvent = r.timeToAct  //下次事件时间(即获取到足够令牌的时刻)
	} else {
		lim.last = last
	}
 
	lim.mu.Unlock()
	return r
}
 
```

```go
func (lim *Limiter) advance(now time.Time) (newNow time.Time, newLast time.Time, newTokens float64) {
	//last不能在当前时间now之后，否则计算出来的elapsed为负数，会导致令牌桶数量减少
	last := lim.last
	if now.Before(last) {
		last = now
	}
 
	// Avoid making delta overflow below when last is very old.
	//根据令牌桶的缺数计算出令牌桶未进行更新的最大时间
	maxElapsed := lim.limit.durationFromTokens(float64(lim.burst) - lim.tokens)
	elapsed := now.Sub(last)  //令牌桶未进行更新的时间段
	if elapsed > maxElapsed {
		elapsed = maxElapsed
	}
 
	// Calculate the new number of tokens, due to time that passed.
	//根据未更新的时间(未向桶中加入令牌的时间段)计算出产生的令牌数。
	delta := lim.limit.tokensFromDuration(elapsed)
	tokens := lim.tokens + delta  //计算出可用的令牌数
	if burst := float64(lim.burst); tokens > burst {
		tokens = burst
	}
 
	return now, last, tokens
}
```
advance方法的作用是更新令牌桶的状态，计算出令牌桶未更新的时间(elapsed)，根据elapsed算出需要向桶中加入的令牌数delta，然后算出桶中可用的令牌数newTokens

### CancelAt
有用户预约到了一个时间，但是他又取消了
```go
func (r *Reservation) Cancel() {
	r.CancelAt(time.Now())
	return
}
 
 
func (r *Reservation) CancelAt(now time.Time) {
	if !r.ok {
		return
	}
 
	r.lim.mu.Lock()
	defer r.lim.mu.Unlock()
 
	/*
	1.如果无需限流
	2. tokens为0 (需要获取的令牌数量为0)
	3. 已经过了截至时间
	以上三种情况无需处理取消操作
	*/
	if r.lim.limit == Inf || r.tokens == 0 || r.timeToAct.Before(now) {
		return
	}
 
	// calculate tokens to restore
	// The duration between lim.lastEvent and r.timeToAct tells us how many tokens were reserved
	// after r was obtained. These tokens should not be restored.
	//计算出出需要还原的令牌数量
	//这里的r.lim.lastEvent可能是本次Reservation的结束时间，也可能是后来的Reservation的结束时间，所以要把本次结束时间点(r.timeToAct)之后产生的令牌数减去
	restoreTokens := float64(r.tokens) - r.limit.tokensFromDuration(r.lim.lastEvent.Sub(r.timeToAct))
	if restoreTokens <= 0 {
		return
	}
 
 
	// advance time to now
	//从新计算令牌桶的状态
	now, _, tokens := r.lim.advance(now)
	// calculate new number of tokens
	//还原当前令牌桶的令牌数量，当前的令牌数tokens加上需要还原的令牌数restoreTokens
	tokens += restoreTokens
	//如果tokens大于桶的最大容量，则将tokens置为桶的最大容量
	if burst := float64(r.lim.burst); tokens > burst {
		tokens = burst
	}
	// update state
	r.lim.last = now  //记录桶的更新时间
	r.lim.tokens = tokens //更新令牌数量
	//还原lastEvent，即上次速率受限制的时间
	if r.timeToAct == r.lim.lastEvent {
		prevEvent := r.timeToAct.Add(r.limit.durationFromTokens(float64(-r.tokens)))
		if !prevEvent.Before(now) {
			r.lim.lastEvent = prevEvent
		}
	}
 
	return
}
```

## 参考

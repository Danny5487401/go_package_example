# 官方扩展库限流算法golang.org/x/time/rate

该限流器也是基于 Token Bucket(令牌桶) 实现的
## 源码分析
time/rate包的Limiter类型对限流器进行了定义，所有限流功能都是通过基于Limiter类型实现的
```go
type Limiter struct {
    mu     sync.Mutex
    limit  Limit
    burst  int // 令牌桶的大小
    tokens float64
    last time.Time // 上次更新tokens的时间
    lastEvent time.Time // 上次发生限速器事件的时间（通过或者限制都是限速器事件）
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
- 
## 适用场景
适合电商抢购或者微博出现热点事件这种场景，因为在限流的同时可以应对一定的突发流量。如果采用漏桶那样的均匀速度处理请求的算法，
在发生热点时间的时候，会造成大量的用户无法访问，对用户体验的损害比较大。

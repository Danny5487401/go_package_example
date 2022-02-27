# go-cache
基于内存的 K/V 存储/缓存 : (类似于Memcached)，适用于单机应用程序 ，支持删除，过期，默认Cache共享锁，

大量key的情况下会造成锁竞争严重

## 为什么选择go-cache？
可以存储任何对象（在给定的持续时间内或永久存储），并且可以由多个goroutine安全地使用缓存。

## 源码分析
主要针对核心的存储结构、Set、Get、Delete、定时清理逻辑进行分析

### 核心的存储结构
```go
// Item 每一个具体缓存值
type Item struct {
    Object     interface{}
    Expiration int64 // 过期时间:设置时间+缓存时长
}

// Cache 整体缓存
type Cache struct {
    *cache
}

// cache 整体缓存
type cache struct {
    defaultExpiration time.Duration // 默认超时时间
    items             map[string]Item // KV对
    mu                sync.RWMutex // 读写锁，在操作（增加，删除）缓存时使用
    onEvicted         func(string, interface{}) // 删除KEY时的CallBack函数
    janitor           *janitor // 定时清空缓存的结构
}

// janitor  定时清空缓存的结构
type janitor struct {
    Interval time.Duration // 多长时间扫描一次缓存
    stop     chan bool // 是否需要停止
}

```

### Set
```go
func (c *cache) Set(k string, x interface{}, d time.Duration) {
    // "Inlining" of set
    var e int64
    if d == DefaultExpiration {
        d = c.defaultExpiration
    }
    if d > 0 {
        e = time.Now().Add(d).UnixNano()
        }
        c.mu.Lock() // 这里可以使用defer？
        c.items[k] = Item{
        Object:     x, // 实际的数据
        Expiration: e, // 下次过期时间
    }
    c.mu.Unlock()
}
```

### get
```go
func (c *cache) Get(k string) (interface{}, bool) {
    c.mu.RLock() // 加锁，限制并发读写
    item, found := c.items[k] // 在 items 这个 map[string]Item 查找数据
    if !found {
        c.mu.RUnlock()
        return nil, false
    }
    if item.Expiration > 0 {
        if time.Now().UnixNano() > item.Expiration { // 已经过期，直接返回nil，为什么在这里不直接就删除了呢？
            c.mu.RUnlock()
            return nil, false
        }
    }
    c.mu.RUnlock()
    return item.Object, true
}
```

### Delete
```go
// Delete an item from the cache. Does nothing if the key is not in the cache.
func (c *cache) Delete(k string) {
    c.mu.Lock()
    v, evicted := c.delete(k)
    c.mu.Unlock()
    if evicted {
        c.onEvicted(k, v) // 删除KEY时的CallBack
    }
}

func (c *cache) delete(k string) (interface{}, bool) {
    if c.onEvicted != nil {
        if v, found := c.items[k]; found {
            delete(c.items, k)
            return v.Object, true
        }
    }
    delete(c.items, k)
    return nil, false
}
```

### 定时清理逻辑

```go
func newCacheWithJanitor(de time.Duration, ci time.Duration, m map[string]Item) *Cache {
    c := newCache(de, m)
    C := &Cache{c}
    if ci > 0 {
        runJanitor(c, ci) // 定时运行清除过期KEY
        runtime.SetFinalizer(C, stopJanitor) // 当C被GC回收时，会停止runJanitor 中的协程
    }
    return C
}

func runJanitor(c *cache, ci time.Duration) {
    j := &janitor{
        Interval: ci,
        stop:     make(chan bool),
    }
    c.janitor = j
    go j.Run(c) // 新的协程做过期删除逻辑
}

func (j *janitor) Run(c *cache) {
    ticker := time.NewTicker(j.Interval)
    for {
        select {
        case <-ticker.C: // 每到一个周期就全部遍历一次
            c.DeleteExpired() // 实际的删除逻辑
        case <-j.stop:
            ticker.Stop()
        return
        }
    }
}

// Delete all expired items from the cache.
func (c *cache) DeleteExpired() {
    var evictedItems []keyAndValue
    now := time.Now().UnixNano()
    c.mu.Lock()
    for k, v := range c.items { // 加锁遍历整个列表
        // "Inlining" of expired
        if v.Expiration > 0 && now > v.Expiration {
            ov, evicted := c.delete(k)
            if evicted {
                evictedItems = append(evictedItems, keyAndValue{k, ov})
            }
        }
    }
    c.mu.Unlock()
    for _, v := range evictedItems {
        c.onEvicted(v.key, v.value)
    }
}
```

## 性能分析
### Lock 的使用
在go-cache中，涉及到读写cache，基本上都用到了锁，而且在遍历的时候也用到锁，当cache的数量非常多时，读写频繁时， 会有严重的锁冲突。
### 使用读写锁？
sync.RWMutex, 在读的时候加RLock, 可以允许多个读。在写的时候加Lock，不允许其他读和写。
### 锁的粒度是否可以变更小？
根据KEY HASH 到不同的map中
### 使用sync.map?
减少锁的使用
### runtime.SetFinalizer
在实际的编程中，我们都希望每个对象释放时执行一个方法，在该方法内执行一些计数、释放或特定的要求， 以往都是在对象指针置nil前调用一个特定的方法， golang提供了runtime.SetFinalizer函数，当GC准备释放对象时，会回调该函数指定的方法，非常方便和有效。

对象可以关联一个SetFinalizer函数， 当gc检测到unreachable对象有关联的SetFinalizer函数时， 会执行关联的SetFinalizer函数， 同时取消关联。这样当下一次gc的时候， 对象重新处于unreachable状态 并且没有SetFinalizer关联， 就会被回收。
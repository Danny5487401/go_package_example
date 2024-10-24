<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [golang-lru](#golang-lru)
  - [源码](#%E6%BA%90%E7%A0%81)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# golang-lru 



golang-lru 一共提供了 3 种 LRU cache 类型：

- Cache 是最简单的 LRU cache，它的实现基于 groupcache
- TwoQueueCache 使用 2Q 缓存淘汰算法，从访问时间和访问频率两个维度来跟踪缓存条目
- ARCCache 使用 ARC 缓存淘汰算法，它也是从访问时间和访问频率两个维度来跟踪缓存条目，同时它也会跟踪缓存条目的淘汰情况，从而动态调整 LRU 队列和 FRU 队列的 size 比例



## 源码

缓存配置
```go
// Cache is a thread-safe fixed size LRU cache.
type Cache[K comparable, V any] struct {
	lru         *simplelru.LRU[K, V]
	evictedKeys []K
	evictedVals []V
	onEvictedCB func(k K, v V)
	lock        sync.RWMutex // Cache 提供线程安全能力
}

```

lRU 
```go
// github.com/hashicorp/golang-lru/v2@v2.0.7/simplelru/lru.go

// LRU implements a non-thread safe fixed size LRU cache
type LRU[K comparable, V any] struct {
	size      int // LRU cache 的大小
	evictList *internal.LruList[K, V] // 缓存条目链表，通过该链表来实现 LRU 算法
	items     map[K]*internal.Entry[K, V] // 缓存条目 map，缓存条目的查找通过该 map 实现。map 的 key 是缓存条目的 key，而 value 则是 evictList 链表元素的指针。链表元素中包含缓存条目的 key 和 value
	onEvict   EvictCallback[K, V] // 缓存条目淘汰时的回调函数
}

```


初始化
```go
// New creates an LRU of the given size.
func New[K comparable, V any](size int) (*Cache[K, V], error) {
	return NewWithEvict[K, V](size, nil)
}

// NewWithEvict constructs a fixed size cache with the given eviction
// callback.
func NewWithEvict[K comparable, V any](size int, onEvicted func(key K, value V)) (c *Cache[K, V], err error) {
	// create a cache with default settings
	c = &Cache[K, V]{
		onEvictedCB: onEvicted,
	}
	if onEvicted != nil {
		c.initEvictBuffers()
		onEvicted = c.onEvicted
	}
	c.lru, err = simplelru.NewLRU(size, onEvicted)
	return
}


func NewLRU[K comparable, V any](size int, onEvict EvictCallback[K, V]) (*LRU[K, V], error) {
	// 缓存大小判断
	if size <= 0 {
		return nil, errors.New("must provide a positive size")
	}

	c := &LRU[K, V]{
		size:      size,
		evictList: internal.NewList[K, V](),
		items:     make(map[K]*internal.Entry[K, V]),
		onEvict:   onEvict,
	}
	return c, nil
}
```


获取缓存
```go
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.lock.Lock()
	value, ok = c.lru.Get(key)
	c.lock.Unlock()
	return value, ok
}


func (c *LRU[K, V]) Get(key K) (value V, ok bool) {
	// 从 map 中根据 key 查找
	if ent, ok := c.items[key]; ok {
		// 如果找到该缓存，需要将其移动到链表头
		c.evictList.MoveToFront(ent)
		// 返回链表元素中保存的缓存 key 对应的 value
		return ent.Value, true
	}
	return
}
```


添加元素
```go
func (c *Cache[K, V]) Add(key K, value V) (evicted bool) {
	var k K
	var v V
	// 操作 lru cache 时需要加锁，保证线程安全
	c.lock.Lock()
	evicted = c.lru.Add(key, value)
	if c.onEvictedCB != nil && evicted {// 如果因为 Add 操作而造成其他缓存条目被淘汰，同时用户指定了回调函数
		k, v = c.evictedKeys[0], c.evictedVals[0]
		c.evictedKeys, c.evictedVals = c.evictedKeys[:0], c.evictedVals[:0]
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil && evicted {
        // 在非临界区调用用户的回调函数
		c.onEvictedCB(k, v)
	}
	return
}

func (c *LRU[K, V]) Add(key K, value V) (evicted bool) {
	// Check for existing item
	if ent, ok := c.items[key]; ok {
		// 发现存在移到前面并覆盖
		c.evictList.MoveToFront(ent)
		ent.Value = value
		return false
	}

	// Add new item
	ent := c.evictList.PushFront(key, value)
	c.items[key] = ent

	// 判断是否需要驱逐
	evict := c.evictList.Length() > c.size
	// Verify size not exceeded
	if evict {
		c.removeOldest()
	}
	return evict
}
```
清除元素
```go
func (c *LRU[K, V]) RemoveOldest() (key K, value V, ok bool) {
	if ent := c.evictList.Back(); ent != nil {
		c.removeElement(ent)
		return ent.Key, ent.Value, true
	}
	return
}

func (c *LRU[K, V]) removeElement(e *internal.Entry[K, V]) {
	c.evictList.Remove(e)
	delete(c.items, e.Key)
	if c.onEvict != nil {
		c.onEvict(e.Key, e.Value)
	}
}

```


## 参考

- [golang-lru 源码分析](https://fuchencong.com/2022/01/26/go-develop-notes-01/)
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [本地缓存](#%E6%9C%AC%E5%9C%B0%E7%BC%93%E5%AD%98)
  - [需求](#%E9%9C%80%E6%B1%82)
  - [Go本身实现背景](#go%E6%9C%AC%E8%BA%AB%E5%AE%9E%E7%8E%B0%E8%83%8C%E6%99%AF)
  - [本地缓存组件优化方式](#%E6%9C%AC%E5%9C%B0%E7%BC%93%E5%AD%98%E7%BB%84%E4%BB%B6%E4%BC%98%E5%8C%96%E6%96%B9%E5%BC%8F)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 本地缓存
参考链接：https://mp.weixin.qq.com/s/UuqkO9UUXjNWNGyI2sjG4Q

## 需求
![](.cache_images/cache_need.png)
1. 我们一般做缓存就是为了能提高系统的读写性能，缓存的命中率越高，也就意味着缓存的效果越好。
2. 其次本地缓存一般都受限于本地内存的大小，所以全量的数据一般存不下。
   - 那基于这样的场景，一方面是想缓存的数据越多，则命中率理论上也会随着缓存数据的增多而提高；
   - 另外一方面是想，既然所有的数据存不下那就想办法利用有限的内存存储有限的数据。这些有限的数据需要是经常访问的，同时有一定时效性（不会频繁改变）的。
     
基于这两个点展开，我们一般对本地缓存会要求其满足支持过期时间、支持淘汰策略。

3. 最后再使用自动管理内存的语言，例如 Go 等开发时，还需要考虑在加入本地缓存后引发的 GC 问题。

## Go本身实现背景
Go 中内置的可以直接用来做本地缓存的无非就是 map 和 sync.Map。
而这两者中，map 是非并发安全的数据结构，在使用时需要加锁；而 sync.Map 虽然是线程安全的。但是需要在并发读写时加锁。
此外二者均无法支持数据的过期和淘汰，同时在存储大量数据时，又会产生比较频繁的 GC 问题，更严重的情况下导致线上服务无法稳定运行


## 本地缓存组件优化方式
1. 实现零 GC 的方案主要就两种：
   a. 无 GC：分配堆外内存（Mmap）
   b. 避免 GC：map 非指针优化（map[uint64]uint32）或者采用 slice 实现一套无指针的 map
   c. 避免 GC：数据存入[]byte slice（可考虑底层采用环形队列封装循环使用空间）

2. 实现高性能的关键在于：
   a. 数据分片（降低锁的粒度)

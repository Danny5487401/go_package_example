<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Consensus algorithm 共识算法](#consensus-algorithm-%E5%85%B1%E8%AF%86%E7%AE%97%E6%B3%95)
  - [基本概念](#%E5%9F%BA%E6%9C%AC%E6%A6%82%E5%BF%B5)
    - [PoW（Proof-of-Work，工作量证明）](#powproof-of-work%E5%B7%A5%E4%BD%9C%E9%87%8F%E8%AF%81%E6%98%8E)
    - [PoS（Proof-of-Stake，权益证明）](#posproof-of-stake%E6%9D%83%E7%9B%8A%E8%AF%81%E6%98%8E)
    - [DPoS（Delegated Proof of Stake，委托权益证明）](#dposdelegated-proof-of-stake%E5%A7%94%E6%89%98%E6%9D%83%E7%9B%8A%E8%AF%81%E6%98%8E)
  - [背景：如何避免单点故障](#%E8%83%8C%E6%99%AF%E5%A6%82%E4%BD%95%E9%81%BF%E5%85%8D%E5%8D%95%E7%82%B9%E6%95%85%E9%9A%9C)
  - [多副本常用的技术方案](#%E5%A4%9A%E5%89%AF%E6%9C%AC%E5%B8%B8%E7%94%A8%E7%9A%84%E6%8A%80%E6%9C%AF%E6%96%B9%E6%A1%88)
    - [主从复制，又分为全同步复制、异步复制、半同步复制，比如 MySQL/Redis 单机主备版就基于主从复制实现的。](#%E4%B8%BB%E4%BB%8E%E5%A4%8D%E5%88%B6%E5%8F%88%E5%88%86%E4%B8%BA%E5%85%A8%E5%90%8C%E6%AD%A5%E5%A4%8D%E5%88%B6%E5%BC%82%E6%AD%A5%E5%A4%8D%E5%88%B6%E5%8D%8A%E5%90%8C%E6%AD%A5%E5%A4%8D%E5%88%B6%E6%AF%94%E5%A6%82-mysqlredis-%E5%8D%95%E6%9C%BA%E4%B8%BB%E5%A4%87%E7%89%88%E5%B0%B1%E5%9F%BA%E4%BA%8E%E4%B8%BB%E4%BB%8E%E5%A4%8D%E5%88%B6%E5%AE%9E%E7%8E%B0%E7%9A%84)
    - [中心化复制](#%E4%B8%AD%E5%BF%83%E5%8C%96%E5%A4%8D%E5%88%B6)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Consensus algorithm 共识算法
共识算法推荐资料：Leslie Lamport和Diego Ongaro的数篇论文、Ongaro在youtube上发的三个视频讲解，以及何登成的ppt。

分布式新系统的共识（consensus）方案需要保证一些基本的的特性：

- safety：多个副本之间对某个决议能够达成一个共识，且不能被篡改，
- liveness：多个副本之间对某个决议能够在有限的时间内形成，


## 基本概念

### PoW（Proof-of-Work，工作量证明）


PoW 算法获取记账权的原理是：利用区块的 index、前一个区块的哈希值、交易的时间戳、区块数据和 nonce 值，通过 SHA256 哈希算法计算出一个哈希值，并判断前 k 个值是否都为 0。如果不是，则递增 nonce 值，重新按照上述方法计算；
如果是，则本次计算的哈希值为要解决的题目的正确答案。谁最先计算出正确答案，谁就获得这个区块的记账权。


### PoS（Proof-of-Stake，权益证明）

PoS 算法中持币越多或持币越久，币龄就会越高，持币人就越容易挖到区块并得到激励，而持币少的人基本没有机会，这样整个系统的安全性实际上会被持币数量较大的一部分人掌握，容易出现垄断现象。

### DPoS（Delegated Proof of Stake，委托权益证明）


DPoS 是在 PoW 和 PoS 的基础上进行改进的，相比于 PoS 算法，DPoS 引入了受托人，优点主要表现在：
- 由投票选举出的若干信誉度更高的受托人记账，解决了所有节点均参与竞争导致消息量大、达成一致的周期长的问题。也就是说，DPoS 能耗更低，具有更快的交易速度。
- 每隔一定周期会调整受托人，避免受托人造假和独权

## 背景：如何避免单点故障

为了解决单点问题，软件系统工程师引入了数据复制技术，实现多副本。通过数据复制方案，一方面我们可以提高服务可用性，避免单点故障。
另一方面，多副本可以提升读吞吐量、甚至就近部署在业务所在的地理位置，降低访问延迟。

## 多副本常用的技术方案

### 主从复制，又分为全同步复制、异步复制、半同步复制，比如 MySQL/Redis 单机主备版就基于主从复制实现的。

1. 全同步复制是指主收到一个写请求后，必须等待全部从节点确认返回后，才能返回给客户端成功。因此如果一个从节点故障，整个系统就会不可用。
   这种方案为了保证多副本的一致性，而牺牲了可用性，一般使用不多。

2. 异步复制是指主收到一个写请求后，可及时返回给 client，异步将请求转发给各个副本，若还未将请求转发到副本前就故障了，则可能导致数据丢失，但是可用性是最高的。

3. 半同步复制介于全同步复制、异步复制之间，它是指主收到一个写请求后，至少有一个副本接收数据后，就可以返回给客户端成功，在数据一致性、可用性上实现了平衡和取舍

### 中心化复制

它是指在一个 n 副本节点集群中，任意节点都可接受写请求，但一个成功的写入需要 w 个节点确认，读取也必须查询至少 r 个节点。
根据实际业务场景对数据一致性的敏感度，设置合适 w/r 参数。比如你希望每次写入后，任意 client 都能读取到新值，如果 n 是 3 个副本，
你可以将 w 和 r 设置为 2，这样当你读两个节点时候，必有一个节点含有最近写入的新值，这种读我们称之为法定票数读（quorum read）。

AWS 的 Dynamo 系统就是基于去中心化的复制算法实现的。它的优点是节点角色都是平等的，降低运维复杂度，可用性更高。
但是缺陷是去中心化复制，势必会导致各种写入冲突，业务需要关注冲突处理。

如何解决以上复制算法的困境呢？

共识算法

## 参考

- [一文读懂11个主流共识算法, 彻底搞懂PoS,PoW,dPoW,PBFT,dBFT](https://cloud.tencent.com/developer/article/1375464)
- [分布式技术原理与算法解析-05|分布式共识：存异求同](https://time.geekbang.org/column/article/144548)
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Consensus algorithm](#consensus-algorithm)
  - [基本概念](#%E5%9F%BA%E6%9C%AC%E6%A6%82%E5%BF%B5)
    - [PoW（Proof-of-Work，工作量证明）](#powproof-of-work%E5%B7%A5%E4%BD%9C%E9%87%8F%E8%AF%81%E6%98%8E)
    - [PoS（Proof-of-Stake，权益证明）](#posproof-of-stake%E6%9D%83%E7%9B%8A%E8%AF%81%E6%98%8E)
    - [DPoS（Delegated Proof of Stake，委托权益证明）](#dposdelegated-proof-of-stake%E5%A7%94%E6%89%98%E6%9D%83%E7%9B%8A%E8%AF%81%E6%98%8E)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Consensus algorithm


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

## 参考

- [一文读懂11个主流共识算法, 彻底搞懂PoS,PoW,dPoW,PBFT,dBFT](https://cloud.tencent.com/developer/article/1375464)
- [分布式技术原理与算法解析-05|分布式共识：存异求同](https://time.geekbang.org/column/article/144548)
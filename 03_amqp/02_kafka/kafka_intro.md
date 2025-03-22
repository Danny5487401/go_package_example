<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Kafka](#kafka)
  - [一. 基本概念](#%E4%B8%80-%E5%9F%BA%E6%9C%AC%E6%A6%82%E5%BF%B5)
    - [AR, ISR, OSR](#ar-isr-osr)
    - [HW(High Watermark 高水位)](#hwhigh-watermark-%E9%AB%98%E6%B0%B4%E4%BD%8D)
    - [LEO(Log End Offset 日志末端位移)](#leolog-end-offset-%E6%97%A5%E5%BF%97%E6%9C%AB%E7%AB%AF%E4%BD%8D%E7%A7%BB)
    - [ISR，HW, LEO之间的关系](#isrhw-leo%E4%B9%8B%E9%97%B4%E7%9A%84%E5%85%B3%E7%B3%BB)
  - [二. producer发布消息](#%E4%BA%8C-producer%E5%8F%91%E5%B8%83%E6%B6%88%E6%81%AF)
    - [producer delivery guarantee生产者发送保证](#producer-delivery-guarantee%E7%94%9F%E4%BA%A7%E8%80%85%E5%8F%91%E9%80%81%E4%BF%9D%E8%AF%81)
  - [三. broker 保存消息](#%E4%B8%89-broker-%E4%BF%9D%E5%AD%98%E6%B6%88%E6%81%AF)
  - [四. 消费者](#%E5%9B%9B-%E6%B6%88%E8%B4%B9%E8%80%85)
    - [两套 consumer API](#%E4%B8%A4%E5%A5%97-consumer-api)
    - [消费组: 消费者是以消费者组的形式对外消费的。](#%E6%B6%88%E8%B4%B9%E7%BB%84-%E6%B6%88%E8%B4%B9%E8%80%85%E6%98%AF%E4%BB%A5%E6%B6%88%E8%B4%B9%E8%80%85%E7%BB%84%E7%9A%84%E5%BD%A2%E5%BC%8F%E5%AF%B9%E5%A4%96%E6%B6%88%E8%B4%B9%E7%9A%84)
    - [消费者重平衡 consumer rebalance](#%E6%B6%88%E8%B4%B9%E8%80%85%E9%87%8D%E5%B9%B3%E8%A1%A1-consumer-rebalance)
  - [五. kafka 高可用 HA](#%E4%BA%94-kafka-%E9%AB%98%E5%8F%AF%E7%94%A8-ha)
    - [副本 replication](#%E5%89%AF%E6%9C%AC-replication)
    - [宕机的场景](#%E5%AE%95%E6%9C%BA%E7%9A%84%E5%9C%BA%E6%99%AF)
  - [六. 开发注意事项](#%E5%85%AD-%E5%BC%80%E5%8F%91%E6%B3%A8%E6%84%8F%E4%BA%8B%E9%A1%B9)
  - [七. 位移主题](#%E4%B8%83-%E4%BD%8D%E7%A7%BB%E4%B8%BB%E9%A2%98)
    - [位移主题的消息格式](#%E4%BD%8D%E7%A7%BB%E4%B8%BB%E9%A2%98%E7%9A%84%E6%B6%88%E6%81%AF%E6%A0%BC%E5%BC%8F)
    - [消费者在消费了信息之后，该把位移发送到哪呢？](#%E6%B6%88%E8%B4%B9%E8%80%85%E5%9C%A8%E6%B6%88%E8%B4%B9%E4%BA%86%E4%BF%A1%E6%81%AF%E4%B9%8B%E5%90%8E%E8%AF%A5%E6%8A%8A%E4%BD%8D%E7%A7%BB%E5%8F%91%E9%80%81%E5%88%B0%E5%93%AA%E5%91%A2)
    - [位移的提交](#%E4%BD%8D%E7%A7%BB%E7%9A%84%E6%8F%90%E4%BA%A4)
  - [八. 安全认证](#%E5%85%AB-%E5%AE%89%E5%85%A8%E8%AE%A4%E8%AF%81)
  - [九. 消息重复和消费幂等](#%E4%B9%9D-%E6%B6%88%E6%81%AF%E9%87%8D%E5%A4%8D%E5%92%8C%E6%B6%88%E8%B4%B9%E5%B9%82%E7%AD%89)
    - [消费阻塞以及堆积](#%E6%B6%88%E8%B4%B9%E9%98%BB%E5%A1%9E%E4%BB%A5%E5%8F%8A%E5%A0%86%E7%A7%AF)
    - [分区个数](#%E5%88%86%E5%8C%BA%E4%B8%AA%E6%95%B0)
  - [十. 面试问题](#%E5%8D%81-%E9%9D%A2%E8%AF%95%E9%97%AE%E9%A2%98)
    - [为什么Kafka不支持读写分离？](#%E4%B8%BA%E4%BB%80%E4%B9%88kafka%E4%B8%8D%E6%94%AF%E6%8C%81%E8%AF%BB%E5%86%99%E5%88%86%E7%A6%BB)
    - [Kafka的那些设计让它有如此高的性能](#kafka%E7%9A%84%E9%82%A3%E4%BA%9B%E8%AE%BE%E8%AE%A1%E8%AE%A9%E5%AE%83%E6%9C%89%E5%A6%82%E6%AD%A4%E9%AB%98%E7%9A%84%E6%80%A7%E8%83%BD)
    - [kafka中的事务](#kafka%E4%B8%AD%E7%9A%84%E4%BA%8B%E5%8A%A1)
  - [十一 kafka主流 golang 客户端](#%E5%8D%81%E4%B8%80-kafka%E4%B8%BB%E6%B5%81-golang-%E5%AE%A2%E6%88%B7%E7%AB%AF)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Kafka
Apache Kafka是消息引擎系统，也是一个分布式流处理平台（Distributed Streaming Platform）.

再强调一遍，Kafka是消息引擎系统，也是分布式流处理平台。

## 一. 基本概念
3.0 之前版本   
![](.kafka_intro_images/kafka_structure.png)
服务端：Broker相当于Kafka的服务端，你可以理解为是队列存在的地方，生产者把消息发送到Broker中，消费者从Broker中获取消息
客户端:生产者与消费者

![](.kafka_intro_images/kafka_3.0_structure.png)

kafka 3 的版本:当中已经彻底去掉了对zk的依赖,KIP-500议案提出了在Kafka中处理元数据的更好方法。基本思想是"Kafka on Kafka"，将Kafka的元数据存储在Kafka本身中，无需增加额外的外部存储比如ZooKeeper等。
在kafka3.0的新的版本当中，使用了新的KRaft协议，使用该协议来保证在元数据仲裁中准确的复制元数据，这个协议类似于zk当中的zab协议以及类似于Raft协议，但是KRaft协议使用的是基于事件驱动的模式，与ZAB协议和Raft协议还有点不一样.

![](.kafka_intro_images/kafka_definition.png)
1. producer： 消息生产者，发布消息到 kafka 集群的终端或服务。
2. broker：kafka 集群中包含的服务器。
3. topic： 每条发布到 kafka 集群的消息属于的类别，即 kafka 是面向 topic 的。
4. partition分区：partition 是物理上的概念，每个 topic 包含一个或多个 partition。kafka 分配的单位是 partition。Broker所在的机器的IO速度，可能会使得这个消息队列存在性能上的瓶颈。
    这里的分区，是扩展的概念，而不是副本的概念。所以，这里的分区，也被称为数据分区、数据分片。这些分区可以部署在不同的机器上，性能也就提升了好几倍。
5. consumer： 从 kafka 集群中消费消息的终端或服务。
6. Consumer group： high-level consumer API 中，每个 consumer 都属于一个 consumer group，每条消息只能被 consumer group 中的一个 Consumer 消费，但可以被多个 consumer group 消费。
7. replica：partition 的副本，保障 partition 的高可用。
8. leader： replica 中的一个角色， producer 和 consumer 只跟 leader 交互。
9. follower： replica 中的一个角色，从 leader 中复制数据。
10. controller： kafka 集群中的其中一个服务器，用来进行 leader election 以及 各种 failover。
11. zookeeper： kafka 通过 zookeeper 来存储集群的 meta 信息，kafka3.0去除了。
12. leader epoch:  leader 的纪元信息（epoch），初始值为0。每当 leader 变更一次，leader epoch 的值就会加1，相当于为 leader 增设了一个版本号。

### AR, ISR, OSR
追随者副本是不对外提供服务,方便实现“Read-your-writes”,当你使用生产者API向Kafka成功写入消息后，马上使用消费者API去读取刚才生产的消息;
方便实现单调读（Monotonic Reads）,在多次消费消息时，它不会看到某条消息一会儿存在一会儿不存在。
1. AR（Assigned Replicas): 分区中的所有副本统称。
2. ISR（In-Sync Replicas ISR副本): 所有与 leader 副本保持一定程度同步的副本（包括 leader 副本在内），ISR 集合是 AR 集合中的一个子集。
3. OSR（Out-of-Sync Replicas 非同步副本）：与leader副本同步滞后过多的副本（不包括leader副本）
4. leader 副本: 负责维护和跟踪 ISR 集合中所有 follower 副本的滞后状态，当 follower 副本落后太多或失效时，leader 副本会把它从 ISR 集合中剔除。

如果 OSR 集合中有 follower 副本“追上”了 leader 副本，那么 leader 副本会把它从 OSR 集合转移至 ISR 集合 。
默认情况下，当 leader 副本发生故障时，只有在 ISR 集合中的副本才有资格被选举为新的 leader，而在 OSR 集合中的副本则没有任何机会（不过这个原则也可以通过修改相应的参数配置来改变）。

![](.kafka_intro_images/hw.png)
### HW(High Watermark 高水位)
高水位的作用主要有2个。

1. 定义消息可见性，即用来标识分区下的哪些消息是可以被消费者消费的。
2. 帮助Kafka完成副本同步

它标识了一个特定的消息偏移量（offset），消费者只能拉取到这个offset之前的消息。

它代表一个日志文件，这个日志文件中有 9 条消息，第一条消息的 offset（LogStartOffset）为0，最后一条消息的offset为8，offset为9的消息用虚线框表示，代表下一条待写入的消息。
日志文件的HW为6，表示消费者只能拉取到offset在0至5之间的消息，而offset为6的消息对消费者而言是不可见的。

### LEO(Log End Offset 日志末端位移)
它标识当前日志文件中下一条待写入消息的offset，图中offset为9的位置即为当前日志文件的LEO，LEO的大小相当于当前日志分区中最后一条消息的offset值加1。
分区ISR集合中的每个副本都会维护自身的LEO，而ISR集合中最小的LEO即为分区的HW，对消费者而言只能消费HW之前的消息。

### ISR，HW, LEO之间的关系
假设某个分区的ISR集合中有3个副本，即一个leader副本和2个follower副本，此时分区的LEO和HW都为3。消息3和消息4从生产者发出之后会被先存入leader副本

![](.kafka_intro_images/isr_process1.png)
![](.kafka_intro_images/isr_process2.png)
![](.kafka_intro_images/isr_process3.png)
![](.kafka_intro_images/isr_process4.png)

Broker端参数replica.lag.time.max.ms参数值:当前默认值是10秒。这就是说，只要一个Follower副本落后Leader副本的时间不连续超过10秒，那么Kafka就认为该Follower副本与Leader是同步的，即使此时Follower副本中保存的消息明显少于Leader副本中的消息。

## 二. producer发布消息

1. 写入方式: 
    producer 采用 push 模式将消息发布到 broker，每条消息都被 append 到 partition 中，属于顺序写磁盘（顺序写磁盘效率比随机写内存要高，保障 kafka 吞吐率）
2. 消息路由:选区原则
   - a. 指定了 partition，则直接使用；
   - b. 未指定 partition 但指定 key，通过对 key 的 value 进行hash 选出一个 partition
   - c. partition 和 key 都未指定，使用轮询选出一个 partition。
3. 写入流程
   - a. producer 先从 zookeeper 的 "/brokers/.../state" 节点找到该 partition 的 leader
   - b. producer 将消息发送给该 leader
   - c. leader 将消息写入本地 log
   - d. followers 从 leader pull 消息，写入本地 log 后 leader 发送 ACK
   - e. leader 收到所有 ISR 中的 replica 的 ACK 后，增加 HW（high watermark，最后 commit 的 offset） 并向 producer 发送 ACK
### producer delivery guarantee生产者发送保证
Kafka默认提供的交付可靠性保障是第二种，即至少一次
- a. At most once 消息可能会丢，但绝不会重复传输
- b. At least one 消息绝不会丢，但可能会重复传输
- c. Exactly once 每条消息肯定会被传输一次且仅传输一次

## 三. broker 保存消息

所理解的“消息”，在Kafka中被称为日志。
在每一个broker中，保存了多个名字为{Topic}-{Partition}的文件夹，例如Test-1、Test-2.这里的意思是，这个broker中能够处理topic为Test，分区为1和2的消息

1. 存储方式   
![](./img/topic_partition.png)   
    物理上把 topic 分成一个或多个 partition（对应 server.properties 中的 num.partitions=3 配置），每个 partition 物理上对应一个文件夹（该文件夹存储该 partition 的所有消息和索引文件）
    在这张图中，一个分区，包含了多个Log Segment。注意，这里的Log Segment也是逻辑上的概念，只有具体到具体的日志文件，才是物理上的概念.
    我们看图片最右边的部分，文件名都是20位的整数，这个数字称为消息的“基准偏移量”。
    例如我们第二个Log Segment是从121开始的，那么代表了这个日志段的第一条消息的偏移量是从121开始的，也代表了在这之前有121条日志记录.
    注意，因为我们的偏移量是从0开始的，所以在121这个偏移量之前有121条数据，而不是120条

2. 存储策略   
    无论消息是否被消费，kafka 都会保留所有消息。有两种策略可以删除旧数据：
    - 2.1  基于时间：log.retention.hours=168
    - 2.2  基于大小：log.retention.bytes=1073741824
    注意：因为Kafka读取特定消息的时间复杂度为O(1)，即与文件大小无关，所以这里删除过期文件与提高 Kafka 性能无关
   
3. Topic创建与删除top   
   - 3.1 controller 在 ZooKeeper 的 /brokers/topics 节点上注册 watcher，当 topic 被创建，则 controller 会通过 watch 得到该 topic 的 partition/replica 分配。
   - 3.2 controller从 /brokers/ids 读取当前所有可用的 broker 列表，对于 set_p 中的每一个 partition：
     - a 从分配给该 partition 的所有 replica（称为AR）中任选一个可用的 broker 作为新的 leader，并将AR设置为新的 ISR
     - b 将新的 leader 和 ISR 写入 /brokers/topics/[topic]/partitions/[partition]/state
   - 3.3 controller 通过 RPC 向相关的 broker 发送 LeaderAndISRRequest。

4. 删除topic    
    - 4.1 controller 在 zooKeeper 的 /brokers/topics 节点上注册 watcher，当 topic 被删除，则 controller 会通过 watch 得到该 topic 的 partition/replica 分配。
    - 4.2 若 delete.topic.enable=false，结束；否则 controller 注册在 /admin/delete_topics 上的 watch 被 fire，controller 通过回调向对应的 broker 发送 StopReplicaRequest
5. 文件格式    
    三种类型的文件，*.log、*.index、*.timeindex。
    - log格式的文件记录了消息
    - index是偏移量索引
    - timeindex 是时间戳索引。  
    broker在接收到生产者发过来的消息的时候，需要将消息写在最后的Log Segment中。这样还带来了一个好处，消息的写入是顺序的IO。也因为如此，最后的一个Log Segment，被称为“active Log Segment”

## 四. 消费者
### 两套 consumer API
- a.  The high-level Consumer API
- b.  The SimpleConsumer API   
  
具体描述
- 1.1 使用 high-level consumer API 可以是多线程的应用,注意事项
    - a. 如果消费线程大于 partition 数量，则有些线程将收不到消息
    - b. 如果 partition 数量大于线程数，则有些线程多收到多个 partition 的消息
    - c. 如果一个线程消费多个 partition，则无法保证你收到的消息的顺序，而一个 partition 内的消息是有序的
- 1.2 SimpleConsumer API,对 partition 有更多的控制权，那就应该使用
    - a. 多次读取一个消息
    - b. 只消费一个 partition 中的部分消息
    - c. 使用事务来保证一个消息仅被消费一次
      
使用此 API 时，partition、offset、broker、leader 等对你不再透明，需要自己去管理。你需要做大量的额外工作：

- a. 必须在应用程序中跟踪 offset，从而确定下一条应该消费哪条消息
- b. 应用程序需要通过程序获知每个 Partition 的 leader 是谁
- c. 需要处理 leader 的变更
      
SimpleConsumer API 的一般流程如下
- 1 查找到一个“活着”的 broker，并且找出每个 partition 的 leader
- 2 找出每个 partition 的 follower
- 3 定义好请求，该请求应该能描述应用程序需要哪些数据
- 4 fetch 数据
- 5 识别 leader 的变化，并对之作出必要的响应

### 消费组: 消费者是以消费者组的形式对外消费的。
原因：   
我们作一个假设，假设没有消费者组这种概念，我们现在有10个消费者订阅了同一个主题，那么当这个主题有新的消息之后，我们这10个消费者是不是应该去“抢消息”进行消费呢？  
   1. 这是一种浪费资源的表现。所以消费者组，也可以认为是一种更加合理分配资源，进行负载均衡的设计。
                 假设有5个消费者属于同一个消费者组，这个消费者组订阅了一个具有10个分区的主题，那么组内的每一个消费者，都会负责处理2个分区的消息
  
   2. 消费者组的设计还能够令我们很方便的横向扩展系统的消费能力。设想一下在我们发觉系统中消息堆积越来越多，消费速度跟不上生产速度的时候，只需要新增消费者，
                    并且将这个消费者划入原来的消费者组中，Kafka会自动调整组内消费者对分区的分配，这个过程称为《重平衡》。
  
   3. kafka 的分配单位是 partition。每个 consumer 都属于一个 group，一个 partition 只能被同一个 group 内的一个 consumer 所消费
                       （也就保障了一个消息只能被 group 内的一个 consumer 所消费），但是多个 group 可以同时消费这个 partition
  
   4. Kafka还支持多个消费者组订阅同一个主题，这样，相同的消息将被发送到所有订阅了这个主题的消费者组中。
                            注意：我们说到了同一分区只能被同一个消费者消费，但是这个说法的前提是这些消费者位于同一个消费者组。也就是说，不同消费者组内的消费者，是可以消费同一个主题分区的
  
   5. 消费方式   
               consumer 采用 pull 模式从 broker 中读取数据。
               push 模式很难适应消费速率不同的消费者，因为消息发送速率是由 broker 决定的。它的目标是尽可能以最快速度传递消息，
               但是这样很容易造成 consumer 来不及处理消息，典型的表现就是拒绝服务以及网络拥塞。而 pull 模式则可以根据 consumer 的消费能力以适当的速率消费消息。
               对于 Kafka 而言，pull 模式更合适，它可简化 broker 的设计，consumer 可自主控制消费消息的速率，同时 consumer 可以自己控制消费方式——即可批量消费也可逐条消费，同时还能选择不同的提交方式从而实现不同的传输语义
   6. 消费保证consumer delivery guarantee  
               如果将 consumer 设置为 autocommit，consumer 一旦读到数据立即自动 commit。如果只讨论这一读取消息的过程，那 Kafka 确保了 Exactly once。
               但实际使用中应用程序并非在 consumer 读取完数据就结束了，而是要进行进一步处理，而数据处理与 commit 的顺序在很大程度上决定了consumer delivery guarantee
      - a. 读完消息先 commit 再处理消息。
                      这种模式下，如果 consumer 在 commit 后还没来得及处理消息就 crash 了，下次重新开始工作后就无法读到刚刚已提交而未处理的消息，这就对应于 At most once
      - b. 读完消息先处理再 commit。
                         这种模式下，如果在处理完消息之后 commit 之前 consumer crash 了，下次重新开始工作时还会处理刚刚未 commit 的消息，实际上该消息已经被处理过了。这就对应于 At least once。
      - c. 如果一定要做到 Exactly once，就需要协调 offset 和实际操作的输出。
                            精典的做法是引入两阶段提交。如果能让 offset 和操作输入存在同一个地方，会更简洁和通用。这种方式可能更好，因为许多输出系统可能不支持两阶段提交。
                            比如，consumer 拿到数据后可能把数据放到 HDFS，如果把最新的 offset 和数据本身一起写到 HDFS，那就可以保证数据的输出和 offset 的更新要么都完成，要么都不完成，
                            间接实现 Exactly once。（目前就 high-level API而言，offset 是存于Zookeeper 中的，无法存于HDFS，而SimpleConsuemr API的 offset 是由自己去维护的，可以将之存于 HDFS 中）

   7. 消费者重平衡 consumer rebalance    
### 消费者重平衡 consumer rebalance
定义：某个消费组内的消费者就如何消费某个主题的所有分区达成一个共识的过程,
但是这个过程对Kafka的吞吐率影响是巨大的，因为这个过程有点像GC中的STW（世界停止），在Rebalance的时候，所有的消费者只能去做重平衡这一件事情，不能消费任何的消息。 
下面我们来说说哪些情况可能会导致Rebalance：

- 1 组内成员数量发生了变化.比如有新的Consumer实例加入组或者离开组，抑或是有Consumer实例崩溃被“踢出”组。
- 2 订阅主题的数量发生了变化.比如consumer.subscribe(Pattern.compile(“t.*c”))就表明该Group订阅所有以字母t开头、字母c结尾的主题。在Consumer Group的运行过程中，你新创建了一个满足这样条件的主题，那么该Group就会发生Rebalance
- 3 订阅主题的分区数量发生了变化.Kafka当前只能允许增加一个主题的分区数。
  
而且在Rebalance的时候，假设有消费者退出了，导致多出了一些分区，Kafka并不是把这几个多出来的分区分配给原来的那些消费者，而是所有的消费者一起参与重新分配所有的分区
当有新的消费者加入的时候，也不是原本的每个消费者分出一些分区给新的消费者，而是所有的消费者一起参与重新分配所有的分区。
这样的分配策略听起来就很奇怪且影响效率，但是没有办法。

不过社区于0.11.0.0版本推出了StickyAssignor（粘性分配）策略，就可以做到我们上面假设的情况，但是目前还存在一些bug    
- a. 将目标 topic 下的所有 partition 排序，存于PT
- b. 对某 consumer group 下所有 consumer 排序，存于 CG，第 i 个consumer 记为 Ci
- c. N=size(PT)/size(CG)，向上取整
- d. 解除 Ci 对原来分配的 partition 的消费权（i从0开始）
- e. 将第i*N到（i+1）*N-1个 partition 分配给 Ci


如何减少非必要重平衡?
后面两个通常都是运维的主动操作，所以它们引发的 Rebalance 大都是不可避免的。接下来，我们主要关心因为组成员数量变化而引发的Rebalance该如何避免?
1. 第一类非必要 Rebalance 是因为未能及时发送心跳，导致Consumer被“踢出”Group而引发的. 设置session.timeout.ms和heartbeat.interval.ms的值
2. 第二类非必要 Rebalance 是Consumer消费时间过长导致的.设置 max.poll.interval.ms参数值

## 五. kafka 高可用 HA
### 副本 replication   
同一个 partition 可能会有多个 replica（对应 server.properties 配置中的 default.replication.factor=N）。
没有 replica 的情况下，一旦 broker 宕机，其上所有 partition 的数据都不可被消费，同时 producer 也不能再将数据存于其上的 partition。
引入replication 之后，同一个 partition 可能会有多个 replica，而这时需要在这些 replica 之间选出一个 leader，
producer 和 consumer 只与这个 leader 交互，其它 replica 作为 follower 从 leader 中复制数据   

复制算法：
- a. 将所有 broker（假设共 n 个 broker）和待分配的 partition 排序
- b. 将第 i 个 partition 分配到第（i mod n）个 broker 上
- c. 将第 i 个 partition 的第 j 个 replica 分配到第（(i + j) mode n）个 broker上

kafka 多副本: 
Kafka 为分区引入了多副本（Replica）机制，通过增加副本数量可以提升容灾能力。
同一分区的不同副本中保存的是相同的消息（当然在同一时刻，副本之间可能并非完全一样），副本之间是“一主多从”的关系，其中leader副本负责处理读写请求，follower副本只负责与leader副本的消息同步。
副本处于不同的broker中，当leader副本出现故障时，从follower副本中重新选举新的leader副本对外提供服务。Kafka通过多副本机制实现了故障的自动转移，当Kafka集群中某个broker失效时仍然能保证服务可用。

Kafka集群中有4个broker，某个主题中有3个分区，且副本因子（即副本个数）也为3，如此每个分区便有1个leader副本和2个follower副本。
生产者和消费者只与leader副本进行交互，而follower副本只负责消息的同步，很多时候follower副本中的消息相对leader副本而言会有一定的滞后。

### 宕机的场景
1. leader宕机
   当 partition 对应的 leader 宕机时，需要从 follower 中选举出新 leader。在选举新leader时，一个基本的原则是，新的 leader 必须拥有旧 leader commit 过的所有消息。
   kafka 在 zookeeper 中（/brokers/.../state）动态维护了一个 ISR（in-sync replicas），
   ISR 里面的所有 replica 都跟上了 leader，只有 ISR 里面的成员才能选为 leader。对于 f+1 个 replica，一个 partition 可以在容忍 f 个 replica 失效的情况下保证消息不丢失。
   当所有 replica 都不工作时，有两种可行的方案
    - a. 等待 ISR 中的任一个 replica 活过来，并选它作为 leader。可保障数据不丢失，但时间可能相对较长。
    - b. 选择第一个活过来的 replica（不一定是 ISR 成员）作为 leader。无法保障数据不丢失，但相对不可用时间较短

2. 中间人broker宕机
    - 3.1 controller 在 zookeeper 的 /brokers/ids/[brokerId] 节点注册 Watcher，当 broker 宕机时 zookeeper 会 fire watch
    - 3.2 controller 从 /brokers/ids 节点读取可用broker
    - 3.3 controller决定set_p，该集合包含宕机 broker 上的所有 partition
    - 3.4 对 set_p 中的每一个 partition
        - 3.4.1 从/brokers/topics/[topic]/partitions/[partition]/state 节点读取 ISR
        - 3.4.2 决定新 leader
        - 3.4.3 将新 leader、ISR、controller_epoch 和 leader_epoch 等信息写入 state 节点
    - 3.5 通过 RPC 向相关 broker 发送 leaderAndISRRequest 命令

3. 控制器controller宕机   
   当 controller 宕机时会触发 controller fail over.
   每个 broker 都会在 zookeeper 的 "/controller" 节点注册 watcher，当 controller 宕机时 zookeeper 中的临时节点消失，
   所有存活的 broker 收到 fire 的通知，每个 broker 都尝试创建新的 controller path，只有一个竞选成功并当选为 controller。
   当新的 controller 当选时，会触发 KafkaController.onControllerFailover 方法
    - 4.1. 读取并增加 Controller Epoch。
    - 4.2. 在 reassignedPartitions Patch(/admin/reassign_partitions) 上注册 watcher。
    - 4.3. 在 preferredReplicaElection Path(/admin/preferred_replica_election) 上注册 watcher。
    - 4.4. 通过 partitionStateMachine 在 broker Topics Patch(/brokers/topics) 上注册 watcher。
    - 4.5. 若 delete.topic.enable=true（默认值是 false），则 partitionStateMachine 在 Delete Topic Patch(/admin/delete_topics) 上注册 watcher。
    - 4.6. 通过 replicaStateMachine在 Broker Ids Patch(/brokers/ids)上注册Watch。
    - 4.7. 初始化 ControllerContext 对象，设置当前所有 topic，“活”着的 broker 列表，所有 partition 的 leader 及 ISR等。
    - 4.8. 启动 replicaStateMachine 和 partitionStateMachine。
    - 4.9. 将 brokerState 状态设置为 RunningAsController。
    - 4.10. 将每个 partition 的 Leadership 信息发送给所有“活”着的 broker。
    - 4.11. 若 auto.leader.rebalance.enable=true（默认值是true），则启动 partition-rebalance 线程。
    - 4.12. 若 delete.topic.enable=true 且Delete Topic Patch(/admin/delete_topics)中有值，则删除相应的Topic。
  
## 六. 开发注意事项  
1. producer 无法发布消息到 broker（奇怪也没有抛错）   
解决方式：server.properties 配置   
advertised.listeners 是 broker 给 producer 和 consumer 连接使用的，如果没有设置，就使用 listeners，而如果 host_name 没有设置的话，就使用 java.net.InetAddress.getCanonicalHostName() 方法返回的主机名

## 七. 位移主题   
在Kafka中的主题名称是__consumer_offsets。因为位移主题也是一个主题。
### 位移主题的消息格式
因为我们是希望保存位移，所以很容易会想到这是一个KV结构。那么Key中应该保存哪些消息呢？
结论： key-value结构

    Key中包含了主题名，分区名，消费者组名。其实在这里是不需要保存消费者id之类的信息的，也就是说只需要具体到是哪个消费者组在哪个主题的哪个分区消费了多少数据，就足够。
    Value中，只需要保存消费位移，就足够了.消费者是可能发生变动的，我们的目的是让消费者发生变动后，能知道从哪里继续消费

### 消费者在消费了信息之后，该把位移发送到哪呢？
Kafka中的位移主题会在第一个消费者被创建的时候创建，默认分区数是50，副本数是3。消费者在提交位移的时候，会根据自己组id的hash值模位移主题的分区数，
所得到的结果就是位移信息该提交的分区id，然后找到这个分区id的leader节点，将位移信息提交到这个leader节点所在的broker中

### 位移的提交
位移的提交与消息丢失和重复消费
虽然有了位移主题这样的设计，但是并不代表了消息一定不会被重复消费，也不代表消息一定不会丢失。
另外，Kafka会严格的执行位移主题中所提交的信息。例如已经消费了0-20的消息，如果你提交的位移是100，那么下一次拉取的信息一定是从100开始的，20-99的消息将会丢失。又比如你提交的位移是10，那么10-20的消息将会被重复消费

提交方式：一种是自动提交，一种是手动提交  

自动提交： enable.auto.commit  
1. 在某一时刻提交了位移100，随后你拉取了100-150的消息，但是还没有到下一次提交位移的时候，消费者宕机了。可能这个时候只消费了100-120的消息，那么在消费者重启后，因为120的位移没有提交，所以这部分的消息会被重复消费一次。
2. 再设想一种情况，你拉取了100-150的消息，这个时候到了自动提交的时间，提交了150的这个位移，而这个时候消费者宕机了，重启之后会从150开始拉取信息处理，那么在这之前的信息将会丢失

手动提交：  
手动提交又分为同步提交和异步提交两种提交方式。
1. 同步提交会直到消息被写入了位移主题，才会返回，这样是安全的，但是可能造成的问题是TPS降低。
2. 异步提交是触发了提交这个操作，就会返回。这样速度是很快的，但是可能会造成提交失败的情况

## 八. 安全认证
SASL(Simple Authentication and Security Layer 简单验证和安全层):用来认证 C/S 模式也就是服务器与客户端的一种认证机制,通俗的话来讲就是让服务器知道连接进来的客户端的身份是谁。
而 SASL 只是一种模式，需要依赖于具体的连接媒介，比如 JAAS(java Authentication Authorization Service)客户端、GSSAPI(Kerberos)、PLAIN、SCRAM-SHA-256、SCRAM-SHA-512、OAuthBearer 等等 
1. GSSAPI： 使用的Kerberos认证，可以集成目录服务，比如AD。从Kafka0.9版本开始支持
2. PLAIN： 使用简单用户名和密码形式。从Kafka0.10版本开始支持
3. SCRAM(Salted Challenge Response Authentication Mechanism)： 主要解决PLAIN动态更新问题以及安全机制，从Kafka0.10.2开始支持
4. OAUTHBEARER：基于OAuth 2认证框架，从Kafka2.0版本开始支持

## 九. 消息重复和消费幂等
消息队列Kafka版消费的语义是at least once， 也就是至少投递一次，保证消息不丢失，但是无法保证消息不重复。在出现网络问题、客户端重启时均有可能造成少量重复消息，此时应用消费端如果对消息重复比较敏感（例如订单交易类），则应该做消息幂等。
以数据库类应用为例，常用做法是：
1.发送消息时，传入key作为唯一流水号ID。
2.消费消息时，判断key是否已经消费过，如果已经消费过了，则忽略，如果没消费过，则消费一次

### 消费阻塞以及堆积
消费端最常见的问题就是消费堆积，最常造成堆积的原因是：
1.消费速度跟不上生产速度，此时应该提高消费速度，详情请参见提高消费速度。
2.消费端产生了阻塞。
3.消费端拿到消息后，执行消费逻辑，通常会执行一些远程调用，如果这个时候同步等待结果，则有可能造成一直等待，消费进程无法向前推进。

消费端应该竭力避免堵塞消费线程，如果存在等待调用结果的情况，建议设置等待的超时时间，超时后作为消费失败进行处理

### 分区个数
分区个数主要影响的是消费者的并发数量。

对于同一个Consumer Group内的消费者来说，一个分区最多只能被一个消费者消费。因此，消费实例的个数不要大于分区的数量，否则会有消费实例分配不到任何分区而处于空跑状态。
一般来说，不建议分区数小于12，否则可能影响消费发送性能； 也不建议超过100个，否则易引发消费端Rebalance。
控制台的默认分区个数是12，可以满足绝大部分场景的需求。您可以根据业务使用量进行增加。


## 十. 面试问题
### 为什么Kafka不支持读写分离？
因为这样有两个明显的缺点：

1. 数据一致性问题。数据从主节点转到从节点必然会有一个延时的时间窗口，这个时间窗口会导致主从节点之间的数据不一致。
2. 延时问题。数据从写入主节点到同步至从节点中的过程需要经历网络→主节点内存→主节点磁盘→网络→从节点内存→从节点磁盘这几个阶段。对延时敏感的应用而言，主写从读的功能并不太适用。
对于Kafka来说，必要性不是很高，因为在Kafka集群中，如果存在多个副本，经过合理的配置，可以让leader副本均匀的分布在各个broker上面，使每个 broker 上的读写负载都是一样的。


### Kafka的那些设计让它有如此高的性能
1. 分区   
   
    kafka是个分布式集群的系统，整个系统可以包含多个broker，也就是多个服务器实例。每个主题topic会有多个分区，kafka将分区均匀地分配到整个集群中，当生产者向对应主题传递消息，消息通过负载均衡机制传递到不同的分区以减轻单个服务器实例的压力。
    一个Consumer Group中可以有多个consumer，多个consumer可以同时消费不同分区的消息，大大的提高了消费者的并行消费能力。但是一个分区中的消息只能被一个Consumer Group中的一个consumer消费。

2. 网络传输上减少开销
   批量发送：
   在发送消息的时候，kafka不会直接将少量数据发送出去，否则每次发送少量的数据会增加网络传输频率，降低网络传输效率。kafka会先将消息缓存在内存中，当超过一个的大小或者超过一定的时间，那么会将这些消息进行批量发送。
   端到端压缩：
   当然网络传输时数据量小也可以减小网络负载，kafka会将这些批量的数据进行压缩，将一批消息打包后进行压缩，发送broker服务器后，最终这些数据还是提供给消费者用，所以数据在服务器上还是保持压缩状态，不会进行解压，而且频繁的压缩和解压也会降低性能，最终还是以压缩的方式传递到消费者的手上。

3. 顺序读写
    kafka将消息追加到日志文件中，利用了磁盘的顺序读写，来提高读写效率。
   
4. 零拷贝技术
![](.kafka_intro_images/copy_before.png)
![](.kafka_intro_images/copy_after.png)  
   
    零拷贝将文件内容从磁盘通过DMA引擎复制到内核缓冲区，而且没有把数据复制到socket缓冲区，只是将数据位置和长度信息的描述符复制到了socket缓存区，然后直接将数据传输到网络接口，最后发送。这样大大减小了拷贝的次数，提高了效率。
   kafka正是调用linux系统给出的sendfile系统调用来使用零拷贝。Java中的系统调用给出的是FileChannel.transferTo接口。

5. 优秀的文件存储机制

   如果分区规则设置得合理，那么所有的消息可以均匀地分布到不同的分区中，这样就可以实现水平扩展。不考虑多副本的情况，一个分区对应一个日志（Log）。为了防止 Log 过大，Kafka 又引入了日志分段（LogSegment）的概念，将 Log 切分为多个 LogSegment，相当于一个巨型文件被平均分配为多个相对较小的文件，这样也便于消息的维护和清理。

   Kafka 中的索引文件以稀疏索引（sparse index）的方式构造消息的索引，它并不保证每个消息在索引文件中都有对应的索引项。每当写入一定量（由 broker 端参数 log.index.interval.bytes 指定，默认值为4096，即 4KB）的消息时，偏移量索引文件和时间戳索引文件分别增加一个偏移量索引项和时间戳索引项，增大或减小 log.index.interval.bytes 的值，对应地可以增加或缩小索引项的密度。

### kafka中的事务

Kafka中的事务可以使应用程序将消费消息、生产消息、提交消费位移当作原子操作来处理，同时成功或失败，即使该生产或消费会跨多个分区。

设置事务型Producer的方法也很简单，满足两个要求即可：

- 和幂等性Producer一样，开启enable.idempotence = true。
- 设置Producer端参数transactional. id。最好为其设置一个有意义的名字。


生产者必须提供唯一的transactionalId，启动后请求事务协调器获取一个PID，transactionalId与PID一一对应。

每次发送数据给<Topic, Partition>前，需要先向事务协调器发送AddPartitionsToTxnRequest，事务协调器会将该<Transaction, Topic, Partition>存于__transaction_state内，并将其状态置为BEGIN。

在处理完 AddOffsetsToTxnRequest 之后，生产者还会发送 TxnOffsetCommitRequest 请求给 GroupCoordinator，从而将本次事务中包含的消费位移信息 offsets 存储到主题 __consumer_offsets 中

一旦上述数据写入操作完成，应用程序必须调用KafkaProducer的commitTransaction方法或者abortTransaction方法以结束当前事务。
无论调用 commitTransaction() 方法还是 abortTransaction() 方法，生产者都会向 TransactionCoordinator 发送 EndTxnRequest 请求。
TransactionCoordinator 在收到 EndTxnRequest 请求后会执行如下操作：

1. 将 PREPARE_COMMIT 或 PREPARE_ABORT 消息写入主题 __transaction_state
2. 通过 WriteTxnMarkersRequest 请求将 COMMIT 或 ABORT 信息写入用户所使用的普通主题和 __consumer_offsets
3. 将 COMPLETE_COMMIT 或 COMPLETE_ABORT 信息写入内部主题 __transaction_state标明该事务结束
   在消费端有一个参数isolation.level，设置为“read_committed”，表示消费端应用不可以看到尚未提交的事务内的消息。
   如果生产者开启事务并向某个分区值发送3条消息 msg1、msg2 和 msg3，在执行 commitTransaction() 或 abortTransaction() 方法前，
   设置为“read_committed”的消费端应用是消费不到这些消息的，不过在 KafkaConsumer 内部会缓存这些消息，直到生产者执行 commitTransaction() 方法之后它才能将这些消息推送给消费端应用。
   反之，如果生产者执行了 abortTransaction() 方法，那么 KafkaConsumer 会将这些缓存的消息丢弃而不推送给消费端应用。


## 十一 kafka主流 golang 客户端


- github.com/confluentinc/confluent-kafka-go 是confluent公司开发的kafka golang sdk，由于confluent公司维护,cgo 包装librdkafka
- github.com/IBM/sarama IBM组织 , 暂不支持contexts, 值作为指针分配会造成内存动态分配,从而更高的内存占用率
- github.com/segmentio/kafka-go 
- github.com/lovoo/goka 

## 参考

- [Kafka 核心技术与实战](https://time.geekbang.org/column/intro/191?tab=catalog)

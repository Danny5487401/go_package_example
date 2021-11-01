#Redis:

![](img/redis_info.png)

	Remote Dictionary Server, 翻译为远程字典服务. Redis是一个完全开源的基于Key-Value的NoSQL存储系统，他是一个使用ANSIC语言编写的，
	遵守BSD协议，支持网络、可基于内存的可持久化的日志型、Key-Value数据库,并提供多种语言的API.
执行原理：

	# 1. 客户端发送命令后，Redis服务器将为这个客户端链接创造一个'输入缓存'，将命令放到里面
	# 2. 再由Redis服务器进行分配挨个执行，顺序是随机的，这将不会产生并发冲突问题，也就不需要事务了.
	# 3. 再将结果返回到客户端的'输出缓存'中，输出缓存先存到'固定缓冲区',如果存满了，就放入'动态缓冲区',客户端再获得信息结果

	# 如果数据时写入命令，例如set name:1  zhangsan 方式添加一个字符串.
	# Redis将根据策略，将这对key:value来用内部编码格式存储，好处是改变内部编码不会对外有影响，正常操作即可,
	# 同时不同情况下存储格式不一样，发挥优势.
##为什么需要Redis?

	传统数据库在存储数据时，需要先定义schema，以确定类型(字节宽度)，并以行存储，所以每行所占的字节宽度是一致的（便于索引数据）。
	数据库内部数据分为多个datapage(一般是4kb)存储，随着数据量的增大，数据库查询速度会越来越慢，其主要瓶颈在于磁盘I/O。
	由于数据量增大查找datapage的时间也会变长，所以索引出现了。索引是一个B+T，存储在内存中，根据索引记录的信息，可以快速定位到datapage的位置。
    索引虽然会大大加快查询速度，但是因为增删改需要维护索引的B+T，所以会把增删改的速度拖慢，所以索引不适合频繁写的表。
    
	此外，当数据库高并发查询的情况下，单位时间内所需的数据量也是很大的，此时可能会受到磁盘带宽的影响，影响磁盘的查询速度。
	在I/O上，内存相比较于磁盘，拥有较好的性能;
	出现了一批基于内存的关系型数据库，比如SAP HAHA数据库，其物理机器内存2T，包含软件以及服务，购买需要1亿元,由于内存关系型数据库的昂贵价格，
	所以大部分公司采用了折中的方案,使用磁盘关系型数据库+内存缓存,比如 Oracle+Memcached,Mysql+Redis

##Redis协议
![](.redis_images/redis_scheme.png)
Redis客户端使用RESP（Redis的序列化协议）协议与Redis的服务器端进行通信。 虽然该协议是专门为Redis设计的，但是该协议也可以用于其他 客户端-服务器 （Client-Server）软件项目。RESP是对以下几件事情的折中实现：
    
    1、实现简单
    
    2、解析快速
    
    3、人类可读
RESP实际上是一个支持以下数据类型的序列化协议：简单字符串（Simple Strings），错误（Errors），整数（Integers），块字符串（Bulk Strings）和数组（Arrays）

    RESP可以序列化不同的数据类型，如整数（integers），字符串（strings），数组（arrays）。它还使用了一个特殊的类型来表示错误（errors）。
    请求以字符串数组的形式来表示要执行命令的参数从客户端发送到Redis服务器。Redis使用命令特有（command-specific）数据类型作为回复。
    
    RESP协议是二进制安全的，并且不需要处理从一个进程传输到另一个进程的块数据的大小，因为它使用前缀长度（prefixed-length）的方式来传输块数据的
在Redis中,RESP用作 请求-响应 协议的方式如下：

    1、客户端将命令作为批量字符串的RESP数组发送到Redis服务器。
    
    2、服务器（Server）根据命令执行的情况返回一个具体的RESP类型作为回复。

在RESP协议中，有些的数据类型取决于第一个字节：

    1、对于简单字符串，回复的第一个字节是“+”
    
    2、对于错误，回复的第一个字节是“ - ”
    
    3、对于整数，回复的第一个字节是“：”
    
    4、对于批量字符串，回复的第一个字节是“$”
    
    5、对于数组，回复的第一个字节是“*”


##Reactor 单线程的Redis为什么这么快?
![](.redis_images/reason_why_redis_so_fast.png)

	# 1. 基于内存的访问，非阻塞I/O，Redis使用事件驱动模型epoll多路复用实现，连接，读写，关闭都转换为事件不在网络I/O上浪费过多的时间.
	# 2. 单线程避免高并发的时候，多线程有锁的问题和线程切换的CPU开销问题.《虽然是单线程，但可以开多实例弥补》
	# 3. 使用C语言编写，更好的发挥服务器性能，并且代码简介，性能高
	
Redis6.0后引入多线程提速：
![](.redis_images/reactor_model.png)

    读写网络的read/write系统耗时 >> Redis运行执行耗时，Redis的瓶颈主要在于网络的 IO 消耗, 优化主要有两个方向:
    
    a.提高网络 IO 性能，典型的实现比如使用 DPDK 来替代内核网络栈的方式
        使用多线程充分利用多核，典型的实现比如 Memcached
        
    b.可以充分利用服务器 CPU 资源，目前主线程只能利用一个核
        多线程任务可以分摊 Redis 同步 IO 读写负荷
关于多线程

    Redis 6.0 版本 默认多线程是关闭的 io-threads-do-reads no
    Redis 6.0 版本 开启多线程后 线程数也要 谨慎设置。
    多线程可以使得性能翻倍，但是多线程只是用来处理网络数据的读写和协议解析，执行命令仍然是单线程顺序执行
    
	
##Redis五种数据类型应用场景
![](img/string.png)

	1.String(sds,simple dynamic string简单动态字符串): 常规的set/get操作,因为string 类型是二进制安全的,可以用来存放图片，视频等内容.
	    C 中字符串遇到 '\0' 会结束，那 '\0' 之后的数据就读取不上了。但在 SDS 中，是根据 len 长度来判断字符串结束的,这样二进制安全的问题就解决了.
	    
        另外由于Redis的高性能读写功能，而string类型的value也可以是数字，一般做一些复杂的计数功能的缓存,还可以用作计数器（INCR,DECR），
		比如分布式环境中统计系统的在线人数，秒杀等

![](img/hash.png)

	2. hash(hash table): value 存放的是键值对结构化后的对象，将一些相关的数据存储在一起，比如用户的购物车。
	    比较方便操作其中某个字段，比如可以做单点登录存放用户信息,以cookiele作为key，
		设置30分钟为缓存过期时间，能很好的模拟出类似session的效
		
		这里需要明确一点： Redis中只有一个K，一个V。其中 K 绝对是字符串对象，而 V 可以是String、List、Hash、Set、ZSet任意一种
		
![](img/list.png)

	3. list(deque): 发布和订阅、慢查询、监视器等都用到了链表，Redis服务器本身还是用链表来保存多个客户端的状态信息，以及使用链表来构建客户端输出缓冲区。
    另外可以利用lrange命令，做基于redis的分页功能,性能极佳，用户体验好
	
![](img/set.png)

	4. set(intset+dict):由于底层是字典实现的，查找元素特别快，另外set 数据类型不允许重复，利用这两个特性我们可以进行全局去重，
		比如在用户注册模块，判断用户名是否注册；另外就是利用交集、并集、差集等操作，可以计算共同喜好，全部的喜好，自己独有的喜好等功能

![](img/sort_set.png)		

	5. Zset(skip list + hash table):有序的集合，可以做范围查找，排行榜应用，取 TOP N 操作等,还可以做延时任务
	
扩展：
1.bitmap
    
    BitMap 原本的含义是用一个比特位来映射某个元素的状态。由于一个比特位只能表示 0 和 1 两种状态，所以 BitMap 能映射的状态有限，
    但是使用比特位的优势是能大量的节省内存空间。
    
    在 Redis 中BitMap 底层是基于字符串类型实现的，可以把 Bitmaps 想象成一个以比特位为单位的数组，数组的每个单元只能存储0和1，
    数组的下标在 Bitmaps 中叫做偏移量，BitMap 的 offset 值上限 2^32 - 1。
    
用途

    签到： key = 年份：用户id offset = （今天是一年中的第几天） % （今年的天数）
    统计活跃用户: 使用日期作为 key，然后用户 id 为 offset 设置不同offset为0 1 即可。
    
    
2.HyperLogLog

    是一种概率数据结构，它使用概率算法来统计集合的近似基数。而它算法的最本源则是伯努利过程 + 分桶 + 调和平均数。
    
    功能：误差允许范围内做基数统计 (基数就是指一个集合中不同值的个数) 的时候非常有用，每个HyperLogLog的键可以计算接近2^64不同元素的基数，
    而大小只需要12KB。错误率大概在0.81%。所以如果用做 UV 统计很合适
    
    HyperLogLog底层 一共分了 2^14 个桶，也就是 16384 个桶。每个(registers)桶中是一个 6 bit 的数组，这里有个骚操作就是一般人可能直接用一个字节当桶浪费2个bit空间，但是Redis底层只用6个然后通过前后拼接实现对内存用到了极致，最终就是 16384*6/8/1024 = 12KB。
    
3.Bloom Filter

    使用布隆过滤器得到的判断结果： 不存在的一定不存在，存在的不一定存在。
    
##RedisDB内部结构
![](.redis_images/redis_internal_structure.png)
![](.redis_images/redis_db_structure.png)
String：存储数字的话，采用int类型的编码，如果是非数字的话，采用 raw 编码

list编码：字符串长度且元素个数小于一定范围使用 ziplist 编码，否则转化为 linkedlist 编码

set编码：保存元素为整数及元素个数小于一定范围使用 intset 编码，任意条件不满足，则使用 hashtable 编码。

zset编码：保存的元素个数小于定值且成员长度小于定值使用 ziplist 编码，任意条件不满足，则使用 skiplist 编码


##Redis数据类型底层数据结构
![](.redis_images/redis_data_structure.png)

###1.sds简单动态字符串
![](.redis_images/c_string_structure.png)
c语言自带的字符串,不过是一个以0结束的字符数组.想要获取 「Redis」的长度，需要从头开始遍历，直到遇到 '\0' 为止
![](.redis_images/c_string.png)
想要获取 「Redis」的长度，需要从头开始遍历，直到遇到 '\0' 为止。
而在redis中，想要获取长度只需要获取 len 字段即可


![](.redis_images/string_structure.png)
自动存储int类型，非int类型用raw编码
![](.redis_images/sdshdr.png)


```c   
typedef char *sds;
```
sds字符串根据字符串的长度，划分了五种结构体sdshdr5、sdshdr8、sdshdr16、sdshdr32、sdshdr64,分别对应的类型为SDS_TYPE_5、SDS_TYPE_8、SDS_TYPE_16、SDS_TYPE_32、SDS_TYPE_64
每个sds 所能存取的最大字符串长度为：

    sdshdr5最大为32(2^5)
    sdshdr8最大为0xff(2^8-1)
    sdshdr16最大为0xffff(2^16-1)
    sdshdr32最大为0xffffffff(2^32-1)
    sdshdr64最大为(2^64-1)
SDS_TYPE_8结构体
```c
struct __attribute__ ((__packed__)) sdshdr8 {
    uint8_t len; /* used */
    uint8_t alloc; /* excluding the header and null terminator */
    unsigned char flags; /* 3 lsb of type, 5 unused bits */
    char buf[];
};

```
###2.hash(ziplist+dict)

![](.redis_images/dict.png)
字典的结构体
```c
//dict.h
typedef struct dict {

    // 包括一些自定义函数，这些函数使得key和value能够存储
    dictType *type;

    void *privdata;
    // ht是一个长度为2的数组，对应的是两个哈希表，一般使用使用ht[0],ht[1]主要在扩容和缩容时使用。
    dictht ht[2];

    long rehashidx; /* 是一个标志量，如果为-1说明当前没有扩容，如果不为 -1 则记录扩容位置 */
    unsigned long iterators; /*当前字典正在进行中的迭代器 */
} dict;

```

哈希表结构体,数据 dictEntry 类型的数组，每个数组的item可能都指向一个链表。
```c
typedef struct dictht {
	//哈希表数组,对应的是多个哈希表节点dictEntry
    dictEntry **table;
    //哈希表大小
    unsigned long size;

   	//哈希表大小的掩码,用于计算索引值
   	//总是等于size-1
    unsigned long sizemask;

   	//已有节点的数量
    unsigned long used;
} dictht;
```

哈希表节点key/value结构体定义，真正的数据节点
![](.redis_images/key_value.png)
```c
typedef struct dictEntry {
	//键
    void *key;
    //值
    union {
        void *val;
        uint64_t u64;
        int64_t s64;
        double d;
    } v;
    //指向下一个哈希表节点,形成链表
    struct dictEntry *next;
} dictEntry
```
哈希表每个节点都保存着一个键值对，key就是键值对的键，v属性就是对应键值对的值，v可以是一个指针也可以是uint64_t,整数也可以是int64_t整数。 

next是一个链表，指向着下一个哈希表节点，这个指针可以将多个哈希值相同的键值对连接在一起，以此来解决哈希冲突问题。

###3.set(intset+dict)

![](.redis_images/intset.png)
![](.redis_images/dict.png)
整数集合intset是集合键的底层实现之一，当一个集合只包含整数值元素，并且这个集合的元素数量不多时，Redis就会使用整数集合键的底层实现。
```c
typedef struct intset {
	//编码方式
    uint32_t encoding;
    //集合包含的元素数量
    uint32_t length;
    //保存元素的数组
    int8_t contents[];
} intset;
```
contents数组是整数集合的底层实现：整数集合的每个元素都是contents数组的一个数据项(item),各个项在数组中按值得大小从小到大有序得排列，并且数组中不包含任何重复项

虽然intset结构将contents属性声明为int8_t类型的数组，但实际上contents并不保存任何int8_t类型的值，contents数组得真正类型取决于encoding属性的值:

    如果encoding属性得值为INTSET_ENC_INT16，那么contents就是一个int16_t类型的数组，数组里得每个项都是一个int16_t类型的整数值(最小值为-32768,最大值为32767)。
    
    如果encoding属性得值为INTSET_ENC_INT32，那么contents就是一个int32_t类型的数组，数组里的每个项都是一个int32_t类型的整数值(最小值为-2147483648,最大值为2147483648)。
    
    如果encoding属性的值为INTSET_ENC_INT64,那么contents就是一个int64_t类型的数组，数组里的每个项都是一个int64_t类型得整数值(最小值为-9223372036854775808,9223372036854775808)
length属性记录了整数集合包含得元素数量，也即是contents数组得长度。

###4.list双端链表
![](.redis_images/list_n_node.png)
list结构体
```
typedef struct list {
    listNode *head;//表头节点
    listNode *tail;//表尾节点
    //节点值复制函数
    void *(*dup)(void *ptr);
    //节点值释放函数
    void (*free)(void *ptr);
    //节点值对比函数
    int (*match)(void *ptr, void *key);
    unsigned long len;
} list;
```
节点
```c
//链表节点adlist.h/listNode
typedef struct listNode {
    struct listNode *prev;//前置节点
    struct listNode *next;//后置节点
    void *value;//节点值
} listNode;


```
压缩列表
![](.redis_images/ziplist.png)
![](.redis_images/ziplist_info.png)

    如果在一个链表节点中存储一个小数据，比如一个字节。那么对应的就要保存头节点，前后指针等额外的数据。
 
    这样就浪费了空间，同时由于反复申请与释放也容易导致内存碎片化。这样内存的使用效率就太低了。
    
    并且压缩列表的内存是连续分配的，遍历的速度很快。
###5.Sort Set(hash+skiptable)


![](.redis_images/skiptable.png)
跳跃表的结构
```c
//定义在server.h/zskiplist
typedef struct zskiplistNode {
    sds ele;//成员对象
    double score;//分数
    struct zskiplistNode *backward;
    //层
    struct zskiplistLevel {
    	//前进指针
        struct zskiplistNode *forward;
        //跨度
        unsigned long span;
    } level[];
    
} zskiplistNode;

typedef struct zskiplist {
    struct zskiplistNode *header, *tail;
    unsigned long length; // 跳跃表的长度
    int level; // 记录跳跃表内，层数最大的那个节点的层数
} zskiplist;
```
跳表就是多层链表的结合体，跳表分为许多层(level)，每一层都可以看作是数据的索引，这些索引的意义就是加快跳表查找数据速度

没有跳表查询时,查询数据37
![](.redis_images/without_skiptable.png)
有跳表查询37时
![](.redis_images/search_value_with_skiptable.png)

###6.stream(radix-tree)



PipeLine:

	Redis的pipeline功能的原理是 Client通过一次性将多条redis命令发往Redis Server，减少了每条命令分别传输的IO开销。
	同时减少了系统调用的次数，因此提升了整体的吞吐能力。
	我们在主-从模式的Redis中，pipeline功能应该用的很多，但是Cluster模式下，估计还没有几个人用过。
	我们知道 redis cluster 默认分配了 16384 个slot，当我们set一个key 时，会用CRC16算法来取模得到所属的slot，
	然后将这个key 分到哈希槽区间的节点上，具体算法就是：CRC16(key) % 16384。如果我们使用pipeline功能，
	一个批次中包含的多条命令，每条命令涉及的key可能属于不同的slot
	
##持久化方式

###1.RDB(redis database)：
    以快照的形式将数据持久化到磁盘
![](.redis_images/rdb_format.png)

    手段(会丢一部分最新数据) 
        * 每个Redis实例只会存一份rdb文件
        * 可以通过Save以及BGSAVE 来调用 
        * 二进制文件, lzf
    
###2.AOF(append only file)
    
    以日志的形式记录每个操作，将Redis执行过的所有指令全部记录下来（读操作不记录）， 只许追加文件但不可以修改文件， Redis启动时会读取AOF配置文件重构数据

    每次操作都需要fsync, 前台线程阻塞
    aof的内容就是redis标准协议

意义 
    将同一个key的反复操作，全部转为最后的值或multi集合(<64)

no-appendfsync-on-rewrite yes

    正导出rdb快照的过程中,要不要停止同步fsync 
auto-aof-rewrite-min-size 3000mb

    aof文件,至少超过3000M时, 再执行aof重写 
auto-aof-rewrite-percentage 80

    aof文件大小比起上次重写时的大小,增长率80%时,执行aof重写

Redis集群特点

    1. 多个redis节点网络互联，数据共享  
    2. 所有的节点都是一主一从（也可以是一主多从），其中从不提供服务，仅作为备用  
    3. 不支持同时处理多个key（如MSET/MGET），因为redis需要把key均匀分布在各个节点上，
    并发量很高的情况下同时创建key-value会降低性能并导致不可预测的行为  
    4. 支持在线增加、删除节点  
    5. 客户端可以连任何一个主节点进行读写

集群方案
![](.redis_images/codis_vs_cluster.png)
    1. vip多线程版本 twemproxy(Twitter开源)
    2. codis
![](.redis_images/codis.png)
    3. redis cluster
![](.redis_images/redis_cluster.png)

#高并发缓存
![](.redis_images/high_concurrency_buffer.png)


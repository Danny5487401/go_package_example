<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [表引擎](#%E8%A1%A8%E5%BC%95%E6%93%8E)
  - [MergeTree引擎](#mergetree%E5%BC%95%E6%93%8E)
  - [MergeTree 数据表的存储结构](#mergetree-%E6%95%B0%E6%8D%AE%E8%A1%A8%E7%9A%84%E5%AD%98%E5%82%A8%E7%BB%93%E6%9E%84)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


# 表引擎

表引擎（即表的类型）决定了：

- 数据的存储方式和位置，写到哪里以及从哪里读取数据
- 支持哪些查询以及如何支持。
- 并发数据访问。
- 索引的使用（如果存在）。
- 是否可以执行多线程请求。
- 数据复制参数。
  

ClickHouse 拥有非常庞大的表引擎体系，总共有合并树、外部存储、内存、文件、接口和其它 6 大类 20 多种表引擎，而在这众多的表引擎中，又属合并树（MergeTree）表引擎及其家族系列（*MergeTree）最为强大，在生产环境中绝大部分场景都会使用此引擎。
 

## MergeTree引擎
MergeTree这个名词是在我们耳熟能详的LSM Tree之上做减法而来——去掉了MemTable和Log。也就是说，向MergeTree引擎族的表插入数据时，数据会不经过缓冲而直接写到磁盘。

> MergeTree is not an LSM tree because it doesn’t contain "memtable" and "log": inserted data is written directly to the filesystem.
> This makes it suitable only to INSERT data in batches, not by individual row and not very frequently – about once per second is ok, but a thousand times a second is not.
> We did it this way for simplicity’s sake, and because we are already inserting data in batches in our applications.


社区通过 https://github.com/ClickHouse/ClickHouse/pull/8290  和 https://github.com/ClickHouse/ClickHouse/pull/10697 两个PR实现了名为Polymorphic Parts的特性，使得MergeTree引擎能够更好地处理频繁的小批量写入，但同时也标志着MergeTree的内核开始向真正的LSM Tree靠拢。


MergeTree 作为家族中最为基础的表引擎，提供了主键索引、数据分区、数据副本和数据采样等基本能力，而家族中的其它其它表引擎则在 MergeTree 的基础之上各有所长。
比如 ReplacingMergeTree 表引擎具有删除重复数据的特性，而 SummingMergeTree 表引擎则会按照排序键自动聚合数据。
如果再给合并树系列的表引擎加上 Replicated 前缀，又会得到一组支持数据副本的表引擎，例如 ReplicatedMergeTree、ReplicatedReplacingMergeTree、ReplicatedSummingMergeTree、ReplicatedAggregatingMergeTree 等等


```clickhouse
-- 使用 MergeTree 创建表
CREATE TABLE [IF NOT EXISTS] [db.]table_name [ON CLUSTER cluster]
(
    name1 [type1] [[NOT] NULL] [DEFAULT|MATERIALIZED|ALIAS|EPHEMERAL expr1] [COMMENT ...] [CODEC(codec1)] [STATISTICS(stat1)] [TTL expr1] [PRIMARY KEY] [SETTINGS (name = value, ...)],
    name2 [type2] [[NOT] NULL] [DEFAULT|MATERIALIZED|ALIAS|EPHEMERAL expr2] [COMMENT ...] [CODEC(codec2)] [STATISTICS(stat2)] [TTL expr2] [PRIMARY KEY] [SETTINGS (name = value, ...)],
    ...
    INDEX index_name1 expr1 TYPE type1(...) [GRANULARITY value1],
    INDEX index_name2 expr2 TYPE type2(...) [GRANULARITY value2],
    ...
    PROJECTION projection_name_1 (SELECT <COLUMN LIST EXPR> [GROUP BY] [ORDER BY]),
    PROJECTION projection_name_2 (SELECT <COLUMN LIST EXPR> [GROUP BY] [ORDER BY])
) ENGINE = MergeTree()
ORDER BY expr
[PARTITION BY expr]
[PRIMARY KEY expr]
[SAMPLE BY expr]
[TTL expr
    [DELETE|TO DISK 'xxx'|TO VOLUME 'xxx' [, ...] ]
    [WHERE conditions]
    [GROUP BY key_expr [SET v1 = aggr_func(v1) [, v2 = aggr_func(v2) ...]] ] ]
[SETTINGS name = value, ...]
```

- PARTITON BY：选填，表示分区键，用于指定表数据以何种标准进行分区。分区键既可以是单个字段、也可以通过元组的形式指定多个字段，同时也支持使用列表达式。
- ORDER BY：必填，表示排序键，用于指定在一个分区内，数据以何种标准进行排序。排序键既可以是单个字段，例如 ORDER BY CounterID，也可以是通过元组声明的多个字段，例如 ORDER BY (CounterID, EventDate)。
- setting 配置
  - index_granularity：对于 MergeTree 而言是一个非常重要的参数，它表示索引的粒度，默认值为 8192。所以 ClickHouse 根据主键生成的索引实际上稀疏索引，默认情况下是每隔 8192 行数据才生成一条索引


## MergeTree 数据表的存储结构

数据存储格式  Wide or Compact format. 
- In Wide format 每列在不同的文件
- in Compact format  所有列在一个文件

```shell
-- ClickHouse server version 25.1.4

-- 该表负责存储用户参加过的活动，每参加一个活动，就会生成一条记录
CREATE TABLE IF NOT EXISTS user_activity_event (
    ID UInt64,  -- 表的 ID
    UserName String,  -- 用户名
    ActivityName String,  -- 活动名称
    ActivityType String,  -- 活动类型
    ActivityLevel Enum('Easy' = 0, 'Medium' = 1, 'Hard' = 2),  -- 活动难度等级
    IsSuccess Int8,  -- 是否成功
    JoinTime DATE  -- 参加时间
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(JoinTime)  -- 按照 toYYYYMM(JoinTime) 进行分区
ORDER BY ID;  -- 按照 ID 字段排序

-- 插入一条数据
INSERT INTO user_activity_event VALUES (1, '张三', '寻找遗失的时间', '市场营销', 'Medium', 1, '2020-05-13')

```

```shell
root@635e708a264a:/var/lib/clickhouse# ls data/helloworld/user_activity_event/202005_1_1_0/ -al
total 44
drwxr-x--- 13 root       root       416 Feb 23 10:02 .
drwxr-x---  5 clickhouse clickhouse 160 Feb 23 10:02 ..
-rw-r-----  1 root       root       330 Feb 23 10:02 checksums.txt  # 校验文件，使用二进制的格式进行存储，它保存了余下各类文件（primary.txt、count.txt 等等）的 size 大小以及哈希值，用于快速校验文件的完整性和正确性
-rw-r-----  1 root       root       204 Feb 23 10:02 columns.txt
-rw-r-----  1 root       root         1 Feb 23 10:02 count.txt
-rw-r-----  1 clickhouse clickhouse 237 Feb 23 10:02 data.bin # 数据文件，使用压缩格式存储，默认为 LZ4 格式，用于存储表的数据。
-rw-r-----  1 clickhouse clickhouse  78 Feb 23 10:02 data.cmrk3 # 标记文件，使用二进制格式存储，标记文件中保存了 data.bin 文件中数据的偏移量信息，并且标记文件与稀疏索引对齐，因此 MergeTree 通过标记文件建立了稀疏索引（primary.idx）与数据文件（data.bin）之间的映射关系。而在读取数据的时候，首先会通过稀疏索引（primary.idx）找到对应数据的偏移量信息（data.mrk），因为两者是对齐的，然后再根据偏移量信息直接从 data.bin 文件中读取数据。
-rw-r-----  1 root       root        10 Feb 23 10:02 default_compression_codec.txt
-rw-r-----  1 root       root         1 Feb 23 10:02 metadata_version.txt
-rw-r-----  1 root       root         4 Feb 23 10:02 minmax_JoinTime.idx
-rw-r-----  1 root       root         4 Feb 23 10:02 partition.dat
-rw-r-----  1 root       root        50 Feb 23 10:02 primary.cidx # 一级索引文件，使用二进制格式存储，用于存储稀疏索引，一张 MergeTree 表只能声明一次一级索引（通过 ORDER BY 或 PRIMARY KEY）。
-rw-r-----  1 root       root       502 Feb 23 10:02 serialization.json

# columns.txt：列信息文件，使用明文格式存储，用于保存此分区下的列字段信息
root@635e708a264a:/var/lib/clickhouse# cat data/helloworld/user_activity_event/202005_1_1_0/columns.txt
columns format version: 1
7 columns:
`ID` UInt64
`UserName` String
`ActivityName` String
`ActivityType` String
`ActivityLevel` Enum8('Easy' = 0, 'Medium' = 1, 'Hard' = 2)
`IsSuccess` Int8
`JoinTime` Date

# count.txt：计数文件，使用明文格式存储，用于记录当前分区下的数据总数。
root@635e708a264a:/var/lib/clickhouse# cat data/helloworld/user_activity_event/202005_1_1_0/count.txt
1
```

- partition.dat 和 minmax_[Column].idx：如果使用了分区键，例如上面的 PARTITION BY toYYYYMM(JoinTime)，则会额外生成 partition.dat 与 minmax_JoinTime.idx 索引文件，它们均使用二进制格式存储。partition.dat 用于保存当前分区下分区表达式最终生成的值，而 minmax_[Column].idx 则负责记录当前分区下分区字段对应原始数据的最小值和最大值。举个栗子，假设我们往上面的 user_activity_event 表中插入了 5 条数据，JoinTime 分别 2020-05-05、2020-05-15、2020-05-31、2020-05-03、2020-05-24，显然这 5 条都会进入到同一个分区，因为 toYYYMM 之后它们的结果是相同的，都是 2020-05，而 partition.dat 中存储的就是 2020-05，也就是分区表达式最终生成的值；同时还会有一个 minmax_JoinTime.idx 文件，里面存储的就是 2020-05-03 2020-05-31，也就是分区字段对应的原始数据的最小值和最大值



## 参考
- [MergeTree 的深度原理解析](https://www.cnblogs.com/traditional/p/15218743.html)
- [官方文档:表引擎](https://clickhouse.com/docs/zh/engines/table-engines)
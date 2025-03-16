<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [curd](#curd)
  - [create](#create)
  - [insert](#insert)
  - [update](#update)
  - [delete](#delete)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# curd

## create
```clickhouse
CREATE DATABASE [IF NOT EXISTS] db_name [ON CLUSTER cluster] [ENGINE = engine(...)] [COMMENT 'Comment']

-- 创建本地表
CREATE TABLE [IF NOT EXISTS] [db.]table_name ON CLUSTER cluster
(
  name1 [type1] [DEFAULT|MATERIALIZED|ALIAS expr1],
  name2 [type2] [DEFAULT|MATERIALIZED|ALIAS expr2],
  ...
  INDEX index_name1 expr1 TYPE type1(...) GRANULARITY value1,
  INDEX index_name2 expr2 TYPE type2(...) GRANULARITY value2
  ) ENGINE = engine_name()
  [PARTITION BY expr]
  [ORDER BY expr]
  [PRIMARY KEY expr]
  [SAMPLE BY expr]
  [SETTINGS name=value, ...];

```

- MATERIALIZED：物化列表达式，表示该列不能被INSERT，是被计算出来的； 在INSERT语句中，不需要写入该列；在SELECT * 查询语句结果集不包含该列；需要指定列表来查询（虚拟列）
- ALIAS ：别名列。这样的列不会存储在表中。 它的值不能够通过INSERT写入，同时SELECT查询使用星号时，这些列也不会被用来替换星号。 但是它们可以用于SELECT中，在这种情况下，在查询分析中别名将被替换

物化列与别名列的区别： 物化列是会保存数据，查询的时候不需要计算，而别名列不会保存数据，查询的时候需要计算，查询时候返回表达式的计算结果

```clickhouse
-- 创建一个本地表
CREATE TABLE ontime_local ON CLUSTER default -- 表名为 ontime_local
(
    Year UInt16,
    Quarter UInt8,
    Month UInt8,
    DayofMonth UInt8,
    DayOfWeek UInt8,
    FlightDate Date,
    FlightNum String,
    Div5WheelsOff String,
    Div5TailNum String
)ENGINE = ReplicatedMergeTree(--表引擎用ReplicatedMergeTree，开启数据副本的合并树表引擎）
    '/clickhouse/tables/ontime_local/{shard}', -- 指定存储路径
    '{replica}')           
 PARTITION BY toYYYYMM(FlightDate)  -- 指定分区键，按FlightDate日期转年+月维度，每月做一个分区
 PRIMARY KEY (intHash32(FlightDate)) -- 指定主键，FlightDate日期转hash值
 ORDER BY (intHash32(FlightDate),FlightNum) -- 指定排序键，包含两列：FlightDate日期转hash值、FlightNunm字符串。
 SAMPLE BY intHash32(FlightDate)  -- 抽样表达式，采用FlightDate日期转hash值
SETTINGS index_granularity= 8192 ;  -- 指定index_granularity指数，每个分区再次划分的数量

```

```clickhouse
-- 基于本地表创建一个分布式表
CREATE TABLE  [db.]table_name  ON CLUSTER default
 AS db.local_table_name
ENGINE = Distributed(<cluster>, <database>, <shard table> [, sharding_key])

```
- sharding_key：分片表达式。可以是一个字段，例如user_id（integer类型），通过对余数值进行取余分片；也可以是一个表达式，例如rand()，通过rand()函数返回值/shards总权重分片；为了分片更均匀，可以加上hash函数，如intHash64(user_id)
```clickhouse
CREATE TABLE ontime_distributed ON CLUSTER default   -- 指定分布式表的表名，所在集群
 AS db_name.ontime_local                             -- 指定对应的 本地表的表名
ENGINE = Distributed(default, db_name, ontime_local, rand());  -- 指定表引擎为Distributed（固定）

```


## insert


插入最佳
- 建议一次插入一千行数据,理想是1000~10000行.
- 插入可以重试,因为是幂等的.对于 MergeTree engine ,会自动去重.



## update and delete 
clickhouse 没有直接支持 update and delete
注意：

- 索引列不支持更新、删除
- 分布式表不支持更新、删除
```clickhouse
ALTER TABLE [<database>.]<table> UPDATE <column> = <expression> WHERE <filter_expr>
```




```clickhouse
ALTER TABLE [<database>.]<table> DELETE WHERE <filter_expr>
```

轻量删除
- 默认异步.Set mutations_sync equal to 1 to wait for one replica to process the statement, and set mutations_sync to 2 to wait for all replicas
- 只适合MergeTree
```clickhouse
DELETE FROM [db.]table [ON CLUSTER cluster] [WHERE expr]
```
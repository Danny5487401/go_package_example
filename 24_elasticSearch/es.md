#ElasticSearch
## VS 关系型数据库
![](es/img/es_n_mysql.png)
## es架构
![](es/img/distribution.png)
## es倒排索引原理
![](es/img/inverted_index.png)
## CRUD增删改查
![](es/img/crud.png)

##索引简介
对于日志或指标（metric）类时序性强的ES索引，因为数据量大，并且写入和查询大多都是近期时间内的数据。

    我们可以采用hot-warm-cold架构将索引数据切分成hot/warm/cold的索引。hot索引负责最新数据的读写，可使用内存存储；
    warm索引负责较旧数据的读取，可使用内存或SSD存储；cold索引很少被读取，可使用大容量磁盘存储
ES从6.7版本推出了索引生命周期管理（Index Lifecycle Management ，简称ILM)机制，能帮我们自动管理一个索引策略（Policy）下索引集群的生命周期。索引策略将一个索引的生命周期定义为四个阶段：

    * Hot：索引可写入，也可查询。
    * Warm：索引不可写入，但可查询。
    * Cold：索引不可写入，但很少被查询，查询的慢点也可接受。
    * Delete：索引可被安全的删除

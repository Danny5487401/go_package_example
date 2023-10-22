<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [关系型数据库和非关系型数据库的应用场景对比](#%E5%85%B3%E7%B3%BB%E5%9E%8B%E6%95%B0%E6%8D%AE%E5%BA%93%E5%92%8C%E9%9D%9E%E5%85%B3%E7%B3%BB%E5%9E%8B%E6%95%B0%E6%8D%AE%E5%BA%93%E7%9A%84%E5%BA%94%E7%94%A8%E5%9C%BA%E6%99%AF%E5%AF%B9%E6%AF%94)
  - [Mysql和MongoDB内存结构](#mysql%E5%92%8Cmongodb%E5%86%85%E5%AD%98%E7%BB%93%E6%9E%84)
    - [1、InnoDb内存使用机制](#1innodb%E5%86%85%E5%AD%98%E4%BD%BF%E7%94%A8%E6%9C%BA%E5%88%B6)
    - [2、MongoDB的存储结构及数据模型](#2mongodb%E7%9A%84%E5%AD%98%E5%82%A8%E7%BB%93%E6%9E%84%E5%8F%8A%E6%95%B0%E6%8D%AE%E6%A8%A1%E5%9E%8B)
      - [1）MongoDB使用的储存引擎是WiredTiger，WiredTiger的结构如图所示](#1mongodb%E4%BD%BF%E7%94%A8%E7%9A%84%E5%82%A8%E5%AD%98%E5%BC%95%E6%93%8E%E6%98%AFwiredtigerwiredtiger%E7%9A%84%E7%BB%93%E6%9E%84%E5%A6%82%E5%9B%BE%E6%89%80%E7%A4%BA)
      - [2）数据模型](#2%E6%95%B0%E6%8D%AE%E6%A8%A1%E5%9E%8B)
        - [内嵌](#%E5%86%85%E5%B5%8C)
        - [引用模型又称规格化模型（Normalized data models)](#%E5%BC%95%E7%94%A8%E6%A8%A1%E5%9E%8B%E5%8F%88%E7%A7%B0%E8%A7%84%E6%A0%BC%E5%8C%96%E6%A8%A1%E5%9E%8Bnormalized-data-models)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 关系型数据库和非关系型数据库的应用场景对比

关系型数据库适合存储结构化数据，如用户的帐号、地址：

    1）这些数据通常需要做结构化查询，比如join，这时候，关系型数据库就要胜出一筹
    2）这些数据的规模、增长的速度通常是可以预期的
    3）事务性、一致性

NoSQL适合存储非结构化数据，如文章、评论：

    1）这些数据通常用于模糊处理，如全文搜索、机器学习
    2）这些数据是海量的，而且增长的速度是难以预期的，
    3）根据数据的特点，NoSQL数据库通常具有无限（至少接近）伸缩性
    4）按key获取数据效率很高，但是对join或其他结构化查询的支持就比较差

优点：

    1）社区活跃，用户较多，应用广泛。
    2）MongoDB在内存充足的情况下数据都放入内存且有完整的索引支持，查询效率较高。
    3）MongoDB的分片机制，支持海量数据的存储和扩展。
缺点：

    1）不支持事务，在MongoDB 4.0版本之前是没有事务管理，但是今天年初MongoDB发布了4.0版本，其中最重大的标志就是引入了事务管理。
    2）不支持join、复杂查询


## Mysql和MongoDB内存结构

### 1、InnoDb内存使用机制
![](09_Nosql/01_mongo/.mongo_images/innodb_model.png)

Innodb关于查询效率有影响的两个比较重要的参数分别是innodb_buffer_pool_size，innodb_read_ahead_threshold。

innodb_buffer_pool_size指的是Innodb缓冲池的大小，本例中Innodb缓冲池大小为20G，该参数的大小可通过命令指定innodb_buffer_pool_size 20G。
缓冲池使用改进的LRU算法进行管理，维护一个LRU列表、一个FREE列表，FREE列表存放空闲页，数据库启动时LRU列表是空的，
当需要从缓冲池分页时，首先从FREE列表查找空闲页，有则放入LRU列表，否则LRU执行淘汰，淘汰尾部的页分配给新页。

innodb_read_ahead_threshold相对应的是数据预加载机制，innodb_read_ahead_threshold 30表示的是如果一个extent中的被顺序读取的page超过或者等于该参数变量的，
Innodb将会异步的将下一个extent读取到buffer pool中，比如该参数的值为30，那么当该extent中有30个pages被sequentially的读取，则会触发innodb linear预读，将下一个extent读到内存中；
在没有该变量之前，当访问到extent的最后一个page的时候，Innodb会决定是否将下一个extent放入到buffer pool中；可以在Mysql服务端通过show innodb status中的Pages read ahead和evicted without access两个值来观察预读的情况：

Innodb_buffer_pool_read_ahead：表示通过预读请求到buffer pool的pages；
Innodb_buffer_pool_read_ahead_evicted：表示由于请求到buffer pool中没有被访问，而驱逐出内存的页数。

可以看出来，Mysql的缓冲池机制是能充分利用内存且有预加载机制，在某些条件下目标数据完全在内存中，也能够具备非常好的查询性能

### 2、MongoDB的存储结构及数据模型
#### 1）MongoDB使用的储存引擎是WiredTiger，WiredTiger的结构如图所示
![](.mongo_images/wiredTiger.png)
![](.mongo_images/wireTiger_cache.png)

Wiredtiger的Cache采用Btree的方式组织，每个Btree节点为一个page，root page是btree的根节点，internal page是btree的中间索引节点，leaf page是真正存储数据的叶子节点；btree的数据以page为单位按需从磁盘加载或写入磁盘。
可以通过在配置文件中指定storage.wiredTiger.engineConfig.cacheSizeGB参数设定引擎使用的内存量。此内存用于缓存工作集数据（索引、namespace，未提交的write，query缓冲等）。

#### 2）数据模型
##### 内嵌
![](.mongo_images/embeded.png)
内嵌类型支持一组相关的数据存储在一个文档中，这样的好处就是，应用程序可以通过比较少的的查询和更新操作来完成一些常规的数据的查询和更新工作。
当遇到以下情况的时候，我们应该考虑使用内嵌类型：

    如果数据关系是一种一对一的包含关系，例如下面的文档，每个人都有一个contact字段来描述这个人的联系方式。像这种一对一的关系，使用内嵌类型可以很方便的进行数据的查询和更新。
```json
{
    "_id": 1,
    "name": "Wilber",
    "contact": {
    "phone": "12345678",
    "email": "wilber@shanghai.com"
   }
}
```

如果数据的关系是一对多，那么也可以考虑使用内嵌模型。例如下面的文档，用posts字段记录所有用户发布的博客。
在这种情况中，如果应用程序会经常通过用户名字段来查询改用户发布的博客信息。那么，把posts作为内嵌字段会是一个比较好的选择，这样就可以减少很多查询的操作.

```json
{
"_id":1,
"name": "Wilber",
"contact": {
    "phone": "12345678",
    "email": "wilber@shanghai.com"
},
"posts": [
    {
        "title": "Indexes in MongoDB",
        "created": "12/01/2014",
        "link": "www.linuxidc.com"
    },
    {
        "title": "Replication in MongoDB",
        "created": "12/02/2014",
        "link": "www.linuxidc.com"
    },
    {
        "title": "Sharding in MongoDB",
        "created": "12/03/2014",
        "link": "www.linuxidc.com"
    }
]
}
```

内嵌模型可以给应用程序提供很好的数据查询性能，因为基于内嵌模型，可以通过一次数据库操作得到所有相关的数据。同时，内嵌模型可以使数据更新操作变成一个原子写操作。然而，内嵌模型也可能引入一些问题，比如说文档会越来越大，这样就可能会影响数据库写操作的性能，还可能会产生数据碎片（data fragmentation）

##### 引用模型又称规格化模型（Normalized data models)

当我们遇到以下情况的时候，就可以考虑使用引用模型了：
![](.mongo_images/refer_model.png)

使用内嵌模型往往会带来数据的冗余，却可以提升数据查询的效率。但是，当应用程序基本上不通过内嵌模型查询，或者说查询效率的提升不足以弥补数据冗余带来的问题时，我们就应该考虑引用模型了。
当需要实现复杂的多对多关系的时候，可以考虑引用模型。比如我们熟知的例子，学生-课程-老师关系，如果用引用模型来实现三者的关系，可能会比内嵌模型更清晰直观，同时会减少很多冗余数据。
当需要实现复杂的树形关系的时候，可以考虑引用模型
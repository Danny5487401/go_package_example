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
```



## insert


插入最佳⌚️
- 建议一次插入一千行数据,理想是1000~10000行.
- 插入可以重试,因为是幂等的.对于 MergeTree engine ,会自动去重.



## update
```clickhouse
ALTER TABLE [<database>.]<table> UPDATE <column> = <expression> WHERE <filter_expr>
```


## delete

```clickhouse
ALTER TABLE [<database>.]<table> DELETE WHERE <filter_expr>
```

轻量删除
- 默认异步.Set mutations_sync equal to 1 to wait for one replica to process the statement, and set mutations_sync to 2 to wait for all replicas
- 只适合MergeTree
```clickhouse
DELETE FROM [db.]table [ON CLUSTER cluster] [WHERE expr]
```
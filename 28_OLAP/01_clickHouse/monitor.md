# 监控


## 内置监控 
https://clickhouse.com/docs/zh/operations/monitoring ,通过 $HOST:$PORT/dashboard 访问

显示以下指标：

- 每秒查询数
- CPU 使用率（核数）
- 正在运行的查询
- 正在进行的合并
- 每秒选定字节数
- IO 等待
- CPU 等待
- 操作系统 CPU 使用率（用户空间）
- 操作系统 CPU 使用率（内核）
- 从磁盘读取
- 从文件系统读取
- 内存（跟踪）
- 每秒插入行数
- 总 MergeTree 部件
- 每个分区的最大部件数
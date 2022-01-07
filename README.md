# go_grpc_example
![grpc](./img/golang.jpeg)
- [MakeFile介绍](makefile.md)
## 第零章 rpc实现选项
- 1 手动实现rpc
- 2 手动实现stub
- 3 json_rpc
- 4 http_rpc
## [第一章 服务注册中心consul](01_consul/consul.md)
- [consul架构](01_consul/consul.md)
- [Raft协议](01_consul/raft.md)
- [raft在consul实现](01_consul/raft_in_consul.md)
- 1 服务注册，过滤，获取
- 2 [分布式锁(consul实现方式)](01_consul/distributed_lock.md)
## 第二章 日志库
- 1 [zerolog](02_log/01_zerolog/zerolog.md)
- 2 [zap使用及源码分析](02_log/02_zap/zap.md)
  - 2.1 控制台输出
  - 2.2 文件输出
  - 2.3 并发安全logger
  - 2.4 配合日志归档库lumberjack库实现定制化log

## 第三章 消息队列
- 1 [rabbitmq](03_amqp/01_rabbitmq/introduction.md)
  - 1.1 消费者：推拉模式
  - 1.1 生产者
- 2 [kafka](03_amqp/02_kafka/kafka_intro.md)
  - 2.1 客户端sarama
  - 2.2 [客户端confluent-kafka-go源码分析](03_amqp/02_kafka/02_confluent-kafka/confluent_kafka_source_code.md)
    - 2.2.1 生产者
    - 2.2.2 消费者
- 3 rocketmq
  - 3.1 消费者：简单消费,延迟消费
  - 3.2 生产者：简单消息，延迟消息，事务消息
## 第四章 [服务注册及配置文件中心Nacos](04_nacos/nacos.md)
- 1 获取配置及监听文件变化
- 2 服务注册，监听，获取
## 第五章 关系型数据库
- [go-mysql-driver插件源码分析](05_rds/go_mysql_driver.md)
- 1 GORM
  - 1.1 GORM原理及实现 
  - 1.2 连接池使用
- 2 XORM
  - 2.1 主从连接
  - 2.2 调用mysql函数
  - 2.3 事务处理
  - 2.4 crud
## 第六章 获取对外可用IP和端口
## 第七章 Gin前端form验证器
- 1 错误英转中
- 2 前端数据校验
## [第八章 GRPC编程 ](08_grpc/grpc.md)
- [protobuf及工具介绍](08_grpc/proto/protobuf_n_tools.md)
  - 引入其他proto文件
  - 编码原理
  - protoc,protoc-gen-go,protoc-gen-go-grpc,protoc-gen-gofast等工具
  
- 1  HelloWorld入门
  - 1.1 [客户端Grpc源码](08_grpc/01_grpc_helloworld/client/client.md)
  - 1.2 [服务端Grpc源码](08_grpc/01_grpc_helloworld/server/server.md)
- 2  [context中的元数据metadata](08_grpc/02_metadata/grpc_context.md)
- 3  流式GRPC
- 4  protobuf的jsonpb包序列化和反序列化
- 5  负载均衡
  - 5.1 [客户端(Resolver接口和Builder接口)](08_grpc/05_grpc_load_balance/client/builder_n_resolver.md)
    - 第三方consul实现Resolver接口和Builder接口
    - 自定义实现Resolver接口和Builder接口
  - 5.2 服务端
- 6  拦截器 
- [7  grpc错误抛出与捕获](08_grpc/07_grpc_error/error.md)
- 8  auth认证 
- 9  Grpc插件-proto字段验证器 
- 10 Grpc插件-grpc网关直接对外http服务 
- 11 Grpc插件-gogoprotobuf
- 12 GRPC生态中间件(拦截器扩展)
- 13 channelz调试
- 14 multiplex多路复用
- 15 自定义grpc插件
## 第九章 Nosql非关系型数据库
- 1 MongoDB
  - [mongo和mysql储存引擎及内存结构](09_Nosql/01_mongo/mongo.md)
  - 1.1 增删改查
- 2 [Redis(协议，原理，数据结构分析)](09_Nosql/02_redis/redis.md)
  - [redis底层数据结构对象源码分析](09_Nosql/02_redis/redis_obj.md)
  - 2.1 redigo使用
  - 2.2 go-redis使用
    - 2.2.1 连接池分析
    - 2.2.2 连接初始化及命令执行流程
    - 2.2.3 protocol协议封装
    - 2.2.4 批处理pipeline分析
## 第十章 链路追踪(Distributed Tracing)
- 1 Jaeger
  - 1.1 结合XORM
  - 1.2 结合redis
## 第十一章 依赖注入
- 1 dig依赖注入及http服务分层 
- 2 wire依赖注入
## 第十二章 测试框架testify(gin使用)
- 1 assert断言
- 2 mock测试替身
- 3 suite测试套件
## [第十三章 序列化反序列化](13_serialize/serialize.md)
- 1 Jsoniter(完全兼容标准库json，性能较好)
  - 1.1 序列化
    - 结构体成员为基本类型,嵌套结构体，及tag标签使用
    - 结构体成员为interface{}
  - 1.2 反序列化
    - 基本使用
    - json字符串数组
    - json.RawMessage二次反序列化
- 2 mapstructure使用（性能低但是方便）
  - 无tag标签
  - 带tag标签
  - embeded内嵌
  - 字段保留
  - 省略字段
  - 元数据
  - 错误
  - 弱解析
  - 自定义解析器
## 第十四章 系统监控
  - 1 systemstat包(适合linux系统，已断更)
  - 2 [gopsutil](14_system_monitor/02_gopsutil/gopsutil.md)
    - 进程信息获取
      - 物理机和虚拟机
      - 容器环境
    - cpu,mem,disk
## [第十五章 分布式事务](15_distributed_transaction/distributed_transaction.md)
- 1 两阶段提交2pc
## 第十六章 数据复制
- 1 [copier(不同类型数据复制)](16_dataCopy/copier/copier.md)
## 第十七章 数据加解密
- 1 phpserialize
## 第十八章 日志收集项目 log_collect
- 1 动态选择文件
- 2 文件内容读取发送
## [第十九章 熔断,限流及降级](19_fuse_currentLimiting_degradation/rate_limit.md)
- 0 熔断，降级，限流(官方包实现)
- 1 Sentinel
  - 1 流量控制
  - 2 熔断
- [2 Hystrix](19_fuse_currentLimiting_degradation/hystrix.md)
## 第二十章 [命令行框架Cobra](20_cobra/introdoction.md)
- 1 介绍及功能使用
- 2 [在k8s中的应用](20_cobra/cobra_in_k8s.md)
## [第二十一章 配置文件获取工具viper(依赖mapstructure,fsnotify)](21_viper/viper.md)
- 1 获取本地文件内容
- [2 监听文件变化(fsnotify)](21_viper/02_fsnotify/fsnotify.md)
- [3 远程读取nacos配置(源码分析)](21_viper/03_remote_config/remote_viper_config.md)
## 第二十二章 ETCD
- 1 CRUD及watch
- 2 [读和写流程分析](22_etcd/etcd_read_n_write.md)

## 第二十三章 Go-Micro框架 (不推荐使用)
- 1 [Config配置加载包](23_micro/01ConfigTest/config.md)

## [ 第二十四章 搜索引擎es](24_elasticSearch/es.md)
- [es索引及索引生命周期管理](24_elasticSearch/es_index.md)
- 1 官方包
- 2 第三方包oliver
  - 2.1 V6版本 
  - 2.2 V7版本 

## 第二十五章 监控sentry
- 结合gin使用

## [第二十六章 图数据库Neo4j](26_neo4j/neo4j.md)
- [cypher语句](26_neo4j/cypher.md)
-  1 CRUD在web服务中

## 第二十七章 Mysql的binlog
- [binlog](27_mysql_binlog/binlog.md)
- canal

## 第二十八章 OLAP(Online Analytical Processing联机分析处理)
- 1 [列数据库ClickHouse](28_OLAP/01_clickHouse/clickHouse.md)
  - 1.1 标准库sql操作clickHouse
  - 1.2 扩展包sqlx操作clickHouse

## 第二十九章 分布式锁
- 1 redis实现



# go_grpc_example
![grpc](./img/golang.jpeg)
## 第零章 rpc实现选项
- 1 手动实现rpc
- 2 手动实现stub
- 3 json_rpc
- 4 http_rpc
## [第一章 服务注册中心consul](01_consul/consul.md)
- 1 服务注册，过滤，获取
## [第二章 日志库zap](02_zap/zap.md)
- 1 源码结构
- 2 控制台输出
- 3 文件输出
## 第三章 消息队列
- 1 [rabbitmq](03_amqp/01_rabbitmq/introduction.md)
  - 1.1 消费者：推拉模式
  - 1.1 生产者
- 2 [kafka](03_amqp/02_kafka/introduction.md)
  - 2.1 客户端sarama
  - 2.2 客户端confluent-kafka-go
- 3 rocketmq
  - 3.1 消费者：简单消费,延迟消费
  - 3.2 生产者：简单消息，延迟消息，事务消息
## 第四章 [服务注册及配置文件中心Nacos](04_nacos/nacos.md)
- 1 获取配置及监听文件变化
- 2 服务注册
## 第五章 数据操作
- 1 GORM
  - 1.1 GORM原理及实现 
  - 1.2 连接池使用
- 2 XORM
  - 2.1 主从连接
- 3 MongoDB
  - mongo和mysql储存引擎及内存结构
  - 3.1 增删改查
- 4 Redis(协议，原理，数据结构分析)
  - 4.1 redigo使用
  - 4.2 go-redis使用
    - 4.2.1 连接池分析
    - 4.2.2 连接初始化及命令执行流程
    - 4.2.3 protocol协议封装
    - 4.2.4 批处理pipeline分析
## 第六章 获取对外可用IP和端口
## 第七章 Gin前端form验证器
- 1 错误英转中
- 2 前端数据校验
## 第八章 GRPC编程 
- 1  HelloWorld入门
  - 1.1 [客户端Grpc源码](08_grpc/01_grpc_helloworld/client/client.md)
  - 1.2 [服务端Grpc源码](08_grpc/01_grpc_helloworld/server/server.md)
- 2  元数据metada
- 3  流式GRPC
- 4  protobuf的jsonpb包序列化和反序列化
- 5  负载均衡 
- 6  拦截器 
- 7  grpc错误抛出与捕获 
- 8  auth认证 
- 9  proto字段验证器 
- 10 grpc网关-直接对外http服务 
- 11 Grpc插件-gogoprotobuf
- 12 GRPC生态中间件
- 13 channelz调试
- 14 multiplex多路复用
## 第九章 DDD领域驱动
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
## 第十三章 序列化反序列化
- 1 Jsoniter(完全兼容标准库json，性能较好)
  - 1.1 序列化
  - 1.2 反序列化
- 2 mapstructure使用（性能低但是方便）
## 第十四章 系统监控指标
## 第十五章 分布式事务
- 1 两阶段提交2pc
## 第十六章 数据复制
- 1 copier(不同类型数据复制)
## 第十七章 数据加解密
- 1 phpserialize
## 第十八章 日志收集项目 log_collect
- 1 动态选择文件
- 2 文件内容读取发送
## 第十九章 熔断和限流Sentinel
- 0 熔断，降级，限流的方法及官方包实现
- 1 流量控制
- 2 熔断
## 第二十章 [命令行框架Cobra](20_cobra/introdoction.md)
- 1 介绍及功能使用
- 2 在k8s中的应用
## 第二十一章 配置文件获取工具viper
- 1 获取文件内容
- 2 监听文件变化(fsNotify)
## 第二十二章 ETCD
- 1 CRUD及watch
- 2 [读和写流程分析](22_etcd/etcd_read_n_write.md)
- 3 [Raft协议](22_etcd/raft.md)
## 第二十三章 Go-Micro框架 (不推荐使用)
- 1 [Config配置加载包](23_micro/01ConfigTest/config.md)

## [ 第二十四章 搜索引擎es](24_elasticSearch/es.md)

- 1 V6版本 
- 2 V7版本 

## 第二十五章 监控sentry
- 结合gin使用




# go_grpc_example
![grpc](./img/golang.jpeg)
## 第零章 rpc实现选项
    1. 手动实现rpc
    2. 手动实现stub
    3. json_rpc
    3. http_rpc
## 第一章 服务注册中心consul
    1. 注册，过滤，获取
## 第二章 日志库zap
    a. 源码结构
    b. 控制台输出
    c. 文件输出
## 第三章 消息队列
    1. rabbitmq
        a. 消费者：推拉模式
        b. 生产者
    2. kafka
        2.1 客户端sarama
        2.2 客户端confluent-kafka-go
    3. rocketmq
        a. 消费者：简单消费,延迟消费
        b. 生产者：简单消息，延迟消息，事务消息
## 第四章 配置文件中心nacos
    1. 获取配置及监听文件变化
## 第五章 数据操作
    5.1 GORM
        a. GORM原理及实现 
        b. 连接池使用
    5.2 XORM
        a. 主从连接
    5.3 MongoDB
    5.4 Redis
        a. redigo使用
        b. go-redis使用
## 第六章 获取对外可用IP和端口
## 第七章 Gin前端form验证器
    1. 错误英转中
    2. 前端数据校验
## 第八章 GRPC编程 
    8.1 Grpc基本使用HelloWorld
    8.2 元数据metada
    8.3 流式GRPC
    8.4 protobuf的jsonpb包序列化和反序列化
    8.5 负载均衡 
    8.6 拦截器 
    8.7 grpc错误抛出与捕获 
    8.8 auth认证 
    8.9 proto字段验证器 
    8.9 grpc网关-直接对外http服务 
    8.9 Grpc插件-gogoprotobuf
## 第九章 错误及异常
## 第十章 链路追踪Jaeger
    10.1 结合XORM
    10.2 结合redis
## 第十一章 依赖注入
    11.1 dig依赖注入及http服务分层
    11.2 wire依赖注入
## 第十一章 GRPC生态中间件
## 第十二章 测试框架testify
## 第十三章 map和structure相互转化
    13.1 mapstructure使用（性能低但是方便）
## 第十九章 熔断和限流Sentinel
    0. 熔断，降级，限流的方法及官方包实现
    1. 流量控制
    2. 熔断
## 第十五章 日志收集项目 log_collect
    1. es操作
    2. etcd操作
    3. kafka操作
## 第二十章 命令行框架Cobra
    1. 介绍及功能使用
## 第二十一章 配置文件获取工具viper
    1. 获取文件内容
    2. 监听文件变化
## 第二十二章 ETCD
    1. CRUD及watch
## 第二十三章 Go-Micro框架 





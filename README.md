# go_grpc_example
![grpc](./img/golang.jpeg)
## 第零章 rpc实现选项
    1. 手动实现rpc
    2. 手动实现stub
    3. json_rpc
    4. http_rpc
## 第一章 服务注册中心consul
    1. 注册，过滤，获取
## 第二章 日志库zap
    1. 源码结构
    2. 控制台输出
    3. 文件输出
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
## 第四章 配置文件中心Nacos
    4.1. 获取配置及监听文件变化
## 第五章 数据操作
    1 GORM
        1.1 GORM原理及实现 
        1.2 连接池使用
    2 XORM
        2.1 主从连接
    3 MongoDB
        mongo和mysql储存引擎及内存结构
    4 Redis
        4.1 redigo使用
        4.2 go-redis使用
## 第六章 获取对外可用IP和端口
## 第七章 Gin前端form验证器
    1. 错误英转中
    2. 前端数据校验
## 第八章 GRPC编程 
    1 HelloWorld分析Grpc源码
    2 元数据metada
    3 流式GRPC
    4 protobuf的jsonpb包序列化和反序列化
    5 负载均衡 
    6 拦截器 
    7 grpc错误抛出与捕获 
    8 auth认证 
    9 proto字段验证器 
    10 grpc网关-直接对外http服务 
    11 Grpc插件-gogoprotobuf
## 第九章 错误及异常
## 第十章 链路追踪(Distributed Tracing)
    1.Jaeger
        1.1 结合XORM
        1.2 结合redis
## 第十一章 依赖注入
    1 dig依赖注入及http服务分层
    2 wire依赖注入
## 第十一章 GRPC生态中间件
## 第十二章 测试框架testify
## 第十三章 序列化反序列化
    1 Jsoniter(完全兼容标准库json，性能较好)
        1.1 序列化
        1.2 反序列化
    2 mapstructure使用（性能低但是方便）
## 第十四章 系统监控指标
## 第十五章 分布式事务
    1 两阶段提交2pc
## 第十六章 不同结构数据之间copy
## 第十九章 熔断和限流Sentinel
    0. 熔断，降级，限流的方法及官方包实现
    1. 流量控制
    2. 熔断
## 第十五章 分布式事务
    1. 两阶段提交
## 第十六章 数据复制
    1. copier(不同类型数据复制)
## 第十七章 数据加解密
    1. phpserialize
## 第十八章 日志收集项目 log_collect
    1. 动态选择文件
    2. 文件内容读取发送

## 第二十章 命令行框架Cobra
    1. 介绍及功能使用
## 第二十一章 配置文件获取工具viper
    1. 获取文件内容
    2. 监听文件变化
## 第二十二章 ETCD
    1. CRUD及watch
    2. 读和写流程分析
    3. Raft协议
## 第二十三章 Go-Micro框架 (不推荐使用)
    1. Config配置加载包
## 第二十四章 搜索引擎es
    1. V6版本 
    2. V7版本 





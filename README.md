<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [go_package_example(Go常用包)](#go_package_examplego%E5%B8%B8%E7%94%A8%E5%8C%85)
  - [第零章 rpc实现选项](#%E7%AC%AC%E9%9B%B6%E7%AB%A0-rpc%E5%AE%9E%E7%8E%B0%E9%80%89%E9%A1%B9)
  - [第一章 服务注册中心consul](#%E7%AC%AC%E4%B8%80%E7%AB%A0-%E6%9C%8D%E5%8A%A1%E6%B3%A8%E5%86%8C%E4%B8%AD%E5%BF%83consul)
  - [第二章 日志库](#%E7%AC%AC%E4%BA%8C%E7%AB%A0-%E6%97%A5%E5%BF%97%E5%BA%93)
  - [第三章 消息队列](#%E7%AC%AC%E4%B8%89%E7%AB%A0-%E6%B6%88%E6%81%AF%E9%98%9F%E5%88%97)
  - [第四章 服务注册及配置文件中心 Nacos](#%E7%AC%AC%E5%9B%9B%E7%AB%A0-%E6%9C%8D%E5%8A%A1%E6%B3%A8%E5%86%8C%E5%8F%8A%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E4%B8%AD%E5%BF%83-nacos)
  - [第五章 关系型数据库](#%E7%AC%AC%E4%BA%94%E7%AB%A0-%E5%85%B3%E7%B3%BB%E5%9E%8B%E6%95%B0%E6%8D%AE%E5%BA%93)
  - [第六章 获取对外可用IP和端口](#%E7%AC%AC%E5%85%AD%E7%AB%A0-%E8%8E%B7%E5%8F%96%E5%AF%B9%E5%A4%96%E5%8F%AF%E7%94%A8ip%E5%92%8C%E7%AB%AF%E5%8F%A3)
  - [第七章 验证器 go-playground/validator](#%E7%AC%AC%E4%B8%83%E7%AB%A0-%E9%AA%8C%E8%AF%81%E5%99%A8-go-playgroundvalidator)
  - [第八章 GRPC编程及调优](#%E7%AC%AC%E5%85%AB%E7%AB%A0-grpc%E7%BC%96%E7%A8%8B%E5%8F%8A%E8%B0%83%E4%BC%98)
  - [第九章 Nosql 非关系型数据库](#%E7%AC%AC%E4%B9%9D%E7%AB%A0-nosql-%E9%9D%9E%E5%85%B3%E7%B3%BB%E5%9E%8B%E6%95%B0%E6%8D%AE%E5%BA%93)
  - [第十章 链路追踪(Distributed Tracing)](#%E7%AC%AC%E5%8D%81%E7%AB%A0-%E9%93%BE%E8%B7%AF%E8%BF%BD%E8%B8%AAdistributed-tracing)
  - [第十一章 依赖注入容器(Dependency Injection Container)](#%E7%AC%AC%E5%8D%81%E4%B8%80%E7%AB%A0-%E4%BE%9D%E8%B5%96%E6%B3%A8%E5%85%A5%E5%AE%B9%E5%99%A8dependency-injection-container)
  - [第十二章 clockwork 虚拟时钟库-->etcd使用](#%E7%AC%AC%E5%8D%81%E4%BA%8C%E7%AB%A0-clockwork-%E8%99%9A%E6%8B%9F%E6%97%B6%E9%92%9F%E5%BA%93--etcd%E4%BD%BF%E7%94%A8)
  - [第十三章 序列化反序列化-涉及多种协议](#%E7%AC%AC%E5%8D%81%E4%B8%89%E7%AB%A0-%E5%BA%8F%E5%88%97%E5%8C%96%E5%8F%8D%E5%BA%8F%E5%88%97%E5%8C%96-%E6%B6%89%E5%8F%8A%E5%A4%9A%E7%A7%8D%E5%8D%8F%E8%AE%AE)
  - [第十四章 系统监控](#%E7%AC%AC%E5%8D%81%E5%9B%9B%E7%AB%A0-%E7%B3%BB%E7%BB%9F%E7%9B%91%E6%8E%A7)
  - [第十五章 分布式事务](#%E7%AC%AC%E5%8D%81%E4%BA%94%E7%AB%A0-%E5%88%86%E5%B8%83%E5%BC%8F%E4%BA%8B%E5%8A%A1)
  - [第十六章 copier(不同类型数据复制)](#%E7%AC%AC%E5%8D%81%E5%85%AD%E7%AB%A0-copier%E4%B8%8D%E5%90%8C%E7%B1%BB%E5%9E%8B%E6%95%B0%E6%8D%AE%E5%A4%8D%E5%88%B6)
  - [第十七章 数据加解密](#%E7%AC%AC%E5%8D%81%E4%B8%83%E7%AB%A0-%E6%95%B0%E6%8D%AE%E5%8A%A0%E8%A7%A3%E5%AF%86)
  - [第十八章 日志收集项目 log_collect](#%E7%AC%AC%E5%8D%81%E5%85%AB%E7%AB%A0-%E6%97%A5%E5%BF%97%E6%94%B6%E9%9B%86%E9%A1%B9%E7%9B%AE-log_collect)
  - [第十九章 熔断,限流及降级](#%E7%AC%AC%E5%8D%81%E4%B9%9D%E7%AB%A0-%E7%86%94%E6%96%AD%E9%99%90%E6%B5%81%E5%8F%8A%E9%99%8D%E7%BA%A7)
  - [第二十章 应用的命令行框架](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E7%AB%A0-%E5%BA%94%E7%94%A8%E7%9A%84%E5%91%BD%E4%BB%A4%E8%A1%8C%E6%A1%86%E6%9E%B6)
  - [第二十一章 配置文件解析:viper(依赖mapstructure,fsnotify,yaml,toml)](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E4%B8%80%E7%AB%A0-%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E8%A7%A3%E6%9E%90viper%E4%BE%9D%E8%B5%96mapstructurefsnotifyyamltoml)
  - [第二十二章 ETCD](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E4%BA%8C%E7%AB%A0-etcd)
  - [第二十三章 Go-Micro框架(不推荐)](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E4%B8%89%E7%AB%A0-go-micro%E6%A1%86%E6%9E%B6%E4%B8%8D%E6%8E%A8%E8%8D%90)
  - [第二十四章 搜索引擎es](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E5%9B%9B%E7%AB%A0-%E6%90%9C%E7%B4%A2%E5%BC%95%E6%93%8Ees)
  - [第二十五章 监控sentry](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E4%BA%94%E7%AB%A0-%E7%9B%91%E6%8E%A7sentry)
  - [第二十六章 图数据库Neo4j](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E5%85%AD%E7%AB%A0-%E5%9B%BE%E6%95%B0%E6%8D%AE%E5%BA%93neo4j)
  - [第二十七章 Mysql的binlog](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E4%B8%83%E7%AB%A0-mysql%E7%9A%84binlog)
  - [第二十八章 OLAP(Online Analytical Processing 联机分析处理)](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E5%85%AB%E7%AB%A0-olaponline-analytical-processing-%E8%81%94%E6%9C%BA%E5%88%86%E6%9E%90%E5%A4%84%E7%90%86)
  - [第二十九章 分布式锁及源码分析](#%E7%AC%AC%E4%BA%8C%E5%8D%81%E4%B9%9D%E7%AB%A0-%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81%E5%8F%8A%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [第三十章 Zookeeper](#%E7%AC%AC%E4%B8%89%E5%8D%81%E7%AB%A0-zookeeper)
  - [第三十一章 分布式 Id](#%E7%AC%AC%E4%B8%89%E5%8D%81%E4%B8%80%E7%AB%A0-%E5%88%86%E5%B8%83%E5%BC%8F-id)
  - [第三十二章 Consensus algorithm 共识算法](#%E7%AC%AC%E4%B8%89%E5%8D%81%E4%BA%8C%E7%AB%A0-consensus-algorithm-%E5%85%B1%E8%AF%86%E7%AE%97%E6%B3%95)
  - [第三十三章 压缩](#%E7%AC%AC%E4%B8%89%E5%8D%81%E4%B8%89%E7%AB%A0-%E5%8E%8B%E7%BC%A9)
  - [第三十四章 本地缓存](#%E7%AC%AC%E4%B8%89%E5%8D%81%E5%9B%9B%E7%AB%A0-%E6%9C%AC%E5%9C%B0%E7%BC%93%E5%AD%98)
  - [第三十五章 sonar静态代码质量分析-涉及与golangci-lint对比使用](#%E7%AC%AC%E4%B8%89%E5%8D%81%E4%BA%94%E7%AB%A0-sonar%E9%9D%99%E6%80%81%E4%BB%A3%E7%A0%81%E8%B4%A8%E9%87%8F%E5%88%86%E6%9E%90-%E6%B6%89%E5%8F%8A%E4%B8%8Egolangci-lint%E5%AF%B9%E6%AF%94%E4%BD%BF%E7%94%A8)
  - [第三十六章 Proto管理工具 Buf](#%E7%AC%AC%E4%B8%89%E5%8D%81%E5%85%AD%E7%AB%A0-proto%E7%AE%A1%E7%90%86%E5%B7%A5%E5%85%B7-buf)
  - [第三十七章 CI持续集成](#%E7%AC%AC%E4%B8%89%E5%8D%81%E4%B8%83%E7%AB%A0-ci%E6%8C%81%E7%BB%AD%E9%9B%86%E6%88%90)
  - [第三十八章 Mergo实现 struct 与 map 之间转换-->k8s中应用](#%E7%AC%AC%E4%B8%89%E5%8D%81%E5%85%AB%E7%AB%A0-mergo%E5%AE%9E%E7%8E%B0-struct-%E4%B8%8E-map-%E4%B9%8B%E9%97%B4%E8%BD%AC%E6%8D%A2--k8s%E4%B8%AD%E5%BA%94%E7%94%A8)
  - [第三十九章 权限管理 casbin](#%E7%AC%AC%E4%B8%89%E5%8D%81%E4%B9%9D%E7%AB%A0-%E6%9D%83%E9%99%90%E7%AE%A1%E7%90%86-casbin)
  - [第四十章 规则引擎 rule engine](#%E7%AC%AC%E5%9B%9B%E5%8D%81%E7%AB%A0-%E8%A7%84%E5%88%99%E5%BC%95%E6%93%8E-rule-engine)
  - [第四十一章 hashicorp/go-plugin 插件使用-->httprunner 4.0 使用](#%E7%AC%AC%E5%9B%9B%E5%8D%81%E4%B8%80%E7%AB%A0-hashicorpgo-plugin-%E6%8F%92%E4%BB%B6%E4%BD%BF%E7%94%A8--httprunner-40-%E4%BD%BF%E7%94%A8)
  - [第四十二章 open-api](#%E7%AC%AC%E5%9B%9B%E5%8D%81%E4%BA%8C%E7%AB%A0-open-api)
  - [第四十三章 go-systemd-->k8s 中使用](#%E7%AC%AC%E5%9B%9B%E5%8D%81%E4%B8%89%E7%AB%A0-go-systemd--k8s-%E4%B8%AD%E4%BD%BF%E7%94%A8)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# go_package_example(Go常用包)

![grpc](img/golang.jpeg)

## 第零章 rpc实现选项
- 1 手动实现rpc
  - [客户端](00_rpc_options/00_helloworld_without_stub/client/client.go)
  - [服务端](00_rpc_options/00_helloworld_without_stub/server/server.go)
- 2 手动实现stub中间人
  - [客户端](00_rpc_options/02_new_helloworld_withStub/client/client.go)
  - [客户端中间人stub](00_rpc_options/02_new_helloworld_withStub/client_proxy/client_proxy.go)
  - [业务方法](00_rpc_options/02_new_helloworld_withStub/handler/handler.go)
  - [服务端](00_rpc_options/02_new_helloworld_withStub/server/server.go)
  - [服务端中间人stub](00_rpc_options/02_new_helloworld_withStub/server_proxy/server_proxy.go)
- 3 json_rpc
  - [客户端](00_rpc_options/03_json_rpc_test/client/client.go)
  - [服务端](00_rpc_options/03_json_rpc_test/server/server.go)
- 4 http_rpc
  - [服务端](00_rpc_options/04_http_rpc_test/server/server.go)

## [第一章 服务注册中心consul](01_consul/consul.md)
- [consul架构](01_consul/consul.md)
- [分布式锁-->consul实现)](01_consul/distributed_lock.md)

- [1 http服务注册发现加健康检查](01_consul/01_http/test/consul_registry_test.go)
- [2 grpc服务注册发现加健康检查](01_consul/01_http/test/consul_registry_test.go)

## 第二章 日志库
- [1 标准库 log](02_log/01_log/log.md) 
- [2 slog-->go 1.21 引入](02_log/02_slog/slog.md)
  - [2.1 自定义 level](02_log/02_slog/01_level/main.go)
  - [2.2 使用 handler](02_log/02_slog/02_handler/main.go)
  - [2.3 使用 group 汇总多个属性](02_log/02_slog/03_group/main.go)
- [3 zap使用及源码分析](02_log/03_zap/zap.md)
  - [3.1 两种打印风格](02_log/03_zap/01_cosole/main.go)
  - [3.2 定义多种输出位置: 控制台输出及文件输出](02_log/03_zap/02_file_stdout/main.go)
  - [3.3 并发安全logger](02_log/03_zap/03_concurrency_safe/main.go)
  - [3.4 zap(配合 lumberjack 库按文件大小切割或 go-file-rotatelogs 库按日期切割)实现定制化log日志归档](02_log/03_zap/04_customized_log/main.go)
    - [3.4.1 lumberjack 日志切割](02_log/03_zap/04_customized_log/lumberjack.md)
  - [3.5 简单的基于Entry实现的hook函数-->无法接收到Fields的相关参数](02_log/03_zap/05_hook/main.go)
- [4 logrus-->兼容 log 库](02_log/04_logrus/logrus.md)
- [5 zerolog](02_log/05_zerolog/zerolog.md)
- [6 C++日志库glog的golang版本](02_log/06_glog/glog.md)
  - [6.1 vlevel 日志级别小于或等于 level 的日志打印处理](02_log/06_glog/01_vlevel)
  - [6.2 vmodule: 单独定制模块级别, log_backtrace_at: 打印堆栈](02_log/06_glog/02_vmodule_and_trace_location/main.go)
- [7 klog-->glog 的 fork 版本,应用 k8s](02_log/07_klog/klog.md) 
## [第三章 消息队列](03_amqp/amqp.md)
- [1 rabbitmq](03_amqp/01_rabbitmq/introduction.md)
  - 1.1 消费者：推拉模式
  - 1.1 生产者
- [2 kafka](03_amqp/02_kafka/kafka_intro.md)
  - [2.1 客户端 sarama](03_amqp/02_kafka/01_sarama/sarama.md)
    - 2.1.1 生产者
    - 2.1.2 消费者
  - [2.2 客户端 confluent-kafka-go](03_amqp/02_kafka/02_confluent-kafka/confluent_kafka.md)
    - 2.2.1 生产者
    - 2.2.2 消费者
- [3 rocketmq](03_amqp/03_rocketmq/rocketmq.md)
  - 3.1 消费者：简单消费,延迟消费
  - 3.2 生产者：简单消息，延迟消息，事务消息
- [4 Asynq 分布式延迟队列](03_amqp/04_asynq/asynq.md)
  - [4.1 生产者](03_amqp/04_asynq/producer)
  - [4.2 消费者](03_amqp/04_asynq/server)

## [第四章 服务注册及配置文件中心 Nacos](04_nacos/nacos.md)
- 1 [获取配置及监听文件变化](04_nacos/config_center/main.go)
- 2 服务注册，监听，获取
  - [V1版本](04_nacos/service_center/v1/main.go)
  - [V2版本](04_nacos/service_center/v2/main.go)

## 第五章 关系型数据库
- [MySQL的主流驱动 go-mysql-driver插件源码分析](05_rds/go_mysql_driver.md)
- 1 GORM
  - 1.1 GORM原理及实现
  - 1.2 连接池使用
- [2 XORM](05_rds/02_xorm/xorm.md)
  - [2.1 主从连接](05_rds/02_xorm/util/util.go)
  - [2.2 调用mysql函数](05_rds/02_xorm/function/sum.go)
  - [2.3 事务处理](05_rds/02_xorm/transaction/transaction.go)
  - 2.4 crud
    - 插入Insert
    - 原生 sql
    - 获取 retrieve
    - 更新 update
- 3 database/sql 源码分析
  - [3.1 converter 把普通的值转化成 driver.Value的接口](05_rds/03_database_sql/01_converter/converter.md)
- [4 SQL生成库 Squirrel](05_rds/04_squirrel/quirrel.md)

## 第六章 获取对外可用IP和端口
- [1 通过google, 国内移动、电信和联通通用的DNS获取对外 Ip 和 port ](06_get_available_ip_port/get_ip/outboundIp_test.go)
- [2 获取本地可用端口](06_get_available_ip_port/get_port/main.go)

## [第七章 验证器 go-playground/validator](07_gin_form_validator/validator.md)
- [1 dive 递归结构体字段验证](07_gin_form_validator/01_dive_validate/main.go)
- [2 前端数据校验](07_gin_form_validator/02_gin_form/main.go)
- [3 错误处理:验证器校验错误英转中](07_gin_form_validator/03_translate_err_from_en_to_ch/main.go)

## [第八章 GRPC编程及调优](08_grpc/grpc.md)
*前置知识*
- [makefile 在 protobuf 中应用,生成 Pb 文件](08_grpc/makefile)
- [protobuf](08_grpc/protobuf.md)
  - 引入其他proto文件,支持编译多个proto文件
  - 编码原理
- [protobuf 工具:protoc,protoc-gen-go,protoc-gen-go-grpc,protoc-gen-gofast 等](08_grpc/proto_tools.md)

- 1  HelloWorld 入门使用及源码分析
  - [1.1 客户端 Grpc 源码](08_grpc/01_grpc_helloworld/client/client.md)
  - [1.2 服务端 Grpc 源码](08_grpc/01_grpc_helloworld/server/server.md)
- [2  context 中的元数据 metadata](08_grpc/02_metadata/grpc_context.md)
- 3  流式GRPC
- [4  jsonpb 包序列化和反序列化: protobuf 转 json](08_grpc/04_jsonpb/jsonpb.md)
- [5  负载均衡](08_grpc/05_grpc_load_balance/load_balance.md)
  - [5.1 客户端负载均衡(Resolver接口和Builder接口)](08_grpc/05_grpc_load_balance/client/builder_n_resolver_n_balancer.md)
    - [第三方consul实现Resolver接口和Builder接口](01_consul/02_grpc/consul_client/main.go)
    - [自定义实现Resolver接口和Builder接口](08_grpc/05_grpc_load_balance/client/customized_resolver_client/client.go)
    - [自定义实现nacos服务注册与发现](08_grpc/05_grpc_load_balance/client/nacos_client)
  - 5.2 服务端
- [6  retry机制](08_grpc/06_grpc_retry/retry.md)
- [7  grpc错误抛出与捕获](08_grpc/07_grpc_error/error.md)
- [8  auth自定义认证](08_grpc/08_grpc_token_auth/credentials.md)
- [9  Grpc插件-引入第三方proto实现字段验证器-->](08_grpc/09_grpc_validate/proto/helloworld.proto)
- [10 Grpc插件-使用 gRPC 转码（RPC Transcoding）实现暴露 http服务-->grpc网关在etcd中应用](08_grpc/10_grpc_gateway/grpc_gateway.md)
  -[buf 在 protobuf 中应用,生成 Pb 文件-->推荐](08_grpc/10_grpc_gateway/buf.work.yaml)
  -[原始不使用 buf 生成pb文件](08_grpc/10_grpc_gateway/proto_without_buf)
- [11 Grpc插件-gogo/protobuf](08_grpc/11_protoc_gogofast/gogoprotobuf.md)
- [12 GRPC生态中间件(拦截器扩展)](08_grpc/12_grpc_middleware/01_grpc_interceptor/server/server.go)
  - 实现基于 CA 的 TLS 证书认证
  - go-grpc-middleware 实现多个中间件：异常保护，日志
- [13 channelz 分析问题](08_grpc/13_channelz/channelz.md)
- 14 multiplex多路复用
- [15 自定义 grpc 插件](08_grpc/15_customized_protobuf_plugin/protobuf_extend.md)
- [16 同目录 proto 文件引入](08_grpc/16_import_proto/proto)
- [17 field masks](08_grpc/17_fieldmask/fieldmask.md)

## 第九章 Nosql 非关系型数据库
- [1 MongoDB](09_Nosql/01_mongo/mongo.md)
  - [mongo和mysql对比：储存引擎及内存结构](09_Nosql/01_mongo/nosql_vs_rds.md)
  - 1.1 增删改查
- [2 Redis(协议，原理，持久化方式)](09_Nosql/02_redis/redis.md)
  - [redis底层数据结构对象源码分析](09_Nosql/02_redis/redis_obj.md)
  - [redis 集群](09_Nosql/02_redis/redis_cluster.md)
  - 2.1 redigo使用
  - 2.2 go-redis使用(官方)
    - [2.2.1 连接池分析](09_Nosql/02_redis/02_go-redis/go-redis_pool.md)
    - [2.2.2 连接初始化及命令执行流程](09_Nosql/02_redis/02_go-redis/go-redis_init_n_excute.md)
    - [2.2.3 protocol协议封装](09_Nosql/02_redis/02_go-redis/go-redis_protocol.md)
    - [2.2.4 批处理pipeline分析](09_Nosql/02_redis/02_go-redis/go-redis_pipeline.md)

## [第十章 链路追踪(Distributed Tracing)](10_distributed_tracing/introduction.md)
- [1 OpenTracing->Jaeger](10_distributed_tracing/01_jaeger/jaeger.md)
  - [1.1 结合XORM](10_distributed_tracing/01_jaeger/01_jaeger_xorm/main_test.go)
  - [1.2 结合redis](10_distributed_tracing/01_jaeger/02_jaeger_redis/hook.go)
- [2 OpenTelemetry 两大开源社区合并](10_distributed_tracing/02_openTelemetry/openTelemetry.md)
  - 跨服务组合tracer代码展示:需开启svc1和svc2两个http服务(url可以是zipkin或则jaeger)

## [第十一章 依赖注入容器(Dependency Injection Container)](11_dependency_injection/dependency_injection.md)
- [1 dig依赖注入及http服务分层->不推荐](11_dependency_injection/00_dig/dig.go)
- 2 wire依赖注入->推荐
  - [2.1 不使用wire现状](11_dependency_injection/01_wire/01_without_wire/main.go)
  - [2.2 使用wire优化](11_dependency_injection/01_wire/02_wire)
  - [2.3 wire使用-带err返回](11_dependency_injection/01_wire/03_wire_return_err/wire)
  - [2.4 wire使用-带参数初始化](11_dependency_injection/01_wire/04_wire_pass_params/wire)


## [第十二章 clockwork 虚拟时钟库-->etcd使用](12_clockwork/clockwork.md)

## [第十三章 序列化反序列化-涉及多种协议](13_serialize/serialize.md)
- [1 标准库 json](13_serialize/01_std_json/json.md)
  - [1.1 omitempty 标签: 忽略空值字段,忽略嵌套结构体空值字段,不修改原结构体忽略空值字段](13_serialize/01_std_json/01_omitempty/main.go)
  - [1.2 string 标签: 处理字符串格式的数字,json.Number 处理json字符串中的数字](13_serialize/01_std_json/02_number/main.go)
  - [1.3 自定义的时间格式解析](13_serialize/01_std_json/03_time/main.go)
  - [1.4 自定义的MarshalJSON方法 和 UnmarshalJSON](13_serialize/01_std_json/04_custom_marshal_unmarshal/main.go)
  - [1.5 使用匿名结构体添加字段, 使用匿名结构体组合多个结构体](13_serialize/01_std_json/05_anonymous_struct/main.go)
  - [1.6 inline 标签: 将嵌套结构体字段展开到父结构体中](13_serialize/01_std_json/06_inline/main.go)
- [2 Jsoniter(完全兼容标准库json，性能较好)-涉及标准库 encoding/json 分析](13_serialize/02_jsoniter/jsoniter.md)
  - 2.1 序列化
    - [2.1.1 指针变量，序列化时自动转换为它所指向的值](13_serialize/02_jsoniter/Marshal/01_pointer/main.go)
    - [2.1.2 结构体成员为interface{}](13_serialize/02_jsoniter/Marshal/02_Interface/main.go)
    - [2.1.3 extra.SetNamingStrategy 统一更改字段的命名风格](13_serialize/02_jsoniter/Marshal/03_name_field/main.go)
  - 2.2 反序列化
    - [2.2.1 反序列化匹配规则](13_serialize/02_jsoniter/Unmarshal/01_json_basic/main.go)
    - [2.2.2 json 字符串数组](13_serialize/02_jsoniter/Unmarshal/02_jsonArray/main.go)
    - [2.2.3 json.RawMessage 二次反序列化](13_serialize/02_jsoniter/Unmarshal/03_RawMessage/main.go)
    - [2.2.4 extra.SupportPrivateFields() 解析私有的字段](13_serialize/02_jsoniter/Unmarshal/04_private_field/main.go)
- [3 mapstructure 将通用的 map 值解码为 struct ](13_serialize/03_mapstructure/mapstructure.md)
  - [3.1 无tag标签](13_serialize/03_mapstructure/01_without_tag/main.go)
  - [3.2 带tag标签mapstructure](13_serialize/03_mapstructure/02_tag/main.go)
  - [3.3 embeded内嵌标签squash](13_serialize/03_mapstructure/03_embeded/main.go)
  - [3.4 未映射字段保留标签remain](13_serialize/03_mapstructure/04_remain/main.go)
  - [3.5 省略字段标签omitempty](13_serialize/03_mapstructure/05_omitempty/main.go)
  - [3.6 元数据展示源数据未映射字段](13_serialize/03_mapstructure/06_metadata/main.go)
  - [3.7 错误](13_serialize/03_mapstructure/07_error/main.go)
  - [3.8 弱解析](13_serialize/03_mapstructure/08_weekDecode/main.go)
  - [3.9 自定义解析器](13_serialize/03_mapstructure/09_decoder/main.go)
  - [3.10 time 类型 DecodeHookFunc ](13_serialize/03_mapstructure/10_time_decode_hook/main.go)
- [4 json patch 两种标准](13_serialize/04_json_patch/json_patch.md)
  - [4.1 github.com/evanphx/json-patch/v5 使用](13_serialize/04_json_patch/main.go)

## 第十四章 系统监控
- [1 systemstat包(适合linux系统，已断更)](14_system_monitor/01_systemstat/main.go)
- [2 gopsutil](14_system_monitor/02_gopsutil/gopsutil.md)
  - [2.1 cpu,mem,disk](14_system_monitor/02_gopsutil/01_disk_n_cpu_n_mem/main.go)
  - 2.2 进程信息获取
    - [物理机和虚拟机](14_system_monitor/02_gopsutil/02_process/01_in_host/main.go)
    - [容器环境](14_system_monitor/02_gopsutil/02_process/02_in_container/main.go)

- [3 prometheus](14_system_monitor/03_prometheus/prometheus.md)
  - [3.1 exporter](14_system_monitor/03_prometheus/01_exporter/exporter.md)
    - [3.1.1 内置 collector](14_system_monitor/03_prometheus/01_exporter/01_embeded_collector/main.go)
    - [3.1.2 使用自定义 collector](14_system_monitor/03_prometheus/01_exporter/02_customized_collector/main.go)
  - [3.2 client](14_system_monitor/03_prometheus/02_client/client.md) 
  - [3.3 k8s 部署](14_system_monitor/03_prometheus/03_k8s_deploy/deploy.md)
    - [3.2.1 原始 yaml --> 测试环境](14_system_monitor/03_prometheus/03_k8s_deploy/01_manual)
    - [3.2.2 Prometheus Operator --> 生产环境](14_system_monitor/03_prometheus/03_k8s_deploy/02_operator)
  - [3.4 PromQL(Prometheus Query Language)](14_system_monitor/03_prometheus/PromQL.md)
  - [3.5 存储模型及监控指标查询性能调优](14_system_monitor/03_prometheus/query.md)
- [4 AlertManager](14_system_monitor/04_alertmanager/alert_manager.md)

## [第十五章 分布式事务](15_distributed_transaction/distributed_transaction.md)
- Note: 使用DTM的代码作为案例 
- [1 两阶段提交2pc/XA](15_distributed_transaction/01_2pc_n_3pc/two_phase_commit.md)
- [2 saga事务](15_distributed_transaction/02_saga/saga.md)
- [3 TCC事务](15_distributed_transaction/03_tcc/tcc.md)
- [4 etcd的STM](15_distributed_transaction/04_stm/stm.md)

## [第十六章 copier(不同类型数据复制)](16_dataCopy/copier.md)

## 第十七章 数据加解密
- 1 phpserialize(不推荐)

## [第十八章 日志收集项目 log_collect](18_log_collect/log_collect.md)
- 1 动态选择文件
- 2 文件内容读取发送

## [第十九章 熔断,限流及降级](19_fuse_currentLimiting_degradation/rate_limit.md)
- [0 令牌桶官方包 x/time/rate](19_fuse_currentLimiting_degradation/00_tokenBucket/time_rate.md)
- 1 Sentinel-->滑动窗口
  - 1.1 基于流量QPS控制
    - [流量控制器的Token计算策略:direct](19_fuse_currentLimiting_degradation/01_sentinel/01_flow/direct/main.go)
    - [流量控制器的Token计算策略:warmUp](19_fuse_currentLimiting_degradation/01_sentinel/01_flow/warm_up/main.go)
  - 1.2 熔断
    - [ErrorCount](19_fuse_currentLimiting_degradation/01_sentinel/02_circuit_breaker/error_count/main.go)
    - [ErrorRatio](19_fuse_currentLimiting_degradation/01_sentinel/02_circuit_breaker/error_ratio/main.go)
    - [SlowRequestRatio](19_fuse_currentLimiting_degradation/01_sentinel/02_circuit_breaker/slow_request_ratio/main.go)
- [2 Hystrix-->滑动窗口](19_fuse_currentLimiting_degradation/02_hystrix/hystrix.md)
  - [2.1 客户端](19_fuse_currentLimiting_degradation/02_hystrix/client/client.go)
  - [2.2 服务端](19_fuse_currentLimiting_degradation/02_hystrix/server/server.go)
- [3 uber-go/ratelimit-->Leaky Bucket(漏桶)](19_fuse_currentLimiting_degradation/03_ubergo_ratelimit/uber-ratelimit.md)
- 4 envoyproxy/ratelimit-->计数器

## 第二十章 应用的命令行框架
- [1 Cobra -->在 k8s 中的应用](20_cli_frame/01_cobra/introdoction.md)
  - [1.1 cobra 构建 time 展示及解析,flag 使用](20_cli_frame/01_cobra/main.go)
- [2 Urfave Cli -->在 buildkit 中的应用](20_cli_frame/02_urfave_cli/urfave_cli.md)
- [3 alecthomas/kingpin -->在 node_exporter 使用](20_cli_frame/03_kingpin/kingpin.md)

## [第二十一章 配置文件解析:viper(依赖mapstructure,fsnotify,yaml,toml)](21_viper/viper.md)
- [1 viper获取本地文件内容](21_viper/01_read_n_watch_config/main.go)
- [2 监听文件变化(fsnotify)原理分析](21_viper/02_fsnotify/fsnotify.md)
- [3 远程读取nacos配置(源码分析)](21_viper/03_remote_config/remote_viper_config.md)
- [4 gopkg.in/yaml.v3 使用](21_viper/04_yaml/yaml.md)
- [5 pelletier/go-toml 使用](21_viper/05_toml/toml.md)
- [6 gopkg.in/ini.v1 使用](21_viper/06_ini/ini.md)

## 第二十二章 ETCD
- [服务端server--读和写流程分析](22_etcd/etcd_read_n_write.md)
- [服务端server--鉴权分析](22_etcd/ectd_auth.md)
- [服务端server--mvcc并发控制](22_etcd/etcd_mvcc.md)
- [服务端server--watch机制](22_etcd/03_watch/etcd_watch.md)
- [etcd 指标](22_etcd/etcd-metrics.md)
- [1 基本操作CRUD及watch监听](22_etcd/01_CRUD/main.go)
- [2 boltdb基本操作及在etcd中的源码分析](22_etcd/04_boltdb/boltdb.md)
- [3 bbolt改善boldb](22_etcd/05_bbolt/bbolt.md)

## 第二十三章 Go-Micro框架(不推荐)
- [1 Config配置加载包](23_micro/01_Config/config.md)

## [第二十四章 搜索引擎es](24_elasticSearch/es.md)
- [es索引及索引生命周期管理](24_elasticSearch/es_index.md)
- [1 go-elasticsearch 官方包](24_elasticSearch/01_official_pkg/go_elasticseach.md)
  - 1.1 批量写入Bulk
  - 1.2 es日志
  - 1.3 并发批量BulkIndexer

## [第二十五章 监控sentry](25_sentry/sentry.md)
- [1 结合gin基本shiyong](25_sentry/gin/main.go)
- [2 自定义zap core模块收集error级别日志上报sentry](25_sentry/zap_sentry/main.go)

## [第二十六章 图数据库Neo4j](26_neo4j/neo4j.md)
- [cypher语句](26_neo4j/cypher.md)
- [1 CRUD在web服务中](26_neo4j/main.go)

## 第二十七章 Mysql的binlog
- [binlog,gtid介绍](27_mysql_binlog/binlog.md)
- [canal使用及源码分析](27_mysql_binlog/canal/canal.md)

## [第二十八章 OLAP(Online Analytical Processing 联机分析处理)](28_OLAP/OLAP.md)
- [1 列数据库 ClickHouse](28_OLAP/01_clickHouse/clickHouse.md)
  - [1.1 database/sql 接口操作 clickHouse](28_OLAP/01_clickHouse/01_database_sql/main.go)
  - [1.2 原声接口操作 clickHouse](28_OLAP/01_clickHouse/02_native_interface/main.go)
  - [驱动 go-clickHouse 源码分析](28_OLAP/01_clickHouse/go-clickHouse.md)
  - [clickhouse 表引擎](28_OLAP/01_clickHouse/engine.md)
  - [clickhouse 基本命令](28_OLAP/01_clickHouse/curd.md)

## [第二十九章 分布式锁及源码分析](29_distributed_lock/distributed_lock.md)
- [1 redsync(RedLock 算法官方实现)](29_distributed_lock/01_redis_distributed_lock/main.go)
- [2 etcd实现分布式锁](29_distributed_lock/02_etcd_distributed_lock/main.go)

## [第三十章 Zookeeper](30_zookeeper/zookeeper.md)
- [1 github.com/go-zookeeper/zk 使用](30_zookeeper/zookeeper.go)

## [第三十一章 分布式 Id](31_distributed_Id/distribued_id.md)
- 1 UUID( Universally Unique Identifier 通用唯一标识码)
  - [1.1 github.com/google/uuid 8个版本使用](31_distributed_Id/01_uuid/main.go)
- 2 雪花算法
  - [2.1 bwmarrin/snowflake-->原生 twitter实现](31_distributed_Id/02_snowflake/01_bwmarrin_snowflake/main.gos)
  - [2.2 SonyFlake-->解决原生算法时间回拨问题](31_distributed_Id/02_snowflake/02_sony_snowflake/main.go)

## [第三十二章 Consensus algorithm 共识算法](32_consensus_algorithm/consensusAlgorithm.md)
- [1 Paxos 协议](32_consensus_algorithm/01_paxos/paxos.md)
- [2 Raft 协议](32_consensus_algorithm/02_raft/raft.md)
  - [hashicorp/raft 在consul实现](32_consensus_algorithm/02_raft/raft_in_consul.md)
  - [raft 在etcd实现原理分析](32_consensus_algorithm/02_raft/raft_in_etcd.md)
  - [2.1 使用 hashicorp/raft 调试应用](32_consensus_algorithm/02_raft/main.go)
- [3 gossip 协议](32_consensus_algorithm/03_gossip/gossip.md)
  - [3.1 github.com/hashicorp/memberlist 使用](32_consensus_algorithm/03_gossip/main.go)

## [第三十三章 压缩](33_compress/compress.md)
- [1 snappy 压缩库-->prometheus 使用](33_compress/01_snappy/snappy.md)


## [第三十四章 本地缓存](34_local_cache/cache.md)
- [1 go-cache源码分析及性能分析](34_local_cache/01_go_cache/go_cache.md)
- [2 free-cache源码分析及性能分析](34_local_cache/02_free_cache/free_cache.md)
- [3 hashicorp/golang-lru实现及变体](34_local_cache/03_lru/lru.md)

## [第三十五章 sonar静态代码质量分析-涉及与golangci-lint对比使用](35_sonar/sonar.md)

## [第三十六章 Proto管理工具 Buf](36_buf/buf_intro.md)

## [第三十七章 CI持续集成](37_CI/gitlabCI.md)
- [1 gitlab-runner 源码分析](37_CI/01_runner/runner.md)

## [第三十八章 Mergo实现 struct 与 map 之间转换-->k8s中应用](38_mergo/mergo.md)
- [1 struct 与 map 之间转换](38_mergo/01_map_to_struct/main.go)
- [2 override 覆盖选项](38_mergo/02_with_override/main.go)
- [3 结构体中的切片使用](38_mergo/03_slice/main.go)
- [4 类型检查](38_mergo/04_type_check/main.go)

## [第三十九章 权限管理 casbin](39_casbin/casbin.md)
- [1 ACL（access-control-list，访问控制列表)](39_casbin/01_acl/main.go)
- [2 RBAC (role-based-access-control 基于角色的权限访问控制)](39_casbin/02_rbac/rbac.md)
- [3 基于domain或tenant租户实现RBAC](39_casbin/03_domain_rbac/main.go)
- [4 ABAC(Attribute-based access control 基于属性的权限验证)使用 eval()功能构造来实现基于自定义规则](39_casbin/04_abac/main.go)


## [第四十章 规则引擎 rule engine](40_rules_engine/rule_engine.md)
- [1 govaluate-->casbin 使用](40_rules_engine/01_govaluate/govaluate.md)
- [2 bilibili/gengine](40_rules_engine/02_gengine/gengine.md)
- [3 expr-lang/expr-->argo-rollouts 使用](40_rules_engine/03_expr/expr.md)

## [第四十一章 hashicorp/go-plugin 插件使用-->httprunner 4.0 使用](41_go_plugin/go-plugin.md)

## [第四十二章 open-api](42_go-openapi/open-api.md)

## [第四十三章 go-systemd-->k8s 中使用](43_systemd/systemd.md)

## 参考
- [awesome-go](https://github.com/avelino/awesome-go)
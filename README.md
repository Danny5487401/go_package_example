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
- [1 grpc服务注册发现加健康检查](01_consul/01_http/test/consul_registry_test.go)

## 第二章 日志库
- 1 [zerolog](02_log/01_zerolog/zerolog.md)
- 2 [zap使用及源码分析](02_log/02_zap/zap.md)
  - [2.1 两种打印风格](02_log/02_zap/01_cosole/main.go)
  - [2.2 定义多种输出位置: 控制台输出及文件输出](02_log/02_zap/02_file_stdout/main.go)
  - [2.3 并发安全logger](02_log/02_zap/03_concurrency_safe/main.go)
  - [2.4 zap(配合lumberjack库或go-file-rotatelogs库)实现定制化log日志归档](02_log/02_zap/04_customized_log/lumberjack.md)
  - [2.5 简单的基于Entry实现的hook函数-->无法接收到Fields的相关参数](02_log/02_zap/05_hook/main.go)

## 第三章 消息队列
- [1 rabbitmq](03_amqp/01_rabbitmq/introduction.md)
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
- [4 Asynq分布式延迟队列](03_amqp/04_asynq/asynq.md)
  - [4.1 生产者](03_amqp/04_asynq/producer)
  - [4.2 消费者](03_amqp/04_asynq/server)

## [第四章 服务注册及配置文件中心Nacos](04_nacos/nacos.md)
- 1 [获取配置及监听文件变化](04_nacos/config_center/main.go)
- 2 服务注册，监听，获取
  - [V1版本](04_nacos/service_center/v1/main.go)
  - [V2版本](04_nacos/service_center/v2/main.go)

## 第五章 关系型数据库
- [go-mysql-driver插件源码分析](05_rds/go_mysql_driver.md)
- 1 GORM
  - 1.1 GORM原理及实现
  - 1.2 连接池使用
- [2 XORM](05_rds/02_xorm/xorm.md)
  - [2.1 主从连接](05_rds/02_xorm/util/util.go)
  - [2.2 调用mysql函数](05_rds/02_xorm/function/sum.go)
  - [2.3 事务处理](05_rds/02_xorm/transaction/transaction.go)
  - 2.4 crud
    - 插入Insert
    - 原生sql
    - 获取retrieve
    - 更新update

## 第六章 获取对外可用IP和端口
- [通过google, 国内移动、电信和联通通用的DNS获取对外Ip和端口](06_get_available_ip_port/get_ip/main.go)

## 第七章 Gin前端form验证器
- [1 验证器校验错误英转中](07_gin_form_validator/err_en_to_ch_translate/main.go)
- [2 前端数据校验](07_gin_form_validator/simpleForm/main.go)

## [第八章 GRPC编程 ](08_grpc/grpc.md)
- [MakeFile介绍](08_grpc/makefile.md)
  - [makefile在protobuf中应用](08_grpc/makefile)
- [grpc前置知识：http知识介绍](08_grpc/http.md)
- [protobuf及工具介绍](08_grpc/16_import_proto/protobuf_n_tools.md)
  - 引入其他proto文件,支持编译多个proto文件
  - 编码原理
  - protoc,protoc-gen-go,protoc-gen-go-grpc,protoc-gen-gofast等工具
- 1  HelloWorld入门
  - [1.1 客户端Grpc源码](08_grpc/01_grpc_helloworld/client/client.md)
  - [1.2 服务端Grpc源码](08_grpc/01_grpc_helloworld/server/server.md)
- 2  [context中的元数据metadata](08_grpc/02_metadata/grpc_context.md)
- 3  流式GRPC
- 4  protobuf的jsonpb包序列化和反序列化
- [5  负载均衡](08_grpc/05_grpc_load_balance/load_balance.md)
  - [5.1 客户端负载均衡(Resolver接口和Builder接口)](08_grpc/05_grpc_load_balance/client/builder_n_resolver_n_balancer.md)
    - [第三方consul实现Resolver接口和Builder接口](01_consul/02_grpc/consul_client/main.go)
    - [自定义实现Resolver接口和Builder接口](08_grpc/05_grpc_load_balance/client/customized_resolver_client/client.go)
    - [自定义实现nacos服务注册与发现](08_grpc/05_grpc_load_balance/client/nacos_client)
  - 5.2 服务端
- [6  retry机制](08_grpc/06_grpc_retry/retry.md)
- [7  grpc错误抛出与捕获](08_grpc/07_grpc_error/error.md)
- [8  auth自定义认证](08_grpc/08_grpc_token_auth/credentials.md)
- 9  Grpc插件-proto字段验证器
- [10 Grpc插件-grpc网关直接对外http服务-->etcd中应用](08_grpc/10_grpc_gateway/grpc_gateway.md)
- [11 Grpc插件-gogo/protobuf](08_grpc/11_protoc_gogofast/gogoprotobuf.md)
- [12 GRPC生态中间件(拦截器扩展)](08_grpc/12_grpc_middleware/01_grpc_interceptor/server/server.go)
  - 实现基于 CA 的 TLS 证书认证
  - go-grpc-middleware实现多个中间间：异常保护，日志
- 13 channelz调试
- 14 multiplex多路复用
- [15 自定义grpc插件](08_grpc/15_customized_protobuf_plugin/protobuf_extend.md)

## 第九章 Nosql非关系型数据库
- [1 MongoDB](09_Nosql/01_mongo/mongo.md)
  - [mongo和mysql对比：储存引擎及内存结构](09_Nosql/01_mongo/nosql_vs_rds.md)
  - 1.1 增删改查
- 2 [Redis(协议，原理，数据结构分析)](09_Nosql/02_redis/redis.md)
  - [redis底层数据结构对象源码分析](09_Nosql/02_redis/redis_obj.md)
  - 2.1 redigo使用
  - 2.2 go-redis使用
    - [2.2.1 连接池分析](09_Nosql/02_redis/02_go-redis/go-redis_pool.md)
    - [2.2.2 连接初始化及命令执行流程](09_Nosql/02_redis/02_go-redis/go-redis_init_n_excute.md)
    - [2.2.3 protocol协议封装](09_Nosql/02_redis/02_go-redis/go-redis_protocol.md)
    - [2.2.4 批处理pipeline分析](09_Nosql/02_redis/02_go-redis/go-redis_pipeline.md)

## [第十章 链路追踪(Distributed Tracing)](10_distributed_tracing/introduction.md)
- [1 Jaeger](10_distributed_tracing/01_jaeger/jaeger.md)
  - [1.1 结合XORM](10_distributed_tracing/01_jaeger/01_jaeger_xorm/main_test.go)
  - [1.2 结合redis](10_distributed_tracing/01_jaeger/02_jaeger_redis/hook.go)
- [2 OpenTelemetry](10_distributed_tracing/02_openTelemetry/openTelemetry.md)
  - 跨服务组合tracer代码展示:需开启svc1和svc2两个http服务(url可以是zipkin或则jaeger)

## 第十一章 依赖注入
- [1 dig依赖注入及http服务分层](11_dependency_injection/00_dig/dig.go)
- 2 wire依赖注入
  - [2.1 不使用wire现状](11_dependency_injection/01_wire/01_without_wire/main.go)
  - [使用wire优化](11_dependency_injection/01_wire/02_wire)
  - [wire使用-带err返回](11_dependency_injection/01_wire/03_wire_return_err/wire)
  - [wire使用-带参数初始化](11_dependency_injection/01_wire/04_wire_pass_params/wire)

## [第十二章 测试框架testify(gin使用)](12_testify/testify.md)
- [1 assert断言](12_testify/01_assert/calculate_test.go)
- [2 mock测试替身](12_testify/02_mock/main_test.go)
- [3 suite测试套件](12_testify/03_suite/suite_test.go)

## [第十三章 序列化反序列化-包含标准库源码分析](13_serialize/serialize.md)
- 1 Jsoniter(完全兼容标准库json，性能较好)
  - 1.1 序列化
    - [结构体成员为基本类型,嵌套结构体，及tag标签使用](13_serialize/01_jsoniter/Marshal/Basic/main.go)
    - [结构体成员为interface{}](13_serialize/01_jsoniter/Marshal/Interface/main.go)
  - 1.2 反序列化
    - [基本使用](13_serialize/01_jsoniter/Unmarshal/json/main.go)
    - [json字符串数组](13_serialize/01_jsoniter/Unmarshal/jsonArray/main.go)
    - [json.RawMessage二次反序列化](13_serialize/01_jsoniter/Unmarshal/RawMessage/main.go)
- 2 mapstructure使用（性能低但是方便）
  - [2.1 无tag标签](13_serialize/02_map2structure/01_without_tag/main.go)
  - [2.2 带tag标签mapstructure](13_serialize/02_map2structure/02_tag/main.go)
  - [2.3 embeded内嵌标签squash](13_serialize/02_map2structure/03_embeded/main.go)
  - [2.4 未映射字段保留标签remain](13_serialize/02_map2structure/04_remain/main.go)
  - [2.5 省略字段标签omitempty](13_serialize/02_map2structure/05_omitempty/main.go)
  - [2.6 元数据展示源数据未映射字段](13_serialize/02_map2structure/06_metadata/main.go)
  - [2.7 错误](13_serialize/02_map2structure/07_error/main.go)
  - [2.8 弱解析](13_serialize/02_map2structure/08_weekDecode/main.go)
  - [2.9 自定义解析器](13_serialize/02_map2structure/09_decoder/main.go)

## 第十四章 系统监控
- [1 systemstat包(适合linux系统，已断更)](14_system_monitor/01_systemstat/main.go)
- [2 gopsutil](14_system_monitor/02_gopsutil/gopsutil.md)
  - 进程信息获取
    - [物理机和虚拟机](14_system_monitor/02_gopsutil/process/in_host/main.go)
    - [容器环境](14_system_monitor/02_gopsutil/process/in_container/main.go)
  - [cpu,mem,disk](14_system_monitor/02_gopsutil/disk_n_cpu_n_mem/main.go)

## [第十五章 分布式事务](15_distributed_transaction/distributed_transaction.md)
- Note: 使用DTM的代码作为案例 
- [1 两阶段提交2pc/XA](15_distributed_transaction/01_2pc_n_3pc/two_phase_commit.md)
- [2 saga事务](15_distributed_transaction/02_saga/saga.md)
- [3 TCC事务](15_distributed_transaction/03_tcc/tcc.md)
- [4 etcd的STM](15_distributed_transaction/04_stm/stm.md)

## 第十六章 数据复制
- [1 copier(不同类型数据复制)](16_dataCopy/copier/copier.md)

## 第十七章 数据加解密
- 1 phpserialize

## [第十八章 日志收集项目 log_collect](18_log_collect/log_collect.md)
- 1 动态选择文件
- 2 文件内容读取发送

## [第十九章 熔断,限流及降级](19_fuse_currentLimiting_degradation/rate_limit.md)
- [0 熔断，降级，限流(官方包x/time/rate)](19_fuse_currentLimiting_degradation/00_tokenBucket/time_rate.md)
- 1 Sentinel
  - 1.1 基于流量QPS控制
    - [流量控制器的Token计算策略:direct](19_fuse_currentLimiting_degradation/01_sentinel/01_flow/direct/main.go)
    - [流量控制器的Token计算策略:warmUp](19_fuse_currentLimiting_degradation/01_sentinel/01_flow/warm_up/main.go)
  - 1.2 熔断
    - [ErrorCount](19_fuse_currentLimiting_degradation/01_sentinel/02_circuit_breaker/error_count/main.go)
    - [ErrorRatio](19_fuse_currentLimiting_degradation/01_sentinel/02_circuit_breaker/error_ratio/main.go)
    - [SlowRequestRatio](19_fuse_currentLimiting_degradation/01_sentinel/02_circuit_breaker/slow_request_ratio/main.go)
- [2 Hystrix](19_fuse_currentLimiting_degradation/02_hystrix/hystrix.md)
  - [2.1 客户端](19_fuse_currentLimiting_degradation/02_hystrix/client/client.go)
  - [2.2 服务端](19_fuse_currentLimiting_degradation/02_hystrix/server/server.go)

## [第二十章 命令行框架Cobra](20_cobra/introdoction.md)
- 1 介绍及功能使用
- 2 [在k8s中的应用](20_cobra/cobra_in_k8s.md)

## [第二十一章 配置文件获取工具viper(依赖mapstructure,fsnotify,yaml)](21_viper/viper.md)
- [1 viper获取本地文件内容](21_viper/01_read_n_watch_config/main.go)
- [2 监听文件变化(fsnotify)原理分析](21_viper/02_fsnotify/fsnotify.md)
- [3 远程读取nacos配置(源码分析)](21_viper/03_remote_config/remote_viper_config.md)
- [4 yaml.v3](21_viper/04_yaml/yaml.md)

## 第二十二章 ETCD
- [服务端server--读和写流程分析](22_etcd/etcd_read_n_write.md)
- [服务端server--鉴权分析](22_etcd/ectd_auth.md)
- [服务端server--mvcc并发控制](22_etcd/etcd_mvcc.md)
- [服务端server--watch机制](22_etcd/03_watch/etcd_watch.md)
- [1 基本操作CRUD及watch监听](22_etcd/01_CRUD/main.go)
- [2 boltdb基本操作及在etcd中的源码分析](22_etcd/04_boltdb/boltdb.md)
- [3 bbolt改善boldb](22_etcd/05_bbolt/bbolt.md)

## 第二十三章 Go-Micro框架 (不推荐使用)
- [1 Config配置加载包](23_micro/01_Config/config.md)

## [第二十四章 搜索引擎es](24_elasticSearch/es.md)
- [es索引及索引生命周期管理](24_elasticSearch/es_index.md)
- [1 官方包](24_elasticSearch/official_pkg/go_elasticseach.md)
  - 1.1 批量写入Bulk
  - 1.2 es日志
  - 1.3 并发批量BulkIndexer
- 2 第三方包oliver
  - 2.1 V6版本
  - 2.2 V7版本

## [第二十五章 监控sentry](25_sentry/sentry.md)
- [1 结合gin基本shiyong](25_sentry/gin/main.go)
- [2 自定义zap core模块收集error级别日志上报sentry](25_sentry/zap_sentry/main.go)

## [第二十六章 图数据库Neo4j](26_neo4j/neo4j.md)
- [cypher语句](26_neo4j/cypher.md)
- [1 CRUD在web服务中](26_neo4j/main.go)

## 第二十七章 Mysql的binlog
- [binlog,gtid介绍](27_mysql_binlog/binlog.md)
- [canal使用及源码分析](27_mysql_binlog/canal/canal.md)

## 第二十八章 OLAP(Online Analytical Processing联机分析处理)
- [1 列数据库ClickHouse](28_OLAP/01_clickHouse/clickHouse.md)
  - [1.1 标准库sql操作clickHouse](28_OLAP/01_clickHouse/01_database_sql/main.go)
  - [1.2 扩展包sqlx操作clickHouse](28_OLAP/01_clickHouse/02_sqlx/main.go)
  - [go-clickHouse源码分析](28_OLAP/01_clickHouse/go-clickHouse.md)

## [第二十九章 分布式锁及源码分析](29_distributed_lock/distributed_lock.md)
- [1 redsync(RedLock 算法官方实现)](29_distributed_lock/01_redis_distributed_lock/main.go)
- [2 etcd实现分布式锁](29_distributed_lock/02_etcd_distributed_lock/main.go)

## [第三十章 Zookeeper](30_zookeeper/zookeeper.md)

## [第三十一章 分布式Id](31_distributed_Id/distribued_id.md)
- 雪花算法
  - [bwmarrin/snowflake库](31_distributed_Id/snowflake/main.go)
  - [SonyFlake(解决时间回拨问题)](31_distributed_Id/sony_snowflake/main.go)

## [第三十二章 多副本常用的技术方案及Raft协议](32_raft/raft.md)
  - [raft在consul实现](32_raft/raft_in_consul.md)
  - [raft在etcd实现原理分析](32_raft/raft_in_etcd.md)
  - [1 使用hashicorp/raft调试应用](32_raft/main.go)

## [第三十三章 多副本常用的技术方案及Paxos协议](33_paxos/paxos.md)

## [第三十四章 本地缓存](34_local_cache/cache.md)
- [1 go-cache源码分析及性能分析](34_local_cache/01_go_cache/go_cache.md)
- [2 free-cache源码分析及性能分析](34_local_cache/02_free_cache/free_cache.md)

## [第三十五章 sonar静态代码质量分析-涉及与golangci-lint对比使用](35_sonar/sonar.md)

## [第三十六章 Proto管理工具Buf](36_buf/buf_intro.md)

## [第三十六章 CI持续集成](37_CI/gitlabCI.md)




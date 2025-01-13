<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [etcd 指标](#etcd-%E6%8C%87%E6%A0%87)
  - [指标分类](#%E6%8C%87%E6%A0%87%E5%88%86%E7%B1%BB)
    - [服务端 server](#%E6%9C%8D%E5%8A%A1%E7%AB%AF-server)
    - [磁盘](#%E7%A3%81%E7%9B%98)
    - [网络](#%E7%BD%91%E7%BB%9C)
    - [调试及不稳定指标](#%E8%B0%83%E8%AF%95%E5%8F%8A%E4%B8%8D%E7%A8%B3%E5%AE%9A%E6%8C%87%E6%A0%87)
    - [快照](#%E5%BF%AB%E7%85%A7)
    - [Prometheus 库相关指标](#prometheus-%E5%BA%93%E7%9B%B8%E5%85%B3%E6%8C%87%E6%A0%87)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# etcd 指标

## 指标分类

etcd 指标名称符合 prometheus 最佳实践


### 服务端 server

以 etcd_server_ 开头.

| 指标 |         含义          |                     异常数值说明                      |
|:--:|:-------------------:|:-----------------------------------------------:|
| has_leader |        可用状态         |                    0代表集群不可用                     |
| leader_changes_seen_total |        切主次数         |                   频繁更改代表集群不稳定                   |
| proposals_committed_total | 写入持久化存储的proposals次数 |           如果leader和Member持续较大的滞后代表不健康           |
| proposals_applied_total |  写入状态机proposals次数   | 正常applied <= committed差值不大,如果比较大,代表大量查询或则事务导致过载 |
| proposals_pending |      排队commit       |            高数值代表客户端或则member无法commit             |
| proposals_failed_total |        失败计数         |                   与领导者选举相关的临时故障或由于集群中法定人数丢失而导致的较长停机时间                   |

### 磁盘

以 etcd_disk_ 开头.

| 指标 |        含义        |    异常数值说明    |
|:--:|:----------------:|:------------:|
| wal_fsync_duration_seconds | wal 模块fsync的延迟分布 | 高延迟代表有磁盘性能问题 |
| backend_commit_duration_seconds |    增量快照时延迟分布     |        高延迟代表有磁盘性能问题      |



### 网络

以 etcd_network_ 开头.

| 指标 |        含义        |    异常数值说明    |
|:--:|:----------------:|:------------:|
| peer_sent_bytes_total | 发送给 peer 的bytes  | |
| peer_received_bytes_total |  接收 peer 的bytes  | |
| peer_sent_failures_total |  发送到 peer 的失败次数  | |
| peer_received_failures_total |  接收 peer 的失败次数   | |
| peer_round_trip_time_seconds | peers间 RTT 延迟分布  | |
| client_grpc_sent_bytes_total | 发送给grpc客户端的bytes | |
| client_grpc_received_bytes_total | 从grpc客户端接收的bytes | |


### 调试及不稳定指标
以 etcd_debugging 开头


### 快照


| 指标 |        含义        |      异常数值说明      |
|:--:|:----------------:|:----------------:|
| snapshot_save_total_duration_seconds | 发送给 peer 的bytes  | 高时延代表磁盘问题导致集群不稳定 |


### Prometheus 库相关指标

| 指标 |      含义      |            异常数值说明            |
|:--:|:------------:|:----------------------------:|
| process_open_fds |  打开的文件描述符数量  | 与进程限制的描述符process_max_fds进行对比 |
| process_max_fds | 允许打开的文件描述符数量 |                              |



## 参考
- https://etcd.io/docs/v3.4/metrics/

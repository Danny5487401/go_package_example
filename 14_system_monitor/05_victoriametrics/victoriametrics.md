
# VictoriaMetrics 




## 架构
![vms_structure.png](vms_structure.png)
https://docs.victoriametrics.com/victoriametrics/cluster-victoriametrics/#architecture-overview

VictoriaMetrics在保持更简单的架构的同时，还包括几个核心组件：

- vmstorage：数据存储以及查询结果返回，默认端口为 8482
- vminsert：数据录入，可实现类似分片、副本功能，默认端口 8480
- vmselect：数据查询，汇总和数据去重，默认端口 8481
- vmagent：数据指标抓取，支持多种后端存储，会占用本地磁盘缓存，默认端口 8429
- vmalert：报警相关组件，不如果不需要告警功能可以不使用该组件，默认端口为 8880


### vmui 

地址: http://<vmselect>:8481/select/<accountID>/vmui/

url 格式: https://docs.victoriametrics.com/victoriametrics/cluster-victoriametrics/#url-format

```shell
✗ kubectl get svc -n victoria                                                                             
NAME                         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                      AGE
my-vmselect-vmcluster-svc    NodePort    10.233.6.15     <none>        8481:32110/TCP               38m
vmagent-vmagent              ClusterIP   10.233.34.26    <none>        8429/TCP                     5h27m
vmclusterlb-vmcluster        ClusterIP   10.233.29.230   <none>        8427/TCP                     5h27m
vminsert-vmcluster           ClusterIP   10.233.52.18    <none>        8480/TCP                     5h25m
vminsertinternal-vmcluster   ClusterIP   None            <none>        8480/TCP                     5h25m
vmselect-vmcluster           ClusterIP   10.233.28.216   <none>        8481/TCP                     5h25m
vmselectinternal-vmcluster   ClusterIP   None            <none>        8481/TCP                     5h25m
vmstorage-vmcluster          ClusterIP   None            <none>        8482/TCP,8400/TCP,8401/TCP   5h26m

```
- grafana 配置地址: 如果使用 prometheus 接口, http://vmselect-vmcluster.victoria.svc:8481/select/0/prometheus


### vmagent
- https://docs.victoriametrics.com/victoriametrics/vmagent/
- https://docs.victoriametrics.com/operator/resources/vmagent/




## 特点

- 解决高基数问题 high cardinality: up to 7x less RAM than Prometheus
- 高流失率 high churn rate(time series 频繁被替代) 优化:  
- 存储空间降低:7x less storage space is required compared to Prometheus


## 部署
使用 operator 的方式:https://docs.victoriametrics.com/guides/getting-started-with-vm-operator/
```shell
helm install victoria-operator vm/victoria-metrics-operator --version 0.47.0
```

这里使用 local-path storageclass 
```yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMCluster
metadata:
  name: vmcluster
  namespace: victoria
spec:
  retentionPeriod: "7"
  requestsLoadBalancer:
    enabled: true
    spec:
      replicaCount: 1
  vmstorage:
    replicaCount: 2
    image:
      repository: swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/victoriametrics/vmstorage
      tag: v1.115.0-cluster
      pullPolicy: IfNotPresent
    storage:
      volumeClaimTemplate:
        metadata:
          name: data
        spec:
          accessModes: [ "ReadWriteOnce" ]
          storageClassName: local-path
          resources:
            requests:
              storage: 5Gi
    resources:
      limits:
        cpu: "1"
        memory: "1Gi"
      requests:
        cpu: "0.5"
        memory: "500Mi"
              
  vmselect:
    replicaCount: 1
    cacheMountPath: "/select-cache"
    image:
      repository: swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/victoriametrics/vmselect
      tag: v1.115.0-cluster
      pullPolicy: IfNotPresent
    storage:
      volumeClaimTemplate:
        spec:
          storageClassName: local-path
          resources:
            requests:
              storage: "1Gi"
    resources:
      limits:
        cpu: "1"
        memory: "1Gi"
      requests:
        cpu: "0.5"
        memory: "500Mi"
        
  vminsert:
    replicaCount: 1
    image:
      repository: swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/victoriametrics/vminsert
      tag: v1.115.0-cluster
      pullPolicy: IfNotPresent
    resources:
      limits:
        cpu: "1"
        memory: "1Gi"
      requests:
        cpu: "0.5"
        memory: "500Mi"

```

vmagent : 修改 remoteWrite

```yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMAgent
metadata:
  name: vmagent
  namespace: victoria
spec:
  image:
    repository: swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/victoriametrics/vmagent
    tag: v1.116.0
  serviceScrapeNamespaceSelector: {}
  podScrapeNamespaceSelector: {}
  podScrapeSelector: {}
  serviceScrapeSelector: {}
  nodeScrapeSelector: {}
  nodeScrapeNamespaceSelector: {}
  staticScrapeSelector: {}
  staticScrapeNamespaceSelector: {}
  replicaCount: 1
  remoteWrite:
    - url: "http://vminsert-vmcluster.victoria.svc.cluster.local:8480/insert/0/prometheus/api/v1/write"
#  secrets:
#    - etcd-client-cert
```

## VictoriaMetrics 对比 prometheus 

- https://last9.io/blog/prometheus-vs-victoriametrics/


## 数据格式
```shell
/vmstorage-data # tree -L 2 /vmstorage-data
-L [error opening dir]
2 [error opening dir]
.
├── data
│   ├── big
│   │   ├── 2025_06
│   │   └── snapshots
│   └── small
│       ├── 2025_06
│       │   ├── 1848DA53DBC2BFDA
│       │   │   ├── index.bin
│       │   │   ├── metadata.json
│       │   │   ├── metaindex.bin
│       │   │   ├── timestamps.bin
│       │   │   └── values.bin
            // ...
│       │   ├── 1848DA53DBC2C383
│       │   │   ├── index.bin
│       │   │   ├── metadata.json
│       │   │   ├── metaindex.bin
│       │   │   ├── timestamps.bin
│       │   │   └── values.bin
│       │   └── parts.json
│       └── snapshots
├── flock.lock
├── indexdb
│   ├── 1848CA450436BA0C
│   │   └── parts.json
│   ├── 1848CA450436BA0D
│   │   ├── 1848CA45062B2EC6
│   │   │   ├── index.bin
│   │   │   ├── items.bin
│   │   │   ├── lens.bin
│   │   │   ├── metadata.json
│   │   │   └── metaindex.bin
        // ...
│   │   └── parts.json
│   ├── 1848CA450436BA0E
│   │   └── parts.json
│   └── snapshots
├── metadata
│   └── minTimestampForCompositeIndex
└── snapshots
```

最主要的是数据目录data和索引目录indexdb，flock.lock文件为文件锁文件，用于VictoriaMetrics进程锁住文件，不允许别的进程进行修改目录或文件。

### 数据目录 data 

VictoriaMetrics分成small目录和big目录，主要是兼顾近期数据的读取和历史数据的压缩率。

在small目录下，以月为单位不断生成partition目录.

big目录下的数据由small目录下的数据在后台compaction时合并生成.



### 索引目录 indexdb

VictoriaMetrics每次内存Flush或者后台Merge时生成的索引part，主要包含metaindex.bin、index.bin、lens.bin、items.bin等4个文件。

## prometheus 迁移到 VictoriaMetrics


- vmagent 替换  prometheus 的 scraping target
- VMServiceScrape (instead of ServiceMonitor)
- VMPodScrape (instead of PodMonitor
- VMRule (instead of PrometheusRule)


会自动将 Prometheus ServiceMonitor, PodMonitor, PrometheusRule, Probe and ScrapeConfig objects 转换成  VictoriaMetrics Operator objects.



### prometheus 存在的问题
https://zetablogs.medium.com/supercharge-your-monitoring-migrate-from-prometheus-to-victoriametrics-for-scalability-and-speed-e1e9df786145
- 单体结构
- 非分布式,只能垂直伸缩
- wal 回放慢
- 内存占用高

## 参考

- [浅析下开源时序数据库VictoriaMetrics的存储机制](https://zhuanlan.zhihu.com/p/368912946)



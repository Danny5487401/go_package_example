<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [prometheus Operator](#prometheus-operator)
  - [prometheus 部署](#prometheus-%E9%83%A8%E7%BD%B2)
    - [测试环境: 使用 prometheus 的docker环境](#%E6%B5%8B%E8%AF%95%E7%8E%AF%E5%A2%83-%E4%BD%BF%E7%94%A8-prometheus-%E7%9A%84docker%E7%8E%AF%E5%A2%83)
      - [作业和实例](#%E4%BD%9C%E4%B8%9A%E5%92%8C%E5%AE%9E%E4%BE%8B)
    - [生产环境: 使用 prometheus-operator](#%E7%94%9F%E4%BA%A7%E7%8E%AF%E5%A2%83-%E4%BD%BF%E7%94%A8-prometheus-operator)
  - [背景](#%E8%83%8C%E6%99%AF)
  - [工作原理](#%E5%B7%A5%E4%BD%9C%E5%8E%9F%E7%90%86)
  - [Operator 能做什么](#operator-%E8%83%BD%E5%81%9A%E4%BB%80%E4%B9%88)
  - [安装](#%E5%AE%89%E8%A3%85)
  - [操作](#%E6%93%8D%E4%BD%9C)
    - [使用后效果](#%E4%BD%BF%E7%94%A8%E5%90%8E%E6%95%88%E6%9E%9C)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# prometheus Operator

## prometheus 部署

### 测试环境: 使用 prometheus 的docker环境
```yaml
# prometheus.yml
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
    monitor: 'codelab-monitor'

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
- job_name: "go-test"
  scrape_interval: 60s
  scrape_timeout: 60s
  metrics_path: "/metrics"

  static_configs:
  - targets: ["localhost:8888"]

```

可以看到配置文件中指定了一个job_name，所要监控的任务即视为一个job, scrape_interval和scrape_timeout是pro进行数据采集的时间间隔和频率，metrics_path指定了访问数据的http路径，target是目标的ip:port,这里使用的是同一台主机上的8888端口

```shell
docker run -p 9090:9090 -v /Users/python/Desktop/github.com/Danny5487401/go_advanced_code/chapter02_goroutine/02_runtime/07prometheus/client/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```
![](.prometheus_images/prometheus_panel.png)

启动之后可以访问web页面http://localhost:9090/graph,在status下拉菜单中可以看到配置文件和目标的状态，此时目标状态为DOWN，因为我们所需要监控的服务还没有启动起来，那就赶紧步入正文，用pro golang client来实现程序吧。

![](.prometheus_images/server_state.png)

启动后状态
![](.prometheus_images/server_state2.png)


#### 作业和实例
在prometheus.yml配置文件中，添加如下配置
```yaml
scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']
```

当前在每一个Job中主要使用了静态配置(static_configs)的方式定义监控目标。
除了静态配置每一个Job的采集Instance地址以外，Prometheus还支持与DNS、Consul、E2C、Kubernetes等进行集成实现自动发现Instance实例，并从这些Instance上获取监控数据。

在Prometheus配置中，一个可以拉取数据的端点IP:Port叫做一个实例（instance），而具有多个相同类型实例的集合称作一个作业（job）
```yaml
- job: api-server
- instance 1: 1.2.3.4:5670
- instance 2: 1.2.3.4:5671
- instance 3: 5.6.7.8:5670
- instance 4: 5.6.7.8:5671

```
在Prometheus中，每一个暴露监控样本数据的HTTP服务称为一个实例。例如在当前主机上运行的node exporter可以被称为一个实例(Instance)


当Prometheus拉取指标数据时，会自动生成一些标签（label）用于区别抓取的来源：
![](.prometheus_images/target_in_ui.png)
- job：配置的作业名；
- instance：配置的实例名，若没有实例名，则是抓取的IP:Port

对于每一个实例（instance）的抓取，Prometheus会默认保存以下数据：

- up{job="<job>", instance="<instance>"}：如果实例是健康的，即可达，值为1，否则为0；
- scrape_duration_seconds{job="<job>", instance="<instance>"}：抓取耗时；
- scrape_samples_post_metric_relabeling{job="<job>", instance="<instance>"}：指标重新标记后剩余的样本数。
- scrape_samples_scraped{job="<job>", instance="<instance>"}：实例暴露的样本数
  该up指标对于监控实例健康状态很有用。

### 生产环境: 使用 prometheus-operator


## 背景
为了在Kubernetes能够方便的管理和部署Prometheus，我们使用ConfigMap了管理Prometheus配置文件。
每次对Prometheus配置文件进行升级时，我们需要手动移除已经运行的Pod实例，从而让Kubernetes可以使用最新的配置文件创建Prometheus。
而如果当应用实例的数量更多时，通过手动的方式部署和升级Prometheus过程繁琐并且效率低下。
```yaml
    scrape_configs:
      - job_name: 'prometheus'
        static_configs:
        - targets: ['localhost:9090']
```

从本质上来讲Prometheus属于是典型的有状态应用，而其有包含了一些自身特有的运维管理和配置管理方式。而这些都无法通过Kubernetes原生提供的应用管理概念实现自动化。
为了简化这类应用程序的管理复杂度，CoreOS率先引入了Operator的概念，并且首先推出了针对在Kubernetes下运行和管理Etcd的Etcd Operator。并随后推出了Prometheus Operator


## 工作原理
![](../../.operator_images/operator_discipline.png)

Prometheus的本职就是一组用户自定义的CRD资源以及Controller的实现，Prometheus Operator负责监听这些自定义资源的变化，并且根据这些资源的定义自动化的完成如Prometheus Server自身以及配置的自动化管理工作

## Operator 能做什么

Prometheus Operator为我们提供了哪些自定义的Kubernetes资源，列出了Prometheus Operator目前提供的️资源：

- Prometheus：声明式创建和管理Prometheus Server实例；
- Alertmanager：声明式的创建和管理Alertmanager实例。
- ServiceMonitor：负责声明式的管理监控配置；
- PrometheusRule：负责声明式的管理告警配置；


还有thanosRuler,podMonitor,Probe等

## 安装
由于需要对Prometheus Operator进行RBAC授权，而默认的bundle.yaml中使用了default命名空间，因此，在安装Prometheus Operator之前需要先替换一下bundle.yaml文件中所有namespace定义，由default修改为monitoring。
```shell
$ kubectl -n monitoring apply -f bundle.yaml
clusterrolebinding.rbac.authorization.k8s.io/prometheus-operator created
clusterrole.rbac.authorization.k8s.io/prometheus-operator created
deployment.apps/prometheus-operator created
serviceaccount/prometheus-operator created
service/prometheus-operator created
```

```shell
$ kubectl -n monitoring get pods
NAME                                   READY     STATUS    RESTARTS   AGE
prometheus-operator-6db8dbb7dd-2hz55   1/1       Running   0          19s
```

##  操作
1. 部署正常的业务http 服务:deployment-app.yaml
```yaml
kind: Service
apiVersion: v1
metadata:
  name: danny-example-app
  labels:
    app: danny-example-app
spec:
  selector:
    app: danny-example-app
  ports:
    - name: web
      port: 8080
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: danny-example-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: danny-example-app
  template:
    metadata:
      labels:
        app: danny-example-app
    spec:
      containers:
        - name: danny-example-app
          image: fabxc/instrumented_app
          ports:
            - name: web
              containerPort: 8080
```

2. 部署监听服务serviceMonitor.yaml：使用matchlabels中的app: danny-example-app去监听业务服务
```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: danny-service-monitor-example-app
  namespace: monitoring
  labels:
    team: frontend
spec:
  # namespaceSelector定义让其可以跨命名空间
  namespaceSelector:
    matchNames:
      - danny-xia
  selector:
    matchLabels:
      app: danny-example-app
  endpoints:
    - port: web
```
3. 配置serviceAccount方便promethues拉取数据:主要是nonResourceURLs: ["/metrics"]获取资源
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: danny-prometheus
  namespace: monitoring
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: danny-prometheus-cluster-role
rules:
  - apiGroups: [""]
    resources:
      - nodes
      - services
      - endpoints
      - pods
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources:
      - configmaps
    verbs: ["get"]
  - nonResourceURLs: ["/metrics"]
    verbs: ["get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: danny-prometheus-cluster-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: danny-prometheus-cluster-role
subjects:
  - kind: ServiceAccount
    name: danny-prometheus
    namespace: monitoring
```
4. 配置prometheus rules:alert_rules.yaml,其中设置标签prometheus: example，role: alert-rules
```
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  labels:
    prometheus: example
    role: alert-rules
  name: danny-prometheus-example-rules
  namespace: monitoring
spec:
  groups:
    - name: ./example.rules
      rules:
        - alert: ExampleAlert
          expr: vector(1)
```
5. 定义alert manager的全局配置secret:alertmanager.yaml,注意alertmanager-danny-alert-instance这名字是固定的，在默认情况下，会通过alertmanager-{ALERTMANAGER_NAME}的命名规则去查找Secret配置并以文件挂载的方式
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: alertmanager-danny-alert-instance
  namespace: monitoring
type: Opaque
stringData:
  alertmanager.yaml: |-
    global:
      resolve_timeout: 5m
    route:
      group_by: ['job']
      group_wait: 30s
      group_interval: 5m
      repeat_interval: 12h
      receiver: 'webhook'
    receivers:
      - name: 'webhook'
        webhook_configs:
          - url: 'http://alertmanagerwh:30500/'
```
6. 开启alert manager实例:alert-manager-instance.yaml
```yaml
apiVersion: monitoring.coreos.com/v1
kind: Alertmanager
metadata:
  name: danny-alert-instance
  namespace: monitoring
spec:
  replicas: 3
```
7. 暴露alert manager 的service:alertmanager-danny-alert-instance-srv.yaml
```yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"labels":{"alertmanager":"danny-alert-instance"},"name":"alertmanager-danny-alert-instance","namespace":"monitoring"},"spec":{"ports":[{"name":"web","port":9093,"targetPort":"web"}],"selector":{"alertmanager":"danny-alert-instance","app":"alertmanager"},"sessionAffinity":"ClientIP"}}
  labels:
    alertmanager: danny-alert-instance
  name: alertmanager-danny-alert-instance
  namespace: monitoring

spec:
  ports:
    - name: web
      port: 9093
      protocol: TCP
      targetPort: web
  selector:
    alertmanager: danny-alert-instance
    app: alertmanager


```


8. 定义prometheus实例：serviceAccountName指定步骤3的danny-prometheus账号，ruleSelector指定步骤4的规则，使用serviceMonitorSelector中的team: frontend去关联步骤2的monitor实例，alerting找步骤7暴露的endpoint
```yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: danny-instance
  namespace: monitoring
spec:
  serviceAccountName: danny-prometheus
  serviceMonitorSelector:
    matchLabels:
      team: frontend
  ruleSelector:
    matchLabels:
      role: alert-rules
      prometheus: example
  alerting:
    alertmanagers:
      - name: alertmanager-danny-alert-instance
        namespace: monitoring
        port: web
  resources:
    requests:
      memory: 400Mi

```

### 使用后效果
1. 配置变化
![](../../.operator_images/operator_effect_on_config.png)
2. 服务发现变化
![](../../.operator_images/operator_effect_on_discovery.png)
3. rules变化
![](../../.operator_images/operator_effect_on_rule.png)



Note:不使用账号会报没有权限
```shell
(⎈ |teleport.gllue.com-test:danny-xia)➜  github.com/Danny5487401/go_advanced_code git:(feature/monitor) ✗ kubectl logs prometheus-danny-instance-0 prometheus -n monitoring --tail 2 
level=error ts=2022-05-05T08:17:55.023Z caller=klog.go:94 component=k8s_client_runtime func=ErrorDepth msg="/app/discovery/kubernetes/kubernetes.go:263: Failed to list *v1.Pod: pods is forbidden: User \"system:serviceaccount:monitoring:default\" cannot list resource \"pods\" in API group \"\" in the namespace \"danny-xia\""
level=error ts=2022-05-05T08:17:55.030Z caller=klog.go:94 component=k8s_client_runtime func=ErrorDepth msg="/app/discovery/kubernetes/kubernetes.go:262: Failed to list *v1.Service: services is forbidden: User \"system:serviceaccount:monitoring:default\" cannot list resource \"services\" in API group \"\" in the namespace \"danny-xia\""

```
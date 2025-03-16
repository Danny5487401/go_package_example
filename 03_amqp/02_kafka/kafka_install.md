# kafka 安装


```shell
(⎈|kubeasz-test:kafka)➜  kafka helm repo add bitnami  https://charts.bitnami.com/bitnami
"bitnami" has been added to your repositories
(⎈|kubeasz-test:kafka)➜  kafka helm search repo -l bitnami/kafka --version 31.5.0
(⎈|kubeasz-test:kafka)➜  kafka cat values.yaml
global:
  defaultStorageClass: local-path
  security:
    ## @param global.security.allowInsecureImages Allows skipping image verification
    allowInsecureImages: true
controller:
  ## @param controller.replicaCount Number of Kafka controller-eligible nodes
  ## Ignore this section if running in Zookeeper mode.
  ##
  replicaCount: 3
  ## @param controller.controllerOnly If set to true, controller nodes will be deployed as dedicated controllers, instead of controller+broker processes.
  ##
image:
  registry: swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io
  repository: bitnami/kafka
  tag: 3.9

      
(⎈|kubeasz-test:kafka)➜  kafka helm install my-kafka -f values.yaml oci://registry-1.docker.io/bitnamicharts/kafka
(⎈|kubeasz-test:kafka)➜  kafka kubectl get secret my-kafka-user-passwords --namespace kafka -o jsonpath='{.data.client-passwords}' | base64 -d | cut -d , -f 1
qwl2pTlW6e%

(⎈|kubeasz-test:kafka)➜  kafka cat kafka-ui.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kafka-ui
  labels:
    app: kafka-ui
  namespace: kafka    
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kafka-ui
  template:
    metadata:
      labels:
        app: kafka-ui
    spec:
      containers:
      - name: kafka-ui
        image: swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/provectuslabs/kafka-ui:v0.7.2
        env:
        - name: KAFKA_CLUSTERS_0_NAME
          value: 'Kafka Cluster'
        - name: KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS
          value: 'my-kafka.kafka.svc.cluster.local:9092'
        - name: KAFKA_CLUSTERS_0_PROPERTIES_SECURITY_PROTOCOL
          value: 'SASL_PLAINTEXT'
        - name: KAFKA_CLUSTERS_0_PROPERTIES_SASL_MECHANISM
          value: 'PLAIN'
        - name: KAFKA_CLUSTERS_0_PROPERTIES_SASL_JAAS_CONFIG
          value: 'org.apache.kafka.common.security.scram.ScramLoginModule required username="user1" password="qwl2pTlW6e";'
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: kafka-ui
  namespace: kafka     
spec:
  selector:
    app: kafka-ui
  type: NodePort
  ports:
    - protocol: TCP
      port: 8080
      targetPort: 8080

# 使用 kt-connect 进行解析集群 k8s 地址
(⎈|kubeasz-test:kafka)➜  kt-connect git:(master) ✗ sudo ./bin/ktcl connect  -n default -d  --includeDomains "gllue.host,aliyuncs.com,svc.cluster.local" --includeIps '172.25.0.0/16,192.168.0.0/16' --podQuota 0.5c,512m 

```

## 参考
- https://github.com/bitnami/charts/tree/main/bitnami/kafka
- [使用 Bitnami Helm 安装 Kafka](https://juejin.cn/post/7183502453873541177)
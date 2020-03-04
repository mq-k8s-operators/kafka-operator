# Kafka Operator

Kafka Operator 构建于 [Operator framework](https://github.com/operator-framework/operator-sdk) 之上，用于在 [Kubernetes](https://kubernetes.io/) 上快速搭建[Apache Kafka](http://kafka.apache.org/)集群。

本项目为简易版本，为了简单可控性而自研并逐渐完善，可作为Operator framework有状态服务的开发例子。如果需要更完善的Operator，有以下项目可以参考。

- https://github.com/banzaicloud/kafka-operator
- https://github.com/krallistic/kafka-operator

# Features

- 自动创建Zookeeper集群
- 自动创建Apache Kafka集群
- 自动创建Yahoo Kafka Manager
- 自动创建K8S集群内部的LB Service
- 自动创建Kafka Manager Ingress 服务
- 基于Kafka Manager调整 + StatefulSet Replicas的扩容/缩容

# Components

- https://github.com/operator-framework/operator-sdk
- https://github.com/pravega/zookeeper-operator current use Tag 0.2.4
- https://github.com/yahoo/kafka-manager
- https://github.com/wurstmeister/kafka-docker

# Usage

```
# build your own operator image
git clone https://github.com/k8s-operators/kafka-operator.git
cd kafka-operator
operator-sdk generate k8s
operator-sdk build your-docker-name/kafka-operator:v1.0.0
docker push your-docker-name/kafka-operator:v1.0.0

# replace deploy/operator.yaml to use your operator image
# or just use jianzhiunique/kafka-operator:v1.0.0

# deploy files under "deploy" dir to server
kubectl apply -f service_account.yaml
kubectl apply -f role.yaml
kubectl apply -f role_binding.yaml
kubectl apply -f operator.yaml
kubectl apply -f deploy/crds/jianzhiunique.github.io_kafkas_crd.yaml

# config cr.yaml
# for all fields, see pkg/apis/jianzhiunique/v1/kafka_types.go
# we apply some default values for you, see pkg/utils/check_cr.go

# the only config that you must specify is storage_class_name, 
# for kafka, we recommend users to use local pv

```
# TODO

- 对K8S集群外部的服务支持
- TLS支持
- ACL支持
- JMX监控及Grafana Dashboard
- Lag监控
- 滚动升级
- 平滑扩容缩容

 # Others
 
 - http://kafka.apache.org/documentation/

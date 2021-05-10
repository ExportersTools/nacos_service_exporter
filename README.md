NacosServiceExporter
====================

NacosServiceExporter 是一个监控 nacos 注册中心服务数量的 exporter.

该项目工程主要用于微服务建设初期, 监控比较薄弱时期, 后期数据微服务监控完善之后, 该组建意义不大. (K8S环境下, 基本无意义.)


建议:

    由于在 发布的时候/主机宕机/程序迁移/等 的时候会有一定的变更, 所以在配置告警的时候需要谨慎配置,不能因为正常行为导致触发告警.

主要采集项(metrics):

    服务是否正常访问Nacos: up
    微服务注册实例数量: serviceInstanceCount
    所有服务数量: allServiceCount
    所有服务实例数量: allServiceInstanceCount
    每个节点(主机)服务个数: endpointServiceCount

StartUP
---------

### 环境变量方法

```shell script
export endPoint="http://127.0.0.1:8848"
export nameSpaceId="xxx-xxx"
./nacosServiceExporter
```

### 传参方式

```shell script
./nacosServiceExporter -endPoint http://127.0.0.1:8848 -nameSpaceId xx-xx-xx
```

Prometheus
-----------

### prometheus config

```yaml
scrape_configs:
  - job_name: 'nacosServiceExporter'
    scrape_interval: 30s
    metrics_path: '/metrics'
    static_configs:
      - targets:
        - 127.0.0.1:11111
```

### alert rules

up
endpointServiceCount
allServiceInstanceCount
allServiceCount
serviceInstanceCount



```yaml
groups:
- name: nacosServiceExporter
  rules:
  # alive 检测是否存活
  - alert: nacosServiceExporterDown
    expr: up{job="nacosServiceExporter"} != 1
    for: 30s
    annotations:
      summary: "实例 {{ $labels.instance }} 失联30秒."
      description: "任务 {{ $labels.job }} 下的实例 {{ $labels.instance }} 可能失联. 当前值 {{ $value }}"

#  - alert: endpointServiceCount
#    expr: endpointServiceCount == 0
#    for: 30s
#    annotations:
#      summary: "主机 {{ $labels.instance }} ."
#      description: "任务 {{ $labels.job }} 下的实例 {{ $labels.instance }} . 当前值 {{ $value }}"

# 当前数 / offset 数 -> 扩容为正, 缩容为负
  - alert: allServiceInstanceCount
    expr: allServiceInstanceCount{job="nacosServiceExporter"} / allServiceInstanceCount{job="nacosServiceExporter"} offset 1m <= 0.8
    for: 20s
    annotations:
      summary: "{{ $labels.instance }} 服务实例数下降较多."
      description: "{{ $labels.job }} 下实例 {{ $labels.instance }} 服务实例数下降较多,请检查是否正常. 当前值 {{ $value }}"
# 当前数 - offset 数 -> 新增为正, 减少为负

  - alert: allServiceCount
    expr: allServiceCount{job="nacosServiceExporter"} - allServiceCount{job="nacosServiceExporter"} offset 1m <= -1
    for: 30s
    annotations:
      summary: "{{ $labels.instance }} 服务数减少."
      description: "{{ $labels.job }} 下实例 {{ $labels.instance }} 服务数减少, 请确认是否有服务下线. 当前值 {{ $value }}"

  - alert: serviceInstanceCount
    expr: serviceInstanceCount{job="nacosServiceExporter"} == 0
    for: 10s
    annotations:
      summary: "{{ $labels.instance }} 无可用实例."
      description: "{{ $labels.job }} 下实例 {{ $labels.instance }} 服务 {{ $labels.service }} 无可用实例, 请确认是否有服务下线. 当前值 {{ $value }}"
```
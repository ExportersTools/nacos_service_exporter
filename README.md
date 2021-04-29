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


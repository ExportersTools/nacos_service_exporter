FROM registry.cn-hangzhou.aliyuncs.com/startops-base/suse:sle15.2.8.2.811

ADD dist/linux/nacosServiceExporter  /nacosServiceExporter

CMD ["/nacosServiceExporter"]
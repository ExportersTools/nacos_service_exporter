FROM registry.cn-hangzhou.aliyuncs.com/startops-base/suse:sle15.2.8.2.811

ADD /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ADD dist/linux/nacosServiceExporter  /nacosServiceExporter

EXPOSE 11111

CMD ["/nacosServiceExporter"]
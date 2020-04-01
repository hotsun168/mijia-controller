[中文文档](README_zh_CN.md)

# Mijia Controller

This program can view the web page which contains information and status of some sub devices which belongs to a Mijia gateway. It also supply the way that can controll or switch the status of some sub devices (such as the wall switch). There is a forwarder in this program which can transfer the gateway's broadcast message to the specified URL. 

All the API is standard HTTP API so that it can be integrated to another system conveniently. 

**注：获取网关信息时必须使用安卓手机的米家客户端进行获取。**

获取网关信息步骤：

![1.gif][1]

## 开发计划：

1. 状态有变化时通知。（已完成）
2. 对接更多设备。


## 构建：

1. git clone 本仓库。
2. 检查并go get未安装的包。
3. go build main.go bindata.go

## 运行：

1. 按照示例填写所需信息。其中，“appInfo”为网关信息，“subDeviceInfo”为子设备信息，可直接从APP复制并粘贴到“config.json”中。
2. 运行编译好的二进制文件。
3. 使用浏览器访问localhost的webServer.port，查看管理界面。


[1]: https://github.com/hotsun168/mijia-controller/raw/master/readme_images/1.gif

## Tips：
1. 如果部署在Linux上的服务收不到whois包的回应，可以尝试在/etc/sysctl.conf中加入如下配置（其中eth0代表接收UDP组播包的网卡）：
```
net.ipv4.all.rp_filter = 0
net.ipv4.eth0.rp_filter = 0
```

修改完配置后需要重启网络服务：
```
service networking restart
```

或直接重启机器。

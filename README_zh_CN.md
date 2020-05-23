# 米家控制器

本程序可以在页面上查看网关和属于该网关的某些子设备的信息和状态，并对某些子设备（如墙壁开关等）进行状态切换，还可网关广播的信息转发到指定URL，以实现第三方系统的对接。

所有的接口都是标准的HTTP接口，可以非常容易的进行第三方系统的接入。

内置WebSocket服务器，可将网关发出的所有数据包广播到WebSocket客户端。

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
3. 使用浏览器访问"http://localhost:${webServer.port}"，查看管理界面。

## WebSocket接入：
使用各种标准WebSocket客户端连接至“ws://localhost:${webServer.port}/ws”，即可收到JSON格式的数据包。

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

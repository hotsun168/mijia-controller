# 米家控制器

本程序可以在页面上查看网关和属于该网关的某些子设备的信息和状态，并对某些子设备（如墙壁开关等）进行状态切换。

所有的接口都是标准的HTTP接口，可以非常容易的进行第三方系统的接入。

界面：

![1.png][1]

## 开发计划：

1. 状态有变化时通知。

## 构建：

1. git clone 本仓库。
2. 检查并go get未安装的包。
3. go build main.go bindata.go

## 运行：

1. 按照示例填写所需信息（iOS版米家无法获取通信密码，请使用Android客户端获取）。其中，“appInfo”为网关信息，“subDeviceInfo”为子设备信息，可直接从APP复制并粘贴到“config.json”中。
2. 运行编译好的二进制文件。
3. 使用浏览器访问“http://localhost:%webServer.port%/”，查看管理界面。


[1]: https://github.com/hotsun168/mijia-controller/raw/master/readme_images/1.png
[2]: https://github.com/hotsun168/mijia-controller/raw/master/readme_images/2.gif
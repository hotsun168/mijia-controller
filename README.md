# 米家控制器

本程序可以在页面上查看网关和该网关下部分子设备的信息和状态，并对某些子设备（如墙壁开关等）进行状态切换。

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

1.

[1]: https://github.com/hotsun168/mijia-controller/raw/master/readme_images/1.png
[2]: https://github.com/hotsun168/mijia-controller/raw/master/readme_images/2.gif
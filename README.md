[中文文档](README_zh_CN.md)

# Mijia Controller

This program contains a web page which contains Mijia gateway's information and status of some sub devices which belongs to the gateway. It also supply the way that can control or switch the status of some sub devices (such as the wall switch). There is a forwarder in this program which can transfer the gateway's broadcast message to the specified URL. 

All the API is standard HTTP API so that it can be integrated to another system conveniently. 

It supply the WebSocket server to broadcasting all gateway's data packages to each WebSocket client.  

**Tip: You have to use Android phone to get the gateway information in the Mijia APP.**

The steps to get the gateway information:

![1.gif][1]

## Development plan: 

1. Make a notification when the status of any sub devices changed. (completed)
2. Support more sub devices. 


## Building steps: 

1. use "git clone" to clone this repository. 
2. use "go get" to setup all the package which are not installed. 
3. use "go build main.go bindata.go" to build this program. 

## Launching steps: 

1. Write the configration content into the config file. In the config file, "appInfo" is gateway information, "subDeviceInfo" is sub devices information. These two parts of config can be copied from the Mijia app and pasted into the "config.json" file. 
2. Launch the builded binary file.
3. visit "http://localhost:${webServer.port}", then you can see the dashboard page. 

## WebSocket:
Connect to "ws://localhost:${webServer.port}/ws" with WebSocket client, then it will received the data packages in JSON format. 

[1]: https://github.com/hotsun168/mijia-controller/raw/master/readme_images/1.gif

## Tips：
1. If this program cannot receive the "whois" response package, you can try to add the below config lines into "/etc/sysctl.conf", the "eth0" is the name of which network adapter is specified to receive the UDP multicasting package. 
```
net.ipv4.all.rp_filter = 0
net.ipv4.eth0.rp_filter = 0
```

You should restart the networking service when finishing changing the config. 
```
service networking restart
```

Also you can reboot the OS directly. 

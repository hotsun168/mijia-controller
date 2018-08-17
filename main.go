package main

import (
	"./config"
	"./device"
	"./webServer"
	"time"
)

func main() {
	config.LoadConfig("config.json")
	device.StartUdpWorker()
	webServer.StartWebServer(assetFS(), config.GetWebServerPort())
	for {
		time.Sleep(1 * time.Second)
	}
}

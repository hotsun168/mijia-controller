package main

import (
	"./config"
	"./device"
	"./notify"
	"./webServer"
	"time"
)

func main() {
	config.LoadConfig("config.json")
	notify.Init()
	device.StartUdpWorker()
	webServer.StartWebServer(assetFS(), config.GetWebServerPort())
	for {
		time.Sleep(1 * time.Second)
	}
}

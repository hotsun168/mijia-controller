package config

import (
	"../utils"
	"encoding/json"
	"io/ioutil"
)

//存放配置项。
var configInfo = ConfigInfo{}

//配置项。
type ConfigInfo struct {
	//Web服务器配置项。
	WebServer WebServer `json:"webServer"`
	//通知URL。
	NotifyUrls []string `json:"notifyUrls"`
	//网关配置项数组。
	Gateways []GatewayInfo `json:"gateways"`
}

//Web服务器配置项。
type WebServer struct {
	//Web端口。
	Port int `json:"port"`
	//Auth用户名。
	UserName string `json:"userName"`
	//Auth密码。
	Password string `json:"password"`
}

//网关配置项。
type GatewayInfo struct {
	//网关名称（自定义）。
	Name string `json:"name"`
	//网关Sid。
	Sid string `json:"sid"`
	//网关通信密码。
	Password string `json:"password"`
	//网关在米家APP中的信息。
	AppInfo interface{} `json:"appInfo"`
	//子设备信息。
	SubDeviceInfo []SubDeviceInfo `json:"subDeviceInfo"`
}

//子设备配置项。
type SubDeviceInfo struct {
	//子设备类型。
	Model string `json:"model"`
	//子设备ID。
	Did string `json:"did"`
	//子设备名称。
	Name string `json:"name"`
}

//执行加载配置项过程。
func LoadConfig(configFileName string) {
	data, err := ioutil.ReadFile(configFileName)
	utils.CheckErrorf(err, "Load config file %s error! ", configFileName)
	err = json.Unmarshal(data, &configInfo)
	utils.CheckError(err, "Parse config JSON error! ")
}

//获取Web服务器端口配置项。
func GetWebServerPort() int {
	return configInfo.WebServer.Port
}

//根据网关Sid，获取网关名称。
func GetGatewayNameBySid(sid string) string {
	for _, g := range configInfo.Gateways {
		if g.Sid == sid {
			return g.Name
		}
	}
	return ""
}

//根据子设备Sid，获取子设备名称。
func GetSubDeviceNameBySid(deviceId string) string {
	for _, g := range configInfo.Gateways {
		for _, d := range g.SubDeviceInfo {
			if d.Did == "lumi."+deviceId {
				return d.Name
			}
		}
	}
	return ""
}

//获取Web服务器基础验证（Auth）的用户名和密码。
func GetWebServerUserNameAndPassword() (string, string) {
	return configInfo.WebServer.UserName, configInfo.WebServer.Password
}

//根据子设备Sid，获取该子设备所属的网关的密码和Sid。
func GetGatewayPasswordAndSid(subDeviceSid string) (string, string) {
	for _, g := range configInfo.Gateways {
		for _, d := range g.SubDeviceInfo {
			if d.Did == "lumi."+subDeviceSid {
				return g.Password, g.Sid
			}
		}
	}
	return "", ""
}

//获取通知URL，可多个。
func GetNotifyUrls() []string {
	urls := make([]string, 0)
	for _, u := range configInfo.NotifyUrls {
		urls = append(urls, u)
	}
	return urls
}

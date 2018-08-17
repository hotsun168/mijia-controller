package device

import (
	"../config"
	"../crypto"
	"../utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	//UDP组播地址。
	broadcastIp = "224.0.0.50"
	//UDP组播端口。
	broadcastPort = 4321
	//UDP监听端口。
	listenPort = 9898
	//UDP数据报最大长度。
	maxDatagramSize = 8192
	//Whois包内容。
	whoisPackage = "{\"cmd\":\"whois\"}"
	//GetIdList包内容。
	getIdListPackage = "{\"cmd\":\"get_id_list\"}"
	//Read包内容。
	readPackage = "{\"cmd\":\"read\",\"sid\":\"%s\"}"
	//Write包内容。
	writePackage = "{\"cmd\":\"write\",\"model\":\"%s\",\"sid\":\"%s\",\"data\":\"{\\\"%s\\\":\\\"%s\\\", \\\"key\\\":\\\"%s\\\"}\"}"
)

//网关。
type Gateway struct {
	Sid           string
	Name          string
	ShortId       int
	Model         string
	//协议版本。
	ProtoVersion  string
	//IP地址。
	Ip            string
	//通信Token，会依据心跳包或Iam包随时更新。
	Token         string
	//通信端口。
	Port          int
	//照度。
	Illumination  int
	//夜灯颜色。
	Rgb           int
	//子设备Sid数组。
	SubDeviceSids []string
}

//墙壁开关单火线单键版。
type CtrlNeutral1 struct {
	Sid      string
	Name     string
	ShortId  int
	Model    string
	//电池电压（始终为3300）。
	Voltage  int
	//开关接通状态。
	Channel0 bool
}

//墙壁开关单火线双键版。
type CtrlNeutral2 struct {
	Sid      string
	Name     string
	ShortId  int
	Model    string
	//电池电压（始终为3300）。
	Voltage  int
	//左键接通状态。
	Channel0 bool
	//右键接通状态。
	Channel1 bool
}

//无线开关
type Switch struct {
	Sid     string
	Name    string
	ShortId int
	Model   string
	//电池电压。
	Voltage int
}

//无线开关贴墙式单键版。
type D86sw1 struct {
	Sid     string
	Name    string
	ShortId int
	Model   string
	//电池电压。
	Voltage int
}

//无线开关贴墙式双键版。
type D86sw2 struct {
	Sid     string
	Name    string
	ShortId int
	Model   string
	//电池电压。
	Voltage int
}

//门窗传感器。
type DoorMagnet struct {
	Sid     string
	Name    string
	ShortId int
	Model   string
	//电池电压。
	Voltage int
	//门窗开合状态。
	IsOpen  bool
}

//人体传感器。
type Motion struct {
	Sid             string
	Name            string
	ShortId         int
	Model           string
	//电池电压。
	Voltage         int
	//未检测到移动的延迟秒数。
	NoMotionSeconds int
	//是否检测到移动。
	IsMotion        bool
}

//Device操作锁。
var lock sync.Mutex

//存放网关与子设备信息，key为网关Sid，value为Gateway实例。
var devices = make(map[string]interface{})

//UDP套接字。
var conn *net.UDPConn

//开启UDP工作Routine。
func StartUdpWorker() {
	conn = ping()
	go recvData(conn)
	go sendWhoisPackage(conn)
}

//连接并暂存UDP套接字。
func ping() *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", broadcastIp, listenPort))
	utils.CheckError(err, "UDP Address build error! ")
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	utils.CheckError(err, "UDP worker listening error! ")
	conn.SetReadBuffer(maxDatagramSize)
	return conn
}

//保持每隔一段时间发送Whois包。
func sendWhoisPackage(conn *net.UDPConn) {
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", broadcastIp, broadcastPort))
	utils.CheckError(err, "UDP worker listening error! ")
	for {
		len, err := conn.WriteToUDP([]byte(whoisPackage), raddr)
		utils.ShowError(err, "Send Whois package error! ")
		log.Printf("Sent whois package. Length is %d. \n", len)
		time.Sleep(2 * 60 * time.Second)
	}
}

func recvData(conn *net.UDPConn) {
	for {
		buf := make([]byte, maxDatagramSize)
		len, remoteAddr, err := conn.ReadFromUDP(buf)
		utils.ShowError(err, "Receive data error! ")
		log.Println(len, "bytes read from", remoteAddr)
		data := string(buf[:len])
		dispatchData(data)
	}
}

func dispatchData(data string) {
	log.Println("Received package: " + data)
	v := make(map[string]interface{})
	json.Unmarshal([]byte(data), &v)
	sid := utils.ParseString(v["sid"])
	cmd := v["cmd"].(string)
	if cmd == "heartbeat" {
		processHeartbeatPackage(v)
	} else if cmd == "iam" {
		processIamPackage(v)
		sendGetIdListPackage(sid)
	} else if cmd == "report" {
		processReportPackage(v)
	} else if cmd == "get_id_list_ack" {
		processGetIdListAckPackage(v)
		sendReadPackage(sid)
	} else if cmd == "read_ack" {
		processReadAckPackage(v)
	}
}

func sendGetIdListPackage(sid string) {
	ip, port := getGatewayIpAndPort(sid)
	if ip != "" {
		raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
		utils.CheckError(err, "Make remote address error! ")
		len, err := conn.WriteToUDP([]byte(getIdListPackage), raddr)
		utils.ShowError(err, "Send get id list package error! ")
		log.Printf("Sent get id list package. Length is %d. \n", len)
	}
}

func sendReadPackage(sid string) {
	subDeviceSids := getGatewaySubDeviceSids(sid)
	if len(subDeviceSids) == 0 {
		return
	}
	ip, port := getGatewayIpAndPort(sid)
	if ip != "" {
		for _, subSid := range subDeviceSids {
			pkg := fmt.Sprintf(readPackage, subSid)
			raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
			utils.CheckError(err, "Make remote address error! ")
			len, err := conn.WriteToUDP([]byte(pkg), raddr)
			utils.ShowError(err, "Send read package error! ")
			log.Printf("Sent read package. Length is %d. \n", len)
		}
	}
}

func processHeartbeatPackage(pkg map[string]interface{}) {
	createOrUpdateDeviceStatus(pkg)
}

func processIamPackage(pkg map[string]interface{}) {
	createOrUpdateDeviceStatus(pkg)
}

func processReportPackage(pkg map[string]interface{}) {
	createOrUpdateDeviceStatus(pkg)
}

func processGetIdListAckPackage(pkg map[string]interface{}) {
	lock.Lock()
	defer lock.Unlock()
	sid := pkg["sid"].(string)
	device := devices[sid]
	if device != nil {
		dataString := utils.ParseStringDef(pkg["data"], "[]")
		subSids := make([]string, 0)
		json.Unmarshal([]byte(dataString), &subSids)
		switch device.(type) {
		case *Gateway:
			{
				gateway := device.(*Gateway)
				gateway.Token = pkg["token"].(string)
				gateway.SubDeviceSids = subSids
				break
			}
		}
	}
}

func processReadAckPackage(pkg map[string]interface{}) {
	createOrUpdateDeviceStatus(pkg)
}

func getGatewayIpAndPort(sid string) (string, int) {
	lock.Lock()
	defer lock.Unlock()
	device := devices[sid]
	if device != nil {
		gateway := device.(*Gateway)
		return gateway.Ip, gateway.Port
	}
	return "", 0
}

func getGatewaySubDeviceSids(sid string) []string {
	lock.Lock()
	defer lock.Unlock()
	device := devices[sid]
	if device != nil {
		gateway := device.(*Gateway)
		return gateway.SubDeviceSids
	}
	return make([]string, 0)
}

//从接收到的数据包，获取设备信息，并创建或暂存到devices中。
func createOrUpdateDeviceStatus(pkg map[string]interface{}) {
	lock.Lock()
	defer lock.Unlock()
	sid := pkg["sid"].(string)
	model := pkg["model"].(string)
	device := devices[sid]
	if device == nil {
		if model == "gateway" {
			device = &Gateway{Sid: sid, Name: config.GetGatewayNameBySid(sid), ShortId: utils.ParseInt(pkg["short_id"]), Model: pkg["model"].(string), SubDeviceSids: make([]string, 0)}
		} else if model == "ctrl_neutral1" {
			device = &CtrlNeutral1{Sid: sid, Name: config.GetSubDeviceNameBySid(sid), ShortId: utils.ParseInt(pkg["short_id"]), Model: pkg["model"].(string)}
		} else if model == "ctrl_neutral2" {
			device = &CtrlNeutral2{Sid: sid, Name: config.GetSubDeviceNameBySid(sid), ShortId: utils.ParseInt(pkg["short_id"]), Model: pkg["model"].(string)}
		} else if model == "switch" {
			device = &Switch{Sid: sid, Name: config.GetSubDeviceNameBySid(sid), ShortId: utils.ParseInt(pkg["short_id"]), Model: pkg["model"].(string)}
		} else if model == "86sw1" {
			device = &D86sw1{Sid: sid, Name: config.GetSubDeviceNameBySid(sid), ShortId: utils.ParseInt(pkg["short_id"]), Model: pkg["model"].(string)}
		} else if model == "86sw2" {
			device = &D86sw2{Sid: sid, Name: config.GetSubDeviceNameBySid(sid), ShortId: utils.ParseInt(pkg["short_id"]), Model: pkg["model"].(string)}
		} else if model == "magnet" {
			device = &DoorMagnet{Sid: sid, Name: config.GetSubDeviceNameBySid(sid), ShortId: utils.ParseInt(pkg["short_id"]), Model: pkg["model"].(string)}
		} else if model == "motion" {
			device = &Motion{Sid: sid, Name: config.GetSubDeviceNameBySid(sid), ShortId: utils.ParseInt(pkg["short_id"]), Model: pkg["model"].(string)}
		}
		devices[sid] = device
	}
	var data map[string]interface{}
	utils.ToJson(pkg["data"], &data)
	switch device.(type) {
	case *Gateway:
		{
			gateway := device.(*Gateway)
			//某些字段可能直接出现在包中，也可能出现在包的data字段中。
			utils.ParseStringAndSetWhenNotBlank(&gateway.Token, pkg["token"])
			utils.ParseStringAndSetWhenNotBlank(&gateway.Ip, pkg["ip"])
			utils.ParseStringAndSetWhenNotBlank(&gateway.ProtoVersion, pkg["proto_version"])
			utils.ParseIntAndSetWhenNotZero(&gateway.Port, pkg["port"])
			utils.ParseStringAndSetWhenNotBlank(&gateway.Ip, data["ip"])
			utils.ParseIntAndSetWhenNotZero(&gateway.Rgb, data["rgb"])
			utils.ParseIntAndSetWhenNotZero(&gateway.Illumination, data["illumination"])
			break
		}
	case *CtrlNeutral1:
		{
			ctrlNeutral1 := device.(*CtrlNeutral1)
			utils.ParseIntAndSetWhenNotZero(&ctrlNeutral1.Voltage, data["voltage"])
			utils.ParseChannelStatusAndSet(&ctrlNeutral1.Channel0, data["channel_0"])
			break
		}
	case *CtrlNeutral2:
		{
			ctrlNeutral2 := device.(*CtrlNeutral2)
			utils.ParseIntAndSetWhenNotZero(&ctrlNeutral2.Voltage, data["voltage"])
			utils.ParseChannelStatusAndSet(&ctrlNeutral2.Channel0, data["channel_0"])
			utils.ParseChannelStatusAndSet(&ctrlNeutral2.Channel1, data["channel_1"])
			break
		}
	case *Switch:
		{
			switchDevice := device.(*Switch)
			utils.ParseIntAndSetWhenNotZero(&switchDevice.Voltage, data["voltage"])
			break
		}
	case *D86sw1:
		{
			d86sw1 := device.(*D86sw1)
			utils.ParseIntAndSetWhenNotZero(&d86sw1.Voltage, data["voltage"])
			break
		}
	case *D86sw2:
		{
			d86sw2 := device.(*D86sw2)
			utils.ParseIntAndSetWhenNotZero(&d86sw2.Voltage, data["voltage"])
			break
		}
	case *DoorMagnet:
		{
			doorMagnet := device.(*DoorMagnet)
			utils.ParseIntAndSetWhenNotZero(&doorMagnet.Voltage, data["voltage"])
			utils.ParseDoorMagnetStatusAndSet(&doorMagnet.IsOpen, data["status"])
			break
		}
	case *Motion:
		{
			motion := device.(*Motion)
			utils.ParseIntAndSetWhenNotZero(&motion.Voltage, data["voltage"])
			status := utils.ParseString(data["status"])
			noMotion := utils.ParseString(data["no_motion"])
			if status == "motion" {
				motion.IsMotion = true
				motion.NoMotionSeconds = 0
			} else if noMotion != "" {
				motion.IsMotion = false
				motion.NoMotionSeconds = utils.ParseInt(noMotion)
			}
			break
		}
	}
}

func GetAllDevices() string {
	lock.Lock()
	defer lock.Unlock()
	data, err := json.Marshal(devices)
	utils.CheckError(err, "Devices to JSON error! ")
	return string(data)
}

func SetSwitchStatus(sid string, btnIndex int, on bool) bool {
	//准备开关、状态、token、密码、sid等参数。
	btnStr := ""
	if btnIndex == 0 {
		btnStr = "channel_0"
	} else {
		btnStr = "channel_1"
	}
	onStr := ""
	if on {
		onStr = "on"
	} else {
		onStr = "off"
	}
	lock.Lock()
	defer lock.Unlock()
	model := ""
	device := devices[sid]
	switch device.(type) {
	case *CtrlNeutral1:
		{
			ctrlNeutral1 := device.(*CtrlNeutral1)
			model = ctrlNeutral1.Model
			break
		}
	case *CtrlNeutral2:
		{
			ctrlNeutral2 := device.(*CtrlNeutral2)
			model = ctrlNeutral2.Model
			break
		}
	}
	password, gatewaySid := config.GetGatewayPasswordAndSid(sid)
	if password == "" || gatewaySid == "" {
		return false
	}
	gatewayDevice := devices[gatewaySid]
	if gatewayDevice == nil {
		return false
	}
	var gateway *Gateway
	switch gatewayDevice.(type) {
	case *Gateway:
		{
			gateway = gatewayDevice.(*Gateway)
			break
		}
	}
	token := gateway.Token
	if token == "" {
		return false
	}
	//基于密码创建aes-128-cbc加密器。
	encryptedBytes, err := crypto.Encrypt(token, []byte(password))
	utils.CheckError(err, "Encrypt key error! ")
	//取密文字节数组对应的16进制字面值，作为write包的key。
	encrypted := hex.EncodeToString(encryptedBytes)
	//组装并发送write包。
	pkg := fmt.Sprintf(writePackage, model, sid, btnStr, onStr, encrypted)
	raddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", gateway.Ip, gateway.Port))
	utils.CheckError(err, "Make remote address error! ")
	len, err := conn.WriteToUDP([]byte(pkg), raddr)
	utils.ShowError(err, "Send set switch status package error! ")
	log.Printf("Sent set switch status package. Length is %d. \n", len)
	return true
}

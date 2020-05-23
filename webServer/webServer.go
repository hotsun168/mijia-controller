package webServer

import (
	"../config"
	"../device"
	"../utils"
	"../webSocket"
	"encoding/base64"
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

func StartWebServer(fs *assetfs.AssetFS, port int) {
	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		go webSocket.HandleWebSocket(c)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !checkAuthentication(w, r) {
			return
		}
		path := r.URL.Path
		if path == "/" || strings.Contains(path, ".html") {
			//fsServer
			if path == "/" {
				path = "/index.html"
			}
			//路由。
			if strings.HasSuffix(path, ".html") {
				data, _ := fs.Asset("static" + path)
				html := string(data)
				w.Header().Add("Content-Type", "text/html; charset=UTF-8")
				w.Write([]byte(html))
			}
		} else if strings.HasSuffix(path, ".do") {
			if path == "/getDeviceStatus.do" {
				w.Header().Add("Content-Type", "application/json; charset=UTF-8")
				w.Write([]byte(device.GetAllDevices()))
			} else if path == "/setSwitchStatus.do" {
				r.ParseForm()
				sid := r.Form["sid"][0]
				index := utils.ParseInt(r.Form["index"][0])
				on := utils.ParseBool(r.Form["on"][0])
				device.SetSwitchStatus(sid, index, on)
			}
		} else {
			if strings.Contains(path, ".js") {
				w.Header().Add("Content-Type", "text/javascript; charset=UTF-8")
			} else if strings.Contains(path, ".css") {
				w.Header().Add("Content-Type", "text/css; charset=UTF-8")
			}

			data, _ := fs.Asset("static" + path)
			w.Write(data)
		}
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	utils.CheckError(err, "Start web server error! ")
}

func checkAuthentication(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		http.Error(w, "Not authorized", 401)
		return false
	}
	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		http.Error(w, err.Error(), 401)
		return false
	}
	arr := strings.SplitN(string(b), ":", 2)
	if len(arr) != 2 {
		http.Error(w, "Not authorized", 401)
		return false
	}
	userName, password := config.GetWebServerUserNameAndPassword()
	if arr[0] != userName || arr[1] != password {
		http.Error(w, "Not authorized", 401)
		return false
	}
	return true
}

package webSocket

import (
	"container/list"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"sync"
)

//WebSocket连接信息List锁。
var webSocketClientsLock sync.Mutex

//WebSocket连接信息。
var webSocketClients = list.List{}

func HandleWebSocket(ws *websocket.Conn) {
	webSocketClientsLock.Lock()
	element := webSocketClients.PushBack(ws)
	log.Printf("[WebSocket] new WebSocket client received. client count: %d\n", webSocketClients.Len())
	webSocketClientsLock.Unlock()
	for {
		b, m := readMessageFromWebSocket(ws)
		if b {
			if m != "" {
				webSocketMessageHandler(ws, m)
			}
		} else {
			webSocketClientsLock.Lock()
			webSocketClients.Remove(element)
			log.Printf("[WebSocket] one WebSocket client disposed. client count: %d\n", webSocketClients.Len())
			webSocketClientsLock.Unlock()
			break
		}
	}
}

func readMessageFromWebSocket(ws *websocket.Conn) (bool, string) {
	var m string
	defer func() (bool, string) {
		e := recover()
		if e != nil {
			log.Println("read message from WebSocket failed, closing...")
			ws.Close()
			return false, ""
		}
		return true, m
	}()
	t, p, err := ws.ReadMessage()
	if err != nil {
		log.Println("read message from WebSocket failed, closing...")
		ws.Close()
		return false, ""
	}
	if t == websocket.TextMessage {
		m = string(p)
		m = strings.TrimSpace(m)
		return true, m
	}
	return true, ""
}

func writeMessageToWebSocket(ws *websocket.Conn, s string) bool {
	defer func() bool {
		e := recover()
		if e != nil {
			log.Println("write message from WebSocket failed, closing...")
			ws.Close()
			return false
		}
		return true
	}()
	err := ws.WriteMessage(websocket.TextMessage, []byte(s))
	if err != nil {
		ws.Close()
		log.Println("write message from WebSocket failed, closing...")
	}
	return true
}

//WebSocket消息处理器。目前无处理逻辑。
func webSocketMessageHandler(ws *websocket.Conn, message string) {

}

//对WebSocket广播。
func WebSocketBroadcast(message string) {
	webSocketClientsLock.Lock()
	total := webSocketClients.Len()
	count := 0
	for e := webSocketClients.Front(); e != nil; e = e.Next() {
		ws := e.Value.(*websocket.Conn)
		b := writeMessageToWebSocket(ws, message)
		if b {
			count++
		}
	}
	webSocketClientsLock.Unlock()
	if count != total {
		log.Printf("[WARN] broadcast not completely successful! size: %d -> succeed: %d\n", total, count)
	}
}

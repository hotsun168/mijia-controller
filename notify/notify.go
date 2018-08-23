package notify

import (
	"../config"
	"../utils"
	"container/list"
	"sync"
	"time"
)

var queue = list.New()
var mutex sync.Mutex
var urls []string

func Init() {
	urls = config.GetNotifyUrls()
	go executeNotify()
}

func executeNotify() {
	for {
		if queue.Len() > 0 {
			mutex.Lock()
			message := queue.Front()
			if message != nil {
				queue.Remove(message)
				doCall(message.Value.(string))
			}
			mutex.Unlock()
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func PushMessage(message string) {
	mutex.Lock()
	queue.PushBack(message)
	mutex.Unlock()
}

func doCall(message string) {
	for _, u := range urls {
		utils.HttpPostJson(u, message)
	}
}

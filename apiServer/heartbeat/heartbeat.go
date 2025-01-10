package heartbeat

import (
	"fmt"
	"goss/pkg/rabbitmq"
	"os"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

func ListenHeartBeat() {
	queue := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer queue.Close()
	// 绑定 交换机 apiServers
	queue.Bind("apiServers")
	consume := queue.Consume()
	// 轮询 移除超时节点
	go removeExpireDataServer()

	for message := range consume {
		// strconv.Unqupte 移除 引号和转义字符
		dataServer, err := strconv.Unquote(string(message.Body))
		if err != nil {
			fmt.Println(err)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

// 移除超时的数据服务节点
func removeExpireDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for server, timer := range dataServers {
			if timer.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, server)
			}
		}
		mutex.Unlock()
	}
}

// 获取所有存活的数据服务节点
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	dataServer := make([]string, 0)
	for server, _ := range dataServers {
		dataServer = append(dataServer, server)
	}
	return dataServer
}

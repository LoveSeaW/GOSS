package locate

import (
	"fmt"
	"goss/pkg/rabbitmq"
	"goss/pkg/types"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var objects = make(map[string]int)

var mutex sync.RWMutex

func Locate(hash string) int {
	mutex.RLock()
	id, ok := objects[hash]
	mutex.RUnlock()
	if !ok {
		return -1
	}

	return id
}

func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
}

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

// 如果 消息队列传递对象文件信息，在本节点，则使用消息单发通知接口服务节点
func StartLocate() {
	queue := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer queue.Close()
	queue.Bind("dataServers")
	consume := queue.Consume()
	for message := range consume {
		hash, err := strconv.Unquote(string(message.Body))
		if err != nil {
			log.Println(err)
			break
		}

		fmt.Println("prepare locate:", hash)
		id := Locate(hash)
		if id != -1 {
			queue.Send(message.ReplyTo, types.LocateMessage{os.Getenv("LISTEN_ADDRESS"), id})
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/")
	for i := range files {
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files[i])
		}
		hash := file[0]
		id, err := strconv.Atoi(file[1])
		if err != nil {
			panic(err)
		}
		objects[hash] = id
	}
}

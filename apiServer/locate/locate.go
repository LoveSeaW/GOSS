package locate

import (
	"encoding/json"
	"fmt"
	"goss/pkg/rabbitmq"
	"goss/pkg/rs"
	"goss/pkg/types"
	"time"
)

// 定位 对象文件 所在 数据节点地址
func Locate(name string) (locateInfo map[int]string) {
	queue := rabbitmq.New("RABBITMQ_SERVER")
	// 向 dataServers 节点发送消息， 绑定 dataServers的 节点服务收到
	queue.Publish("dataServers", name)
	// 接收匿名队列的消息， 当 数据服务的 节点接收到 发送的消息时，判断 name 是否为自己， 为自己则将自身节点地址信息作为消息返回
	consume := queue.Consume()
	go func() {
		time.Sleep(time.Second)
		queue.Close()
	}()

	locateInfo = make(map[int]string)
	for i := 0; i < rs.AllShard; i++ {
		// 阻塞接收 消息队列消息， 接收数量到节点限制
		message := <-consume
		fmt.Println("receive message: ", message.Body)
		if len(message.Body) == 0 {
			return
		}
		var info types.LocateMessage
		json.Unmarshal(message.Body, &info)
		fmt.Println("locate message")
		locateInfo[info.Id] = info.Address
	}
	return
}

// 如果存储该对象的节点 超过数据分片， 则可以被修复，允许存储
func Exist(name string) bool {
	return len(Locate(name)) >= rs.DataShard
}

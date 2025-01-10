package heartbeat

import (
	"goss/pkg/rabbitmq"
	"os"
	"time"
)

// 数据节点 定期向 apiServers 交换机发送自身服务信息， apiServer 服务可以获取这些节点信息后进行操作
func StartHeartBeat() {
	queue := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer queue.Close()
	for {
		queue.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}

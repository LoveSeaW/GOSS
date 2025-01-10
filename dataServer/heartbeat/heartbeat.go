package heartbeat

import (
	"goss/pkg/rabbitmq"
	"os"
	"time"
)

func StartHeartBeat() {
	queue := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer queue.Close()
	for {
		queue.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}

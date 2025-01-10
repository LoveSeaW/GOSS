package main

import (
	"goss/dataServer/heartbeat"
	"goss/dataServer/locate"
	"goss/dataServer/objects"
	"goss/dataServer/temp"
	"log"
	"net/http"
	"os"
)

func main() {
	// 收集 存储文件信息
	locate.CollectObjects()
	// 开启心跳检测
	go heartbeat.StartHeartBeat()
	// 向 rabbitmq 注册节点信息
	go locate.StartLocate()

	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)

	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}

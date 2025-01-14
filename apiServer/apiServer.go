package main

import (
	"goss/apiServer/heartbeat"
	"goss/apiServer/locate"
	"goss/apiServer/objects"
	"goss/apiServer/temp"
	version "goss/apiServer/versions"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.ListenHeartBeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", version.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}

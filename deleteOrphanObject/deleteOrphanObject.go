package main

import (
	"goss/pkg/es7"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 删除没有元数据引用的对象数据的工具
func main() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		hashInMetaData, err := es7.HasHash(hash)
		if err != nil {
			log.Println(err)
			return
		}
		if !hashInMetaData {
			del(hash)
		}
	}
}

func del(hash string) {
	log.Println("delete:", hash)
	url := "http://" + os.Getenv("LISTEN_ADDRESS") + "/objects/" + hash
	request, _ := http.NewRequest("DELETE", url, nil)
	client := http.Client{}
	client.Do(request)
}

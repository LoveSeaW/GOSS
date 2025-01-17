package temp

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// 当接口服务数据验证一致后，调用将临时文件转正
func put(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	temp_info, err := readFromFile(uuid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusFound)
		return
	}

	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	dataFile := infoFile + ".dat"
	file, err := os.Open(dataFile)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	actual := info.Size()
	os.Remove(infoFile)
	if actual != temp_info.Size {
		os.Remove(dataFile)
		log.Println("actual size mismatch, expect ", temp_info.Size, " actual", actual)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commitTempObject(dataFile, temp_info)
}

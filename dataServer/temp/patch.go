package temp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func patch(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	temp_info, err := readFromFile(uuid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	filePath := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	dataFile := filePath + ".dat"
	file, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	actual := info.Size()
	if actual > temp_info.Size {
		os.Remove(dataFile)
		os.Remove(filePath)
		log.Println("actual size, ", actual, "exceeds", temp_info.Size)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

// 获取文件夹 temp 下对应 uuid 文件的内容信息
func readFromFile(uuid string) (*tempInfo, error) {
	file, err := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)
	var info tempInfo
	json.Unmarshal(bytes, &info)
	return &info, nil
}

package temp

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

// POST /temp/<hash> Size: < 需要缓存的对象大小>
// return uuid
// 创建临时对象
// 将这个对象写入 /temp/<uuid>文件中
// 再创建一个 /temp/<uuid>.dat 文件
// 返回uuid
func post(w http.ResponseWriter, r *http.Request) {
	// 调用系统命令 uuidgen 生成一个 uuid， 捕获其输出
	output, _ := exec.Command("uuidgen").Output()
	uuid := strings.TrimSuffix(string(output), "\n")
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	temp := tempInfo{uuid, name, size}
	err = temp.writeToFile()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + temp.Uuid + ".dat")
	w.Write([]byte(uuid))
}

func (t *tempInfo) writeToFile() error {
	file, err := os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid)
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, _ := json.Marshal(t)
	file.Write(bytes)
	return nil
}

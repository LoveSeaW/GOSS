package temp

import (
	"net/http"
	"os"
	"strings"
)

// 在接口服务数据校验不一致时，调用数据服务 del 删除临时文件
func del(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	dataFile := infoFile + ".dat"
	os.Remove(infoFile)
	os.Remove(dataFile)
	w.WriteHeader(http.StatusOK)
}

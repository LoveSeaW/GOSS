package objects

import (
	"goss/pkg/es7"
	"log"
	"net/http"
	"strings"
)

// 接口服务-del 在接口服务数据校验不一致时，删除对象文件
func del(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version, err := es7.SearchLatestVersion(name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 表示删除的特殊版本
	err = es7.PutMetadata(name, version.Version+1, 0, "")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

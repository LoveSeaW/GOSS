package objects

import (
	"goss/apiServer/heartbeat"
	"goss/apiServer/locate"
	"goss/pkg/es7"
	"goss/pkg/rs"
	"goss/pkg/utils"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// 上传临时对象存储
func post(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 获取文件大小
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if locate.Exist(url.PathEscape(hash)) {
		err = es7.AddVersion(name, hash, size)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	// 随机选择 数据服务节点
	dataServers := heartbeat.ChooseRandomDataServers(rs.AllShard, nil)

	if len(dataServers) != rs.AllShard {
		log.Println("cannot find enough dataServer")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	// 使用 rs纠错码 将 对象 生成数据碎片分散到不同的 数据服务节点 中
	stream, err := rs.NewRSResumablePutStream(dataServers, name, url.PathEscape(hash), size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("location", "/temp/"+url.PathEscape(stream.ToToken()))
	w.WriteHeader(http.StatusCreated)
}

package objects

import (
	"compress/gzip"
	"fmt"
	"goss/pkg/es7"
	"goss/pkg/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// 下载对象
func get(w http.ResponseWriter, r *http.Request) {
	// r.URL.EscapedPath() 进行URL解码，获取URL
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionId := r.URL.Query()["version"]
	version := 0
	var err error
	if len(versionId) != 0 {
		version, err = strconv.Atoi(versionId[0])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	// 从 es 获取 文件信息
	meta, err := es7.GetMetadata(name, version)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	hash := url.PathEscape(meta.Hash)
	stream, err := GetStream(hash, meta.Size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 数据偏移量
	offset := utils.GetOffsetFromHeader(r.Header)
	if offset != 0 {
		stream.Seek(offset, io.SeekCurrent)
		// 获取 数据读取范围
		w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", offset, meta.Size-1, meta.Size))
		w.WriteHeader(http.StatusPartialContent)
	}

	acceptGzip := false
	encoding := r.Header["Accept-Encoding"]
	for i := range encoding {
		if encoding[i] == "gzip" {
			acceptGzip = true
			break
		}
	}
	// 将对象写入ResponseWriter中
	if acceptGzip {
		w.Header().Set("content-encoding", "gzip")
		w2 := gzip.NewWriter(w)
		io.Copy(w2, stream)
		w2.Close()
	} else {
		io.Copy(w, stream)
	}
	stream.Close()
}

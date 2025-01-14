package temp

import (
	"fmt"
	"goss/apiServer/locate"
	"goss/pkg/es7"
	"goss/pkg/rs"
	"goss/pkg/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := rs.PutStreamFromToken(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	fmt.Println("current:", current)
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	offset := utils.GetOffsetFromHeader(r.Header)
	fmt.Println("offset:", offset)
	if current != offset {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	bytes := make([]byte, rs.BlockSize)
	for {
		content, err := io.ReadFull(r.Body, bytes)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		current += int64(content)
		if current > stream.Size {
			stream.Commit(false)
			log.Println("resumable put exceed size")
			w.WriteHeader(http.StatusForbidden)
		}
		if content != rs.BlockSize && current != stream.Size {
			return
		}
		stream.Write(bytes[:content])
		if current == stream.Size {
			stream.Flush()
			// 获取 不同数据节点临时对象数据碎片，计算出hash值比较原对象的hash分片
			getStream, _ := rs.NewRSResumableGetStream(stream.Servers, stream.Uuids, stream.Size)
			hash := url.PathEscape(utils.CalculateHash(getStream))
			if hash != stream.Hash {
				stream.Commit(false)
				log.Println("resumable put done but hash mismatch")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if locate.Exist(url.PathEscape(hash)) {
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			err = es7.AddVersion(stream.Name, stream.Hash, stream.Size)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}

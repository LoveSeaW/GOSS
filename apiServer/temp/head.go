package temp

import (
	"fmt"
	"goss/pkg/rs"
	"log"
	"net/http"
	"strings"
)

// 访问 token 获取当前已经上传了多少数据
func head(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := rs.PutStreamFromToken(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-length", fmt.Sprintf("%d", current))
}

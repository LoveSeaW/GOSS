package objects

import (
	"fmt"
	"goss/apiServer/locate"
	"goss/pkg/utils"
	"io"
	"net/http"
	"net/url"
)

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}

	stream, err := putStream(url.PathEscape(hash), size)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// 将输入流（r）的数据同时传递到两个地方：一个是它原本的目标（即 r 本身），另一个是你指定的 stream（它通常是一个 Writer）。
	reader := io.TeeReader(r, stream)
	// 当reader被读取时，同时也会写入stream即r.Body会写入stream
	caculate := utils.CalculateHash(reader)
	if caculate != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch,caculated=%s,requestd=%s", caculate, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}

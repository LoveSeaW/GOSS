package objectstream

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type TempPutStream struct {
	Server string // 数据节点服务
	Uuid   string // 数据节点中临时文件的uuid
}

// 在 数据服务节点上创建临时文件，获取临时文件的uuid
func NewTempPutStream(server, object string, size int64) (*TempPutStream, error) {
	request, err := http.NewRequest("POST", "http://"+server+"/temp/"+object, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	uuid, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &TempPutStream{Server: server, Uuid: string(uuid)}, nil
}

// 将 对象数据碎片 写入临时文件 数据服务
func (w *TempPutStream) Write(p []byte) (n int, err error) {
	request, err := http.NewRequest("PATCH", "http://"+w.Server+"/temp"+w.Uuid, strings.NewReader(string(p)))
	if err != nil {
		return 0, err
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer return http code %d", response.StatusCode)
	}
	return len(p), err
}

// Put 方法 将临时对象转正
// Delete 方法 将临时对象删除
func (w *TempPutStream) Commit(good bool) {
	method := http.MethodDelete
	if good {
		method = http.MethodPut
	}
	request, _ := http.NewRequest(method, "http://"+w.Server+"/temp/"+w.Uuid, nil)
	client := http.Client{}
	client.Do(request)
}

func NewTempGetStream(server, uuid string) (*GetStream, error) {
	return newGetStream("http://" + server + "/temp/" + uuid)
}

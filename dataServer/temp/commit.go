package temp

import (
	"compress/gzip"
	"goss/dataServer/locate"
	"goss/pkg/utils"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func (t *tempInfo) hash() string {
	str := strings.Split(t.Name, ".")
	return str[0]
}

func (t *tempInfo) id() int {
	str := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(str[1])
	return id
}

func commitTempObject(dataFile string, info *tempInfo) {
	file, _ := os.Open(dataFile)
	defer file.Close()

	calculated := url.PathEscape(utils.CalculateHash(file))
	// 将文件的读写指针重置到文件开头，确保后续操作从文件开始位置读取内容
	file.Seek(0, io.SeekStart)
	path, _ := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + info.Name + "." + calculated)
	// 创建一个 gzip 压缩器，将目标文件路径作为输出，文件内容将以压缩形式存储
	writer := gzip.NewWriter(path)
	io.Copy(writer, file)
	writer.Close()
	// 删除原始文件
	os.Remove(dataFile)
	locate.Add(info.hash(), info.id())
}

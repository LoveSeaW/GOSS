package rs

import (
	"fmt"
	objectStream "goss/pkg/objectStream"
	"io"
)

type RSGetStream struct {
	*decoder
}

// 获取RSGetStream 对象，该对象可以通过读取数据碎片，重新生产对象文件
func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	if len(locateInfo)+len(dataServers) != AllShard {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	readers := make([]io.Reader, AllShard)
	// 补全所有数据节点
	for i := 0; i < AllShard; i++ {
		server := locateInfo[i]
		if server == "" {
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}

		reader, err := objectStream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		if err == nil {
			readers[i] = reader
		}
	}

	writers := make([]io.Writer, AllShard)
	// 获取 数据分片大小，向上取整
	perShard := (size + DataShard - 1) / DataShard
	var err error
	for i := range readers {
		if readers[i] == nil {
			// 生成的write 可以用于数据碎片写入数据服务节点
			writers[i], err = objectStream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if err != nil {
				return nil, err
			}
		}
	}
	decode := NewDecoder(readers, writers, size)
	return &RSGetStream{decoder: decode}, nil
}

func (s *RSGetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectStream.TempPutStream).Commit(true)
		}
	}
}

// 跳转至客户端请求位置，输出数据
func (s *RSGetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		panic("only support SeekCurrent")
	}

	if offset < 0 {
		panic("only support forward seek")
	}

	for offset != 0 {
		length := int64(BlockSize)
		if offset < length {
			length = offset
		}

		buffer := make([]byte, length)
		// 将数据 存入缓冲区
		// 从 Readers 切片中依次读取数据
		io.ReadFull(s, buffer)
		offset -= length
	}
	return offset, nil
}

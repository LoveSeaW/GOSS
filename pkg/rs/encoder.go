package rs

import (
	"io"

	"github.com/klauspost/reedsolomon"
)

type encoder struct {
	writers []io.Writer
	encode  reedsolomon.Encoder
	cache   []byte
}

func NewEncoder(writers []io.Writer) *encoder {
	encode, _ := reedsolomon.New(DataShard, ParityShard)
	return &encoder{encode: encode, writers: writers}
}

// 将数据按照 规定数据块的大小 依次写入cache，达到数据块规定大小后调用 Flush() 方法
func (e *encoder) Write(p []byte) (count int, err error) {
	length := len(p)
	current := 0
	for length != 0 {
		next := BlockSize - len(e.cache)
		if next > length {
			next = length
		}
		e.cache = append(e.cache, p[current:current+next]...)
		if len(e.cache) == BlockSize {
			e.Flush()
		}
		current += next
		length -= next
	}
	return len(p), nil
}

// 将 stream 流数据写入
func (e *encoder) Flush() {
	if len(e.cache) == 0 {
		return
	}
	// 将数据切成4个数据片
	shards, _ := e.encode.Split(e.cache)
	e.encode.Encode(shards)
	for i := range shards {
		// 调用 数据服务的 patch /temp/, 写入数据节点的临时文件
		e.writers[i].Write(shards[i])
	}
	e.cache = []byte{}
}

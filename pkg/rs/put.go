package rs

import (
	"fmt"
	objectstream "goss/pkg/objectStream"
	"io"
)

type RSPutStream struct {
	*encoder
}

func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	if len(dataServers) != AllShard {
		return nil, fmt.Errorf("dataServer number mismatch")
	}

	perShard := (size + DataShard - 1) / DataShard
	writers := make([]io.Writer, AllShard)
	var err error
	for i := range writers {
		writers[i], err = objectstream.NewTempPutStream(dataServers[i], fmt.Sprintf("%s.%d", hash, i), perShard)
		if err != nil {
			return nil, err
		}
	}

	encoder := NewEncoder(writers)
	return &RSPutStream{encoder: encoder}, nil
}

func (s *RSPutStream) Commit(success bool) {
	s.Flush()
	for i := range s.writers {
		s.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}

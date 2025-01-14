package objects

import (
	"fmt"
	"goss/apiServer/heartbeat"
	"goss/pkg/rs"
)

func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	servers := heartbeat.ChooseRandomDataServers(rs.AllShard, nil)
	if len(servers) != rs.AllShard {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}
	return rs.NewRSPutStream(servers, hash, size)
}

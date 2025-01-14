package objects

import (
	"fmt"
	"goss/apiServer/heartbeat"
	"goss/apiServer/locate"
	"goss/pkg/rs"
)

func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	locateInfo := locate.Locate(hash)
	// 如果存活数据节点 小于 rs 纠错码可恢复的最小节点即报错
	if len(locateInfo) < rs.DataShard {
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)
	}

	dataServers := make([]string, 0)
	if len(locateInfo) != rs.AllShard {
		// 随机选择剩余的存活节点
		dataServers = heartbeat.ChooseRandomDataServers(rs.AllShard-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}

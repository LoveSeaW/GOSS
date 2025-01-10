package heartbeat

import "math/rand"

// 选择随机的数据服务节点， n 代表选择的节点数量， exclude 代表需要排除的节点， 比如自身节点
func ChooseRandomDataServers(n int, exclude map[int]string) (dataServers []string) {
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	// 去除重复的数据服务节点 address
	for id, address := range exclude {
		reverseExcludeMap[address] = id
	}

	servers := GetDataServers()
	for i := range servers {
		server := servers[i]
		_, exclude := reverseExcludeMap[server]
		// 检查节点是否存活可用，不可排除加入 candidates
		if !exclude {
			candidates = append(candidates, server)
		}
	}
	length := len(candidates)
	if length < n {
		return
	}

	part := rand.Perm(length)
	for i := 0; i < n; i++ {
		dataServers = append(dataServers, candidates[part[i]])
	}
	return
}

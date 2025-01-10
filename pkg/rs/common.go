package rs

const (
	DataShard     = 4 // 数据碎片
	ParityShard   = 2 // 冗余碎片
	AllShard      = DataShard + ParityShard
	BlockPerShard = 8000                      // 每个碎片包含的数据块数量
	BlockSize     = BlockPerShard * DataShard // 所有数据块的大小
)

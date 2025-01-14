package deleteoldmetadata

import (
	"goss/pkg/es7"
	"log"
)

const MinVersionCount = 5

// 删除过期元数据的工具
func main() {
	buckets, err := es7.SearchVersionStatus(MinVersionCount + 1)
	if err != nil {
		log.Println(err)
		return
	}
	for i := range buckets {
		bucket := buckets[i]
		for v := 0; v < bucket.DocCount-MinVersionCount; v++ {
			es7.DelMetadata(bucket.Key, v+int(bucket.MinVersion.Value))
		}
	}
}

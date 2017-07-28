package storage

import (
	"time"
)

// ParsePutTime 提供了将PutTime转换为 time.Time 的功能
func ParsePutTime(putTime int64) (t time.Time) {
	t = time.Unix(0, putTime*100)
	return
}

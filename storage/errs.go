package storage

import (
	"errors"
)

var (
	// ErrBucketNotExist 用户存储空间不存在
	ErrBucketNotExist = errors.New("bucket not exist")
)

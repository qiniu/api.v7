package storage

const (
	DontCheckCrc    = 0
	CalcAndCheckCrc = 1
	CheckCrc        = 2
)

// 表单上传的额外可选项
type PutExtra struct {
	// 可选，用户自定义参数，必须以 "x:" 开头。若不以x:开头，则忽略。
	Params map[string]string

	// 可选，当为 "" 时候，服务端自动判断。
	MimeType string

	Crc32 uint32

	// CheckCrc == 0 (DontCheckCrc): 表示不进行 crc32 校验
	// CheckCrc == 1 (CalcAndCheckCrc): 对于 Put 等同于 CheckCrc = 2；对于 PutFile 会自动计算 crc32 值
	// CheckCrc == 2 (CheckCrc): 表示进行 crc32 校验，且 crc32 值就是上面的 Crc32 变量
	CheckCrc uint32

	// 上传事件：进度通知。这个事件的回调函数应该尽可能快地结束。
	OnProgress func(fsize, uploaded int64)
}

type BlkputRet struct {
	Ctx      string `json:"ctx"`
	Checksum string `json:"checksum"`
	Crc32    uint32 `json:"crc32"`
	Offset   uint32 `json:"offset"`
	Host     string `json:"host"`
}

// 分片上传的额外可选项
type RputExtra struct {
	Params     map[string]string                             // 可选。用户自定义参数，以"x:"开头 否则忽略
	MimeType   string                                        // 可选。
	ChunkSize  int                                           // 可选。每次上传的Chunk大小
	TryTimes   int                                           // 可选。尝试次数
	Progresses []BlkputRet                                   // 可选。上传进度
	Notify     func(blkIdx int, blkSize int, ret *BlkputRet) // 可选。进度提示（注意多个block是并行传输的）
	NotifyErr  func(blkIdx int, blkSize int, err error)
}

// 如果 uptoken 没有指定 ReturnBody，那么返回值是标准的 PutRet 结构
type PutRet struct {
	Hash         string `json:"hash"`
	PersistentId string `json:"persistentId"`
	Key          string `json:"key"`
}

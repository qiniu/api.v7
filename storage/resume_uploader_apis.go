package storage

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/qiniu/api.v7/v7/client"
	"github.com/qiniu/api.v7/v7/conf"
)

type resumeUploaderAPIs struct {
	Client *client.Client
	Cfg    *Config
}

// BlkputRet 表示分片上传每个片上传完毕的返回值
type BlkputRet struct {
	Ctx       string `json:"ctx"`
	Checksum  string `json:"checksum"`
	Crc32     uint32 `json:"crc32"`
	Offset    uint32 `json:"offset"`
	Host      string `json:"host"`
	ExpiredAt int64  `json:"expired_at"`
	blkIdx    int
}

func (p *resumeUploaderAPIs) mkBlk(ctx context.Context, upToken string, upHost string, ret *BlkputRet, blockSize int, body io.Reader, size int) error {
	reqUrl := upHost + "/mkblk/" + strconv.Itoa(blockSize)

	return p.Client.CallWith(ctx, ret, "POST", reqUrl, makeHeadersForUpload(upToken), body, size)
}

func (p *resumeUploaderAPIs) bput(ctx context.Context, upToken string, ret *BlkputRet, body io.Reader, size int) error {
	reqUrl := ret.Host + "/bput/" + ret.Ctx + "/" + strconv.FormatUint(uint64(ret.Offset), 10)

	return p.Client.CallWith(ctx, ret, "POST", reqUrl, makeHeadersForUpload(upToken), body, size)
}

// RputExtra 表示分片上传额外可以指定的参数
type RputExtra struct {
	Params     map[string]string // 可选。用户自定义参数，以"x:"开头，而且值不能为空，否则忽略
	UpHost     string
	MimeType   string                                        // 可选。
	ChunkSize  int                                           // 可选。每次上传的Chunk大小
	TryTimes   int                                           // 可选。尝试次数
	Progresses []BlkputRet                                   // 可选。上传进度
	Notify     func(blkIdx int, blkSize int, ret *BlkputRet) // 可选。进度提示（注意多个block是并行传输的）
	NotifyErr  func(blkIdx int, blkSize int, err error)
}

func (p *resumeUploaderAPIs) mkfile(ctx context.Context, upToken string, upHost string, ret interface{}, key string, hasKey bool, fsize int64, extra *RputExtra) (err error) {
	url := upHost + "/mkfile/" + strconv.FormatInt(fsize, 10)
	if extra == nil {
		extra = &RputExtra{}
	}
	if extra.MimeType != "" {
		url += "/mimeType/" + encode(extra.MimeType)
	}
	if hasKey {
		url += "/key/" + encode(key)
	}
	for k, v := range extra.Params {
		if (strings.HasPrefix(k, "x:") || strings.HasPrefix(k, "x-qn-meta-")) && v != "" {
			url += "/" + k + "/" + encode(v)
		}
	}
	ctxs := make([]string, len(extra.Progresses))
	for i, progress := range extra.Progresses {
		ctxs[i] = progress.Ctx
	}
	buf := strings.Join(ctxs, ",")
	return p.Client.CallWith(ctx, ret, "POST", url, makeHeadersForUpload(upToken), strings.NewReader(buf), len(buf))
}

func (p *resumeUploaderAPIs) upHost(ak, bucket string) (upHost string, err error) {
	return getUpHost(p.Cfg, ak, bucket)
}

func (p *resumeUploaderAPIs) getUpHostFromUploadToken(upToken string) (upHost string, err error) {
	var ak, bucket string

	if ak, bucket, err = getAkBucketFromUploadToken(upToken); err != nil {
		return
	}
	upHost, err = p.upHost(ak, bucket)
	return
}

func makeHeadersForUpload(upToken string) http.Header {
	headers := http.Header{}
	headers.Add("Content-Type", conf.CONTENT_TYPE_OCTET)
	headers.Add("Authorization", "UpToken "+upToken)
	return headers
}

func encode(raw string) string {
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

func (r *RputExtra) init() {
	if r.ChunkSize == 0 {
		r.ChunkSize = settings.ChunkSize
	}
	if r.TryTimes == 0 {
		r.TryTimes = settings.TryTimes
	}
	if r.Notify == nil {
		r.Notify = func(blkIdx, blkSize int, ret *BlkputRet) {}
	}
	if r.NotifyErr == nil {
		r.NotifyErr = func(blkIdx, blkSize int, err error) {}
	}
}

type BlkputRets []BlkputRet

func (rets BlkputRets) Len() int {
	return len(rets)
}

func (rets BlkputRets) Less(i, j int) bool {
	return rets[i].blkIdx < rets[j].blkIdx
}

func (rets BlkputRets) Swap(i, j int) {
	rets[i], rets[j] = rets[j], rets[i]
}

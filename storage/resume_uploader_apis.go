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
	Ctx        string `json:"ctx"`
	Checksum   string `json:"checksum"`
	Crc32      uint32 `json:"crc32"`
	Offset     uint32 `json:"offset"`
	Host       string `json:"host"`
	ExpiredAt  int64  `json:"expired_at"`
	chunkSize  int
	fileOffset int64
	blkIdx     int
}

func (p *resumeUploaderAPIs) mkBlk(ctx context.Context, upToken, upHost string, ret *BlkputRet, blockSize int, body io.Reader, size int) error {
	reqUrl := upHost + "/mkblk/" + strconv.Itoa(blockSize)

	return p.Client.CallWith(ctx, ret, "POST", reqUrl, makeHeadersForUpload(upToken), body, size)
}

func (p *resumeUploaderAPIs) bput(ctx context.Context, upToken string, ret *BlkputRet, body io.Reader, size int) error {
	reqUrl := ret.Host + "/bput/" + ret.Ctx + "/" + strconv.FormatUint(uint64(ret.Offset), 10)

	return p.Client.CallWith(ctx, ret, "POST", reqUrl, makeHeadersForUpload(upToken), body, size)
}

// RputExtra 表示分片上传额外可以指定的参数
type RputExtra struct {
	Recorder   Recorder          // 可选。上传进度记录
	Params     map[string]string // 可选。用户自定义参数，以"x:"开头，而且值不能为空，否则忽略
	UpHost     string
	MimeType   string                                        // 可选。
	ChunkSize  int                                           // 可选。每次上传的Chunk大小
	TryTimes   int                                           // 可选。尝试次数
	Progresses []BlkputRet                                   // 可选。上传进度
	Notify     func(blkIdx int, blkSize int, ret *BlkputRet) // 可选。进度提示（注意多个block是并行传输的）
	NotifyErr  func(blkIdx int, blkSize int, err error)
}

func (p *resumeUploaderAPIs) mkfile(ctx context.Context, upToken, upHost string, ret interface{}, key string, hasKey bool, fsize int64, extra *RputExtra) (err error) {
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

// InitPartsRet 表示分片上传 v2 初始化完毕的返回值
type InitPartsRet struct {
	UploadID string `json:"uploadId"`
}

func (p *resumeUploaderAPIs) initParts(ctx context.Context, upToken, upHost, bucket, key string, hasKey bool, ret *InitPartsRet) error {
	reqUrl := upHost + "/buckets/" + bucket + "/objects/" + encodeV2(key, hasKey) + "/uploads"

	return p.Client.CallWith(ctx, ret, "POST", reqUrl, makeHeadersForUploadEx(upToken, ""), nil, 0)
}

// UploadPartsRet 表示分片上传 v2 每个片上传完毕的返回值
type UploadPartsRet struct {
	Etag string `json:"etag"`
	MD5  string `json:"md5"`
}

func (p *resumeUploaderAPIs) uploadParts(ctx context.Context, upToken, upHost, bucket, key string, hasKey bool, uploadId string, partNumber int64, partMD5 string, ret *UploadPartsRet, body io.Reader, size int) error {
	reqUrl := upHost + "/buckets/" + bucket + "/objects/" + encodeV2(key, hasKey) + "/uploads/" + uploadId + "/" + strconv.FormatInt(partNumber, 10)

	return p.Client.CallWith(ctx, ret, "PUT", reqUrl, makeHeadersForUploadPart(upToken, partMD5), body, size)
}

type uploadPartInfo struct {
	Etag       string `json:"etag"`
	PartNumber int64  `json:"partNumber"`
	partSize   int
	fileOffset int64
}

// RputV2Extra 表示分片上传 v2 额外可以指定的参数
type RputV2Extra struct {
	Recorder   Recorder          // 可选。上传进度记录
	Metadata   map[string]string // 可选。用户自定义文件 metadata 信息
	CustomVars map[string]string // 可选。用户自定义参数，以"x:"开头，而且值不能为空，否则忽略
	UpHost     string
	MimeType   string                                      // 可选。
	PartSize   int64                                       // 可选。每次上传的块大小
	TryTimes   int                                         // 可选。尝试次数
	progresses []uploadPartInfo                            // 上传进度
	Notify     func(partNumber int64, ret *UploadPartsRet) // 可选。进度提示（注意多个block是并行传输的）
	NotifyErr  func(partNumber int64, err error)
}

func (p *resumeUploaderAPIs) completeParts(ctx context.Context, upToken, upHost string, ret interface{}, bucket, key string, hasKey bool, uploadId string, extra *RputV2Extra) (err error) {
	type CompletePartBody struct {
		Parts      []uploadPartInfo  `json:"parts"`
		MimeType   string            `json:"mimeType,omitempty"`
		Metadata   map[string]string `json:"metadata,omitempty"`
		CustomVars map[string]string `json:"customVars,omitempty"`
	}
	if extra == nil {
		extra = &RputV2Extra{}
	}
	completePartBody := CompletePartBody{
		Parts:      extra.progresses,
		MimeType:   extra.MimeType,
		Metadata:   extra.Metadata,
		CustomVars: make(map[string]string),
	}
	for k, v := range extra.CustomVars {
		if strings.HasPrefix(k, "x:") && v != "" {
			completePartBody.CustomVars[k] = v
		}
	}

	reqUrl := upHost + "/buckets/" + bucket + "/objects/" + encodeV2(key, hasKey) + "/uploads/" + uploadId

	return p.Client.CallWithJson(ctx, ret, "POST", reqUrl, makeHeadersForUploadEx(upToken, conf.CONTENT_TYPE_JSON), &completePartBody)
}

func (p *resumeUploaderAPIs) upHost(ak, bucket string) (upHost string, err error) {
	return getUpHost(p.Cfg, ak, bucket)
}

func (p *resumeUploaderAPIs) getUpHostFromUploadToken(upToken string) (upHost string, err error) {
	_, upHost, err = p.getBucketAndUpHostFromUploadToken(upToken)
	return
}

func (p *resumeUploaderAPIs) getBucketAndUpHostFromUploadToken(upToken string) (bucket, upHost string, err error) {
	_, bucket, upHost, err = p.getAkAndBucketAndUpHostFromUploadToken(upToken)
	return
}

func (p *resumeUploaderAPIs) getAkAndBucketAndUpHostFromUploadToken(upToken string) (ak, bucket, upHost string, err error) {
	if ak, bucket, err = getAkBucketFromUploadToken(upToken); err != nil {
		return
	}
	upHost, err = p.upHost(ak, bucket)
	return
}

func makeHeadersForUpload(upToken string) http.Header {
	return makeHeadersForUploadEx(upToken, conf.CONTENT_TYPE_OCTET)
}

func makeHeadersForUploadPart(upToken, partMD5 string) http.Header {
	headers := makeHeadersForUpload(upToken)
	headers.Add("Content-MD5", partMD5)
	return headers
}

func makeHeadersForUploadEx(upToken, contentType string) http.Header {
	headers := http.Header{}
	if contentType != "" {
		headers.Add("Content-Type", contentType)
	}
	headers.Add("Authorization", "UpToken "+upToken)
	return headers
}

func encode(raw string) string {
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

func encodeV2(key string, hasKey bool) string {
	if !hasKey {
		return "~"
	} else {
		return encode(key)
	}
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

func (r *RputV2Extra) init() {
	if r.PartSize == 0 {
		r.PartSize = settings.PartSize
	}
	if r.TryTimes == 0 {
		r.TryTimes = settings.TryTimes
	}
	if r.Notify == nil {
		r.Notify = func(partNumber int64, ret *UploadPartsRet) {}
	}
	if r.NotifyErr == nil {
		r.NotifyErr = func(partNumber int64, err error) {}
	}
}

type blkputRets []BlkputRet

func (rets blkputRets) Len() int {
	return len(rets)
}

func (rets blkputRets) Less(i, j int) bool {
	return rets[i].blkIdx < rets[j].blkIdx
}

func (rets blkputRets) Swap(i, j int) {
	rets[i], rets[j] = rets[j], rets[i]
}

type uploadPartInfos []uploadPartInfo

func (infos uploadPartInfos) Len() int {
	return len(infos)
}

func (infos uploadPartInfos) Less(i, j int) bool {
	return infos[i].PartNumber < infos[j].PartNumber
}

func (infos uploadPartInfos) Swap(i, j int) {
	infos[i], infos[j] = infos[j], infos[i]
}

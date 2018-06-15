package storage

import (
	by "bytes"
	"context"
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"strconv"

	"github.com/qiniu/api.v7/conf"
	"github.com/qiniu/x/bytes.v7"
	"github.com/qiniu/x/xlog.v7"
)

// ResumeUploader 表示一个分片上传的对象
type ResumeUploader struct {
	client *Client
	cfg    *Config
}

// NewResumeUploader 表示构建一个新的分片上传的对象
func NewResumeUploader(cfg *Config) *ResumeUploader {
	if cfg == nil {
		cfg = &Config{}
	}

	return &ResumeUploader{
		cfg:    cfg,
		client: &DefaultClient,
	}
}

// NewResumeUploaderEx 表示构建一个新的分片上传的对象
func NewResumeUploaderEx(cfg *Config, client *Client) *ResumeUploader {
	if cfg == nil {
		cfg = &Config{}
	}

	if client == nil {
		client = &DefaultClient
	}

	return &ResumeUploader{
		client: client,
		cfg:    cfg,
	}
}

// 创建块请求
func (p *ResumeUploader) Mkblk(
	ctx context.Context, upToken string, upHost string, ret *BlkputRet, blockSize int, body io.Reader, size int) error {

	reqUrl := upHost + "/mkblk/" + strconv.Itoa(blockSize)
	headers := http.Header{}
	headers.Add("Content-Type", conf.CONTENT_TYPE_OCTET)
	headers.Add("Authorization", "UpToken "+upToken)

	return p.client.CallWith(ctx, ret, "POST", reqUrl, headers, body, size)
}

// 发送bput请求
func (p *ResumeUploader) Bput(
	ctx context.Context, upToken string, ret *BlkputRet, body io.Reader, size int) error {

	reqUrl := ret.Host + "/bput/" + ret.Ctx + "/" + strconv.FormatUint(uint64(ret.Offset), 10)
	headers := http.Header{}
	headers.Add("Content-Type", conf.CONTENT_TYPE_OCTET)
	headers.Add("Authorization", "UpToken "+upToken)

	return p.client.CallWith(ctx, ret, "POST", reqUrl, headers, body, size)
}

type blockReader struct {
	*by.Reader
	blkIdx  int
	blkSize int
}

func newBlockReader(r io.Reader, blkIdx int) (*blockReader, error) {
	maxBlockSize := 1 << blockBits
	b := make([]byte, maxBlockSize)
	n, e := io.ReadFull(r, b)
	if n != maxBlockSize && e != io.EOF && e != io.ErrUnexpectedEOF {
		return nil, e
	}
	return &blockReader{by.NewReader(b[:n]), blkIdx, n}, nil
}

func (br *blockReader) ReadAt(p []byte, off int64) (int, error) {
	off -= int64(br.blkIdx) << blockBits
	return br.Reader.ReadAt(p, off)
}

func (p *ResumeUploader) resumableBputWithoutSize(
	ctx context.Context, upToken, upHost string, ret *BlkputRet, br *blockReader, extra *RputExtra) error {

	return p.resumableBput(ctx, upToken, upHost, ret, br, br.blkIdx, br.blkSize, extra)
}

// 分片上传请求
func (p *ResumeUploader) resumableBput(
	ctx context.Context, upToken string, upHost string, ret *BlkputRet, f io.ReaderAt, blkIdx, blkSize int, extra *RputExtra) (err error) {

	log := xlog.NewWith(ctx)
	h := crc32.NewIEEE()
	offbase := int64(blkIdx) << blockBits
	chunkSize := extra.ChunkSize

	var bodyLength int

	if ret.Ctx == "" {

		if chunkSize < blkSize {
			bodyLength = chunkSize
		} else {
			bodyLength = blkSize
		}

		body1 := io.NewSectionReader(f, offbase, int64(bodyLength))
		body := io.TeeReader(body1, h)

		err = p.Mkblk(ctx, upToken, upHost, ret, blkSize, body, bodyLength)
		if err != nil {
			return
		}
		if ret.Crc32 != h.Sum32() || int(ret.Offset) != bodyLength {
			err = ErrUnmatchedChecksum
			return
		}
		extra.Notify(blkIdx, blkSize, ret)
	}

	for int(ret.Offset) < blkSize {

		if chunkSize < blkSize-int(ret.Offset) {
			bodyLength = chunkSize
		} else {
			bodyLength = blkSize - int(ret.Offset)
		}

		tryTimes := extra.TryTimes

	lzRetry:
		h.Reset()
		body1 := io.NewSectionReader(f, offbase+int64(ret.Offset), int64(bodyLength))
		body := io.TeeReader(body1, h)

		err = p.Bput(ctx, upToken, ret, body, bodyLength)
		if err == nil {
			if ret.Crc32 == h.Sum32() {
				extra.Notify(blkIdx, blkSize, ret)
				continue
			}
			log.Warn("ResumableBlockput: invalid checksum, retry")
			err = ErrUnmatchedChecksum
		} else {
			if ei, ok := err.(*ErrorInfo); ok && ei.Code == InvalidCtx {
				ret.Ctx = "" // reset
				log.Warn("ResumableBlockput: invalid ctx, please retry")
				return
			}
			log.Warn("ResumableBlockput: bput failed -", err)
		}
		if tryTimes > 1 {
			tryTimes--
			log.Info("ResumableBlockput retrying ...")
			goto lzRetry
		}
		break
	}
	return
}

// 创建文件请求
func (p *ResumeUploader) Mkfile(
	ctx context.Context, upToken string, upHost string, ret interface{}, key string, hasKey bool, fsize int64, extra *RputExtra) (err error) {

	url := upHost + "/mkfile/" + strconv.FormatInt(fsize, 10)

	if extra.MimeType != "" {
		url += "/mimeType/" + encode(extra.MimeType)
	}
	if hasKey {
		url += "/key/" + encode(key)
	}
	for k, v := range extra.Params {
		url += fmt.Sprintf("/%s/%s", k, encode(v))
	}

	buf := make([]byte, 0, 196*len(extra.Progresses))
	for _, prog := range extra.Progresses {
		buf = append(buf, prog.Ctx...)
		buf = append(buf, ',')
	}
	if len(buf) > 0 {
		buf = buf[:len(buf)-1]
	}

	headers := http.Header{}
	headers.Add("Content-Type", conf.CONTENT_TYPE_OCTET)
	headers.Add("Authorization", "UpToken "+upToken)

	return p.client.CallWith(
		ctx, ret, "POST", url, headers, bytes.NewReader(buf), len(buf))
}

func encode(raw string) string {
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

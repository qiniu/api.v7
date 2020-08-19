package storage

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"sort"
	"sync"

	"github.com/qiniu/api.v7/v7/client"
)

// ResumeUploaderV2 表示一个分片上传 v2 的对象
type ResumeUploaderV2 struct {
	Client *client.Client
	Cfg    *Config
}

// NewResumeUploaderV2 表示构建一个新的分片上传的对象
func NewResumeUploaderV2(cfg *Config) *ResumeUploaderV2 {
	return NewResumeUploaderV2Ex(cfg, nil)
}

// NewResumeUploaderV2Ex 表示构建一个新的分片上传 v2 的对象
func NewResumeUploaderV2Ex(cfg *Config, clt *client.Client) *ResumeUploaderV2 {
	if cfg == nil {
		cfg = &Config{}
	}

	if clt == nil {
		clt = &client.DefaultClient
	}

	return &ResumeUploaderV2{
		Client: clt,
		Cfg:    cfg,
	}
}

// Put 方法用来上传一个文件，支持断点续传和分块上传。
//
// ctx     是请求的上下文。
// ret     是上传成功后返回的数据。如果 upToken 中没有设置 CallbackUrl 或 ReturnBody，那么返回的数据结构是 PutRet 结构。
// upToken 是由业务服务器颁发的上传凭证。
// key     是要上传的文件访问路径。比如："foo/bar.jpg"。注意我们建议 key 不要以 '/' 开头。另外，key 为空字符串是合法的。
// f       是文件内容的访问接口。考虑到需要支持分块上传和断点续传，要的是 io.ReaderAt 接口，而不是 io.Reader。
// fsize   是要上传的文件大小。
// extra   是上传的一些可选项。详细见 RputV2Extra 结构的描述。
//
func (p *ResumeUploaderV2) Put(ctx context.Context, ret interface{}, upToken string, key string, f io.ReaderAt, fsize int64, extra *RputV2Extra) error {
	return p.rput(ctx, ret, upToken, key, true, f, fsize, extra)
}

func (p *ResumeUploaderV2) PutWithoutSize(ctx context.Context, ret interface{}, upToken, key string, r io.Reader, extra *RputV2Extra) error {
	return p.rputWithoutSize(ctx, ret, upToken, key, true, r, extra)
}

// PutWithoutKey 方法用来上传一个文件，支持断点续传和分块上传。文件命名方式首先看看
// upToken 中是否设置了 saveKey，如果设置了 saveKey，那么按 saveKey 要求的规则生成 key，否则自动以文件的 hash 做 key。
//
// ctx     是请求的上下文。
// ret     是上传成功后返回的数据。如果 upToken 中没有设置 CallbackUrl 或 ReturnBody，那么返回的数据结构是 PutRet 结构。
// upToken 是由业务服务器颁发的上传凭证。
// f       是文件内容的访问接口。考虑到需要支持分块上传和断点续传，要的是 io.ReaderAt 接口，而不是 io.Reader。
// fsize   是要上传的文件大小。
// extra   是上传的一些可选项。详细见 RputV2Extra 结构的描述。
//
func (p *ResumeUploaderV2) PutWithoutKey(ctx context.Context, ret interface{}, upToken string, f io.ReaderAt, fsize int64, extra *RputV2Extra) error {
	return p.rput(ctx, ret, upToken, "", false, f, fsize, extra)
}

// PutFile 用来上传一个文件，支持断点续传和分块上传。
// 和 Put 不同的只是一个通过提供文件路径来访问文件内容，一个通过 io.ReaderAt 来访问。
//
// ctx       是请求的上下文。
// ret       是上传成功后返回的数据。如果 upToken 中没有设置 CallbackUrl 或 ReturnBody，那么返回的数据结构是 PutRet 结构。
// upToken   是由业务服务器颁发的上传凭证。
// key       是要上传的文件访问路径。比如："foo/bar.jpg"。注意我们建议 key 不要以 '/' 开头。另外，key 为空字符串是合法的。
// localFile 是要上传的文件的本地路径。
// extra     是上传的一些可选项。详细见 RputV2Extra 结构的描述。
//
func (p *ResumeUploaderV2) PutFile(ctx context.Context, ret interface{}, upToken, key, localFile string, extra *RputV2Extra) error {
	return p.rputFile(ctx, ret, upToken, key, true, localFile, extra)
}

// PutFileWithoutKey 上传一个文件，支持断点续传和分块上传。文件命名方式首先看看
// upToken 中是否设置了 saveKey，如果设置了 saveKey，那么按 saveKey 要求的规则生成 key，否则自动以文件的 hash 做 key。
// 和 PutWithoutKey 不同的只是一个通过提供文件路径来访问文件内容，一个通过 io.ReaderAt 来访问。
//
// ctx       是请求的上下文。
// ret       是上传成功后返回的数据。如果 upToken 中没有设置 CallbackUrl 或 ReturnBody，那么返回的数据结构是 PutRet 结构。
// upToken   是由业务服务器颁发的上传凭证。
// localFile 是要上传的文件的本地路径。
// extra     是上传的一些可选项。详细见 RputV2Extra 结构的描述。
//
func (p *ResumeUploaderV2) PutFileWithoutKey(ctx context.Context, ret interface{}, upToken, localFile string, extra *RputV2Extra) error {
	return p.rputFile(ctx, ret, upToken, "", false, localFile, extra)
}

func (p *ResumeUploaderV2) rput(ctx context.Context, ret interface{}, upToken string, key string, hasKey bool, f io.ReaderAt, fsize int64, extra *RputV2Extra) (err error) {
	if extra == nil {
		extra = &RputV2Extra{}
	}
	extra.init()

	var bucket, upHost string
	if bucket, upHost, err = p.resumeUploaderAPIs().getBucketAndUpHostFromUploadToken(upToken); err != nil {
		return
	}
	if extra.UpHost != "" {
		upHost = extra.UpHost
	}

	return uploadByWorkers(newResumeUploaderV2Impl(p, bucket, key, hasKey, upToken, upHost, extra, ret), ctx, newSizedChunkReader(f, fsize, extra.PartSize), extra.TryTimes)
}

func (p *ResumeUploaderV2) rputWithoutSize(ctx context.Context, ret interface{}, upToken string, key string, hasKey bool, r io.Reader, extra *RputV2Extra) (err error) {
	if extra == nil {
		extra = &RputV2Extra{}
	}
	extra.init()

	var bucket, upHost string
	if bucket, upHost, err = p.resumeUploaderAPIs().getBucketAndUpHostFromUploadToken(upToken); err != nil {
		return
	}
	if extra.UpHost != "" {
		upHost = extra.UpHost
	}

	return uploadByWorkers(newResumeUploaderV2Impl(p, bucket, key, hasKey, upToken, upHost, extra, ret), ctx, newUnsizedChunkReader(r, extra.PartSize), extra.TryTimes)
}

func (p *ResumeUploaderV2) rputFile(ctx context.Context, ret interface{}, upToken string, key string, hasKey bool, localFile string, extra *RputV2Extra) (err error) {
	var (
		file     *os.File
		fileInfo os.FileInfo
	)

	if file, err = os.Open(localFile); err != nil {
		return
	}
	defer file.Close()

	if fileInfo, err = file.Stat(); err != nil {
		return
	}

	return p.rput(ctx, ret, upToken, key, hasKey, file, fileInfo.Size(), extra)
}

// 创建块请求
func (p *ResumeUploaderV2) InitParts(ctx context.Context, upToken, upHost, bucket, key string, hasKey bool, ret *InitPartsRet) error {
	return p.resumeUploaderAPIs().initParts(ctx, upToken, upHost, bucket, key, hasKey, ret)
}

// 发送bput请求
func (p *ResumeUploaderV2) UploadParts(ctx context.Context, upToken, upHost, bucket, key string, hasKey bool, uploadId string, partNumber int64, partMD5 string, ret *UploadPartsRet, body io.Reader, size int) error {
	return p.resumeUploaderAPIs().uploadParts(ctx, upToken, upHost, bucket, key, hasKey, uploadId, partNumber, partMD5, ret, body, size)
}

// 创建文件请求
func (p *ResumeUploaderV2) CompleteParts(ctx context.Context, upToken, upHost string, ret interface{}, bucket, key string, hasKey bool, uploadId string, extra *RputV2Extra) (err error) {
	return p.resumeUploaderAPIs().completeParts(ctx, upToken, upHost, ret, bucket, key, hasKey, uploadId, extra)
}

func (p *ResumeUploaderV2) UpHost(ak, bucket string) (upHost string, err error) {
	return p.resumeUploaderAPIs().upHost(ak, bucket)
}

func (p *ResumeUploaderV2) resumeUploaderAPIs() *resumeUploaderAPIs {
	return &resumeUploaderAPIs{Client: p.Client, Cfg: p.Cfg}
}

// 用于实现 resumeUploaderBase 的 V2 分片接口
type resumeUploaderV2Impl struct {
	ctx      context.Context
	client   *client.Client
	cfg      *Config
	bucket   string
	key      string
	hasKey   bool
	uploadId string
	upToken  string
	upHost   string
	extra    *RputV2Extra
	ret      interface{}
	lock     sync.Mutex
}

func newResumeUploaderV2Impl(resumeUploader *ResumeUploaderV2, bucket, key string, hasKey bool, upToken string, upHost string, extra *RputV2Extra, ret interface{}) *resumeUploaderV2Impl {
	return &resumeUploaderV2Impl{
		client:  resumeUploader.Client,
		cfg:     resumeUploader.Cfg,
		bucket:  bucket,
		key:     key,
		hasKey:  hasKey,
		upToken: upToken,
		upHost:  upHost,
		extra:   extra,
		ret:     ret,
	}
}

func (impl *resumeUploaderV2Impl) initUploader(ctx context.Context) error {
	var ret InitPartsRet
	err := impl.resumeUploaderAPIs().initParts(ctx, impl.upToken, impl.upHost, impl.bucket, impl.key, impl.hasKey, &ret)
	if err == nil {
		impl.uploadId = ret.UploadID
	}
	return err
}

func (impl *resumeUploaderV2Impl) uploadChunk(ctx context.Context, c chunk) error {
	var (
		apis = impl.resumeUploaderAPIs()
		ret  UploadPartsRet
		err  error
	)

	md5ByteArray := md5.Sum(c.data)
	md5Value := hex.EncodeToString(md5ByteArray[:])
	partNumber := c.id + 1

	if err = apis.uploadParts(ctx, impl.upToken, impl.upHost, impl.bucket, impl.key, impl.hasKey, impl.uploadId, partNumber, md5Value, &ret, bytes.NewReader(c.data), len(c.data)); err != nil {
		impl.extra.NotifyErr(partNumber, err)
	} else {
		impl.extra.Notify(partNumber, &ret)
		func() {
			impl.lock.Lock()
			defer impl.lock.Unlock()
			impl.extra.progresses = append(impl.extra.progresses, uploadPartInfo{Etag: ret.Etag, PartNumber: partNumber})
		}()
	}
	return err
}

func (impl *resumeUploaderV2Impl) final(ctx context.Context) error {
	sort.Sort(uploadPartInfos(impl.extra.progresses))
	return impl.resumeUploaderAPIs().completeParts(ctx, impl.upToken, impl.upHost, impl.ret, impl.bucket, impl.key, impl.hasKey, impl.uploadId, impl.extra)
}

func (impl *resumeUploaderV2Impl) resumeUploaderAPIs() *resumeUploaderAPIs {
	return &resumeUploaderAPIs{Client: impl.client, Cfg: impl.cfg}
}

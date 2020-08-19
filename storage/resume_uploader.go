package storage

import (
	"bytes"
	"context"
	"hash/crc32"
	"io"
	"os"
	"sort"
	"sync"

	"github.com/qiniu/api.v7/v7/client"
)

// ResumeUploader 表示一个分片上传的对象
type ResumeUploader struct {
	Client *client.Client
	Cfg    *Config
}

// NewResumeUploader 表示构建一个新的分片上传的对象
func NewResumeUploader(cfg *Config) *ResumeUploader {
	return NewResumeUploaderEx(cfg, nil)
}

// NewResumeUploaderEx 表示构建一个新的分片上传的对象
func NewResumeUploaderEx(cfg *Config, clt *client.Client) *ResumeUploader {
	if cfg == nil {
		cfg = &Config{}
	}

	if clt == nil {
		clt = &client.DefaultClient
	}

	return &ResumeUploader{
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
// extra   是上传的一些可选项。详细见 RputExtra 结构的描述。
//
func (p *ResumeUploader) Put(ctx context.Context, ret interface{}, upToken string, key string, f io.ReaderAt, fsize int64, extra *RputExtra) error {
	return p.rput(ctx, ret, upToken, key, true, f, fsize, extra)
}

func (p *ResumeUploader) PutWithoutSize(ctx context.Context, ret interface{}, upToken, key string, r io.Reader, extra *RputExtra) error {
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
// extra   是上传的一些可选项。详细见 RputExtra 结构的描述。
//
func (p *ResumeUploader) PutWithoutKey(ctx context.Context, ret interface{}, upToken string, f io.ReaderAt, fsize int64, extra *RputExtra) error {
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
// extra     是上传的一些可选项。详细见 RputExtra 结构的描述。
//
func (p *ResumeUploader) PutFile(ctx context.Context, ret interface{}, upToken, key, localFile string, extra *RputExtra) error {
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
// extra     是上传的一些可选项。详细见 RputExtra 结构的描述。
//
func (p *ResumeUploader) PutFileWithoutKey(ctx context.Context, ret interface{}, upToken, localFile string, extra *RputExtra) error {
	return p.rputFile(ctx, ret, upToken, "", false, localFile, extra)
}

func (p *ResumeUploader) rput(ctx context.Context, ret interface{}, upToken string, key string, hasKey bool, f io.ReaderAt, fsize int64, extra *RputExtra) (err error) {
	if extra == nil {
		extra = &RputExtra{}
	}
	extra.init()

	var upHost string
	if extra.UpHost != "" {
		upHost = extra.UpHost
	} else if upHost, err = p.resumeUploaderAPIs().getUpHostFromUploadToken(upToken); err != nil {
		return
	}

	return uploadByWorkers(newResumeUploaderImpl(p, key, hasKey, upToken, upHost, extra, ret), ctx, newSizedChunkReader(f, fsize, 1<<blockBits), extra.TryTimes)
}

func (p *ResumeUploader) rputWithoutSize(ctx context.Context, ret interface{}, upToken string, key string, hasKey bool, r io.Reader, extra *RputExtra) (err error) {
	if extra == nil {
		extra = &RputExtra{}
	}
	extra.init()

	var upHost string
	if extra.UpHost != "" {
		upHost = extra.UpHost
	} else if upHost, err = p.resumeUploaderAPIs().getUpHostFromUploadToken(upToken); err != nil {
		return
	}

	return uploadByWorkers(newResumeUploaderImpl(p, key, hasKey, upToken, upHost, extra, ret), ctx, newUnsizedChunkReader(r, 1<<blockBits), extra.TryTimes)
}

func (p *ResumeUploader) rputFile(ctx context.Context, ret interface{}, upToken string, key string, hasKey bool, localFile string, extra *RputExtra) (err error) {
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
func (p *ResumeUploader) Mkblk(ctx context.Context, upToken string, upHost string, ret *BlkputRet, blockSize int, body io.Reader, size int) error {
	return p.resumeUploaderAPIs().mkBlk(ctx, upToken, upHost, ret, blockSize, body, size)
}

// 发送bput请求
func (p *ResumeUploader) Bput(ctx context.Context, upToken string, ret *BlkputRet, body io.Reader, size int) error {
	return p.resumeUploaderAPIs().bput(ctx, upToken, ret, body, size)
}

// 创建文件请求
func (p *ResumeUploader) Mkfile(ctx context.Context, upToken string, upHost string, ret interface{}, key string, hasKey bool, fsize int64, extra *RputExtra) (err error) {
	return p.resumeUploaderAPIs().mkfile(ctx, upToken, upHost, ret, key, hasKey, fsize, extra)
}

func (p *ResumeUploader) UpHost(ak, bucket string) (upHost string, err error) {
	return p.resumeUploaderAPIs().upHost(ak, bucket)
}

func (p *ResumeUploader) resumeUploaderAPIs() *resumeUploaderAPIs {
	return &resumeUploaderAPIs{Client: p.Client, Cfg: p.Cfg}
}

type (
	notifyFunc    func(blkIdx, blkSize int, ret *BlkputRet)
	notifyErrFunc func(blkIdx, blkSize int, err error)

	// 用于实现 resumeUploaderBase 的 V1 分片接口
	resumeUploaderImpl struct {
		ctx      context.Context
		client   *client.Client
		cfg      *Config
		key      string
		hasKey   bool
		upToken  string
		upHost   string
		extra    *RputExtra
		ret      interface{}
		fileSize int64
		lock     sync.Mutex
	}
)

func newResumeUploaderImpl(resumeUploader *ResumeUploader, key string, hasKey bool, upToken string, upHost string, extra *RputExtra, ret interface{}) *resumeUploaderImpl {
	return &resumeUploaderImpl{
		client:   resumeUploader.Client,
		cfg:      resumeUploader.Cfg,
		key:      key,
		hasKey:   hasKey,
		upToken:  upToken,
		upHost:   upHost,
		extra:    extra,
		ret:      ret,
		fileSize: 0,
	}
}

func (impl *resumeUploaderImpl) initUploader(ctx context.Context) error {
	// Do nothing
	return nil
}

func (impl *resumeUploaderImpl) uploadChunk(ctx context.Context, c chunk) error {
	var (
		chunkSize int = impl.extra.ChunkSize
		apis          = impl.resumeUploaderAPIs()
		chunkData []byte
		blkPutRet BlkputRet
		err       error
	)

	for chunkOffset := 0; chunkOffset < len(c.data); chunkOffset += len(chunkData) {
		chunkData = c.data[chunkOffset:]
		actualChunkSize := len(chunkData)
		if actualChunkSize > chunkSize {
			actualChunkSize = chunkSize
		}
		chunkData = chunkData[:actualChunkSize]
	UploadSingleChunk:
		for retried := 0; retried < impl.extra.TryTimes; retried += 1 {
			if chunkOffset == 0 {
				err = apis.mkBlk(ctx, impl.upToken, impl.upHost, &blkPutRet, len(c.data), bytes.NewReader(chunkData), len(chunkData))
			} else {
				err = apis.bput(ctx, impl.upToken, &blkPutRet, bytes.NewReader(chunkData), len(chunkData))
			}
			if err != nil {
				if err == context.Canceled {
					break UploadSingleChunk
				}
				continue UploadSingleChunk
			}
			if blkPutRet.Crc32 != crc32.ChecksumIEEE(chunkData) || int(blkPutRet.Offset) != chunkOffset+len(chunkData) {
				err = ErrUnmatchedChecksum
				continue UploadSingleChunk
			}
			break UploadSingleChunk
		}
		if err != nil {
			impl.extra.NotifyErr(int(c.id), len(c.data), err)
			return err
		}
	}

	blkPutRet.blkIdx = int(c.id)
	impl.extra.Notify(blkPutRet.blkIdx, len(c.data), &blkPutRet)
	func() {
		impl.lock.Lock()
		defer impl.lock.Unlock()
		impl.extra.Progresses = append(impl.extra.Progresses, blkPutRet)
		impl.fileSize += int64(len(c.data))
	}()
	return nil
}

func (impl *resumeUploaderImpl) final(ctx context.Context) error {
	sort.Sort(BlkputRets(impl.extra.Progresses))
	return impl.resumeUploaderAPIs().mkfile(ctx, impl.upToken, impl.upHost, impl.ret, impl.key, impl.hasKey, impl.fileSize, impl.extra)
}

func (impl *resumeUploaderImpl) resumeUploaderAPIs() *resumeUploaderAPIs {
	return &resumeUploaderAPIs{Client: impl.client, Cfg: impl.cfg}
}

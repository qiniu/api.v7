package storage

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"sync"

	"github.com/qiniu/api.v7/v7"
)

// 分片上传过程中可能遇到的错误
var ErrUnmatchedChecksum = errors.New("unmatched checksum")

const (
	// 获取下一个分片Reader失败
	ErrNextReader = "ErrNextReader"
	// 超过了最大的重试上传次数
	ErrMaxUpRetry = "ErrMaxUpRetry"
)

const (
	blockBits = 22
	blockMask = (1 << blockBits) - 1
)

// Settings 为分片上传设置
type Settings struct {
	TaskQsize int   // 可选。任务队列大小。为 0 表示取 Workers * 4。
	Workers   int   // 并行 Goroutine 数目。
	ChunkSize int   // 默认的Chunk大小，不设定则为4M（仅在分片上传 v1 中使用）
	PartSize  int64 // 默认的Part大小，不设定则为4M（仅在分片上传 v2 中使用）
	TryTimes  int   // 默认的尝试次数，不设定则为3
}

// 分片上传默认参数设置
const (
	defaultWorkers   = 4               // 默认的并发上传的块数量
	defaultChunkSize = 4 * 1024 * 1024 // 默认的分片大小，4MB（仅在分片上传 v1 中使用）
	defaultPartSize  = 4 * 1024 * 1024 // 默认的分片大小，4MB（仅在分片上传 v2 中使用）
	defaultTryTimes  = 3               // mkblk / bput / uploadParts 失败重试次数
)

// 分片上传的默认设置
var settings = Settings{
	TaskQsize: defaultWorkers * 4,
	Workers:   defaultWorkers,
	ChunkSize: defaultChunkSize,
	PartSize:  defaultPartSize,
	TryTimes:  defaultTryTimes,
}

// SetSettings 可以用来设置分片上传参数
func SetSettings(v *Settings) {
	settings = *v
	if settings.Workers == 0 {
		settings.Workers = defaultWorkers
	}
	if settings.TaskQsize == 0 {
		settings.TaskQsize = settings.Workers * 4
	}
	if settings.ChunkSize == 0 {
		settings.ChunkSize = defaultChunkSize
	}
	if settings.PartSize == 0 {
		settings.PartSize = defaultPartSize
	}
	if settings.TryTimes == 0 {
		settings.TryTimes = defaultTryTimes
	}
}

var (
	tasks chan func()
	once  sync.Once
)

func worker(tasks chan func()) {
	for task := range tasks {
		task()
	}
}

func _initWorkers() {
	tasks = make(chan func(), settings.TaskQsize)
	for i := 0; i < settings.Workers; i++ {
		go worker(tasks)
	}
}

func initWorkers() {
	once.Do(_initWorkers)
}

// 代表一块分片的基本信息
type (
	chunk struct {
		id      int64
		offset  int64
		data    []byte
		retried int
	}

	// 代表一块分片的上传错误信息
	chunkError struct {
		chunk
		err error
	}

	// 通用分片上传接口，同时适用于分片上传 v1 和 v2 接口
	resumeUploaderBase interface {
		// 开始上传前调用一次用于初始化，在 v1 中该接口不做任何事情，而在 v2 中该接口对应 initParts
		initUploader(context.Context) ([]int64, error)
		// 上传实际的分片数据，允许并发上传。在 v1 中该接口对应 mkblk 和 bput 组合，而在 v2 中该接口对应 uploadParts
		uploadChunk(context.Context, chunk) error
		// 上传所有分片后调用一次用于结束上传，在 v1 中该接口对应 mkfile，而在 v2 中该接口对应 completeParts
		final(context.Context) error
	}

	// 将已知数据流大小的情况和未知数据流大小的情况抽象成一个接口
	chunkReader interface {
		readChunks([]int64, func(chunkID int64, off int64, data []byte) error) error
	}

	// 已知数据流大小的情况下读取数据流
	unsizedChunkReader struct {
		body      io.Reader
		blockSize int64
	}

	// 未知数据流大小的情况下读取数据流
	sizedChunkReader struct {
		body      io.ReaderAt
		totalSize int64
		blockSize int64
	}
)

// 使用并发 Goroutine 上传数据
func uploadByWorkers(uploader resumeUploaderBase, ctx context.Context, body chunkReader, tryTimes int) (err error) {
	var (
		wg           sync.WaitGroup
		failedChunks sync.Map
		recovered    []int64
	)

	initWorkers()

	if recovered, err = uploader.initUploader(ctx); err != nil {
		return
	}

	// 读取 Chunk 并创建任务
	err = body.readChunks(recovered, func(chunkID int64, off int64, data []byte) error {
		newChunk := chunk{id: chunkID, offset: off, data: data, retried: 0}
		wg.Add(1)
		tasks <- func() {
			defer wg.Done()
			if err := uploader.uploadChunk(ctx, newChunk); err != nil {
				newChunk.retried += 1
				failedChunks.LoadOrStore(newChunk.id, chunkError{chunk: newChunk, err: err})
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	for {
		// 等待这一轮任务执行完毕
		wg.Wait()
		// 不断重试先前失败的任务
		failedTasks := 0
		failedChunks.Range(func(key, chunkValue interface{}) bool {
			chunkErr := chunkValue.(chunkError)
			if chunkErr.err == context.Canceled {
				err = chunkErr.err
				return false
			}

			failedChunks.Delete(key)
			if chunkErr.retried < tryTimes {
				failedTasks += 1
				wg.Add(1)
				tasks <- func() {
					defer wg.Done()
					if cerr := uploader.uploadChunk(ctx, chunkErr.chunk); cerr != nil {
						chunkErr.retried += 1
						failedChunks.LoadOrStore(chunkErr.id, chunkErr)
					}
				}
				return true
			} else {
				err = chunkErr.err
				return false
			}
		})
		if err != nil {
			err = api.NewError(ErrMaxUpRetry, err.Error())
			return
		} else if failedTasks == 0 {
			break
		}
	}

	err = uploader.final(ctx)
	return
}

func newUnsizedChunkReader(body io.Reader, blockSize int64) *unsizedChunkReader {
	return &unsizedChunkReader{body: body, blockSize: blockSize}
}

func (r *unsizedChunkReader) readChunks(recovered []int64, f func(chunkID int64, off int64, data []byte) error) error {
	var (
		lastChunk          = false
		chunkID      int64 = 0
		off          int64 = 0
		chunkSize    int
		err          error
		recoveredMap = make(map[int64]struct{}, len(recovered))
	)

	for _, roff := range recovered {
		recoveredMap[roff] = struct{}{}
	}

	for !lastChunk {
		buf := make([]byte, r.blockSize)
		if chunkSize, err = io.ReadFull(r.body, buf); err != nil {
			switch err {
			case io.EOF:
				lastChunk = true
				return nil
			case io.ErrUnexpectedEOF:
				buf = buf[:chunkSize]
				lastChunk = true
			default:
				return api.NewError(ErrNextReader, err.Error())
			}
		}

		if _, ok := recoveredMap[off]; !ok {
			if err = f(chunkID, off, buf); err != nil {
				return err
			}
		}
		chunkID += 1
		off += int64(chunkSize)
	}
	return nil
}

func newSizedChunkReader(body io.ReaderAt, totalSize, blockSize int64) *sizedChunkReader {
	return &sizedChunkReader{body: body, totalSize: totalSize, blockSize: blockSize}
}

func (r *sizedChunkReader) readChunks(recovered []int64, f func(chunkID int64, off int64, data []byte) error) error {
	var (
		chunkID      int64 = 0
		off          int64 = 0
		buf          []byte
		err          error
		recoveredMap = make(map[int64]struct{}, len(recovered))
	)

	for _, roff := range recovered {
		recoveredMap[roff] = struct{}{}
	}

	for off < r.totalSize {
		shouldRead := r.totalSize - off
		if shouldRead > r.blockSize {
			shouldRead = r.blockSize
		}
		if _, ok := recoveredMap[off]; ok {
			off += shouldRead
		} else {
			if buf, err = ioutil.ReadAll(io.NewSectionReader(r.body, off, shouldRead)); err != nil {
				return err
			}
			if len(buf) == 0 { // 理论上应该读到 shouldRead 个字节，实际上为空，将直接结束该方法
				return nil
			}
			if err = f(chunkID, off, buf); err != nil {
				return err
			}
			off += int64(len(buf))
		}
		chunkID += 1
	}
	return nil
}

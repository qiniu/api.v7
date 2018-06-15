package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

var (
	accessKey = os.Getenv("QINIU_ACCESS_KEY")
	secretKey = os.Getenv("QINIU_SECRET_KEY")
	bucket    = os.Getenv("QINIU_TEST_BUCKET")
)

const FILESIZE = 1 << 24

type dummyReader struct{}

func (*dummyReader) Read(p []byte) (int, error) {
	return len(p), nil
}

func main() {

	key := "resumeUploadKey"

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	resumeUploader := storage.NewResumeUploader(&cfg)
	upToken := putPolicy.UploadToken(mac)
	ret := storage.PutRet{}
	fmt.Println("resume uploading, size:", FILESIZE)
	err := resumeUploader.PutWithoutSize(context.Background(), &ret, upToken,
		key, &io.LimitedReader{&dummyReader{}, FILESIZE}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ret.Key, ret.Hash)
}

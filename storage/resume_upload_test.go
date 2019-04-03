package storage

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
)

func TestResumeUploadPutFile(t *testing.T) {
	var putRet PutRet
	ctx := context.TODO()
	putPolicy := PutPolicy{
		Scope:           testBucket,
		DeleteAfterDays: 7,
	}
	upToken := putPolicy.UploadToken(mac)
	testKey := fmt.Sprintf("testRPutFileKey_%d", rand.Int())

	err := resumeUploader.PutFile(ctx, &putRet, upToken, testKey, testLocalFile, nil)
	if err != nil {
		t.Fatalf("ResumeUploader#PutFile() error, %s", err)
	}
	t.Logf("Key: %s, Hash:%s", putRet.Key, putRet.Hash)
}

func TestNextReader(t *testing.T) {

	type testUploader struct {
		length int64
		*uploader
	}

	blkSize := int64(1 << blockBits)

	uploaders := []testUploader{
		{
			uploader: &uploader{
				body:    strings.NewReader("hello world"),
				blkSize: blkSize,
			},
			length: 11,
		},
	}

	for _, up := range uploaders {
		up.init()
		_, n, _, err := up.nextReader()
		if err != io.EOF || int64(n) != up.length {
			t.Fatalf("nextReader(): %q\n", err)
		}
	}
}

func TestPutWithoutSize(t *testing.T) {
	var putRet PutRet
	putPolicy := PutPolicy{
		Scope:           testBucket,
		DeleteAfterDays: 7,
	}
	upToken := putPolicy.UploadToken(mac)
	testKey := fmt.Sprintf("testRPutFileKey_%d", rand.Int())

	err := resumeUploader.PutWithoutSize(context.Background(), &putRet, upToken, testKey, strings.NewReader("hello world"), nil)
	if err != nil {
		t.Fatalf("PutWithoutSize() error, %s", err)
	}
	t.Logf("Key: %s, Hash:%s", putRet.Key, putRet.Hash)
}

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

	tu1 := testUploader{
		uploader: &uploader{
			body:    strings.NewReader("hello world"),
			blkSize: blkSize,
		},
		length: 11,
	}

	tu1.init()
	_, n, _, err := tu1.nextReader()
	if err != io.EOF || int64(n) != tu1.length {
		t.Fatalf("nextReader(): %q\n", err)
	}

	tu2 := testUploader{
		uploader: &uploader{
			body: &NotSeekerReader{
				Reader: strings.NewReader(strings.Repeat("hello", 1<<blockBits)),
			},
			blkSize: blkSize,
		},
		length: 5 * blkSize,
	}
	tu2.init()

	for i := 0; i < 4; i++ {
		_, n, _, err = tu2.nextReader()
		if err != nil || int64(n) != blkSize {
			t.Fatalf("nextReader(): %q, n: %d\n", err, n)
		}
	}
	_, n, _, err = tu2.nextReader()
	if err != io.EOF && err != nil {
		t.Fatalf("nextReader(): %q\n", err)
	}
}

type NotSeekerReader struct {
	io.Reader
}

func (r *NotSeekerReader) Read(p []byte) (n int, err error) {
	return r.Reader.Read(p)
}

func NewNotSeekerReader(r io.Reader) *NotSeekerReader {
	return &NotSeekerReader{
		Reader: r,
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

	rds := []io.Reader{
		strings.NewReader("hello world"),
		strings.NewReader(strings.Repeat("he", 1<<blockBits)),
		NewNotSeekerReader(strings.NewReader(strings.Repeat("test", 1<<blockBits))),
	}

	for _, rd := range rds {

		err := resumeUploader.PutWithoutSize(context.Background(), &putRet, upToken, testKey, rd, nil)
		if err != nil {
			t.Fatalf("PutWithoutSize() error, %s", err)
		}
		t.Logf("Key: %s, Hash:%s", putRet.Key, putRet.Hash)
	}
}

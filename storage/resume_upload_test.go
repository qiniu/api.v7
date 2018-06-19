package storage

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
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
	testKey := fmt.Sprintf("testRPutFile_%d", rand.Int())
	localFile := "/tmp/" + testKey
	f, _ := os.Create(localFile)
	f.WriteString("12345")
	f.Close()

	err := resumeUploader.PutFile(ctx, &putRet, upToken, testKey, localFile, nil)
	os.Remove(localFile)
	if err != nil {
		t.Fatalf("ResumeUploader#PutFile() error, %s", err)
	}
	if putRet.Key != testKey {
		t.Fatal(putRet)
	}
}

func TestResumeUploadPutReader(t *testing.T) {
	var putRet PutRet
	ctx := context.TODO()
	putPolicy := PutPolicy{
		Scope:           testBucket,
		DeleteAfterDays: 7,
	}
	upToken := putPolicy.UploadToken(mac)
	testKey := fmt.Sprintf("testRPutReader_%d", rand.Int())

	err := resumeUploader.PutReader(ctx, &putRet, upToken, testKey, strings.NewReader("test"), nil)
	if err != nil {
		t.Fatalf("ResumeUploader#PutReader() error, %s", err)
	}
	if putRet.Key != testKey {
		t.Fatal(putRet)
	}
}

func TestResumeUploadPutReaderWithoutKey(t *testing.T) {
	var putRet PutRet
	ctx := context.TODO()
	putPolicy := PutPolicy{
		Scope:           testBucket,
		DeleteAfterDays: 7,
	}
	upToken := putPolicy.UploadToken(mac)

	err := resumeUploader.PutReaderWithoutKey(ctx, &putRet, upToken, strings.NewReader("test"), nil)
	if err != nil {
		t.Fatalf("ResumeUploader#PutReaderWithoutKey() error, %s", err)
	}
	if putRet.Key != putRet.Hash {
		t.Fatal(putRet)
	}
}

type readAtAdapter struct {
	io.Reader
}

func (r *readAtAdapter) ReadAt(p []byte, off int64) (int, error) {
	return r.Read(p)
}

func TestResumeUploadPut(t *testing.T) {
	var putRet PutRet
	ctx := context.TODO()
	putPolicy := PutPolicy{
		Scope:           testBucket,
		DeleteAfterDays: 7,
	}
	upToken := putPolicy.UploadToken(mac)
	testKey := fmt.Sprintf("testRPut_%d", rand.Int())

	err := resumeUploader.Put(ctx, &putRet, upToken, testKey, &readAtAdapter{strings.NewReader("test")}, 4, nil)
	if err != nil {
		t.Fatalf("ResumeUploader#Put() error, %s", err)
	}
	if putRet.Key != testKey {
		t.Fatal(putRet)
	}
}

func TestResumeUploadPutWithoutKey(t *testing.T) {
	var putRet PutRet
	ctx := context.TODO()
	putPolicy := PutPolicy{
		Scope:           testBucket,
		DeleteAfterDays: 7,
	}
	upToken := putPolicy.UploadToken(mac)

	err := resumeUploader.PutWithoutKey(ctx, &putRet, upToken, &readAtAdapter{strings.NewReader("test")}, 4, nil)
	if err != nil {
		t.Fatalf("ResumeUploader#PutWithoutKey() error, %s", err)
	}
	if putRet.Key != putRet.Hash {
		t.Fatal(putRet)
	}
}

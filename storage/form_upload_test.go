package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestFormUploadPutFile(t *testing.T) {
	var putRet PutRet
	ctx := context.TODO()
	putPolicy := PutPolicy{
		Scope:           testBucket,
		DeleteAfterDays: 7,
	}

	// prepare file for test uploading
	testLocalFile, err := ioutil.TempFile("", "TestFormUploadPutFile")
	if err != nil {
		t.Fatalf("ioutil.TempFile file failed, err: %v", err)
	}
	defer os.Remove(testLocalFile.Name())

	upToken := putPolicy.UploadToken(mac)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	testKey := fmt.Sprintf("testPutFileKey_%d", r.Int())

	err = formUploader.PutFile(ctx, &putRet, upToken, testKey, testLocalFile.Name(), nil)
	if err != nil {
		t.Fatalf("FormUploader#PutFile() error, %s", err)
	}
	t.Logf("Key: %s, Hash:%s", putRet.Key, putRet.Hash)
}

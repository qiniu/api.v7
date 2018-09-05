package storage

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

var (
	testLocalFile string
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	exPath := filepath.Dir(pwd)
	testLocalFile = filepath.Join(exPath, "Makefile")
}

func TestFormUploadPutFile(t *testing.T) {
	var putRet PutRet
	ctx := context.TODO()
	putPolicy := PutPolicy{
		Scope:           testBucket,
		DeleteAfterDays: 7,
	}
	upToken := putPolicy.UploadToken(mac)
	testKey := fmt.Sprintf("testPutFileKey_%d", rand.Int())

	err := formUploader.PutFile(ctx, &putRet, upToken, testKey, testLocalFile, nil)
	if err != nil {
		t.Fatalf("FormUploader#PutFile() error, %s", err)
	}
	t.Logf("Key: %s, Hash:%s", putRet.Key, putRet.Hash)
}

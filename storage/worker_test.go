package storage

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func doWork(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

func TestWorkerCopy(t *testing.T) {
	wg := sync.WaitGroup{}
	var initOnce sync.Once
	workers := 10
	var tasks chan func()
	initOnce.Do(func() {
		tasks = make(chan func(), workers)
		for i := 0; i < workers; i++ {
			go doWork(tasks)
		}
	})

	m := NewBucketManager(mac, nil)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		tasks <- func() {
			defer wg.Done()
			err := m.Copy(testBucket, "qiniu.png", testBucket, fmt.Sprintf("test_%d", i), true)
			t.Log(err)
		}
	}
	wg.Wait()
}

func TestWorkerUpload(t *testing.T) {
	wg := sync.WaitGroup{}
	var initOnce sync.Once
	workers := 10
	var tasks chan func()
	initOnce.Do(func() {
		tasks = make(chan func(), workers)
		for i := 0; i < workers; i++ {
			go doWork(tasks)
		}
	})

	uploader := NewResumeUploader(nil)
	ctx := context.Background()

	for i := 0; i < 20; i++ {
		wg.Add(1)

		tasks <- func() {
			defer wg.Done()

			testKey := fmt.Sprintf("testPutFileKey_%d", rand.Int())
			t.Logf("start to upload %s ...", testKey)

			var putRet PutRet
			putPolicy := PutPolicy{
				Scope:           testBucket + ":" + testKey,
				DeleteAfterDays: 7,
			}
			upToken := putPolicy.UploadToken(mac)
			err := uploader.PutFile(ctx, &putRet, upToken, testKey, testLocalFile, nil)
			if err != nil {
				t.Fatalf("TestWorkerUpload error, %s", err)
			}

			t.Logf("upload success, key: %s, hash:%s", putRet.Key, putRet.Hash)
		}
	}

	wg.Wait()

}

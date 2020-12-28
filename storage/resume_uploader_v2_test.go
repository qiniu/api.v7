package storage

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestPutWithoutSizeV2(t *testing.T) {
	var putRet PutRet
	putPolicy := PutPolicy{
		Scope:           testBucket,
		DeleteAfterDays: 7,
	}
	upToken := putPolicy.UploadToken(mac)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sizes := []int64{
		64,
		1 << blockBits,
		(1 << blockBits) - 1,
		(1 << blockBits) + 1,
		(1 << (blockBits + 4)) + 1,
	}
	partSizes := []int64{0, 1 << 20, 1 << 24}

	for _, partSize := range partSizes {
		for _, size := range sizes {
			md5Sumer := md5.New()
			rd := io.TeeReader(&io.LimitedReader{R: r, N: size}, md5Sumer)
			testKey := fmt.Sprintf("testRPutFileV2Key_%d", r.Int())
			err := resumeUploaderV2.PutWithoutSize(context.Background(), &putRet, upToken, testKey, rd, &RputV2Extra{
				PartSize: partSize,
				Notify: func(partNumber int64, ret *UploadPartsRet) {
					t.Logf("Notify: partNumber: %d, ret: %#v", partNumber, ret)
				},
				NotifyErr: func(partNumber int64, err error) {
					t.Logf("NotifyErr: partNumber: %d, err: %s", partNumber, err)
				},
			})
			if err != nil {
				t.Fatalf("PutWithoutSize() error, %s", err)
			}
			md5Value := hex.EncodeToString(md5Sumer.Sum(nil))
			validateMD5(t, testKey, md5Value, size)
			t.Logf("Size: %d, Part: %d, Key: %s, Hash:%s", size, partSize, putRet.Key, putRet.Hash)
		}
	}
}

func TestPutWithSizeV2(t *testing.T) {
	var putRet PutRet
	putPolicy := PutPolicy{
		Scope:           testBucket,
		DeleteAfterDays: 7,
	}
	upToken := putPolicy.UploadToken(mac)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sizes := []int64{
		64,
		1 << blockBits,
		(1 << blockBits) - 1,
		(1 << blockBits) + 1,
		(1 << (blockBits + 4)) + 1,
	}
	partSizes := []int64{0, 1 << 20, 1 << 24}

	for _, partSize := range partSizes {
		for _, size := range sizes {
			data := make([]byte, size)
			if _, err := io.ReadFull(r, data); err != nil {
				t.Error(err)
			}
			testKey := fmt.Sprintf("testRPutFileV2Key_%d", r.Int())
			err := resumeUploaderV2.Put(context.Background(), &putRet, upToken, testKey, bytes.NewReader(data), size, &RputV2Extra{
				PartSize: partSize,
				Notify: func(partNumber int64, ret *UploadPartsRet) {
					t.Logf("Notify: partNumber: %d, ret: %#v", partNumber, ret)
				},
				NotifyErr: func(partNumber int64, err error) {
					t.Logf("NotifyErr: partNumber: %d, err: %s", partNumber, err)
				},
			})
			if err != nil {
				t.Fatalf("Put() error, %s", err)
			}
			md5ByteArray := md5.Sum(data)
			md5Value := hex.EncodeToString(md5ByteArray[:])
			validateMD5(t, testKey, md5Value, size)
			t.Logf("Size: %d, Part: %d, Key: %s, Hash:%s", size, partSize, putRet.Key, putRet.Hash)
		}
	}
	for _, partSize := range partSizes {
		for _, size := range sizes {
			md5Sumer := md5.New()
			testKey := fmt.Sprintf("testRPutFileV2Key_%d_*", r.Int())
			tmpFile, err := ioutil.TempFile("", testKey)
			if err != nil {
				t.Error(err)
			}
			if _, err = io.CopyN(tmpFile, io.TeeReader(r, md5Sumer), size); err != nil {
				t.Error(err)
			} else if err = tmpFile.Close(); err != nil {
				t.Error(err)
			}
			err = resumeUploaderV2.PutFile(context.Background(), &putRet, upToken, testKey, tmpFile.Name(), &RputV2Extra{
				PartSize: partSize,
				Notify: func(partNumber int64, ret *UploadPartsRet) {
					t.Logf("Notify: partNumber: %d, ret: %#v", partNumber, ret)
				},
				NotifyErr: func(partNumber int64, err error) {
					t.Logf("NotifyErr: partNumber: %d, err: %s", partNumber, err)
				},
			})
			if err != nil {
				t.Fatalf("PutFile() error, %s", err)
			}
			md5ByteArray := md5Sumer.Sum(nil)
			md5Value := hex.EncodeToString(md5ByteArray[:])
			validateMD5(t, testKey, md5Value, size)
			t.Logf("Size: %d, Part: %d, Key: %s, Hash:%s", size, partSize, putRet.Key, putRet.Hash)
		}
	}
}

func TestPutWithoutKeyV2(t *testing.T) {
	var putRet PutRet
	putPolicy := PutPolicy{
		Scope:           testBucket,
		DeleteAfterDays: 7,
	}
	upToken := putPolicy.UploadToken(mac)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	data := make([]byte, 64)
	if _, err := io.ReadFull(r, data); err != nil {
		t.Error(err)
	}
	err := resumeUploaderV2.PutWithoutKey(context.Background(), &putRet, upToken, bytes.NewReader(data), int64(len(data)), &RputV2Extra{
		Notify: func(partNumber int64, ret *UploadPartsRet) {
			t.Logf("Notify: partNumber: %d, ret: %#v", partNumber, ret)
		},
		NotifyErr: func(partNumber int64, err error) {
			t.Logf("NotifyErr: partNumber: %d, err: %s", partNumber, err)
		},
	})
	if err != nil {
		t.Fatalf("PutWithoutKey() error, %s", err)
	}
	t.Logf("Key: %s, Hash:%s", putRet.Key, putRet.Hash)

	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}
	if _, err = io.Copy(tmpFile, bytes.NewReader(data)); err != nil {
		t.Error(err)
	} else if err = tmpFile.Close(); err != nil {
		t.Error(err)
	}
	err = resumeUploaderV2.PutFileWithoutKey(context.Background(), &putRet, upToken, tmpFile.Name(), &RputV2Extra{
		Notify: func(partNumber int64, ret *UploadPartsRet) {
			t.Logf("Notify: partNumber: %d, ret: %#v", partNumber, ret)
		},
		NotifyErr: func(partNumber int64, err error) {
			t.Logf("NotifyErr: partNumber: %d, err: %s", partNumber, err)
		},
	})
	if err != nil {
		t.Fatalf("PutWithoutKey() error, %s", err)
	}
	t.Logf("Key: %s, Hash:%s", putRet.Key, putRet.Hash)
}
func TestPutWithRecoveryV2(t *testing.T) {
	var putRet PutRet
	putPolicy := PutPolicy{
		Scope:           testBucket,
		DeleteAfterDays: 7,
	}
	upToken := putPolicy.UploadToken(mac)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	testKey := fmt.Sprintf("testRPutFileKey_%d", r.Int())
	dirName, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirName)
	recorder, err := NewFileRecorder(dirName)
	if err != nil {
		t.Fatal(err)
	}

	fileName := filepath.Join(dirName, "originalFile")
	testFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer testFile.Close()

	size := int64(4 * (1 << blockBits))
	io.CopyN(testFile, r, size)

	for i := 0; i < 4; i++ {
		ctx, cancelFunc := context.WithCancel(context.Background())
		counter := uint32(0)
		err = resumeUploaderV2.PutFile(ctx, &putRet, upToken, testKey, fileName, &RputV2Extra{
			Recorder: recorder,
			Notify: func(partNumber int64, ret *UploadPartsRet) {
				t.Logf("Notify: partNumber: %d, ret: %#v", partNumber, ret)
				if atomic.AddUint32(&counter, 1) >= 2 {
					cancelFunc()
				}
			},
			NotifyErr: func(partNumber int64, err error) {
				t.Logf("NotifyErr: partNumber: %d, err: %s", partNumber, err)
			},
		})
		if err == nil {
			return
		}
	}
	t.Fatal(err)
}

func TestPutWithEmptyKeyV2(t *testing.T) {
	bucketManager.Delete(testBucket, "")

	var putRet PutRet
	putPolicy := PutPolicy{
		Scope:           testBucket + ":",
		DeleteAfterDays: 7,
	}
	upToken := putPolicy.UploadToken(mac)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	data := make([]byte, 64)
	if _, err := io.ReadFull(r, data); err != nil {
		t.Error(err)
	}
	err := resumeUploaderV2.Put(context.Background(), &putRet, upToken, "", bytes.NewReader(data), int64(len(data)), &RputV2Extra{
		Notify: func(partNumber int64, ret *UploadPartsRet) {
			t.Logf("Notify: partNumber: %d, ret: %#v", partNumber, ret)
		},
		NotifyErr: func(partNumber int64, err error) {
			t.Logf("NotifyErr: partNumber: %d, err: %s", partNumber, err)
		},
	})
	if err != nil {
		t.Fatalf("Put() error, %s", err)
	}
	t.Logf("Key: %s, Hash:%s", putRet.Key, putRet.Hash)
}

package storage

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
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
		1 << blockBits,
		(1 << blockBits) - 1,
		(1 << blockBits) + 1,
		(1 << blockBits) + 1,
		(1 << (blockBits + 2)) + 1,
		(1 << (blockBits + 4)) + 1,
	}
	partSizes := []int64{0, 1 << 20, 1 << 24}

	for _, partSize := range partSizes {
		for _, size := range sizes {
			md5Sumer := md5.New()
			rd := io.TeeReader(&io.LimitedReader{R: r, N: size}, md5Sumer)
			testKey := fmt.Sprintf("testRPutFileV2Key_%d", rand.Int())
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
		1 << blockBits,
		(1 << blockBits) - 1,
		(1 << blockBits) + 1,
		(1 << blockBits) + 1,
		(1 << (blockBits + 2)) + 1,
		(1 << (blockBits + 4)) + 1,
	}
	partSizes := []int64{0, 1 << 20, 1 << 24}

	for _, partSize := range partSizes {
		for _, size := range sizes {
			data := make([]byte, size)
			if _, err := io.ReadFull(r, data); err != nil {
				t.Error(err)
			}
			testKey := fmt.Sprintf("testRPutFileV2Key_%d", rand.Int())
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

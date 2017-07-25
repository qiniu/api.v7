package storage

import (
	"github.com/qiniu/api.v7/auth/qbox"
	"os"
	"testing"
)

var (
	testAK     = os.Getenv("QINIU_ACCESS_KEY")
	testSK     = os.Getenv("QINIU_SECRET_KEY")
	testBucket = os.Getenv("QINIU_TEST_BUCKET")
	testKey    = "qiniu.jpg"
)

var mac *qbox.Mac
var bucketManager *BucketManager

func init() {
	if testAK == "" || testSK == "" {
		panic("please run ./test-env.sh first")
	}
	mac = qbox.NewMac(testAK, testSK)
	cfg := Config{}
	bucketManager = NewBucketManager(mac, &cfg)
}

//Test get zone
func TestGetZone(t *testing.T) {
	zone, err := GetZone(testAK, testBucket)
	if err != nil {
		t.Fatalf("GetZone() error, %s", err)
	}
	t.Log(zone.String())
}

//Test get bucket list
func TestBuckets(t *testing.T) {
	shared := true
	buckets, err := bucketManager.Buckets(shared)
	if err != nil {
		t.Fatalf("Buckets() error, %s", err)
	}

	for _, bucket := range buckets {
		t.Log(bucket)
	}
}

//Test get file info
func TestStat(t *testing.T) {
	keysToStat := []string{"qiniu.jpg", "qiniu.png", "qiniu.mp4"}
	for _, eachKey := range keysToStat {
		info, err := bucketManager.Stat(testBucket, eachKey)
		if err != nil {
			t.Logf("Stat() error, %s", err)
			t.Fail()
		} else {
			t.Logf("FileInfo:\n %s", info.String())
		}
	}
}

func TestDelete(t *testing.T) {
	keysToDelete := []string{"qiniu_1.jpg", "qiniu_2.jpg", "qiniu_3.jpg"}
	for _, eachKey := range keysToDelete {
		err := bucketManager.Copy(testBucket, testKey, testBucket, eachKey, true)
		if err != nil {
			t.Logf("Copy() error, %s", err)
			t.Fail()
		}
	}

	for _, eachKey := range keysToDelete {
		err := bucketManager.Delete(testBucket, eachKey)
		if err != nil {
			t.Logf("Delete() error, %s", err)
			t.Fail()
		}
	}

}

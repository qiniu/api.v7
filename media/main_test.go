package media

import (
	"qiniupkg.com/api.v7/conf"
	"testing"
	"os"
)

func init()  {
	conf.ACCESS_KEY=os.Getenv("QINIU_ACCESS_KEY")
	conf.SECRET_KEY=os.Getenv("QINIU_SECRET_KEY")
	Bucket=os.Getenv("QINIU_TEST_BUCKET")
	Pipline = "default-pipline"
}

func TestAvthumb(t *testing.T) {
	view := NewAvthumb()
	view.Format = "mp4"
	options := Options{NeedConvertFileName: "320ae9957cb2a82416d9a843903ce32c.m3u8"}
	res, err := view.Avthumb(options)
	if err!=nil {
		t.Log("Avthumb Err:", err)
	}else{
		t.Log("persistentId:",res.PersistentId)
	}
}

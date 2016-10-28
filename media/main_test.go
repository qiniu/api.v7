package media

import (
	"os"
	"qiniupkg.com/api.v7/conf"
	"testing"
)

func init() {
	conf.ACCESS_KEY = os.Getenv("QINIU_ACCESS_KEY")
	conf.SECRET_KEY = os.Getenv("QINIU_SECRET_KEY")
	Bucket = os.Getenv("QINIU_TEST_BUCKET")
	Pipline = "default-pipline"
}

func TestAvthumb(t *testing.T) {
	view := NewAvthumb()
	view.Format = "mp4"
	options := Options{NeedConvertFileName: "320ae9957cb2a82416d9a843903ce32c.mp4"}
	res, err := view.Avthumb(options)
	if err != nil {
		t.Log("Avthumb Err:", err)
	} else {
		t.Log("persistentId:", res.PersistentId)
	}
}

func TestAvSegtime(t *testing.T)  {
	view := NewAvSegtime()
	view.Format = "m3u8"
	options := Options{NeedConvertFileName: "320ae9957cb2a82416d9a843903ce32c.m3u8"}
	res, err := view.AvSegtime(options)
	if err != nil {
		t.Log("Avthumb Err:", err)
	} else {
		t.Log("persistentId:", res.PersistentId)
	}
}

func TestAvConcat(t *testing.T)  {
	view:=NewAvConcat()
	view.Mode="2"
	view.Format="mp4"
	view.Urls=append(view.Urls,"code/v6/api/dora-api/av/avconcat.html")
	view.Urls=append(view.Urls,"code/v6/api/dora-api/av/avconcat.html")
	options:=Options{NeedConvertFileName:"thinking-in-go.1.mp4"}
	res,err:=view.AvConcat(options)
	if err != nil {
		t.Log("AvConcat Err:", err)
	} else {
		t.Log("persistentId:", res.PersistentId)
	}
}

func TestAvinfo(t *testing.T) {
	res, err := Avinfo("http://78re52.com1.z0.glb.clouddn.com/resource%2Fthinking-in-go.mp4")
	if err != nil {
		t.Log("Avinfo Err:", err)
	} else {
		t.Log("Avinfo Res:", res)
	}
}

func TestVFrame(t *testing.T) {
	view:=NewVFrame()
	view.Format="jpg"
	view.Offset="1"
	view.Width="11"
	view.Height="22"
	options:=Options{NeedConvertFileName:"thinking-in-go.1.mp4"}
	res,err:=view.VFrame(options)
	if err != nil {
		t.Log("VFrame Err:", err)
	} else {
		t.Log("persistentId:", res.PersistentId)
	}
}
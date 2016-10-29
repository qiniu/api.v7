package media

import (
	"os"
	"qiniupkg.com/api.v7/conf"
	"testing"
	"time"
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

func TestVSample(t *testing.T) {
	view:=NewVSample()
	view.Format="jpg"
	options:=Options{NeedConvertFileName:"thinking-in-go.1.mp4"}
	res,err:=view.VSample(options)
	if err != nil {
		t.Log("VSample Err:", err)
	} else {
		t.Log("persistentId:", res.PersistentId)
	}
}

func TestPrivateM3U8(t *testing.T) {
	view:=NewPrivateM3U8()
	res,err:=view.Download("http://7xoucz.com1.z0.glb.clouddn.com/qiniutest.m3u8",time.Now().Unix()+1200)
	if err != nil {
		t.Log("VSample Err:", err)
	} else {
		t.Log("body:", res.Body)
	}
}

func TestAvvod(t *testing.T)  {
	view:=NewAvvod()
	res:=view.Avvod("http://7xvilo.com1.z0.glb.clouddn.com/%E4%B8%83%E7%89%9B%E4%BA%91%E5%AD%98%E5%82%A8%E8%A7%86%E9%A2%91%EF%BC%8D%E4%B8%89%E5%91%A8%E5%B9%B4%20.mp4")
	t.Log("res:",res)
}

func TestAdapt(t *testing.T)  {
	view := NewAdapt()
	options := Options{NeedConvertFileName: "4k.mp4"}
	res, err := view.Adapt(options)
	if err != nil {
		t.Log("Adapt Err:", err)
	} else {
		t.Log("persistentId:", res.PersistentId)
	}
}

func TestImageInfo(t *testing.T)  {
	res,err:=GetImageInfo("http://78re52.com1.z0.glb.clouddn.com/resource/gogopher.jpg")
	if err != nil {
		t.Log("ImageInfo Err:", err)
	} else {
		t.Log("res:", res)
	}
}

func TestImageExif(t *testing.T)  {
	res,err:=GetImageExif("http://78re52.com1.z0.glb.clouddn.com/resource/gogopher.jpg")
	if err != nil {
		t.Log("ImageExif Err:", err)
	} else {
		t.Log("res:", res)
	}
}

func TestImageAVE(t *testing.T)  {
	res,err:=GetImageAVE("http://78re52.com1.z0.glb.clouddn.com/resource/gogopher.jpg")
	if err != nil {
		t.Log("ImageAVE Err:", err)
	} else {
		t.Log("res:", res)
	}
}

func TestTextConcat(t *testing.T)  {
	options:=Options{NeedConvertFileName:"a.txt"}
	view:=NewTextConcatView("text")
	view.Urls=append(view.Urls,"http://test.clouddn.com/b.txt")
	view.Urls=append(view.Urls,"http://test.clouddn.com/c.txt")
	res,err:=view.TextConcat(options)
	if err != nil {
		t.Log("TextConcat Err:", err)
	} else {
		t.Log("persistentId:", res.PersistentId)
	}
}
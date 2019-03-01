package linking

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestVodToken(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	device := "sdk-testVodtoken"
	defer c.DeleteDevice(testApp, device)
	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev)
	noError(t, err)
	time.Sleep(time.Second)
	token, err := c.VodToken(testApp, device, time.Now().Add(time.Hour*5).Unix())
	noError(t, err)
	fmt.Println(token)
	// make a request and return code should NOT be 401
	url := c.url("/device/resource/playback.m3u8?dtoken=%s&start=%d&end=%d", token, time.Now().Add(-time.Hour).Unix(), time.Now().Add(time.Hour).Unix())
	resp, err := http.Get(url)
	shouldBeEqual(t, 612, resp.StatusCode)
	noError(t, err)
}

func TestStatusToken(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	device := "sdk-testStatusToken"
	defer c.DeleteDevice(testApp, device)
	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev)
	noError(t, err)
	time.Sleep(time.Second)
	token, err := c.VodToken(testApp, device, time.Now().Add(time.Hour*5).Unix())
	noError(t, err)
	fmt.Println(token)
	// make a request and return code should NOT be 401
	url := c.url("/device/resource/status?dtoken=%s&start=%d&end=%d", token, time.Now().Add(-time.Hour).Unix(), time.Now().Add(time.Hour).Unix())
	resp, err := http.Get(url)
	shouldBeEqual(t, 200, resp.StatusCode)
	noError(t, err)

}

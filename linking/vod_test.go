package linking

import (
	"fmt"
	"testing"
	"time"
)

// 这个测试case需要保证最近1个小时ts文件在上传
func TestSegments(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	// delete if exist
	device := "sdk-testSegmentsDevice"
	defer c.DeleteDevice(testApp, device)
	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev)
	noError(t, err)

	end := time.Now().Unix()
	start := time.Now().Add(-time.Hour).Unix()
	segs, marker, err := c.Segments(testApp, device, int(start), int(end), "", 1000)
	noError(t, err)
	fmt.Printf("segs = %#v, marker = %#v", segs, marker)
}

// 这个测试case需要保证最近1个小时ts文件在上传
func TestSaveas(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	// delete if exist
	device := "sdk-testSaveasDevice"
	defer c.DeleteDevice(testApp, device)
	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev)
	noError(t, err)

	end := time.Now().Unix()
	start := time.Now().Add(-time.Hour).Unix()
	saveasReply, _ := c.Saveas(testApp, device, int(start), int(end), "testSaveas.mp4", "mp4")
	fmt.Printf("saveas reply = %#v", saveasReply)
}

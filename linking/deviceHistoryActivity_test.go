package linking

import (
	"fmt"
	"testing"
	"time"
)

// 这个测试case需要保证最近1个小时设备有上下线操作
func TestHistoryActivity(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	// delete if exist
	device := "sdk-testHistActivityDevice"
	defer c.DeleteDevice(testApp, device)
	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev)
	noError(t, err)
	end := time.Now().Unix()
	start := time.Now().Add(-time.Hour).Unix()
	segs, marker, err := c.ListDeviceHistoryactivity(testApp, device, int(start), int(end), "", 1000)
	noError(t, err)
	fmt.Printf("segs = %#v, marker = %#v", segs, marker)
}

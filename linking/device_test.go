package linking

import (
	"testing"
)

func TestAddDevice(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	// delete if exist
	device := "sdk-testAddDevice"
	defer c.DeleteDevice(testApp, device)
	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev)
	noError(t, err)
}

func TestQueryDevice(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	device := "sdk-testQueryDevice"
	defer c.DeleteDevice(testApp, device)
	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev)
	noError(t, err)
	dev2, err := c.QueryDevice(testApp, device)
	noError(t, err)
	shouldBeEqual(t, 7, dev2.SegmentExpireDays)
}

func TestUpdateDevice(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	device := "sdk-testUpdateDevice"
	defer c.DeleteDevice(testApp, device)

	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev)
	noError(t, err)
	dev2, err := c.QueryDevice(testApp, device)
	noError(t, err)
	shouldBeEqual(t, 7, dev2.SegmentExpireDays)

	// udpate device segmentexpiredays to 30
	ops := []PatchOperation{
		PatchOperation{Op: "replace", Key: "segmentExpireDays", Value: 30},
	}
	dev3, err := c.UpdateDevice(testApp, device, ops)
	noError(t, err)
	shouldBeEqual(t, 30, dev3.SegmentExpireDays)
}

func TestListDevice(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	device1 := "sdk-testListDevice1"
	defer c.DeleteDevice(testApp, device1)

	dev1 := &Device{
		Device:            device1,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev1)
	noError(t, err)

	device2 := "sdk-testListDevice2"
	defer c.DeleteDevice(testApp, device2)

	dev2 := &Device{
		Device:            device2,
		SegmentExpireDays: 7,
	}
	_, err = c.AddDevice(testApp, dev2)
	noError(t, err)

	device3 := "sdk-testListDevice3"
	defer c.DeleteDevice(testApp, device3)

	dev3 := &Device{
		Device:            device3,
		SegmentExpireDays: 7,
	}
	_, err = c.AddDevice(testApp, dev3)
	noError(t, err)

	devices, marker, err := c.ListDevice(testApp, "sdk-testListDevice", "", 2, false, false, 0, "")
	noError(t, err)
	shouldBeEqual(t, 2, len(devices))

	devices, marker, err = c.ListDevice(testApp, "sdk-testListDevice", "", 1000, false, false, 0, "")
	noError(t, err)
	shouldBeEqual(t, 3, len(devices))
	shouldBeEqual(t, "", marker)
}
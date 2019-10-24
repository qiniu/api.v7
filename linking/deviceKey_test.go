package linking

import (
	"testing"
)

func TestAddDeviceKey(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	// delete if exist
	device := "sdk-testAddDeviceKey"
	defer c.DeleteDevice(testApp, device)
	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev)
	noError(t, err)
	keys, err := c.AddDeviceKey(testApp, device)
	noError(t, err)
	shouldBeEqual(t, 2, len(keys))
}

func TestQueryDeviceKey(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	// delete if exist
	device := "sdk-testQueryDeviceKey"
	defer c.DeleteDevice(testApp, device)
	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev)
	if err != nil {
		t.Fatalf("Added Device Failed error = %s", err.Error())
	}

	keys, err := c.QueryDeviceKey(testApp, device)
	if err != nil {
		t.Fatalf("Query DeviceKey error = %s", err.Error())
	}
	if 1 != len(keys) {
		t.Fatal("length of keys should be 1")
	}
}

func TestUpdateDeviceKey(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	// delete if exist
	device := "sdk-testUpdateDeviceKey"
	defer c.DeleteDevice(testApp, device)
	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}

	_, err := c.AddDevice(testApp, dev)
	if err != nil {
		t.Fatalf("Added Device Failed error = %s", err.Error())
	}
	keys, err := c.QueryDeviceKey(testApp, device)
	noError(t, err)
	shouldBeEqual(t, 1, len(keys))
	shouldBeEqual(t, 0, keys[0].State)
	dak := keys[0].AccessKey
	err = c.UpdateDeviceKeyState(testApp, device, dak, 1)
	noError(t, err)
	// check key state == 1
	keys, err = c.QueryDeviceKey(testApp, device)
	noError(t, err)
	shouldBeEqual(t, 1, keys[0].State)
}

func TestDeleteDeviceKey(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	// delete if exist
	device := "sdk-testDeleteDeviceKey"
	defer c.DeleteDevice(testApp, device)
	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev)
	noError(t, err)

	keys, err := c.AddDeviceKey(testApp, device)
	noError(t, err)
	shouldBeEqual(t, 2, len(keys))
	dak := keys[0].AccessKey
	err = c.DeleteDeviceKey(testApp, device, dak)
	keys, err = c.QueryDeviceKey(testApp, device)
	noError(t, err)
	shouldBeEqual(t, 1, len(keys))
}

func TestCloneDeviceKey(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	// delete if exist
	device1 := "sdk-testCloneDeviceKey1"
	defer c.DeleteDevice(testApp, device1)
	dev1 := &Device{
		Device:            device1,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev1)
	noError(t, err)
	device2 := "sdk-testCloneDeviceKey2"
	defer c.DeleteDevice(testApp, device2)
	dev2 := &Device{
		Device:            device2,
		SegmentExpireDays: 7,
	}
	_, err = c.AddDevice(testApp, dev2)
	noError(t, err)

	// query device1 keys
	keys, err := c.QueryDeviceKey(testApp, device1)
	dak1 := keys[0].AccessKey

	// clone device1 key ----> device2
	keys, err = c.QueryDeviceKey(testApp, device2)

	noError(t, err)
	shouldBeEqual(t, 1, len(keys))
	keys, err = c.CloneDeviceKey(testApp, device1, device2, false, false, dak1)
	noError(t, err)
	keys, err = c.QueryDeviceKey(testApp, device2)

	noError(t, err)
	shouldBeEqual(t, 2, len(keys))
	// make sure device1 key is deleted
	keys, err = c.QueryDeviceKey(testApp, device1)
	noError(t, err)
	shouldBeEqual(t, 0, len(keys))
}
func TestQueryAppidDeviceByAccessKey(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	// delete if exist
	device := "sdk-testQueryAppidDeviceNameByAccessKey"
	defer c.DeleteDevice(testApp, device)
	dev := &Device{
		Device:            device,
		SegmentExpireDays: 7,
	}
	_, err := c.AddDevice(testApp, dev)
	noError(t, err)

	keys, err := c.QueryDeviceKey(testApp, device)
	noError(t, err)
	shouldBeEqual(t, 1, len(keys))
	dak := keys[0].AccessKey
	realAppid, realDevice, err := c.QueryAppidDeviceNameByAccessKey(dak)
	noError(t, err)
	shouldBeEqual(t, device, realDevice)
	shouldBeEqual(t, testApp, realAppid)
}

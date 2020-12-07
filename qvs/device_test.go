package qvs

import (
	"fmt"
	"testing"
)

func TestDeviceCRUD(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	ns := &NameSpace{
		Name:       "testNamespace",
		AccessType: "gb28181",
		UrlMode:    1,
		Domains:    []string{"qiniu1.com"},
	}
	ns1, err := c.AddNamespace(ns)
	noError(t, err)

	device := &Device{
		NamespaceId: ns1.ID,
		Type:        2,
		Username:    "username",
		Password:    "password",
	}
	postDevice, err := c.AddDevice(device)
	noError(t, err)

	queryDevice, err := c.QueryDevice(ns1.ID, postDevice.GBId)
	noError(t, err)
	shouldBeEqual(t, postDevice.GBId, queryDevice.GBId)
	shouldBeEqual(t, "username", queryDevice.Username)
	shouldBeEqual(t, "password", queryDevice.Password)

	ops := []PatchOperation{
		{
			Op:    "replace",
			Key:   "desc",
			Value: "test",
		},
	}
	updateDevice, err := c.UpdateDevice(ns1.ID, queryDevice.GBId, ops)
	noError(t, err)
	shouldBeEqual(t, updateDevice.GBId, queryDevice.GBId)
	shouldBeEqual(t, updateDevice.Desc, "test")

	deviceChannels, err := c.ListChannels(ns1.ID, postDevice.GBId, "")
	noError(t, err)
	fmt.Println(*deviceChannels)

	device2 := &Device{
		NamespaceId: ns1.ID,
		Username:    "username",
		Password:    "password",
	}
	postDevice2, err := c.AddDevice(device2)
	noError(t, err)

	device3 := &Device{
		NamespaceId: ns1.ID,
		Username:    "username",
		Password:    "password",
	}
	postDevice3, err := c.AddDevice(device3)
	noError(t, err)

	devices, total, err := c.ListDevice(ns1.ID, 0, 2, "", "", 0)
	noError(t, err)
	shouldBeEqual(t, int64(3), total)
	shouldBeEqual(t, 2, len(devices))

	devices, total, err = c.ListDevice(ns1.ID, 2, 2, "", "", 0)
	noError(t, err)
	shouldBeEqual(t, int64(3), total)
	shouldBeEqual(t, 1, len(devices))

	err = c.StartDevice(ns1.ID, postDevice.GBId, nil)
	fmt.Println(err)

	err = c.StopDevice(ns1.ID, postDevice.GBId, nil)
	fmt.Println(err)

	c.DeleteDevice(ns1.ID, postDevice.GBId)
	c.DeleteDevice(ns1.ID, postDevice2.GBId)
	c.DeleteDevice(ns1.ID, postDevice3.GBId)

	c.DeleteNamespace(ns1.ID)

	channel, err := c.QueryChannel("3nm4x0vyz7xlu", "31011500991180000270", "31011500991180000135")
	noError(t, err)
	fmt.Println(*channel)

	//err = c.DeleteChannel("3nm4x0vyz7xlu", "31011500991180000270", "31011500991180000135")
	//noError(t, err)

	err = c.FetchCatalog("3nm4x0vyz7xlu", "31011500991180000270")
	noError(t, err)
}

func TestQueryGBRecordHistories(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	ret, err := c.QueryGBRecordHistories("3nm4x0vyz7xlu", "31011500991180000270", "34020000001310000020", 1604817540, 1604903940)
	noError(t, err)
	for _, v := range ret.Items {
		fmt.Println(v)
	}
}

func TestDevicePTZ(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	err := c.DevicePtz("3nm4x0vyz7xlu", "31011500991180000270", "34020000001310000034", "left", 10)
	noError(t, err)

	err = c.DeviceFocus("3nm4x0vyz7xlu", "31011500991180000270", "34020000001310000034", "focusnear", 10)
	noError(t, err)

	err = c.DeviceIris("3nm4x0vyz7xlu", "31011500991180000270", "34020000001310000034", "irisin", 10)
	noError(t, err)

	items, err := c.QueryDevicePresets("3nm4x0vyz7xlu", "31011500991180000270", "34020000001310000034", true)
	noError(t, err)
	fmt.Println("items", items)
}

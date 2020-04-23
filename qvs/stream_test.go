package qvs

import (
	"fmt"
	"testing"
)

func TestStreamCRUD(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	ns := &NameSpace{
		Name:        "testNamespace",
		AccessType:  "rtmp",
		RTMPURLType: 1,
		Domains:     []string{"qtest.com"},
	}
	ns1, err := c.AddNamespace(ns)
	noError(t, err)

	stream := &Stream{
		StreamID: "test001",
	}
	_, err = c.AddStream(ns1.ID, stream)
	noError(t, err)

	stream2, err := c.QueryStream(ns1.ID, "test001")
	noError(t, err)
	shouldBeEqual(t, stream.StreamID, stream2.StreamID)

	ops := []PatchOperation{
		{
			Op:    "replace",
			Key:   "desc",
			Value: "test",
		},
	}
	stream3, err := c.UpdateStream(ns1.ID, stream.StreamID, ops)
	noError(t, err)
	shouldBeEqual(t, stream3.StreamID, stream2.StreamID)
	shouldBeEqual(t, stream3.Desc, "test")

	stream4 := &Stream{
		StreamID: "test002",
	}
	_, err = c.AddStream(ns1.ID, stream4)
	noError(t, err)

	stream5 := &Stream{
		StreamID: "test003",
	}
	_, err = c.AddStream(ns1.ID, stream5)
	noError(t, err)

	streams, total, err := c.ListStream(ns1.ID, 0, 2, "", "", 0)
	noError(t, err)
	shouldBeEqual(t, int64(3), total)
	shouldBeEqual(t, 2, len(streams))

	streams, total, err = c.ListStream(ns1.ID, 2, 2, "", "", 0)
	noError(t, err)
	shouldBeEqual(t, int64(3), total)
	shouldBeEqual(t, 1, len(streams))

	err = c.DisableStream(ns1.ID, stream.StreamID)
	noError(t, err)
	ret, err := c.QueryStream(ns1.ID, stream.StreamID)
	noError(t, err)
	shouldBeEqual(t, true, ret.Disabled)

	err = c.EnableStream(ns1.ID, stream.StreamID)
	noError(t, err)
	ret, err = c.QueryStream(ns1.ID, stream.StreamID)
	noError(t, err)
	shouldBeEqual(t, false, ret.Disabled)

	c.DeleteStream(ns1.ID, "test001")
	c.DeleteStream(ns1.ID, "test002")
	c.DeleteStream(ns1.ID, "test003")

	c.DeleteNamespace(ns1.ID)
}

func TestDynamicPublishPlayURL(t *testing.T) {

	c := getTestManager()
	ret, err := c.DynamicPublishPlayURL("2akrarsdkltth", "device005", &DynamicLiveRoute{PublishIP: "127.0.0.1", PlayIP: "127.0.0.1", UrlExpireSec: 0})
	fmt.Println(err, ret)
}

func TestStaticPublishPlayURL(t *testing.T) {

	c := getTestManager()
	ret, err := c.StaticPublishPlayURL("test.com", "2akrarsdkltth", "device005", 0, 0)
	fmt.Println(err, ret)
}

func TestStreamsSnapshots(t *testing.T) {
	c := getTestManager()
	ret, err := c.StreamsSnapshots("2akrarsdkltth", "device005", 1585565152, 1585568752, 0, 0, "")
	fmt.Println(err, string(ret))
}

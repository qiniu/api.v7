package qvs

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestQueryStreamRecordHistories(t *testing.T) {

	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	items, _, err := c.QueryStreamRecordHistories("2xenzw5o81ods", "31011500991320000356", 1604851200, 1604894400, "", 0, "flv")
	noError(t, err)
	for _, v := range items {
		fmt.Println(v)
	}
}

func TestOndemandRecordAndSnap(t *testing.T) {

	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	err := c.OndemandSnap("2xenzw5o81ods", "31011500991320000356")
	noError(t, err)
	err = c.StartRecord("2xenzw5o81ods", "31011500991320000356")
	noError(t, err)
	time.Sleep(15 * 60 * time.Second)
	err = c.StopRecord("2xenzw5o81ods", "31011500991320000356")
	noError(t, err)
}

func TestDeleteSnapshots(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	var res = struct {
		Items []struct {
			Snap string `json:"snap"`
			Time int    `json:"time"`
		} `json:"items"`
		Marker string `json:"marker"`
	}{}
	ret, err := c.StreamsSnapshots("2xenzw5o81ods", "31011500991320000356", 1604988765, 1605002214, 1, 20, "")
	noError(t, err)
	fmt.Println(string(ret))
	err = json.Unmarshal(ret, &res)
	noError(t, err)
	fmt.Println(len(res.Items))
	if len(res.Items) > 0 {
		err = c.DeleteSnapshots("2xenzw5o81ods", "31011500991320000356", []string{res.Items[0].Snap[strings.Index(res.Items[0].Snap, "snapshot"):strings.Index(res.Items[0].Snap, "?")]})
		noError(t, err)
	}
	ret, err = c.StreamsSnapshots("2xenzw5o81ods", "31011500991320000356", 1604988765, 1605002214, 1, 20, "")
	noError(t, err)
	err = json.Unmarshal(ret, &res)
	noError(t, err)
	fmt.Println(len(res.Items))
}

func TestRecordClipsSaveasAndDeleteRecord(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	ret, err := c.RecordClipsSaveas("2xenzw5o81ods", "31011500991320000356", &saveasArgs{
		Format: "m3u8",
		Start: 1604989846,
		End: 1604990735,
	})
	noError(t, err)
	shouldBeEqual(t, ret.Fname, "record/2xenzw5o81ods/31011500991320000356/1604989846152-1604990735281-852640.m3u8")

	err = c.DeleteStreamRecordHistories("2xenzw5o81ods", "31011500991320000356", []string{"record/2xenzw5o81ods/31011500991320000356/1604989846152-1604990735281-852640.m3u8"})
	noError(t, err)
}



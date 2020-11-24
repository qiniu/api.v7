package qvs

import (
	"testing"
	"time"
)

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

func TestRecordClipsSaveasAndDeleteRecord(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()

	ret, err := c.RecordClipsSaveas("2xenzw5o81ods", "31011500991320000356", &saveasArgs{
		Format: "m3u8",
		Start:  1604989846,
		End:    1604990735,
	})
	noError(t, err)
	shouldBeEqual(t, ret.Fname, "record/2xenzw5o81ods/31011500991320000356/1604989846152-1604990735281-852640.m3u8")

	err = c.DeleteStreamRecordHistories("2xenzw5o81ods", "31011500991320000356", []string{"record/2xenzw5o81ods/31011500991320000356/1604989846152-1604990735281-852640.m3u8"})
	noError(t, err)
}

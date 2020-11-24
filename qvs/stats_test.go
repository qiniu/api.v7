package qvs

import (
	"fmt"
	"testing"
)

func TestQueryStats(t *testing.T) {
	if skipTest() {
		t.SkipNow()
	}
	c := getTestManager()
	ret, err := c.QueryFlow("", "", "5min", 20200901, 20200902)
	noError(t, err)
	fmt.Println(*ret)

	ret, err = c.QueryBandwidth("", "", "hour", 20200901, 20200902)
	noError(t, err)
	fmt.Println(*ret)
}

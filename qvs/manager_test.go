package qvs

import (
	"runtime/debug"
	"testing"

	"github.com/qiniu/api.v7/v7/auth"
)

var (
	testAccessKey = "2PSNcce7OqV05g9sF2Ngvsp-h4pDDqMhzBQRgXid"
	testSecretKey = "-iPXzzP0oAXYvGMseW71FfT4N5-rK8f5mR73LNZQ"
)

func skipTest() bool {
	return testAccessKey == "" || testSecretKey == ""
}

func getTestManager() *Manager {
	mac := auth.New(testAccessKey, testSecretKey)
	return NewManager(mac, nil)
}
func noError(t *testing.T, err error) {
	if err != nil {
		debug.PrintStack()
		t.Fatalf("should be nil, err = %s", err.Error())
	}
}

func shouldBeEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		debug.PrintStack()
		t.Fatalf("should be equal, expect = %#v, but got  = %#v", a, b)
	}
}

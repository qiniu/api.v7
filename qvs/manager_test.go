package qvs

import (
	"runtime/debug"
	"testing"

	"github.com/qiniu/api.v7/v7/auth"
)

var (
	testAccessKey = ""
	testSecretKey = ""
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
		t.Helper()
		debug.PrintStack()
		t.Fatalf("should be nil, err = %s", err.Error())
	}
}

func shouldBeEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Helper()
		debug.PrintStack()
		t.Fatalf("should be equal, expect = %#v, but got  = %#v", a, b)
	}
}

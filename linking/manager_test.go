package linking

import (
	"runtime/debug"
	"testing"

	"github.com/qiniu/api.v7/v7/auth"
)

var (
	testAccessKey = ""
	testSecretKey = ""
	testApp       = "2xenzvm26zx2b"
)

func skipTest() bool {
	return testAccessKey == "" || testSecretKey == "" || testApp == ""
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

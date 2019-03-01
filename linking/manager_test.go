package linking

import (
	"runtime/debug"
	"testing"

	"github.com/qiniu/api.v7/auth"
)

var (
	testAccessKey = ""
	testSecretKey = ""
	testApp       = ""
)

func skipTest() bool {
	return testAccessKey == "" || testSecretKey == "" || testApp == ""
}

func getTestManager() *Manager {
	mac := auth.New(testAccessKey, testSecretKey)
	return New(mac, nil)
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

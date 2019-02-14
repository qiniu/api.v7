package storage

import (
	"fmt"
	"github.com/qiniu/api.v7/conf"
	"runtime"
	"testing"
)

func TestVariable(t *testing.T) {
	appName := "test"

	SetAppName(appName)

	want := fmt.Sprintf("QiniuGo/%s (%s; %s; %s) %s", conf.Version, runtime.GOOS, runtime.GOARCH, appName, runtime.Version())

	if UserAgent != want {
		t.Fail()
	}
}

package storage

import (
	"testing"
)

func TestReqHost(t *testing.T) {
	zoneHuadong := Zone{
		SrcUpHosts: []string{
			"up.qiniup.com",
			"up-nb.qiniup.com",
			"up-xs.qiniup.com",
		},
		CdnUpHosts: []string{
			"upload.qiniup.com",
			"upload-nb.qiniup.com",
			"upload-xs.qiniup.com",
		},
		RsHost:    "rs.qbox.me",
		RsfHost:   "rsf.qbox.me",
		ApiHost:   "api.qiniu.com",
		IovipHost: "iovip.qbox.me",
	}
	cfgs := []Config{
		{UseHTTPS: true, Zone: &zoneHuadong},
		{UseHTTPS: true, RsHost: "http://rshost.com"},
		{UseHTTPS: true, RsHost: "https://rshost.com"},
		{UseHTTPS: true, Zone: &zoneHuadong, RsHost: "http://rshost.com"},
		{UseHTTPS: false, Zone: &zoneHuadong},
		{UseHTTPS: false, RsHost: "http://rshost.com"},
		{UseHTTPS: false, RsHost: "https://rshost.com"},
		{UseHTTPS: false, Zone: &zoneHuadong, RsHost: "http://rshost.com"},
		{UseHTTPS: false, Region: &zoneHuadong, RsHost: "http://rshost.com"},
	}
	wantRsHosts := []string{
		"https://rs.qbox.me",
		"https://rshost.com",
		"https://rshost.com",
		"https://rs.qbox.me",
		"http://rs.qbox.me",
		"http://rshost.com",
		"http://rshost.com",
		"http://rs.qbox.me",
		"http://rs.qbox.me",
	}

	for ind, cfg := range cfgs {
		got := cfg.RsReqHost()
		want := wantRsHosts[ind]
		if got != want {
			t.Errorf("ind = %d, want = %q, got = %q\n", ind, want, got)
		}
	}
}

package cdn

import (
	"os"
	"testing"
	"time"

	"fmt"
	"github.com/qiniu/api.v7/v7/auth"
)

//global variables

var (
	ak     = os.Getenv("accessKey")
	sk     = os.Getenv("secretKey")
	domain = os.Getenv("QINIU_TEST_DOMAIN")

	layout    = "2006-01-02"
	now       = time.Now()
	startDate = now.AddDate(0, 0, -2).Format(layout)
	endDate   = now.AddDate(0, 0, -1).Format(layout)
	logDate   = now.AddDate(0, 0, -1).Format(layout)

	testUrls = []string{
		fmt.Sprintf("http://%s/qiniu.png", domain),
		fmt.Sprintf("http://%s/qiniu-fetch.png", domain),
	}
	testDirs = []string{
		fmt.Sprintf("http://%s/gosdkintegration/", domain),
		fmt.Sprintf("http://%s/gosdkintegration1/", domain),
	}
)

var mac *auth.Credentials
var cdnManager *CdnManager

func init() {
	if ak == "" || sk == "" {
		panic("ak/sk should not be empty")
	}
	mac = auth.New(ak, sk)
	cdnManager = NewCdnManager(mac)
}

//TestGetBandwidthData
func TestGetBandwidthData(t *testing.T) {
	type args struct {
		startDate   string
		endDate     string
		granularity string
		domainList  []string
	}

	testCases := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "CdnManager_TestGetBandwidthData",
			args: args{
				startDate,
				endDate,
				"5min",
				[]string{domain},
			},
			wantCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret, err := cdnManager.GetBandwidthData(tc.args.startDate, tc.args.endDate,
				tc.args.granularity, tc.args.domainList)
			if err != nil || ret.Code != tc.wantCode {
				t.Errorf("GetBandwidth() error = %v, %v", err, ret.Error)
				return
			}
		})
	}
}

//TestGetFluxData
func TestGetFluxData(t *testing.T) {
	type args struct {
		startDate   string
		endDate     string
		granularity string
		domainList  []string
	}

	testCases := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "CdnManager_TestGetFluxData",
			args: args{
				startDate,
				endDate,
				"5min",
				[]string{domain},
			},
			wantCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret, err := cdnManager.GetFluxData(tc.args.startDate, tc.args.endDate,
				tc.args.granularity, tc.args.domainList)
			if err != nil || ret.Code != tc.wantCode {
				t.Errorf("GetFlux() error = %v, %v", err, ret.Error)
				return
			}
		})
	}
}

//TestRefreshUrls
func TestRefreshUrls(t *testing.T) {
	type args struct {
		urls []string
	}

	testCases := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "CdnManager_TestRefresUrls",
			args: args{
				urls: testUrls,
			},
			wantCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret, err := cdnManager.RefreshUrls(tc.args.urls)
			if err != nil || ret.Code != tc.wantCode {
				t.Errorf("RefreshUrls() %v error = %v, %v", tc.args.urls, err, ret.Error)
				return
			}
		})
	}
}

//TestRefreshDirs
func TestRefreshDirs(t *testing.T) {
	type args struct {
		dirs []string
	}

	testCases := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "CdnManager_TestRefreshDirs",
			args: args{
				dirs: testDirs,
			},
			wantCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret, err := cdnManager.RefreshDirs(tc.args.dirs)
			if err != nil || ret.Code != tc.wantCode {
				if ret.Error == "refresh dir limit error" {
					t.Logf("RefreshDirs() error=%v", ret.Error)
				} else {
					t.Errorf("RefreshDirs() error = %v, %v", err, ret.Error)
				}
				return
			}
		})
	}
}

/* 预取有额度限制
//TestPrefetchUrls
func TestPrefetchUrls(t *testing.T) {
	type args struct {
		urls []string
	}

	testCases := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "CdnManager_PrefetchUrls",
			args: args{
				urls: testUrls,
			},
			wantCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret, err := cdnManager.PrefetchUrls(tc.args.urls)
			if err != nil || ret.Code != tc.wantCode {
				t.Errorf("PrefetchUrls() error = %v, %v", err, ret.Error)
				return
			}
		})
	}
}
*/

//TestGetCdnLogList
func TestGetCdnLogList(t *testing.T) {
	type args struct {
		date    string
		domains []string
	}

	testCases := []struct {
		name string
		args args
	}{
		{
			name: "CdnManager_TestGetCdnLogList",
			args: args{
				date:    logDate,
				domains: []string{domain},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := cdnManager.GetCdnLogList(tc.args.date, tc.args.domains)
			if err != nil {
				t.Errorf("GetCdnLogList() error = %v", err)
				return
			}
		})
	}
}

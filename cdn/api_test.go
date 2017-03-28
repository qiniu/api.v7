package cdn

import (
	"os"
	"reflect"
	"testing"

	"qiniupkg.com/api.v7/kodo"
)

var (
	ak = os.Getenv("QINIU_ACCESS_KEY")
	sk = os.Getenv("QINIU_SECRET_KEY")
)

func TestGetBandWidthData(t *testing.T) {
	type args struct {
		startDate   string
		endDate     string
		granularity string
		domainList  []string
	}
	tests := []struct {
		name        string
		args        args
		wantTraffic TrafficResp
		wantErr     bool
	}{
		{
			name: "BandWidthTest_1",
			args: args{
				"2016-12-20",
				"2016-12-20",
				"5min",
				[]string{"abc.def.com"},
			},
		},
	}
	kodo.SetMac(ak, sk)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBandWidthData(tt.args.startDate, tt.args.endDate, tt.args.granularity, tt.args.domainList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBandWidthData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestGetFluxData(t *testing.T) {
	type args struct {
		startDate   string
		endDate     string
		granularity string
		domainList  []string
	}
	tests := []struct {
		name        string
		args        args
		wantTraffic TrafficResp
		wantErr     bool
	}{
		{
			name: "BandWidthTest_1",
			args: args{
				"2016-12-20",
				"2016-12-20",
				"5min",
				[]string{"abc.def.com"},
			},
		},
	}
	kodo.SetMac(ak, sk)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetFluxData(tt.args.startDate, tt.args.endDate, tt.args.granularity, tt.args.domainList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFluxData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRefreshUrlsAndDirs(t *testing.T) {
	type args struct {
		urls []string
		dirs []string
	}
	tests := []struct {
		name       string
		args       args
		wantResult RefreshResp
		wantErr    bool
	}{
		{
			name: "refresh_test_1",
			args: args{
				urls: []string{""},
			},
			wantErr: false,
		},
	}
	kodo.SetMac(ak, sk)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RefreshUrlsAndDirs(tt.args.urls, tt.args.dirs)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshUrlsAndDirs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestRefreshUrls(t *testing.T) {
	type args struct {
		urls []string
	}
	tests := []struct {
		name       string
		args       args
		wantResult RefreshResp
		wantErr    bool
	}{
		{
			name: "refresh_test_1",
			args: args{
				urls: []string{""},
			},
			wantErr: false,
		},
	}
	kodo.SetMac(ak, sk)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RefreshUrls(tt.args.urls)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshUrls() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRefreshDirs(t *testing.T) {
	type args struct {
		dirs []string
	}
	tests := []struct {
		name       string
		args       args
		wantResult RefreshResp
		wantErr    bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := RefreshDirs(tt.args.dirs)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshDirs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("RefreshDirs() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestPrefetchUrls(t *testing.T) {
	type args struct {
		urls []string
	}
	tests := []struct {
		name       string
		args       args
		wantResult PrefetchResp
		wantErr    bool
	}{
		{
			name: "refresh_test_1",
			args: args{
				urls: []string{""},
			},
			wantErr: false,
		},
	}
	kodo.SetMac(ak, sk)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := PrefetchUrls(tt.args.urls)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrefetchUrls() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

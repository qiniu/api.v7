package cdn

import (
	"reflect"
	"testing"
)

func TestGetCdnLogList(t *testing.T) {
	type args struct {
		date    string
		domains string
	}
	tests := []struct {
		name           string
		args           args
		wantDomainLogs []LogDomainInfo
		wantErr        bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDomainLogs, err := GetCdnLogList(tt.args.date, tt.args.domains)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCdnLogList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDomainLogs, tt.wantDomainLogs) {
				t.Errorf("GetCdnLogList() = %v, want %v", gotDomainLogs, tt.wantDomainLogs)
			}
		})
	}
}

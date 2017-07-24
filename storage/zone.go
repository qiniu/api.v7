package storage

import (
	"context"
	"fmt"
	"github.com/qiniu/x/rpc.v7"
	"sync"
)

type Zone struct {
	SrcUpHosts []string
	CdnUpHosts []string
	RsHost     string
	RsfHost    string
	ApiHost    string
	IovipHost  string
}

//z0
var Zone_z0 = Zone{
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
	RsHost:    "rs.qiniu.com",
	RsfHost:   "rsf.qiniu.com",
	ApiHost:   "api.qiniu.com",
	IovipHost: "iovip.qbox.me",
}

//z1
var Zone_z1 = Zone{
	SrcUpHosts: []string{
		"up-z1.qiniup.com",
	},
	CdnUpHosts: []string{
		"upload-z1.qiniup.com",
	},
	RsHost:    "rs-z1.qiniu.com",
	RsfHost:   "rsf-z1.qiniu.com",
	ApiHost:   "api-z1.qiniu.com",
	IovipHost: "iovip-z1.qbox.me",
}

//z2
var Zone_z2 = Zone{
	SrcUpHosts: []string{
		"up-z2.qiniup.com",
		"up-gz.qiniup.com",
		"up-fs.qiniup.com",
	},
	CdnUpHosts: []string{
		"upload-z2.qiniup.com",
		"upload-gz.qiniup.com",
		"upload-fs.qiniup.com",
	},
	RsHost:    "rs-z2.qiniu.com",
	RsfHost:   "rsf-z2.qiniu.com",
	IovipHost: "iovip-z2.qbox.me",
}

//na0
var Zone_na0 = Zone{
	SrcUpHosts: []string{
		"up-na0.qiniu.com",
	},
	CdnUpHosts: []string{
		"upload-na0.qiniu.com",
	},
	RsHost:    "rs-na0.qiniu.com",
	RsfHost:   "rsf-na0.qiniu.com",
	IovipHost: "iovip-na0.qbox.me",
}

//////////////////////////////

const UC_HOST = "https://uc.qbox.me"

type UcQueryRet struct {
	Ttl int                            `json:"ttl"`
	Io  map[string]map[string][]string `json:"io"`
	Up  map[string]map[string][]string `json:"up"`
}

var (
	zoneMutext sync.RWMutex
	zoneCache  = make([string]Zone)
)

//v2 version
func QueryZone(ak, bucket string) (zone Zone, err error) {
	//check from cache

	//query from server
	reqUrl := fmt.Sprintf("%s/v2/query?ak=%s&bucket=%s", UC_HOST, ak, bucket)
	var ret UcQueryRet
	ctx := context.Background()
	qErr := rpc.CallWithForm(ctx, &ret, "GET", reqUrl, nil)
	if qErr != nil {
		err = fmt.Errorf("query zone error, %s", qErr.Error())
		return
	}

}

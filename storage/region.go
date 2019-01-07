package storage

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// 存储所在的地区，例如华东，华南，华北
// 每个存储区域可能有多个机房信息，每个机房可能有多个上传入口
type Region struct {
	// 上传入口
	SrcUpHosts []string

	// 加速上传入口
	CdnUpHosts []string

	// 获取文件信息入口
	RsHost string

	// bucket列举入口
	RsfHost string

	ApiHost string

	// 存储io 入口
	IovipHost string
}

func (z *Region) String() string {
	str := ""
	str += fmt.Sprintf("SrcUpHosts: %v\n", z.SrcUpHosts)
	str += fmt.Sprintf("CdnUpHosts: %v\n", z.CdnUpHosts)
	str += fmt.Sprintf("IovipHost: %s\n", z.IovipHost)
	str += fmt.Sprintf("RsHost: %s\n", z.RsHost)
	str += fmt.Sprintf("RsfHost: %s\n", z.RsfHost)
	str += fmt.Sprintf("ApiHost: %s\n", z.ApiHost)
	return str
}

func endpoint(useHttps bool, host string) string {
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		return host
	}
	scheme := "http://"
	if useHttps {
		scheme = "https://"
	}
	return fmt.Sprintf("%s%s", scheme, host)
}

// 获取rsfHost
func (z *Region) GetRsfHost(useHttps bool) string {
	return endpoint(useHttps, z.RsfHost)
}

// 获取io host
func (z *Region) GetIoHost(useHttps bool) string {
	return endpoint(useHttps, z.IovipHost)
}

// 获取RsHost
func (z *Region) GetRsHost(useHttps bool) string {
	return endpoint(useHttps, z.RsHost)
}

// 获取api host
func (z *Region) GetApiHost(useHttps bool) string {
	return endpoint(useHttps, z.ApiHost)
}

// RegionHuadong 表示华东机房
var RegionHuadong = Region{
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

// RegionHuabei 表示华北机房
var RegionHuabei = Region{
	SrcUpHosts: []string{
		"up-z1.qiniup.com",
	},
	CdnUpHosts: []string{
		"upload-z1.qiniup.com",
	},
	RsHost:    "rs-z1.qbox.me",
	RsfHost:   "rsf-z1.qbox.me",
	ApiHost:   "api-z1.qiniu.com",
	IovipHost: "iovip-z1.qbox.me",
}

// RegionHuanan 表示华南机房
var RegionHuanan = Region{
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
	RsHost:    "rs-z2.qbox.me",
	RsfHost:   "rsf-z2.qbox.me",
	ApiHost:   "api-z2.qiniu.com",
	IovipHost: "iovip-z2.qbox.me",
}

// RegionBeimei 表示北美机房
var RegionBeimei = Region{
	SrcUpHosts: []string{
		"up-na0.qiniup.com",
	},
	CdnUpHosts: []string{
		"upload-na0.qiniup.com",
	},
	RsHost:    "rs-na0.qbox.me",
	RsfHost:   "rsf-na0.qbox.me",
	ApiHost:   "api-na0.qiniu.com",
	IovipHost: "iovip-na0.qbox.me",
}

// RegionXinjiapo 表示新加坡机房
var RegionXinjiapo = Region{
	SrcUpHosts: []string{
		"up-as0.qiniup.com",
	},
	CdnUpHosts: []string{
		"upload-as0.qiniup.com",
	},
	RsHost:    "rs-as0.qbox.me",
	RsfHost:   "rsf-as0.qbox.me",
	ApiHost:   "api-as0.qiniu.com",
	IovipHost: "iovip-as0.qbox.me",
}

// for programmers
var Region_z0 = RegionHuadong
var Region_z1 = RegionHuabei
var Region_z2 = RegionHuanan
var Region_na0 = RegionBeimei
var Region_as0 = RegionXinjiapo

// UcHost 为查询空间相关域名的API服务地址
const UcHost = "https://uc.qbox.me"

// UcQueryRet 为查询请求的回复
type UcQueryRet struct {
	TTL int                            `json:"ttl"`
	Io  map[string]map[string][]string `json:"io"`
	Up  map[string]UcQueryUp           `json:"up"`
}

// UcQueryUp 为查询请求回复中的上传域名信息
type UcQueryUp struct {
	Main   []string `json:"main,omitempty"`
	Backup []string `json:"backup,omitempty"`
	Info   string   `json:"info,omitempty"`
}

var (
	regionMutext sync.RWMutex
	regionCache  = make(map[string]*Region)
)

// GetRegion 用来根据ak和bucket来获取空间相关的机房信息
func GetRegion(ak, bucket string) (region *Region, err error) {
	regionID := fmt.Sprintf("%s:%s", ak, bucket)
	//check from cache
	regionMutext.RLock()
	if v, ok := regionCache[regionID]; ok {
		region = v
	}
	regionMutext.RUnlock()
	if region != nil {
		return
	}

	//query from server
	reqURL := fmt.Sprintf("%s/v2/query?ak=%s&bucket=%s", UcHost, ak, bucket)
	var ret UcQueryRet
	ctx := context.TODO()
	qErr := DefaultClient.CallWithForm(ctx, &ret, "GET", reqURL, nil, nil)
	if qErr != nil {
		err = fmt.Errorf("query region error, %s", qErr.Error())
		return
	}

	ioHost := ret.Io["src"]["main"][0]
	srcUpHosts := ret.Up["src"].Main
	if ret.Up["src"].Backup != nil {
		srcUpHosts = append(srcUpHosts, ret.Up["src"].Backup...)
	}
	cdnUpHosts := ret.Up["acc"].Main
	if ret.Up["acc"].Backup != nil {
		cdnUpHosts = append(cdnUpHosts, ret.Up["acc"].Backup...)
	}

	region = &Region{
		SrcUpHosts: srcUpHosts,
		CdnUpHosts: cdnUpHosts,
		IovipHost:  ioHost,
		RsHost:     DefaultRsHost,
		RsfHost:    DefaultRsfHost,
		ApiHost:    DefaultAPIHost,
	}

	//set specific hosts if possible
	setSpecificHosts(ioHost, region)

	regionMutext.Lock()
	regionCache[regionID] = region
	regionMutext.Unlock()
	return
}

func setSpecificHosts(ioHost string, region *Region) {
	if strings.Contains(ioHost, "-z1") {
		region.RsHost = "rs-z1.qbox.me"
		region.RsfHost = "rsf-z1.qbox.me"
		region.ApiHost = "api-z1.qiniu.com"
	} else if strings.Contains(ioHost, "-z2") {
		region.RsHost = "rs-z2.qbox.me"
		region.RsfHost = "rsf-z2.qbox.me"
		region.ApiHost = "api-z2.qiniu.com"
	} else if strings.Contains(ioHost, "-na0") {
		region.RsHost = "rs-na0.qbox.me"
		region.RsfHost = "rsf-na0.qbox.me"
		region.ApiHost = "api-na0.qiniu.com"
	} else if strings.Contains(ioHost, "-as0") {
		region.RsHost = "rs-as0.qbox.me"
		region.RsfHost = "rsf-as0.qbox.me"
		region.ApiHost = "api-as0.qiniu.com"
	}
}

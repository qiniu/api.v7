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

type RegionCode string

// GetDefaultReion 根据RegionID获取对应的Region信息
func GetDefaultRegion(regionCode RegionCode) (Region, bool) {
	if r, ok := regionMap[regionCode]; ok {
		return r, ok
	}
	return Region{}, false
}

// GetDefaultRegionStr 根据region code string返回Region信息
func GetDefaultRegionStr(regionCodeStr string) (Region, bool) {
	switch regionCodeStr {
	case "z0":
		return GetDefaultRegion(RCodeHuadong)
	case "z1":
		return GetDefaultRegion(RCodeHuabei)
	case "z2":
		return GetDefaultRegion(RCodeHuanan)
	case "na0":
		return GetDefaultRegion(RCodeBeimei)
	case "as0":
		return GetDefaultRegion(RCodeAsia)
	}
	return Region{}, false
}

func (r *Region) String() string {
	str := ""
	str += fmt.Sprintf("SrcUpHosts: %v\n", r.SrcUpHosts)
	str += fmt.Sprintf("CdnUpHosts: %v\n", r.CdnUpHosts)
	str += fmt.Sprintf("IovipHost: %s\n", r.IovipHost)
	str += fmt.Sprintf("RsHost: %s\n", r.RsHost)
	str += fmt.Sprintf("RsfHost: %s\n", r.RsfHost)
	str += fmt.Sprintf("ApiHost: %s\n", r.ApiHost)
	return str
}

func endpoint(useHttps bool, host string) string {
	host = strings.TrimSpace(host)
	host = strings.TrimLeft(host, "http://")
	host = strings.TrimLeft(host, "https://")
	if host == "" {
		return ""
	}
	scheme := "http://"
	if useHttps {
		scheme = "https://"
	}
	return fmt.Sprintf("%s%s", scheme, host)
}

// 获取rsfHost
func (r *Region) GetRsfHost(useHttps bool) string {
	return endpoint(useHttps, r.RsfHost)
}

// 获取io host
func (r *Region) GetIoHost(useHttps bool) string {
	return endpoint(useHttps, r.IovipHost)
}

// 获取RsHost
func (r *Region) GetRsHost(useHttps bool) string {
	return endpoint(useHttps, r.RsHost)
}

// 获取api host
func (r *Region) GetApiHost(useHttps bool) string {
	return endpoint(useHttps, r.ApiHost)
}

var (
	// regionHuadong 表示华东机房
	regionHuadong = Region{
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

	// regionHuabei 表示华北机房
	regionHuabei = Region{
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
	// regionHuanan 表示华南机房
	regionHuanan = Region{
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

	// regionBeimei 表示北美机房
	regionBeimei = Region{
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
	// regionXinjiapo 表示新加坡机房
	regionXinjiapo = Region{
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
)

const (
	// region code
	RCodeHuadong = RegionCode("z0")
	RCodeHuabei  = RegionCode("z1")
	RCodeHuanan  = RegionCode("z2")
	RCodeBeimei  = RegionCode("na0")
	RCodeAsia    = RegionCode("as0")
)

// regionMap 是RegionID到具体的Region的映射
var regionMap = map[RegionCode]Region{
	RCodeHuadong: regionHuadong,
	RCodeHuanan:  regionHuanan,
	RCodeHuabei:  regionHuabei,
	RCodeAsia:    regionXinjiapo,
	RCodeBeimei:  regionBeimei,
}

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

func regionFromHost(ioHost string) (Region, bool) {
	if strings.Contains(ioHost, "-z1") {
		return GetDefaultRegion(RCodeHuabei)
	}
	if strings.Contains(ioHost, "-z2") {
		return GetDefaultRegion(RCodeHuanan)
	}

	if strings.Contains(ioHost, "-na0") {
		return GetDefaultRegion(RCodeBeimei)
	}
	if strings.Contains(ioHost, "-as0") {
		return GetDefaultRegion(RCodeAsia)
	}
	return Region{}, false
}

func setSpecificHosts(ioHost string, region *Region) {
	r, ok := regionFromHost(ioHost)
	if ok {
		region.RsHost = r.RsHost
		region.RsfHost = r.RsfHost
		region.ApiHost = r.ApiHost
	}
}

package storage

import (
	"strings"
)

// Config 为文件上传，资源管理等配置
type Config struct {
	Zone          *Region //空间所在的机房
	UseHTTPS      bool    //是否使用https域名
	UseCdnDomains bool    //是否使用cdn加速域名
	CentralRsHost string  //中心机房的RsHost，用于list bucket

	// 兼容保留
	RsHost  string
	RsfHost string
	UpHost  string
	ApiHost string
	IoHost  string
}

func emptyHost(useHttps bool) string {
	return ""
}

func reqHost(useHttps bool, fn func(bool) string, host, defaultHost string) (endp string) {
	endp = fn(useHttps)
	endp = strings.TrimSpace(endp)
	if endp != "" && endp != "http://" && endp != "https://" {
		return
	}
	if host == "" {
		host = defaultHost
	}
	return endpoint(useHttps, host)
}

// 获取RsHost
// 优先使用Zone中的Host信息，如果Zone中的host信息没有配置，那么使用Config中的Host信息
func (c *Config) RsReqHost() string {
	var method func(bool) string
	if c.Zone != nil {
		method = c.Zone.GetRsHost
	} else {
		method = emptyHost
	}
	return reqHost(c.UseHTTPS, method, c.RsHost, DefaultRsHost)
}

// 获取rsfHost
// 优先使用Zone中的Host信息，如果Zone中的host信息没有配置，那么使用Config中的Host信息
func (c *Config) RsfReqHost() string {
	var method func(bool) string
	if c.Zone != nil {
		method = c.Zone.GetRsfHost
	} else {
		method = emptyHost
	}
	return reqHost(c.UseHTTPS, method, c.RsfHost, DefaultRsfHost)
}

// 获取apiHost
// 优先使用Zone中的Host信息，如果Zone中的host信息没有配置，那么使用Config中的Host信息
func (c *Config) ApiReqHost() string {
	var method func(bool) string
	if c.Zone != nil {
		method = c.Zone.GetApiHost
	} else {
		method = emptyHost
	}
	return reqHost(c.UseHTTPS, method, c.ApiHost, DefaultAPIHost)
}

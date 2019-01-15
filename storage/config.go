package storage

import (
	"strings"
)

// Config 为文件上传，资源管理等配置
type Config struct {
	//兼容保留
	Zone *Region //空间所在的机房

	Region        *Region
	UseHTTPS      bool   //是否使用https域名
	UseCdnDomains bool   //是否使用cdn加速域名
	CentralRsHost string //中心机房的RsHost，用于list bucket

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

// reqHost 返回一个Host链接
// 如果getHost返回的Host是合法的，那么优先使用该返回值； 否则使用参数host, 当host 为空时，使用默认的host
// 主要用于Config 中Host的获取，Region优先级最高， Zone次之， 最后才使用设置的Host信息
func reqHost(useHttps bool, getHost func(bool) string, host, defaultHost string) (endp string) {
	endp = getHost(useHttps)
	endp = strings.TrimSpace(endp)

	// 返回的host合法
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
	if c.Region != nil {
		method = c.Region.GetRsHost
	} else if c.Zone != nil {
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
	if c.Region != nil {
		method = c.Region.GetRsfHost
	} else if c.Zone != nil {
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
	if c.Region != nil {
		method = c.Region.GetApiHost
	} else if c.Zone != nil {
		method = c.Zone.GetApiHost
	} else {
		method = emptyHost
	}
	return reqHost(c.UseHTTPS, method, c.ApiHost, DefaultAPIHost)
}

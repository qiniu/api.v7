package kodo

import (
	"net/http"

	"qiniupkg.com/api.v7/auth/qbox"
	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/x/rpc.v7"
)

var (
	RS_HOST  = "http://rs.qbox.me"
	RSF_HOST = "http://rsf.qbox.me"
)

// ----------------------------------------------------------

type zoneConfig struct {
	UpHosts []string
}

var zones = []zoneConfig{
	// z0:
	{
		UpHosts: []string{
			"http://up.qiniu.com",
			"http://upload.qiniu.com",
			"-H up.qiniu.com http://183.136.139.16",
		},
	},
	// z1:
	{
		UpHosts: []string{
			"http://up-z1.qiniu.com",
			"http://upload-z1.qiniu.com",
			"-H up-z1.qiniu.com http://106.38.227.27",
		},
	},
}

// ----------------------------------------------------------

type Config struct {
	AccessKey string
	SecretKey string
	RSHost    string
	RSFHost   string
	UpHosts   []string
	Transport http.RoundTripper
}

// ----------------------------------------------------------

type Client struct {
	rpc.Client
	mac *qbox.Mac
	Config
}

func New(zone int, cfg *Config) (p *Client) {

	p = new(Client)
	if cfg != nil {
		p.Config = *cfg
	}

	p.mac = qbox.NewMac(p.AccessKey, p.SecretKey)
	p.Client = rpc.Client{qbox.NewClient(p.mac, p.Transport)}

	if p.RSHost == "" {
		p.RSHost = RS_HOST
	}
	if p.RSFHost == "" {
		p.RSFHost = RSF_HOST
	}
	if len(p.UpHosts) == 0 {
		if zone < 0 || zone >= len(zones) {
			panic("invalid config: invalid zone")
		}
		p.UpHosts = zones[zone].UpHosts
	}
	return
}

// ----------------------------------------------------------

// userApp should be [A-Za-z0-9_\ \-\.]*
//
func SetAppName(userApp string) error {

	return conf.SetAppName(userApp)
}

// ----------------------------------------------------------


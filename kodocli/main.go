package kodocli

import (
	"net/http"

	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/x/rpc.v7"
)

// ----------------------------------------------------------

type ZoneConfig struct {
	UpHosts []string
}

var zones = []ZoneConfig{
	// z0:
	{
		UpHosts: []string{
			"http://upload.qiniu.com",
			"http://up.qiniu.com",
			"-H up.qiniu.com http://183.136.139.16",
		},
	},
	// z1:
	{
		UpHosts: []string{
			"http://upload-z1.qiniu.com",
			"http://up-z1.qiniu.com",
			"-H up-z1.qiniu.com http://106.38.227.27",
		},
	},
}

// ----------------------------------------------------------

type UploadConfig struct {
	UpHosts   []string
	Transport http.RoundTripper
	Zone      int
}

type Uploader struct {
	Conn    rpc.Client
	UpHosts []string
}

func NewUploader(cfg *UploadConfig) (p Uploader) {

	var uc UploadConfig
	if cfg != nil {
		uc = *cfg
	}
	if len(uc.UpHosts) == 0 {
		if uc.Zone >= len(zones) {
			panic("invalid upload config: invalid zone")
		}
		uc.UpHosts = zones[uc.Zone].UpHosts
	}

	p.UpHosts = uc.UpHosts
	p.Conn.Client = &http.Client{Transport: uc.Transport}
	return
}

// ----------------------------------------------------------

// userApp should be [A-Za-z0-9_\ \-\.]*
//
func SetAppName(userApp string) error {

	return conf.SetAppName(userApp)
}

// ----------------------------------------------------------


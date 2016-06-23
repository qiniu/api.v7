package api

import (
	"time"

	"qiniupkg.com/x/rpc.v7"

	. "golang.org/x/net/context"
)

const DefaultApiHost string = "http://uc.qbox.me"

var hostCache = make(map[string]BucketInfo)

type Client struct {
	*rpc.Client
	host   string
	scheme string
}

func NewClient(host string, scheme string) *Client {
	if host == "" {
		host = DefaultApiHost
	}
	client := rpc.DefaultClient
	return &Client{&client, host, scheme}
}

type BucketInfo struct {
	UpHosts []string `json:"up"`
	IoHost  string   `json:"io"`
	Expire  int64    `json:"expire"` // expire == 0 means no expire
}

func (p *Client) GetBucketInfo(ak, bucketName string) (ret BucketInfo, err error) {
	key := ak + ":" + bucketName + ":" + p.scheme
	if info, ok := hostCache[key]; ok && (info.Expire == 0 || info.Expire > time.Now().Unix()) {
		ret = info
		return
	}
	info, err := p.bucketHosts(ak, bucketName)
	if err != nil {
		return
	}
	ret.Expire = time.Now().Unix()
	if p.scheme == "https" {
		ret.UpHosts = info.Https["up"]
		if iohosts, ok := info.Https["io"]; ok && len(iohosts) != 0 {
			ret.IoHost = iohosts[0]
		}
	} else {
		ret.UpHosts = info.Http["up"]
		if iohosts, ok := info.Http["io"]; ok && len(iohosts) != 0 {
			ret.IoHost = iohosts[0]
		}
	}
	hostCache[key] = ret
	return
}

type HostsInfo struct {
	Ttl   int64               `json:"ttl"`
	Http  map[string][]string `json:"http"`
	Https map[string][]string `json:"https"`
}

/*
请求包：
	GET /v1/query?ak=<ak>&&bucket=<bucket>
返回包：
	200 OK {
	  "ttl": <ttl>,              // 有效时间
	  "http": {
	    "up": [],
	    "io": [],                // 当bucket为global时，我们不需要iohost, io缺省
	  },
	  "https": {
	    "up": [],
	    "io": [],                // 当bucket为global时，我们不需要iohost, io缺省
	  }
	}
*/

func (p *Client) bucketHosts(ak, bucket string) (info HostsInfo, err error) {
	ctx := Background()
	err = p.CallWithForm(ctx, &info, "GET", p.host+"/v1/query", map[string][]string{
		"ak":     []string{ak},
		"bucket": []string{bucket},
	})
	return
}

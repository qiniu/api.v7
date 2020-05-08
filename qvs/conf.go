package qvs

import (
	"net/http"

	"github.com/qiniu/api.v7/v7/auth"
)

// ---------------------------------------------------------------------------------------

// APIHost 指定了 API 服务器的地址
var APIHost = "qvs.qiniuapi.com/v1"

// APIHTTPScheme 指定了在请求 API 服务器时使用的 HTTP 模式.
var APIHTTPScheme = "http://"

type transport struct {
	http.RoundTripper
	mac *auth.Credentials
}

func newTransport(mac *auth.Credentials, tr http.RoundTripper) *transport {
	if tr == nil {
		tr = http.DefaultTransport
	}
	return &transport{tr, mac}
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	token, err := t.mac.SignRequestV2(req)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Qiniu "+token)
	return t.RoundTripper.RoundTrip(req)
}
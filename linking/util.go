package linking

import (
	"net/http"

	"github.com/qiniu/api.v7/v7/auth"
)

// ---------------------------------------------------------------------------------------

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

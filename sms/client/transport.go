package client

import (
	"encoding/base64"
	"net/http"
)

// Mac qiniu mac type
type Mac struct {
	AccessKey string
	SecretKey []byte
}

// Transport with qiniu mac
type Transport struct {
	mac       Mac
	Transport http.RoundTripper
}

// RoundTrip transport round trip method
func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	sign, err := SignRequest(t.mac.SecretKey, req)
	if err != nil {
		return
	}

	auth := "Qiniu " + t.mac.AccessKey + ":" + base64.URLEncoding.EncodeToString(sign)
	req.Header.Set("Authorization", auth)
	return t.Transport.RoundTrip(req)
}

// NestedObject return transport
func (t *Transport) NestedObject() interface{} {

	return t.Transport
}

// NewTransport return transport with qiniu mac
func NewTransport(mac *Mac, transport http.RoundTripper) *Transport {

	if transport == nil {
		transport = http.DefaultTransport
	}

	t := &Transport{Transport: transport}
	t.mac = *mac

	return t
}

// NewClient return qiniu mac client
func NewClient(mac *Mac, transport http.RoundTripper) *http.Client {

	t := NewTransport(mac, transport)
	return &http.Client{Transport: t}
}

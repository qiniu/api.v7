package rpc

import (
	"net"
	"net/http"
	"time"
)

// DefaultDailTimeout 默认超时时间: 5秒
const DefaultDailTimeout = time.Duration(5) * time.Second

// DefaultTransport 默认 HTTP Transport
var DefaultTransport = NewTransportTimeout(DefaultDailTimeout, 0)

// NewTransportTimeout 返回指定超时时间的 Transport 对象
func NewTransportTimeout(dial, resp time.Duration) http.RoundTripper {
	t := &http.Transport{ // DefaultTransport
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	t.Dial = (&net.Dialer{
		Timeout:   dial,
		KeepAlive: 30 * time.Second,
	}).Dial
	t.ResponseHeaderTimeout = resp
	return t
}

// NewTransportTimeoutWithConnsPool 返回指定超时时间和最大连接主机数的 Transport 对象
func NewTransportTimeoutWithConnsPool(dial, resp time.Duration, poolSize int) http.RoundTripper {

	t := &http.Transport{ // DefaultTransport
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConnsPerHost: poolSize,
	}
	t.Dial = (&net.Dialer{
		Timeout:   dial,
		KeepAlive: 30 * time.Second,
	}).Dial
	t.ResponseHeaderTimeout = resp
	return t
}

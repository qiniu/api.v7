package api

import (
	"io"
	"io/ioutil"
	"net/http"
)

// BytesFromRequest 读取http.Request.Body的内容到slice中
func BytesFromRequest(r *http.Request) (b []byte, err error) {
	if r.ContentLength == 0 {
		return
	}
	if r.ContentLength > 0 {
		b = make([]byte, int(r.ContentLength))
		_, err = io.ReadFull(r.Body, b)
		return
	}
	return ioutil.ReadAll(r.Body)
}

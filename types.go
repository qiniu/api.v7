package api

import (
	"io"
	"io/ioutil"
	"net/http"
)

// Bool 获取一个bool类型的指针
func Bool(v bool) *bool {
	return &v
}

// String获取一个字符串指针
func String(v string) *string {
	return &v
}

// Int64 获取一个整形的指针
func Int64(v int64) *int64 {
	i := int64(v)
	return &i
}

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

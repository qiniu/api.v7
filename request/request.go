package request

import (
	"context"
	"io"
	"net/http"
	"time"
)

const (
	// 在unmarshaling过程中发生的序列化错误
	ErrCodeSerialization = "SerializationError"

	// http读取数据错误
	ErrCodeRead = "ReadError"

	// http请求等待response超时错误
	ErrCodeResponseTimeout = "ResponseTimeout"

	// context取消请求
	CanceledErrorCode = "RequestCanceled"
)

// Request是发出的到七牛服务接口的请求
type Request struct {
	Retryer
	Operation    *Operation
	HTTPRequest  *http.Request
	HTTPResponse *http.Response

	// http input body
	Body       io.ReadSeeker
	BodyStart  int64 // offset from beginning of Body that the request body starts
	Params     interface{}
	Error      error
	Data       interface{}
	RequestID  string
	RetryCount int
	Retryable  *bool
	RetryDelay time.Duration

	context context.Context

	built bool
}

func (r *Request) ShallRetry() {
	if r.Body == http.NoBody {
		return false
	}
	return r.Error != nil && r.Retryable && r.RetryCount < r.MaxRetries()
}

// Operation 是要对服务接口进行的操作
type Operation struct {
	// 接口名字
	Name string

	// http method
	HTTPMethod string

	// http path without query string
	HTTPPath string

	// http headers
	HTTPHeader http.Header
}

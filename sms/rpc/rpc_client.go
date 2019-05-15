package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// UserAgent user agent
var UserAgent = "Golang qiniu/rpc package"

// --------------------------------------------------------------------

// Client a golang http client
type Client struct {
	*http.Client
}

// DefaultClient a golang default http client
var DefaultClient = Client{&http.Client{Transport: DefaultTransport}}

// NewClientTimeout return a golang http client
func NewClientTimeout(dial, resp time.Duration) Client {
	return Client{&http.Client{Transport: NewTransportTimeout(dial, resp)}}
}

// --------------------------------------------------------------------

// Logger interface
type Logger interface {
	ReqID() string
	Xput(logs []string)
}

// --------------------------------------------------------------------

// Head send head method request
func (r Client) Head(l Logger, url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// Get send get method request
func (r Client) Get(l Logger, url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// Delete send delete method request
func (r Client) Delete(l Logger, url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// PostEx send post method request with url
func (r Client) PostEx(l Logger, url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// PostWith send post method request with url, bodyType, body and bodyLength
func (r Client) PostWith(l Logger, url1 string, bodyType string, body io.Reader, bodyLength int) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url1, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = int64(bodyLength)
	return r.Do(l, req)
}

// PostWith64 send post method request with url, bodyType, body and bodyLength(64)
func (r Client) PostWith64(l Logger, url1 string, bodyType string, body io.Reader, bodyLength int64) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url1, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = bodyLength
	return r.Do(l, req)
}

// PostWithForm send post method request with url and form data
func (r Client) PostWithForm(l Logger, url1 string, data map[string][]string) (resp *http.Response, err error) {
	msg := url.Values(data).Encode()
	return r.PostWith(l, url1, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
}

// PostWithJSON send post method request with url and application/json data
func (r Client) PostWithJSON(l Logger, url1 string, data interface{}) (resp *http.Response, err error) {
	msg, err := json.Marshal(data)
	if err != nil {
		return
	}
	return r.PostWith(l, url1, "application/json", bytes.NewReader(msg), len(msg))
}

// PutEx send put request with url
func (r Client) PutEx(l Logger, url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

// PutWith send put method request with url, bodyType, body and bodyLength
func (r Client) PutWith(l Logger, url1 string, bodyType string, body io.Reader, bodyLength int) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url1, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = int64(bodyLength)
	return r.Do(l, req)
}

// PutWith64 send put method request with url, bodyType, body and bodyLength(64)
func (r Client) PutWith64(l Logger, url1 string, bodyType string, body io.Reader, bodyLength int64) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url1, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = bodyLength
	return r.Do(l, req)
}

// PutWithForm send put method request with url and form data
func (r Client) PutWithForm(l Logger, url1 string, data map[string][]string) (resp *http.Response, err error) {
	msg := url.Values(data).Encode()
	return r.PutWith(l, url1, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
}

// PutWithJSON send put method request with url and application/json data
func (r Client) PutWithJSON(l Logger, url1 string, data interface{}) (resp *http.Response, err error) {
	msg, err := json.Marshal(data)
	if err != nil {
		return
	}
	return r.PutWith(l, url1, "application/json", bytes.NewReader(msg), len(msg))
}

// --------------------------------------------------------------------

// Do 发送 HTTP Request, 并返回 HTTP Response
func (r Client) Do(l Logger, req *http.Request) (resp *http.Response, err error) {
	// debug
	start := time.Now()
	defer func() {
		end := time.Now()
		latency := end.Sub(start)
		if err != nil {

		} else {
			reqid := req.Header.Get("X-Request-Id")
			if rid := resp.Header.Get("X-Request-Id"); len(rid) > 0 {
				reqid = rid
			}
			method := req.Method
			methodColor := colorForMethod(method)
			statusCode := resp.StatusCode
			statusColor := colorForStatus(statusCode)

			fmt.Printf("[RPC] [%s] %v |%s %3d %s| %13v |%s %s %s|\t %s\n",
				reqid,
				start.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, reset,
				latency,
				methodColor, method, reset,
				req.URL.String(),
			)
		}
	}()

	if l != nil {
		req.Header.Set("X-Request-Id", l.ReqID())
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", UserAgent)
	}
	resp, err = r.Client.Do(req)
	if err != nil {
		return
	}

	if l != nil {
		details := resp.Header["X-Log"]
		if len(details) > 0 {
			l.Xput(details)
		}
	}
	return
}

// --------------------------------------------------------------------

// RespError interface
type RespError interface {
	ErrorDetail() string
	Error() string
	HttpCode() int
}

// ErrorInfo type
type ErrorInfo struct {
	Err     string   `json:"error"`
	Reqid   string   `json:"reqid"`
	Details []string `json:"details"`
	Code    int      `json:"code"`
}

// ErrorDetail return error detail
func (r *ErrorInfo) ErrorDetail() string {
	msg, _ := json.Marshal(r)
	return string(msg)
}

// Error return error message
func (r *ErrorInfo) Error() string {
	if r.Err != "" {
		return r.Err
	}
	return http.StatusText(r.Code)
}

// HTTPCode return rpc http StatusCode
func (r *ErrorInfo) HTTPCode() int {
	return r.Code
}

// --------------------------------------------------------------------

type errorRet struct {
	Error string `json:"error"`
}

// ResponseError return response error
func ResponseError(resp *http.Response) (err error) {
	e := &ErrorInfo{
		Details: resp.Header["X-Log"],
		Reqid:   resp.Header.Get("X-Request-Id"),
		Code:    resp.StatusCode,
	}
	if resp.StatusCode > 299 {
		if resp.ContentLength != 0 {
			if ct := resp.Header.Get("Content-Type"); strings.TrimSpace(strings.SplitN(ct, ";", 2)[0]) == "application/json" {
				var ret1 errorRet
				json.NewDecoder(resp.Body).Decode(&ret1)
				e.Err = ret1.Error
			}
		}
	}
	return e
}

// CallRet parse http response
func CallRet(l Logger, ret interface{}, resp *http.Response) (err error) {
	return callRet(l, ret, resp)
}

// callRet parse http response
func callRet(l Logger, ret interface{}, resp *http.Response) (err error) {
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode/100 == 2 || resp.StatusCode/100 == 3 {
		if ret != nil && resp.ContentLength != 0 {
			err = json.NewDecoder(resp.Body).Decode(ret)
			if err != nil {
				return
			}
		}
		return nil
	}
	return ResponseError(resp)
}

// CallWithForm send post method request with url and form data then parse response
func (r Client) CallWithForm(l Logger, ret interface{}, url1 string, param map[string][]string) (err error) {
	resp, err := r.PostWithForm(l, url1, param)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// CallWithJSON send post method request with url and application/json data then parse response
func (r Client) CallWithJSON(l Logger, ret interface{}, url1 string, param interface{}) (err error) {
	resp, err := r.PostWithJSON(l, url1, param)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// CallWith send post method request with url, bodyType, body and bodyLength then parse response
func (r Client) CallWith(l Logger, ret interface{}, url1 string, bodyType string, body io.Reader, bodyLength int) (err error) {
	resp, err := r.PostWith(l, url1, bodyType, body, bodyLength)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// CallWith64 send post method request with url, bodyType, body and bodyLength(64) then parse response
func (r Client) CallWith64(l Logger, ret interface{}, url1 string, bodyType string, body io.Reader, bodyLength int64) (err error) {
	resp, err := r.PostWith64(l, url1, bodyType, body, bodyLength)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// Call send post method request with url then parse response
func (r Client) Call(l Logger, ret interface{}, url1 string) (err error) {
	resp, err := r.PostWith(l, url1, "application/x-www-form-urlencoded", nil, 0)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// PutCallWithForm send put method request with url and param then parse response
func (r Client) PutCallWithForm(l Logger, ret interface{}, url1 string, param map[string][]string) (err error) {
	resp, err := r.PutWithForm(l, url1, param)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// PutCallWithJSON send put method request with url and param then parse response
func (r Client) PutCallWithJSON(l Logger, ret interface{}, url1 string, param interface{}) (err error) {
	resp, err := r.PutWithJSON(l, url1, param)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// PutCallWith send put method request with url, bodyType, body and bodyLength then parse response
func (r Client) PutCallWith(l Logger, ret interface{}, url1 string, bodyType string, body io.Reader, bodyLength int) (err error) {
	resp, err := r.PutWith(l, url1, bodyType, body, bodyLength)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// PutCallWith64 send post method request with url, bodyType, body and bodyLength(64) then parse response
func (r Client) PutCallWith64(l Logger, ret interface{}, url1 string, bodyType string, body io.Reader, bodyLength int64) (err error) {
	resp, err := r.PutWith64(l, url1, bodyType, body, bodyLength)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// PutCall send put method request with url then parse response
func (r Client) PutCall(l Logger, ret interface{}, url1 string) (err error) {
	resp, err := r.PutWith(l, url1, "application/x-www-form-urlencoded", nil, 0)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// GetCall send get method request with url then parse response
func (r Client) GetCall(l Logger, ret interface{}, url1 string) (err error) {
	resp, err := r.Get(l, url1)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// GetCallWithForm send get method request with url and param then parse response
func (r Client) GetCallWithForm(l Logger, ret interface{}, url1 string, param map[string][]string) (err error) {
	payload := url.Values(param).Encode()
	if strings.ContainsRune(url1, '?') {
		url1 += "&"
	} else {
		url1 += "?"
	}
	url1 += payload
	resp, err := r.Get(l, url1)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// DeleteCall send delete method request with url
func (r Client) DeleteCall(l Logger, ret interface{}, url string) (err error) {
	resp, err := r.Delete(l, url)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// --------------------------------------------------------------------

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

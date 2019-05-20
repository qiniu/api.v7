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

// Head send head method request
func (r Client) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return
	}
	return r.Do(req)
}

// Get send get method request
func (r Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	return r.Do(req)
}

// Delete send delete method request
func (r Client) Delete(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}
	return r.Do(req)
}

// PostEx send post method request with url
func (r Client) PostEx(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}
	return r.Do(req)
}

// PostWith send post method request with url, bodyType, body and bodyLength
func (r Client) PostWith(url1 string, bodyType string, body io.Reader, bodyLength int) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url1, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = int64(bodyLength)
	return r.Do(req)
}

// PostWith64 send post method request with url, bodyType, body and bodyLength(64)
func (r Client) PostWith64(url1 string, bodyType string, body io.Reader, bodyLength int64) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url1, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = bodyLength
	return r.Do(req)
}

// PostWithForm send post method request with url and form data
func (r Client) PostWithForm(url1 string, data map[string][]string) (resp *http.Response, err error) {
	msg := url.Values(data).Encode()
	return r.PostWith(url1, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
}

// PostWithJSON send post method request with url and application/json data
func (r Client) PostWithJSON(url1 string, data interface{}) (resp *http.Response, err error) {
	msg, err := json.Marshal(data)
	if err != nil {
		return
	}
	return r.PostWith(url1, "application/json", bytes.NewReader(msg), len(msg))
}

// PutEx send put request with url
func (r Client) PutEx(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return
	}
	return r.Do(req)
}

// PutWith send put method request with url, bodyType, body and bodyLength
func (r Client) PutWith(url1 string, bodyType string, body io.Reader, bodyLength int) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url1, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = int64(bodyLength)
	return r.Do(req)
}

// PutWith64 send put method request with url, bodyType, body and bodyLength(64)
func (r Client) PutWith64(url1 string, bodyType string, body io.Reader, bodyLength int64) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url1, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = bodyLength
	return r.Do(req)
}

// PutWithForm send put method request with url and form data
func (r Client) PutWithForm(url1 string, data map[string][]string) (resp *http.Response, err error) {
	msg := url.Values(data).Encode()
	return r.PutWith(url1, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
}

// PutWithJSON send put method request with url and application/json data
func (r Client) PutWithJSON(url1 string, data interface{}) (resp *http.Response, err error) {
	msg, err := json.Marshal(data)
	if err != nil {
		return
	}
	return r.PutWith(url1, "application/json", bytes.NewReader(msg), len(msg))
}

// --------------------------------------------------------------------

// Do 发送 HTTP Request, 并返回 HTTP Response
func (r Client) Do(req *http.Request) (resp *http.Response, err error) {
	// debug
	start := time.Now()
	defer func() {
		end := time.Now()
		latency := end.Sub(start)
		if err != nil {

		} else {
			reqid := req.Header.Get("X-Reqid")
			if rid := resp.Header.Get("X-Reqid"); len(rid) > 0 {
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

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", UserAgent)
	}

	resp, err = r.Client.Do(req)
	return
}

// --------------------------------------------------------------------

// ErrorInfo type
type ErrorInfo struct {
	Err       string `json:"error"`
	RequestID string `json:"reqid"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
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
		RequestID: resp.Header.Get("X-Reqid"),
		Code:      resp.StatusCode,
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
func CallRet(ret interface{}, resp *http.Response) (err error) {
	return callRet(ret, resp)
}

// callRet parse http response
func callRet(ret interface{}, resp *http.Response) (err error) {
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
func (r Client) CallWithForm(ret interface{}, url1 string, param map[string][]string) (err error) {
	resp, err := r.PostWithForm(url1, param)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
}

// CallWithJSON send post method request with url and application/json data then parse response
func (r Client) CallWithJSON(ret interface{}, url1 string, param interface{}) (err error) {
	resp, err := r.PostWithJSON(url1, param)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
}

// CallWith send post method request with url, bodyType, body and bodyLength then parse response
func (r Client) CallWith(ret interface{}, url1 string, bodyType string, body io.Reader, bodyLength int) (err error) {
	resp, err := r.PostWith(url1, bodyType, body, bodyLength)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
}

// CallWith64 send post method request with url, bodyType, body and bodyLength(64) then parse response
func (r Client) CallWith64(ret interface{}, url1 string, bodyType string, body io.Reader, bodyLength int64) (err error) {
	resp, err := r.PostWith64(url1, bodyType, body, bodyLength)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
}

// Call send post method request with url then parse response
func (r Client) Call(ret interface{}, url1 string) (err error) {
	resp, err := r.PostWith(url1, "application/x-www-form-urlencoded", nil, 0)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
}

// PutCallWithForm send put method request with url and param then parse response
func (r Client) PutCallWithForm(ret interface{}, url1 string, param map[string][]string) (err error) {
	resp, err := r.PutWithForm(url1, param)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
}

// PutCallWithJSON send put method request with url and param then parse response
func (r Client) PutCallWithJSON(ret interface{}, url1 string, param interface{}) (err error) {
	resp, err := r.PutWithJSON(url1, param)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
}

// PutCallWith send put method request with url, bodyType, body and bodyLength then parse response
func (r Client) PutCallWith(ret interface{}, url1 string, bodyType string, body io.Reader, bodyLength int) (err error) {
	resp, err := r.PutWith(url1, bodyType, body, bodyLength)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
}

// PutCallWith64 send post method request with url, bodyType, body and bodyLength(64) then parse response
func (r Client) PutCallWith64(ret interface{}, url1 string, bodyType string, body io.Reader, bodyLength int64) (err error) {
	resp, err := r.PutWith64(url1, bodyType, body, bodyLength)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
}

// PutCall send put method request with url then parse response
func (r Client) PutCall(ret interface{}, url1 string) (err error) {
	resp, err := r.PutWith(url1, "application/x-www-form-urlencoded", nil, 0)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
}

// GetCall send get method request with url then parse response
func (r Client) GetCall(ret interface{}, url1 string) (err error) {
	resp, err := r.Get(url1)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
}

// GetCallWithForm send get method request with url and param then parse response
func (r Client) GetCallWithForm(ret interface{}, url1 string, param map[string][]string) (err error) {
	payload := url.Values(param).Encode()
	if strings.ContainsRune(url1, '?') {
		url1 += "&"
	} else {
		url1 += "?"
	}
	url1 += payload
	resp, err := r.Get(url1)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
}

// DeleteCall send delete method request with url
func (r Client) DeleteCall(ret interface{}, url string) (err error) {
	resp, err := r.Delete(url)
	if err != nil {
		return err
	}
	return callRet(ret, resp)
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

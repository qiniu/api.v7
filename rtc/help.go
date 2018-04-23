package rtc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/qiniu/api.v7/auth/qbox"
)

// ResInfo is httpresponse infomation
type ResInfo struct {
	URL    string
	Method string
	Code   int
	Err    error
	Msg    string
	Header map[string][]string
}

func NewResInfo() ResInfo {
	info := ResInfo{}
	info.Header = make(map[string][]string)
	return info
}

func CopyHttpHeader(src *http.Header, t *ResInfo) {
	for k, v := range *src {
		K := strings.Title(k)
		if strings.Contains(K, "Reqid") || K == "Content-Length" {
			t.Header[k] = v
		}
	}
}

func buildURL(path string) string {
	if strings.Index(path, "/") != 0 {
		path = "/" + path
	}
	return "https://" + RtcHost + path
}

func postReq(httpClient *http.Client, mac *qbox.Mac, url string,
	reqParam interface{}, ret interface{}) *ResInfo {
	info := NewResInfo()
	var reqData []byte
	var err error

	switch v := reqParam.(type) {
	case *string:
		reqData = []byte(*v)
	case string:
		reqData = []byte(v)
	case *[]byte:
		reqData = *v
	case []byte:
		reqData = v
	default:
		reqData, err = json.Marshal(reqParam)
	}

	if err != nil {
		info.Err = err
		return &info
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqData))
	if err != nil {
		info.Err = err
		return &info
	}
	req.Header.Add("Content-Type", "application/json")
	return callReq(httpClient, req, mac, &info, ret)
}

func getReq(httpClient *http.Client, mac *qbox.Mac, url string, ret interface{}) *ResInfo {
	info := NewResInfo()
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		info.Err = err
		return &info
	}
	return callReq(httpClient, req, mac, &info, ret)
}

func delReq(httpClient *http.Client, mac *qbox.Mac, url string, ret interface{}) *ResInfo {
	info := NewResInfo()
	req, err := http.NewRequest("DELETE", url, strings.NewReader(""))
	if err != nil {
		info.Err = err
		return &info
	}
	return callReq(httpClient, req, mac, &info, ret)
}

func callReq(httpClient *http.Client, req *http.Request, mac *qbox.Mac,
	info *ResInfo, ret interface{}) (oinfo *ResInfo) {
	oinfo = info
	accessToken, err := mac.SignRequestV2(req)
	if err != nil {
		info.Err = err
		return
	}
	req.Header.Add("Authorization", "Qiniu "+accessToken)
	client := httpClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		info.Err = err
		return
	}

	defer resp.Body.Close()
	info.Method = req.Method
	info.URL = req.URL.RequestURI()
	info.Code = resp.StatusCode
	CopyHttpHeader(&resp.Header, info)
	if resp.ContentLength > 2*1024*1024 {
		err = fmt.Errorf("response is too long. Content-Length: %v", resp.ContentLength)
		info.Err = err
		return
	}
	resData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		info.Err = err
		return
	}
	if info.Code != 200 {
		info.Msg = string(resData)
		return
	}
	if ret != nil {
		err = json.Unmarshal(resData, ret)
		info.Err = err
		if err != nil {
			info.Msg = string(resData)
		}
	}
	return
}

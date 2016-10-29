package media

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
	"net/url"
	"qiniupkg.com/api.v7/auth/qbox"
	"qiniupkg.com/api.v7/conf"
	"regexp"
	"strings"
)

const (
	domain = "http://api.qiniu.com/pfop/" //pfop地址
)

var (
	Bucket  = "" //bucket
	Pipline = ""
)

type Options struct {
	Bucket              string `url:"bucket,omitempty"`
	NotifyURL           string `url:"notifyURL,omitempty"`
	Pipeline            string `url:"pipeline,omitempty"`
	NeedConvertFileName string `url:"key"`
	Fops                string `url:"fops"`
}

type Result struct {
	Error        string `json:"error,omitempty"`
	PersistentId string `json:"persistentId,omitempty"`
}

func get(params ...interface{}) (body []byte, err error) {
	if len(params) == 0 {
		return
	}
	urlStr, ok := params[0].(string)
	if !ok {
		return
	}
	req, err := http.NewRequest("GET", urlStr, strings.NewReader(""))
	if err != nil {
		return
	}
	if len(params) >= 2 {
		withToken, ok := params[1].(bool)
		if !ok {
			return
		}
		if withToken {
			mac := qbox.NewMac(conf.ACCESS_KEY, conf.SECRET_KEY)
			token, _ := mac.SignRequest(req, false)
			req.URL, _ = url.Parse(urlStr + fmt.Sprintf("&token=%s", token))
		}
	}
	body, err = request(req)
	return
}

func post(params Options) (result Result, err error) {
	if len(params.Bucket) == 0 {
		params.Bucket = Bucket
	}
	if len(params.Pipeline) == 0 {
		params.Pipeline = Pipline
	}
	if conf.ACCESS_KEY == "" || conf.SECRET_KEY == "" || len(params.Bucket) == 0 || len(params.Pipeline) == 0 {
		err = errors.New("missing some required parameters")
		return
	}
	v, _ := query.Values(params)
	req, err := http.NewRequest("POST", domain, strings.NewReader(v.Encode()))
	if err != nil {
		return
	}
	mac := qbox.NewMac(conf.ACCESS_KEY, conf.SECRET_KEY)
	token, _ := mac.SignRequest(req, true)
	req.Header.Add("Authorization", fmt.Sprintf("QBox %s", token))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	body, err := request(req)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}
	if len(result.Error) > 0 {
		err = errors.New(result.Error)
	}
	return
}

func request(req *http.Request) (body []byte, err error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return

}

func makeFops(queryStruct interface{}) (fops string) {
	v, _ := query.Values(queryStruct)
	queryStr := v.Encode()
	fops = convertQueryStrToFopStr(queryStr)
	return fops
}

func convertQueryStrToFopStr(str string) (fops string) {
	fops = regexp.MustCompile(`[&=]`).ReplaceAllString(str, "/")
	fops = regexp.MustCompile(`\d*-`).ReplaceAllString(fops, "")
	fops = fixedDollarToEqual(fops)
	return
}

/**
only for media package
 */
func UrlBase64Encode(str string) (base64Str string) {
	base64Str = base64.StdEncoding.EncodeToString([]byte(str))
	base64Str = fixedEqualToDollar(base64Str)
	return
}

/**
fixed base64 equals conflict with url query encode equals
 */
func fixedEqualToDollar(str string) string {
	return regexp.MustCompile(`[=]`).ReplaceAllString(str, "$")
}

/**
fixed base64 equals from dollar after query encode
 */
func fixedDollarToEqual(str string) string {
	return regexp.MustCompile(`[$]`).ReplaceAllString(str, "=")
}

/**
only for media package
 */
func UrlSafeBase64Encode(str string) (base64Str string) {
	base64Str = base64.StdEncoding.EncodeToString([]byte(str))
	base64Str = regexp.MustCompile(`[+]`).ReplaceAllString(base64Str, "-")
	base64Str = regexp.MustCompile(`[/]`).ReplaceAllString(base64Str, "_")
	base64Str = fixedEqualToDollar(base64Str)
	return
}

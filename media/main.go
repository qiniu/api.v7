package media

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/conf"
	"io/ioutil"
	"net/http"
	"qiniupkg.com/x/errors.v7"
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

func put(params Options) (result Result, err error) {
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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}
	if len(result.Error) > 0 {
		err = errors.New(result.Error)
	}
	return
}

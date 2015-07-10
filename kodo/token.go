package kodo

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"qiniupkg.com/api.v7/auth/qbox"
	"qiniupkg.com/x/url.v7"
)

// ----------------------------------------------------------

// 根据空间(Bucket)的域名，以及文件的 key，获得 baseUrl。
// 如果空间是 public 的，那么通过 baseUrl 可以直接下载文件内容。
// 如果空间是 private 的，那么需要对 baseUrl 进行私有签名得到一个临时有效的 privateUrl 进行下载。
//
func MakeBaseUrl(domain, key string) (baseUrl string) {

	return "http://" + domain + "/" + url.Escape(key)
}

// ----------------------------------------------------------

type GetPolicy struct {
	Expires uint32
}

func (p *Client) MakePrivateUrl(baseUrl string, policy *GetPolicy) (privateUrl string) {

	var expires int64
	if policy == nil || policy.Expires == 0 {
		expires = 3600
	} else {
		expires = int64(policy.Expires)
	}
	deadline := time.Now().Unix() + expires

	if strings.Contains(baseUrl, "?") {
		baseUrl += "&e="
	} else {
		baseUrl += "?e="
	}
	baseUrl += strconv.FormatInt(deadline, 10)

	token := qbox.Sign(p.mac, []byte(baseUrl))
	return baseUrl + "&token=" + token
}

// --------------------------------------------------------------------------------

type PutPolicy struct {
	Scope               string `json:"scope"`
	Expires             uint32 `json:"deadline"`             // 截止时间（以秒为单位）
	InsertOnly          uint16 `json:"insertOnly,omitempty"` // 若非0, 即使Scope为 Bucket:Key 的形式也是insert only
	DetectMime          uint8  `json:"detectMime,omitempty"` // 若非0, 则服务端根据内容自动确定 MimeType
	CallbackFetchKey    uint8  `json:"callbackFetchKey,omitempty"`
	FsizeLimit          int64  `json:"fsizeLimit,omitempty"`
	MimeLimit           string `json:"mimeLimit,omitempty"`
	SaveKey             string `json:"saveKey,omitempty"`
	CallbackUrl         string `json:"callbackUrl,omitempty"`
	CallbackHost        string `json:"callbackHost,omitempty"`
	CallbackBody        string `json:"callbackBody,omitempty"`
	CallbackBodyType    string `json:"callbackBodyType,omitempty"`
	ReturnUrl           string `json:"returnUrl,omitempty"`
	ReturnBody          string `json:"returnBody,omitempty"`
	PersistentOps       string `json:"persistentOps,omitempty"`
	PersistentNotifyUrl string `json:"persistentNotifyUrl,omitempty"`
	PersistentPipeline  string `json:"persistentPipeline,omitempty"`
	AsyncOps            string `json:"asyncOps,omitempty"`
	EndUser             string `json:"endUser,omitempty"`
	Checksum            string `json:"checksum,omitempty"` // 格式：<HashName>:<HexHashValue>，目前支持 MD5/SHA1。
}

func (p *Client) MakeUptoken(policy *PutPolicy) string {

	var rr = *policy
	if rr.Expires == 0 {
		rr.Expires = 3600
	}
	rr.Expires += uint32(time.Now().Unix())
	b, _ := json.Marshal(&rr)
	return qbox.SignWithData(p.mac, b)
}

// ----------------------------------------------------------

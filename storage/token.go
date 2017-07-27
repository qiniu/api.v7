package storage

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/qiniu/api.v7/auth/qbox"
	"strings"
	"time"
)

type PutPolicy struct {
	Scope               string `json:"scope"`
	isPrefixalScope     int    `json:"isPrefixalScope"`
	Expires             uint32 `json:"deadline"`             // 截止时间（以秒为单位）
	InsertOnly          uint16 `json:"insertOnly,omitempty"` // 若非0, 即使Scope为 Bucket:Key 的形式也是insert only
	DetectMime          uint8  `json:"detectMime,omitempty"` // 若非0, 则服务端根据内容自动确定 MimeType
	FsizeLimit          int64  `json:"fsizeLimit,omitempty"`
	MimeLimit           string `json:"mimeLimit,omitempty"`
	SaveKey             string `json:"saveKey,omitempty"`
	CallbackFetchKey    uint8  `json:"callbackFetchKey,omitempty"`
	CallbackUrl         string `json:"callbackUrl,omitempty"`
	CallbackHost        string `json:"callbackHost,omitempty"`
	CallbackBody        string `json:"callbackBody,omitempty"`
	CallbackBodyType    string `json:"callbackBodyType,omitempty"`
	ReturnUrl           string `json:"returnUrl,omitempty"`
	ReturnBody          string `json:"returnBody,omitempty"`
	PersistentOps       string `json:"persistentOps,omitempty"`
	PersistentNotifyUrl string `json:"persistentNotifyUrl,omitempty"`
	PersistentPipeline  string `json:"persistentPipeline,omitempty"`
	EndUser             string `json:"endUser,omitempty"`
	DeleteAfterDays     int    `json:"deleteAfterDays,omitempty"`
	FileType            int    `json:"fileType,omitempty"`
}

func (p *PutPolicy) UploadToken(mac *qbox.Mac) (token string) {
	if p.Expires == 0 {
		p.Expires = 3600 // 1 hour
	}
	p.Expires += uint32(time.Now().Unix())

	putPolicyJson, _ := json.Marshal(p)
	token = mac.SignWithData(putPolicyJson)
	return
}

func getAkBucketFromUploadToken(token string) (ak, bucket string, err error) {
	items := strings.Split(token, ":")
	if len(items) != 3 {
		err = errors.New("invalid upload token, format error")
		return
	}

	ak = items[0]
	policyBytes, dErr := base64.URLEncoding.DecodeString(items[2])
	if dErr != nil {
		err = errors.New("invalid upload token, invalid put policy")
		return
	}

	putPolicy := PutPolicy{}
	uErr := json.Unmarshal(policyBytes, &putPolicy)
	if uErr != nil {
		err = errors.New("invalid upload token, invalid put policy")
		return
	}

	bucket = strings.Split(putPolicy.Scope, ":")[0]
	return
}

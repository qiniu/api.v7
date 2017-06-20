package media

import (
	"encoding/json"
	"errors"
	"fmt"
)

type PrivateM3U8 struct {
	Mode     string `url:"11-pm3u8,omitempty"`
	Expires  string `url:"12-expires,omitempty"`
	Deadline string `url:"13-deadline,omitempty"`
}

type M3U8Result struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
	Body  string `json:"body,omitempty"`
}

func NewPrivateM3U8() PrivateM3U8 {
	return PrivateM3U8{Mode: "0", Expires: "43200"}
}

func (this PrivateM3U8) Download(M3U8DownloadURI string,expiresTs int64) (result M3U8Result, err error) {
	fops := makeFops(this)
	M3U8DownloadURI = fmt.Sprintf("%s?%s&e=%d", M3U8DownloadURI, fops, expiresTs)
	body, err := get(M3U8DownloadURI, true)
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
	result.Body = string(body)
	return
}

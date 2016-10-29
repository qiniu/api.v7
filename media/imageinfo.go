package media

import (
	"encoding/json"
	"fmt"
	"errors"
)

type ImageInfoResult struct {
	Code        json.Number `json:"code,omitempty"`
	Error       string      `json:"error,omitempty"`
	Format      string      `json:"format,omitempty"`
	Width       json.Number `json:"width,omitempty"`
	Height      json.Number `json:"height,omitempty"`
	ColorModel  string      `json:"colorModel,omitempty"`
	FrameNumber string      `json:"frameNumber,omitempty"`
}

/**
http://developer.qiniu.com/code/v6/api/kodo-api/image/imageinfo.html
 */
func GetImageInfo(imgUrl string) (result ImageInfoResult,err error) {
	body,err:=get(fmt.Sprintf("%s?imageInfo",imgUrl))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}
	if len(result.Error) > 0 {
		err = errors.New(result.Error)
	}
	return
}
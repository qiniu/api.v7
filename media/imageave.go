package media

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ImageAVEResult struct {
	Code  json.Number `json:"code,omitempty"`
	Error string      `json:"error,omitempty"`
	RGB   string      `json:"RGB,omitempty"`
}

/**
http://developer.qiniu.com/code/v6/api/kodo-api/image/imageave.html
 */
func GetImageAVE(imgUrl string) (result ImageAVEResult, err error) {
	body, err := get(fmt.Sprintf("%s?imageAve", imgUrl))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}
	if len(result.Error) > 0 {
		err = errors.New(result.Error)
	}
	return
}

package media

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ImageExifResult struct {
	Code   json.Number `json:"code,omitempty"`
	Error  string      `json:"error,omitempty"`
	Info   interface{} `json:"info,omitempty"`
}

/**
http://developer.qiniu.com/code/v6/api/kodo-api/image/exif.html
*/
func GetImageExif(imgUrl string) (result ImageExifResult, err error) {
	body, err := get(fmt.Sprintf("%s?exif", imgUrl))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}
	if len(result.Error) > 0 {
		err = errors.New(result.Error)
		return
	}
	jsonStrWarp:=fmt.Sprintf(`{"info":%s}`,string(body))
	err = json.Unmarshal([]byte(jsonStrWarp), &result)
	if err != nil {
		return
	}
	return
}

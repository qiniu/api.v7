package media

import (
	"fmt"
	"encoding/json"
	"errors"
)

type AvInfoResult struct {
	Streams []interface{} `json:"streams,omitempty"`
	Format  interface{}   `json:"format,omitempty"`
	Error   string        `json:"error,omitempty"`
}

func Avinfo(AvDownloadURI string)(result AvInfoResult,err error)  {
	body,err:=get(fmt.Sprintf("%s?avinfo",AvDownloadURI))
	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}
	if len(result.Error) > 0 {
		err = errors.New(result.Error)
	}
	return
}
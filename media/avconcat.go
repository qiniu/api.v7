package media

import (
	"errors"
	"github.com/google/go-querystring/query"
)

type AvConcatView struct {
	Mode   string   `url:"11-avconcat,omitempty"`
	Format string   `url:"12-format,omitempty"`
	Urls   []string `url:"-"`
}

func NewAvConcat() AvConcatView {
	view := AvConcatView{}
	view.Urls = make([]string, 0,20)
	return AvConcatView{}
}

/**
http://developer.qiniu.com/code/v6/api/dora-api/av/avconcat.html
*/
func (this AvConcatView) makeFops() string {
	v, _ := query.Values(this)
	queryStr := v.Encode()
	for _, urlStr := range this.Urls {
		queryStr += "/" + UrlSafeBase64Encode(urlStr)
	}
	fops := convertQueryStrToFopStr(queryStr)
	return fops
}

func (this AvConcatView) AvConcat(params Options) (result Result, err error) {
	if len(this.Urls) == 0 {
		err = errors.New("need Urls parameters")
		return
	}
	fops := this.makeFops()
	params.Fops = fops
	result, err = post(params)
	return
}

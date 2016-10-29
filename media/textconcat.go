package media

import (
	"errors"
	"github.com/google/go-querystring/query"
	"qiniupkg.com/x/url.v7"
)

type TextConcatView struct {
	EncodedMimeType string   `url:"11-concat/mimeType,omitempty"`
	Urls            []string `url:"-"`
}

func NewTextConcatView(encodedMimeType string) TextConcatView {
	view := TextConcatView{EncodedMimeType: UrlSafeBase64Encode(encodedMimeType)}
	view.Urls = make([]string, 0, 1000)
	return view
}

/**
http://developer.qiniu.com/code/v6/api/dora-api/concat.html#concat-specification
*/
func (this TextConcatView) makeFops() string {
	v, _ := query.Values(this)
	queryStr,_ := url.Unescape( v.Encode())
	for _, urlStr := range this.Urls {
		queryStr += "/" + UrlSafeBase64Encode(urlStr)
	}
	fops := convertQueryStrToFopStr(queryStr)
	return fops
}

func (this TextConcatView) TextConcat(params Options) (result Result, err error) {
	if len(this.Urls) == 0 {
		err = errors.New("need Urls parameters")
		return
	}
	params.Fops = this.makeFops()
	result, err = post(params)
	return
}

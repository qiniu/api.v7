package media

type VSampleView struct {
	Format     string `url:"11-vsample,omitempty"`
	StartTime  string `url:"12-ss,omitempty"`
	Duration   string `url:"13-t,omitempty"`
	Resolution string `url:"14-s,omitempty"`
	Degree     string `url:"15-rotate,omitempty"`
	Interval   string `url:"16-interval,omitempty"`
	Pattern    string `url:"17-pattern,omitempty"` //need use UrlSafeBase64Encode at external
}

func NewVSample() VSampleView {
	return VSampleView{}
}

/**
http://developer.qiniu.com/code/v6/api/dora-api/av/vsample.html
 */

func (this VSampleView) VSample(params Options) (result Result, err error) {
	fops := makeFops(this)
	params.Fops = fops
	result, err = post(params)
	return
}

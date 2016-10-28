package media

type AvSegtimeView struct {
	Format     string `url:"11-avthumb,omitempty"`
	NoDomain   string `url:"12-noDomain,omitempty"`
	Domain     string `url:"13-domain,omitempty"`
	SegTime    string `url:"14-segtime,omitempty"`
	AB         string `url:"15-ab,omitempty"`
	AQ         string `url:"16-aq,omitempty"`
	AR         string `url:"17-ar,omitempty"`
	R          string `url:"18-r,omitempty"`
	VB         string `url:"19-vb,omitempty"`
	VCodec     string `url:"20-vcodec,omitempty"`
	ACodec     string `url:"21-acodec,omitempty"`
	SCodec     string `url:"22-scodec,omitempty"`
	Subtitle   string `url:"23-subtitle,omitempty"`
	SS         string `url:"24-ss,omitempty"`
	T          string `url:"25-t,omitempty"`
	S          string `url:"26-s,omitempty"`
	StripMeta  string `url:"27-stripmeta,omitempty"`
	Rotate     string `url:"29-rotate,omitempty"`
	HlsKey     string `url:"30-hlsKey,omitempty"`
	HlsKeyType string `url:"31-hlsKeyType,omitempty"`
	HlsKeyUrl  string `url:"32-hlsKeyUrl,omitempty"`
	Pattern    string `url:"33-pattern,omitempty"`
}

func NewAvSegtime() AvSegtimeView {
	return AvSegtimeView{}
}

/**
http://developer.qiniu.com/code/v6/api/dora-api/av/segtime.html
 */

func (this AvSegtimeView) AvSegtime(params Options) (result Result, err error) {
	if len(this.Domain) > 0 {
		this.Domain = UrlBase64Encode(this.Domain)
	}
	if len(this.HlsKeyUrl) > 0 {
		this.HlsKeyUrl = UrlSafeBase64Encode(this.HlsKeyUrl)
	}
	fops := makeFops(this)
	params.Fops = fops
	result, err = post(params)
	return
}

package media


type AuthumbView struct {
	Format        string `url:"11-avthumb,omitempty"`
	AB            string `url:"12-ab,omitempty"`
	AQ            string `url:"13-aq,omitempty"`
	AR            string `url:"14-ar,omitempty"`
	R             string `url:"15-r,omitempty"`
	VB            string `url:"16-vb,omitempty"`
	VCodec        string `url:"17-vcodec,omitempty"`
	ACodec        string `url:"18-acodec,omitempty"`
	AudioProfile  string `url:"19-audioProfile,omitempty"`
	SCodec        string `url:"20-scodec,omitempty"`
	Subtitle      string `url:"21-subtitle,omitempty"`
	SS            string `url:"22-ss,omitempty"`
	T             string `url:"23-t,omitempty"`
	S             string `url:"24-s,omitempty"`
	AutoScale     string `url:"25-autoscale,omitempty"`
	Aspect        string `url:"26-aspect,omitempty"`
	StripMeta     string `url:"27-stripmeta,omitempty"`
	H264Crf       string `url:"28-h264Crf,omitempty"`
	Rotate        string `url:"29-rotate,omitempty"`
	WmImage       string `url:"30-wmImage,omitempty"`
	WmGravity     string `url:"31-wmGravity,omitempty"`
	WmText        string `url:"32-wmText,omitempty"`
	WmGravityText string `url:"33-wmGravityText,omitempty"`
	WmFont        string `url:"34-wmFont,omitempty"`
	WmFontColor   string `url:"35-wmFontColor,omitempty"`
	WmFontSize    string `url:"36-wmFontSize,omitempty"`
	WriteXing     string `url:"37-writeXing,omitempty"`
	AN            string `url:"38-an,omitempty"`
	VN            string `url:"39-vn,omitempty"`
}

func NewAvthumb() AuthumbView {
	return AuthumbView{}
}

/**
http://developer.qiniu.com/code/v6/api/dora-api/av/avthumb.html
 */
func (this AuthumbView) Avthumb(params Options) (result Result, err error) {
	if len(this.Subtitle) > 0 {
		this.Subtitle = UrlSafeBase64Encode(this.Subtitle)
	}
	if len(this.WmImage) > 0 {
		this.WmImage = UrlSafeBase64Encode(this.WmImage)
	}
	if len(this.WmText) > 0 {
		this.WmText = UrlSafeBase64Encode(this.WmText)
	}
	if len(this.WmFont) > 0 {
		this.WmFont = UrlSafeBase64Encode(this.WmFont)
	}
	if len(this.WmFontColor) > 0 {
		this.WmFontColor = UrlSafeBase64Encode(this.WmFontColor)
	}
	fops := makeFops(this)
	params.Fops = fops
	result, err = post(params)
	return
}

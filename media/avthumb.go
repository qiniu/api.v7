package media

import (
	"encoding/base64"
	"github.com/google/go-querystring/query"
	"regexp"
)

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

func (this AuthumbView) makeFops() string {
	if len(this.AB) > 0 {
		this.AB = base64.StdEncoding.EncodeToString([]byte(this.AB))
	}
	v, _ := query.Values(this)
	fops := regexp.MustCompile(`[&=]`).ReplaceAllString(v.Encode(), "/")
	fops = regexp.MustCompile(`\d*-`).ReplaceAllString(fops, "")
	return fops
}

func (this AuthumbView) Avthumb(params Options) (result Result, err error) {
	fops := this.makeFops()
	params.Fops = fops
	result, err = put(params)
	return
}

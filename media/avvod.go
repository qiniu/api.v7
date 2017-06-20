package media

import "fmt"

type AvvodView struct {
	Format       string `url:"11-avvod,omitempty"`
	BitRate      string `url:"12-ab,omitempty"`
	AudioQuality string `url:"13-aq,omitempty"`
	SamplingRate string `url:"14-ar,omitempty"`
	FrameRate    string `url:"15-r,omitempty"`
	VideoBitRate string `url:"16-vb,omitempty"`
	VideoCodec   string `url:"17-vcodec,omitempty"`
	AudioCodec   string `url:"18-acodec,omitempty"`
	Resolution   string `url:"19-s,omitempty"`
}

func NewAvvod() AvvodView {
	return AvvodView{Format: "m3u8",Resolution:"960x640",VideoBitRate:"1000k"}
}

/**
http://developer.qiniu.com/code/v6/api/dora-api/av/avvod.html
 */
func (this AvvodView) Avvod(videoUrl string) (m3u8Url string) {
	fops := makeFops(this)
	m3u8Url = fmt.Sprintf("%s?%s", videoUrl, fops)
	return
}

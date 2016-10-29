package media

type AdaptView struct {
	Format          string `url:"11-adapt,omitempty"`
	EnvBandWidth    string `url:"12-envBandWidth,omitempty"`
	MultiVb         string `url:"13-multiVb,omitempty"`
	MultiAb         string `url:"14-multiAb,omitempty"`
	MultiResolution string `url:"15-multiResolution,omitempty"`
	VideoBitRate    string `url:"16-vb,omitempty"`
	AudioBitRate    string `url:"17-ab,omitempty"`
	Resolution      string `url:"18-resolution,omitempty"`
	HlsTime         string `url:"19-hlstime,omitempty"`
}

func NewAdapt() AdaptView {
	return AdaptView{Format: "m3u8"}
}

/**
http://developer.qiniu.com/code/v6/api/dora-api/av/adapt.html
*/
func (this AdaptView) Adapt(params Options) (result Result, err error) {
	fops := makeFops(this)
	params.Fops = fops
	result, err = post(params)
	return
}

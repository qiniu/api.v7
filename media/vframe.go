package media

type VFrameView struct {
	Format string `url:"11-vframe,omitempty"`
	Offset string `url:"12-offset,omitempty"`
	Width  string `url:"13-w,omitempty"`
	Height string `url:"14-h,omitempty"`
	Degree string `url:"15-rotate,omitempty"`
}

func NewVFrame() VFrameView {
	return VFrameView{}
}

func (this VFrameView) VFrame(params Options) (result Result, err error) {
	fops:=makeFops(this)
	params.Fops = fops
	result, err = post(params)
	return
}
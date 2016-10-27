package media

type AuthumbView struct {
	Format string // 封装格式
}

func (this AuthumbView) MakeFops() string {

	fops := "avthumb/" + this.Format
	return fops
}

func Avthumb(params Options) (result Result, err error) {
	authumbView := AuthumbView{
		Format: "mp4",
	}
	fops := authumbView.MakeFops() //fop命令
	params.Fops = fops
	result, err = put(params)
	return
}
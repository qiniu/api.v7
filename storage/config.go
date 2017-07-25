package storage

type Config struct {
	Zone          *Zone //空间所在的机房
	UseHttps      bool  //是否使用https域名
	UseCdnDomains bool  //是否使用cdn加速域名
}

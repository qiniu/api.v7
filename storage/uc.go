// package storage 提供了用户配置(uc)方面的功能, 定义了UC API 的返回结构体类型
package storage

import (
	"strings"
)

// BucketInfo 存储空间的详细信息
type BucketInfo struct {
	// 镜像回源地址， 接口返回的多个地址以；分割
	Source string `json:"source"`

	// 镜像回源的时候请求头中的HOST
	Host string `json:"host"`

	// 镜像回源地址过期时间(秒数)， 现在这个功能没有实现，因此这个字段现在是没有意义的
	Expires int `json:"expires"`

	// 是否开启了原图保护
	Protected bool `json:"protected"`

	// 是否是私有空间
	Private bool `json:"private"`

	// 如果NoIndexPage是false表示开启了空间根目录index.html
	// 如果是true, 表示没有开启
	// 开启了根目录下的index.html, 文件将会被作为默认首页展示
	NoIndexPage bool `json:"no_index_page"`

	// 图片样式分隔符， 接口返回的可能有多个
	Separator string `json:"separator"`

	// 图片样式， map中的key表示图片样式命令名字
	// map中的value表示图片样式命令的内容
	Styles map[string]string `json:"styles"`

	// 该字段已经废弃
	RefreshTime string `json:"refresh_time"`

	// 防盗链模式
	// 1 - 表示设置了防盗链的referer白名单
	// 2 - 表示设置了防盗链的referer黑名单
	AntiLeechMode int `json:"anti_leech_mode"`

	// 使用token签名进行防盗链
	// 0 - 表示关闭
	// 1 - 表示开启
	TokenAntiLeechMode int `json:"token_anti_leech"`

	// 防盗链referer白名单列表
	ReferWl []string `json:"refer_wl"`

	// 防盗链referer黑名单列表
	ReferBl []string `json:"refer_bl`

	// 是否允许空的referer访问
	NoRefer bool `json:"no_refer"`

	// 用于防盗链token的生成
	MacKey string `json:"mac_key"`

	// 用于防盗链token的生成
	MacKey2 string `json:"mac_key2"`

	// 存储区域， 兼容保留
	Zone string

	// 存储区域
	Region string

	// 是否是全局域名
	Global bool
}

// ReferAntiLeechConfig 是用户存储空间的Refer防盗链配置
type ReferAntiLeechConfig struct {
	// 防盗链模式， 0 - 关闭Refer防盗链, 1 - 开启Referer白名单，2 - 开启Referer黑名单
	Mode int

	// 是否允许空的referer访问
	AllowEmptyReferer bool

	// Pattern 匹配HTTP Referer头, 当模式是1或者2的时候有效
	// Mode为1的时候表示允许Referer符合该Pattern的HTTP请求访问
	// Mode为2的时候表示禁止Referer符合该Pattern的HTTP请求访问
	// 当前允许的匹配字符串格式分为三种:
	// 一种为空主机头域名, 比如 foo.com; 一种是泛域名, 比如 *.bar.com;
	// 一种是完全通配符, 即一个 *;
	// 多个规则之间用;隔开, 比如: foo.com;*.bar.com;sub.foo.com;*.sub.bar.com
	Pattern string

	// 是否开启源站的防盗链， 默认为0， 只开启CDN防盗链， 当设置为1的时候
	// 在源站支持的情况下开启源站的Referer防盗链
	EnableSource bool
}

// SetMode 设置referer防盗链模式
func (r *ReferAntiLeechConfig) SetMode(mode int) *ReferAntiLeechConfig {
	if mode != 0 && mode != 1 && mode != 2 {
		panic("Referer anti_leech_mode must be in [0, 1, 2]")
	}
	r.Mode = mode
	return r
}

// SetEmptyReferer 设置是否允许空Referer访问
func (r *ReferAntiLeechConfig) SetEmptyReferer(enable bool) *ReferAntiLeechConfig {
	r.AllowEmptyReferer = enable
	return r
}

// SetPattern 设置匹配Referer的模式
func (r *ReferAntiLeechConfig) SetPattern(pattern string) *ReferAntiLeechConfig {
	if pattern == "" {
		panic("Empty pattern is not allowed")
	}

	r.Pattern = pattern
	return r
}

// SetEnableSource 设置是否开启源站的防盗链
func (r *ReferAntiLeechConfig) SetEnableSource(enable bool) *ReferAntiLeechConfig {
	r.EnableSource = enable
	return r
}

// ImageSources 返回多个镜像回源地址的列表
func (b *BucketInfo) ImageSources() (srcs []string) {
	srcs = strings.Split(b.Source, ";")
	return
}

// IndexPageOn 返回空间是否开启了根目录下的index.html
func (b *BucketInfo) IndexPageOn() bool {
	return b.NoIndexPage == false
}

// Separators 返回分隔符列表
func (b *BucketInfo) Separators() (ret []rune) {
	for _, r := range b.Separator {
		ret = append(ret, r)
	}
	return
}

// WhiteListSet 是否设置了防盗链白名单
func (b *BucketInfo) WhiteListSet() bool {
	return b.AntiLeechMode == 1
}

// BlackListSet 是否设置了防盗链黑名单
func (b *BucketInfo) BlackListSet() bool {
	return b.AntiLeechMode == 2
}

// TokenAntiLeechModeOn 返回是否使用token签名防盗链开启了
func (b *BucketInfo) TokenAntiLeechModeOn() bool {
	return b.TokenAntiLeechMode == 1
}

// GetBucketInfo 返回BucketInfo结构
func GetBucketInfo(bucketName string) (bucketInfo BucketInfo, err error) {

}

// BucketInfosForRegion 获取指定区域的该用户的所有bucketInfo信息
func BucketInfosInRegion(region RegionID) (bucketInfos []BucketInfo, err error) {

}

// SetReferAntiLeechMode 配置存储空间referer防盗链模式
func SetReferAntiLeechMode(bucketName string, refererAntiLeechConfig *ReferAntiLeechConfig) (err error) {

}

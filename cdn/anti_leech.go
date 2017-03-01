package cdn

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"time"
)

// CreateTimestampAntiLeech
// 构建带时间戳防盗链的链接
// host需要加上 "http://" 或 "https://"
// encryptKey 七牛防盗链key
func CreateTimestampAntiLeechUrl(host, fileName string, queryStr url.Values, encryptKey string, durationInSeconds int64) (string, error) {

	var urlStr string
	if queryStr != nil {
		urlStr = fmt.Sprintf("%s/%s?%s", host, fileName, queryStr.Encode())
	} else {
		urlStr = fmt.Sprintf("%s/%s", host, fileName)
	}

	u, parseErr := url.Parse(urlStr)
	if parseErr != nil {
		err = parseErr
		return
	}
	return createTimestampAntiLeechUrl(u, encryptKey, durationInSeconds), nil

}

func createTimestampAntiLeechUrl(u *url.URL, encryptKey string, duration int64) string {

	expireTime := time.Now().Add(time.Second * time.Duration(duration)).Unix()
	toSignStr := fmt.Sprintf("%s%s%x", encryptKey, u.EscapedPath(), expireTime)
	signedStr := fmt.Sprintf("%x", md5.Sum([]byte(toSignStr)))

	q := u.Query()
	q.Add("sign", signedStr)
	q.Add("t", fmt.Sprintf("%x", expireTime))
	u.RawQuery = q.Encode()

	return fmt.Sprintf("%s://%s%s?%s", u.Scheme, u.Host, u.Path, u.Query().Encode())

}

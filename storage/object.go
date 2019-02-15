package storage

import (
	"context"
	"fmt"
	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/client"
	"github.com/qiniu/api.v7/conf"
	"net/http"
)

// ObjectManager 管理七牛存储空间中的对象
// 封装BucketManager的主要原因是继承一部分原来在BucketManager中实现了的对象管理的功能，比如Stat获取文件信息
// 该类型应该被用来对象管理， 其他的和存储空间相关的操作，比如创建，删除存储空间应该使用BucketManager
type ObjectManager struct {
	*BucketManager
}

// NewObjectManager 返回一个ObjectManager指针
func NewObjectManager(accessKey, secretKey string, cfg *Config) *ObjectManager {
	at := auth.New(accessKey, secretKey)

	return &ObjectManager{
		BucketManager: NewBucketManager(at, cfg),
	}
}

// NewObjectManagerEx 返回一个ObjectManager指针
func NewObjectManagerEx(accessKey, secretKey string, cfg *Config, clt *client.Client) *ObjectManager {
	at := auth.New(accessKey, secretKey)

	return &ObjectManager{
		BucketManager: NewBucketManagerEx(at, cfg, clt),
	}
}

// UpdateObjectStatus 用来修改文件状态, 禁用和启用文件的可访问性

// 请求包：
//
// POST /chstatus/<EncodedEntry>/status/<status>
// status：0表示启用，1表示禁用
// 返回包(JSON)：
//
// 200 OK
// 当<EncodedEntryURI>解析失败，返回400 Bad Request {"error":"invalid argument"}
// 当<EncodedEntryURI>不符合UTF-8编码，返回400 Bad Request {"error":"key must be utf8 encoding"}
// 当文件不存在时，返回612 status code 612 {"error":"no such file or directory"}
// 当文件当前状态和设置的状态已经一致，返回400 {"error":"already enabled"}或400 {"error":"already disabled"}
func (o *ObjectManager) UpdateObjectStatus(bucketName string, key string, enable bool) error {
	var status string
	ee := EncodedEntry(bucketName, key)
	if enable {
		status = "0"
	} else {
		status = "1"
	}
	path := fmt.Sprintf("/chstatus/%s/status/%s", ee, status)

	ctx := context.WithValue(context.TODO(), "mac", o.Mac)
	reqHost, reqErr := o.RsReqHost(bucketName)
	if reqErr != nil {
		return reqErr
	}
	reqURL := fmt.Sprintf("%s%s", reqHost, path)
	headers := http.Header{}
	headers.Add("Content-Type", conf.CONTENT_TYPE_FORM)
	return o.Client.Call(ctx, nil, "POST", reqURL, headers)
}

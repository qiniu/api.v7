package storage

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/x/rpc.v7"
	"net/url"
	"strconv"
	"strings"
)

const (
	DefaultRsHost  = "rs.qiniu.com"
	DefaultRsfHost = "rsf.qiniu.com"
	DefaultApiHost = "api.qiniu.com"
	DefaultPubHost = "pu.qbox.me:10200"
)

type FileInfo struct {
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	Type     int    `json:"type"`
}

func (f *FileInfo) String() string {
	str := ""
	str += fmt.Sprintf("Hash:     %s\n", f.Hash)
	str += fmt.Sprintf("Fsize:    %d\n", f.Fsize)
	str += fmt.Sprintf("PutTime:  %d\n", f.PutTime)
	str += fmt.Sprintf("MimeType: %s\n", f.MimeType)
	str += fmt.Sprintf("Type:     %d\n", f.Type)
	return str
}

type FetchRet struct {
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	MimeType string `json:"mimeType"`
	Key      string `json:"key"`
}

func (r *FetchRet) String() string {
	str := ""
	str += fmt.Sprintf("Key:      %s\n", r.Key)
	str += fmt.Sprintf("Hash:     %s\n", r.Hash)
	str += fmt.Sprintf("Fsize:    %d\n", r.Fsize)
	str += fmt.Sprintf("MimeType: %s\n", r.MimeType)
	return str
}

type ListItem struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	Type     int    `json:"type"`
	EndUser  string `json:"endUser"`
}

func (l *ListItem) String() string {
	str := ""
	str += fmt.Sprintf("Hash:     %s\n", l.Hash)
	str += fmt.Sprintf("Fsize:    %d\n", l.Fsize)
	str += fmt.Sprintf("PutTime:  %d\n", l.PutTime)
	str += fmt.Sprintf("MimeType: %s\n", l.MimeType)
	str += fmt.Sprintf("Type:     %d\n", l.Type)
	str += fmt.Sprintf("EndUser:  %s\n", l.EndUser)
	return str
}

type BucketManager struct {
	client *rpc.Client
	mac    *qbox.Mac
	cfg    *Config
}

func NewBucketManager(mac *qbox.Mac, cfg *Config) *BucketManager {
	if cfg == nil {
		cfg = &Config{}
	}

	return &BucketManager{
		client: NewClient(mac, nil),
		mac:    mac,
		cfg:    cfg,
	}
}

// 获取空间列表
// @param shared - 是否同时列出被授权访问的bucket
func (m *BucketManager) Buckets(shared bool) (buckets []string, err error) {
	ctx := context.TODO()
	var reqHost string
	if m.cfg.UseHttps {
		reqHost = fmt.Sprintf("https://%s", DefaultRsHost)
	} else {
		reqHost = fmt.Sprintf("http://%s", DefaultRsHost)
	}

	reqUrl := fmt.Sprintf("%s/buckets?shared=%v", reqHost, shared)
	err = m.client.Call(ctx, &buckets, "POST", reqUrl)
	return
}

// 获取文件信息
func (m *BucketManager) Stat(bucket, key string) (info FileInfo, err error) {
	ctx := context.TODO()
	reqHost, reqErr := m.RsHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}

	reqUrl := fmt.Sprintf("%s%s", reqHost, URIStat(bucket, key))
	err = m.client.Call(ctx, &info, "POST", reqUrl)
	return
}

// 删除文件
func (m *BucketManager) Delete(bucket, key string) (err error) {
	ctx := context.TODO()
	reqHost, reqErr := m.RsHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}
	reqUrl := fmt.Sprintf("%s%s", reqHost, URIDelete(bucket, key))
	err = m.client.Call(ctx, nil, "POST", reqUrl)
	return
}

// 复制文件
func (m *BucketManager) Copy(srcBucket, srcKey, destBucket, destKey string, force bool) (err error) {
	ctx := context.TODO()
	reqHost, reqErr := m.RsHost(srcBucket)
	if reqErr != nil {
		err = reqErr
		return
	}

	reqUrl := fmt.Sprintf("%s%s", reqHost, URICopy(srcBucket, srcKey, destBucket, destKey, force))
	err = m.client.Call(ctx, nil, "POST", reqUrl)
	return
}

// 移动文件
func (m *BucketManager) Move(srcBucket, srcKey, destBucket, destKey string, force bool) (err error) {
	ctx := context.TODO()
	reqHost, reqErr := m.RsHost(srcBucket)
	if reqErr != nil {
		err = reqErr
		return
	}

	reqUrl := fmt.Sprintf("%s%s", reqHost, URIMove(srcBucket, srcKey, destBucket, destKey, force))
	err = m.client.Call(ctx, nil, "POST", reqUrl)
	return
}

// 更新文件的格式类型
func (m *BucketManager) ChangeMime(bucket, key, newMime string) (err error) {
	ctx := context.TODO()
	reqHost, reqErr := m.RsHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}
	reqUrl := fmt.Sprintf("%s%s", reqHost, URIChangeMime(bucket, key, newMime))
	err = m.client.Call(ctx, nil, "POST", reqUrl)
	return
}

// 更新文件存储类型
func (m *BucketManager) ChangeType(bucket, key string, fileType int) (err error) {
	ctx := context.TODO()
	reqHost, reqErr := m.RsHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}
	reqUrl := fmt.Sprintf("%s%s", reqHost, URIChangeType(bucket, key, fileType))
	err = m.client.Call(ctx, nil, "POST", reqUrl)
	return
}

// 更新文件生命周期
func (m *BucketManager) DeleteAfterDays(bucket, key string, days int) (err error) {
	ctx := context.TODO()
	reqHost, reqErr := m.RsHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}

	reqUrl := fmt.Sprintf("%s%s", reqHost, URIDeleteAfterDays(bucket, key, days))
	err = m.client.Call(ctx, nil, "POST", reqUrl)
	return
}

// 抓取资源
func (m *BucketManager) Fetch(resUrl, bucket, key string) (fetchRet FetchRet, err error) {
	ctx := context.TODO()
	reqHost, reqErr := m.IovipHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}
	reqUrl := fmt.Sprintf("%s%s", reqHost, uriFetch(resUrl, bucket, key))
	err = m.client.Call(ctx, &fetchRet, "POST", reqUrl)
	return
}

// 抓取资源，如果不指定key，则以文件的内容hash作为文件名
func (m *BucketManager) FetchWithoutKey(resUrl, bucket string) (fetchRet FetchRet, err error) {
	ctx := context.TODO()
	reqHost, reqErr := m.IovipHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}
	reqUrl := fmt.Sprintf("%s%s", reqHost, uriFetchWithoutKey(resUrl, bucket))
	err = m.client.Call(ctx, &fetchRet, "POST", reqUrl)
	return
}

// 同步镜像空间的资源和镜像源资源内容
func (m *BucketManager) Prefetch(bucket, key string) (err error) {
	ctx := context.TODO()
	reqHost, reqErr := m.IovipHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}
	reqUrl := fmt.Sprintf("%s%s", reqHost, uriPrefetch(bucket, key))
	err = m.client.Call(ctx, nil, "POST", reqUrl)
	return
}

// 设置空间镜像源
func (m *BucketManager) SetImage(siteUrl, bucket string) (err error) {
	ctx := context.TODO()
	reqUrl := fmt.Sprintf("http://%s%s", DefaultPubHost, uriSetImage(siteUrl, bucket))
	err = m.client.Call(ctx, nil, "POST", reqUrl)
	return
}

// 设置空间镜像源，额外添加回源Host头部
func (m *BucketManager) SetImageWithHost(siteUrl, bucket, host string) (err error) {
	ctx := context.TODO()
	reqUrl := fmt.Sprintf("http://%s%s", DefaultPubHost,
		uriSetImageWithHost(siteUrl, bucket, host))
	err = m.client.Call(ctx, nil, "POST", reqUrl)
	return
}

// 取消空间镜像源设置
func (m *BucketManager) UnsetImage(bucket string) (err error) {
	ctx := context.TODO()
	reqUrl := fmt.Sprintf("http://%s%s", DefaultPubHost, uriUnsetImage(bucket))
	err = m.client.Call(ctx, nil, "POST", reqUrl)
	return err
}

type listFilesRet struct {
	Marker         string     `json:"marker"`
	Items          []ListItem `json:"items"`
	CommonPrefixes []string   `json:"commonPrefixes"`
}

// 获取空间文件列表
// @param bucket
// @param prefix
// @param delimiter
// @param marker
// @param limit
func (m *BucketManager) ListFiles(bucket, prefix, delimiter, marker string,
	limit int) (entries []ListItem, commonPrefixes []string, nextMarker string, hasNext bool, err error) {
	ctx := context.TODO()
	reqHost, reqErr := m.RsfHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}

	ret := listFilesRet{}
	reqUrl := fmt.Sprintf("%s%s", reqHost, uriListFiles(bucket, prefix, delimiter, marker, limit))
	err = m.client.Call(ctx, &ret, "POST", reqUrl)
	if err != nil {
		return
	}

	commonPrefixes = ret.CommonPrefixes
	nextMarker = ret.Marker
	entries = ret.Items
	if ret.Marker != "" {
		hasNext = true
	}

	return
}

// 获取资源管理域名
func (m *BucketManager) RsHost(bucket string) (rsHost string, err error) {
	zone, zoneErr := GetZone(m.mac.AccessKey, bucket)
	if zoneErr != nil {
		err = zoneErr
		return
	}

	if m.cfg.UseHttps {
		rsHost = fmt.Sprintf("https://%s", zone.RsHost)
	} else {
		rsHost = fmt.Sprintf("http://%s", zone.RsHost)
	}
	return
}

func (m *BucketManager) RsfHost(bucket string) (rsfHost string, err error) {
	zone, zoneErr := GetZone(m.mac.AccessKey, bucket)
	if zoneErr != nil {
		err = zoneErr
		return
	}

	if m.cfg.UseHttps {
		rsfHost = fmt.Sprintf("https://%s", zone.RsfHost)
	} else {
		rsfHost = fmt.Sprintf("http://%s", zone.RsfHost)
	}
	return
}

// 获取IOVIP域名
func (m *BucketManager) IovipHost(bucket string) (iovipHost string, err error) {
	zone, zoneErr := GetZone(m.mac.AccessKey, bucket)
	if zoneErr != nil {
		err = zoneErr
		return
	}

	if m.cfg.UseHttps {
		iovipHost = fmt.Sprintf("https://%s", zone.IovipHost)
	} else {
		iovipHost = fmt.Sprintf("http://%s", zone.IovipHost)
	}
	return
}

// 构建op的方法，导出的方法支持在Batch操作中使用
func URIStat(bucket, key string) string {
	return fmt.Sprintf("/stat/%s", EncodedEntry(bucket, key))
}

func URIDelete(bucket, key string) string {
	return fmt.Sprintf("/delete/%s", EncodedEntry(bucket, key))
}

func URICopy(srcBucket, srcKey, destBucket, destKey string, force bool) string {
	return fmt.Sprintf("/copy/%s/%s/force/%v", EncodedEntry(srcBucket, srcKey),
		EncodedEntry(destBucket, destKey), force)
}

func URIMove(srcBucket, srcKey, destBucket, destKey string, force bool) string {
	return fmt.Sprintf("/move/%s/%s/force/%v", EncodedEntry(srcBucket, srcKey),
		EncodedEntry(destBucket, destKey), force)
}

func URIDeleteAfterDays(bucket, key string, days int) string {
	return fmt.Sprintf("/deleteAfterDays/%s/%d", EncodedEntry(bucket, key), days)
}

func URIChangeMime(bucket, key, newMime string) string {
	return fmt.Sprintf("/chgm/%s/mime/%s", EncodedEntry(bucket, key),
		base64.URLEncoding.EncodeToString([]byte(newMime)))
}

func URIChangeType(bucket, key string, fileType int) string {
	return fmt.Sprintf("/chtype/%s/type/%d", EncodedEntry(bucket, key), fileType)
}

// 构建op的方法，非导出的方法无法用在Batch操作中
func uriFetch(resUrl, bucket, key string) string {
	return fmt.Sprintf("/fetch/%s/to/%s",
		base64.URLEncoding.EncodeToString([]byte(resUrl)), EncodedEntry(bucket, key))
}

func uriFetchWithoutKey(resUrl, bucket string) string {
	return fmt.Sprintf("/fetch/%s/to/%s",
		base64.URLEncoding.EncodeToString([]byte(resUrl)), EncodedEntryWithoutKey(bucket))
}

func uriPrefetch(bucket, key string) string {
	return fmt.Sprintf("/prefetch/%s", EncodedEntry(bucket, key))
}

func uriSetImage(siteUrl, bucket string) string {
	return fmt.Sprintf("/image/%s/from/%s", bucket,
		base64.URLEncoding.EncodeToString([]byte(siteUrl)))
}

func uriSetImageWithHost(siteUrl, bucket, host string) string {
	return fmt.Sprintf("/image/%s/from/%s/host/%s", bucket,
		base64.URLEncoding.EncodeToString([]byte(siteUrl)),
		base64.URLEncoding.EncodeToString([]byte(host)))
}

func uriUnsetImage(bucket string) string {
	return fmt.Sprintf("/unimage/%s", bucket)
}

func uriListFiles(bucket, prefix, delimiter, marker string, limit int) string {
	query := make(url.Values)
	query.Add("bucket", bucket)
	if prefix != "" {
		query.Add("prefix", prefix)
	}
	if delimiter != "" {
		query.Add("delimiter", delimiter)
	}
	if marker != "" {
		query.Add("marker", marker)
	}
	if limit > 0 {
		query.Add("limit", strconv.FormatInt(int64(limit), 10))
	}
	return fmt.Sprintf("/list?%s", query.Encode())
}

// EncodedEntry
func EncodedEntry(bucket, key string) string {
	entry := fmt.Sprintf("%s:%s", bucket, key)
	return base64.URLEncoding.EncodeToString([]byte(entry))
}

func EncodedEntryWithoutKey(bucket string) string {
	return base64.URLEncoding.EncodeToString([]byte(bucket))
}

// 公开空间资源下载链接
// @param domain - 下载域名，例如 http://img.example.com
// @param key    - 下载文件名，例如 img/2017/test.png
func MakePublicUrl(domain, key string) (publicUrl string) {
	return fmt.Sprintf("%s/%s", domain, url.QueryEscape(key))
}

// 私有空间资源下载链接
// @param domain   - 下载域名，例如 http://img.example.com
// @param key      - 下载文件名，例如 img/2017/test.png
// @param deadline - 链接过期Unix时间戳
func MakePrivateUrl(mac *qbox.Mac, domain, key string, deadline int64) (privateUrl string) {
	publicUrl := MakePublicUrl(domain, key)
	urlToSign := publicUrl
	if strings.Contains(publicUrl, "?") {
		urlToSign = fmt.Sprintf("%s&e=%d", urlToSign, deadline)
	} else {
		urlToSign = fmt.Sprintf("%s?e=%d", urlToSign, deadline)
	}
	token := mac.Sign([]byte(urlToSign))
	privateUrl = fmt.Sprintf("%s&token=%s", urlToSign, token)
	return
}

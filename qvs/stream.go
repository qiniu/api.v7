package qvs

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Stream struct {
	StreamID string `json:"streamId"`       // 流名称, 流名称在空间中唯一，可包含 字母、数字、中划线、下划线；1 ~ 100 个字符长；创建后将不可修改
	Desc     string `json:"desc,omitempty"` // 关于流的描述信息

	NamespaceId string `json:"nsId"`   // 所属的空间ID
	Namespace   string `json:"nsName"` // 所属的空间名称

	RecordTemplateId   string `json:"recordTemplateId"`   // 录制模版ID，配置流维度的录制模板
	SnapShotTemplateId string `json:"snapshotTemplateId"` // 截图模版ID，配置流维度的截图模板

	Status       bool  `json:"status"`       // 设备是否在线
	Disabled     bool  `json:"disabled"`     // 流是否被禁用
	LastPushedAt int64 `json:"lastPushedAt"` // 最后一次推流时间,0:表示没有推流

	CreatedAt int64 `json:"createdAt,omitempty"` // 流创建时间
	UpdatedAt int64 `json:"updatedAt,omitempty"` // 流更新时间
}

/*
	创建流API
*/
func (manager *Manager) AddStream(nsId string, stream *Stream) (*Stream, error) {

	var ret Stream
	err := manager.client.CallWithJson(context.Background(), &ret, "POST", manager.url("/namespaces/%s/streams", nsId), nil, stream)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
	查询流API
*/
func (manager *Manager) QueryStream(nsId string, streamId string) (*Stream, error) {

	var ret Stream
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/streams/%s", nsId, streamId), nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
	更新流API
*/
func (manager *Manager) UpdateStream(nsId string, streamId string, ops []PatchOperation) (*Stream, error) {

	req := M{"operations": ops}
	var ret Stream
	err := manager.client.CallWithJson(context.Background(), &ret, "PATCH", manager.url("/namespaces/%s/streams/%s", nsId, streamId), nil, req)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
	删除流API
*/
func (manager *Manager) DeleteStream(nsId string, streamId string) error {

	return manager.client.Call(context.Background(), nil, "DELETE", manager.url("/namespaces/%s/streams/%s", nsId, streamId), nil)
}

/*
	查询流列表API
*/
func (manager *Manager) ListStream(nsId string, offset, line int, prefix, sortBy string, qType int) ([]Stream, int64, error) {

	ret := struct {
		Items []Stream `json:"items"`
		Total int64    `json:"total"`
	}{}

	query := url.Values{}
	setQuery(query, "offset", offset)
	setQuery(query, "line", line)
	setQuery(query, "sortBy", sortBy)
	setQuery(query, "prefix", prefix)
	setQuery(query, "qtype", qType)

	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/streams?%v", nsId, query.Encode()), nil)
	return ret.Items, ret.Total, err
}

type DynamicLiveRoute struct {
	PublishIP    string `json:"publishIP"`    // 推流端对外IP地址
	PlayIP       string `json:"playIP"`       // 拉流端对外IP地址
	UrlExpireSec int64  `json:"urlExpireSec"` // 地址过期时间,urlExpireSec:100代表100秒后过期;  默认urlExpireSec:0,永不过期.
}

type RouteRet struct {
	PublishUrl        string        `json:"publishUrl"`        // rtmp推流地址
	PlayUrls          RoutePlayUrls `json:"playUrls"`          // 拉流URLs
	PublishUrlExpired int64         `json:"publishUrlExpired"` // 推拉流地址过期时间点(unix时间戳,单位second)
}

type RoutePlayUrls struct {
	Rtmp string `json:"rtmp"` // rtmp播放地址
	Flv  string `json:"flv"`  // flv播放地址
	Hls  string `json:"hls"`  // hls播放地址
}

/*
	获取流地址API：推拉流IP地址计算最合适的设备端推拉流地址
*/
func (manager *Manager) DynamicPublishPlayURL(nsId string, streamId string, route *DynamicLiveRoute) (*RouteRet, error) {

	var ret RouteRet
	err := manager.client.CallWithJson(context.Background(), &ret, "POST", manager.url("/namespaces/%s/streams/%s/urls", nsId, streamId), nil, route)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
	获取流地址API：
	生成推拉流地址-domain urlExpireSec: urlExpireSec:100代表100秒后过期;  默认urlExpireSec:0,永不过期.  authEnable 0:开启鉴权 1:关闭鉴权
 	domain 如推流地址为qvs-publish.test.com, 拉流地址为qvs-live-rtmp.test.com, qvs-live-hls.test.com, qvs-live-hdl.test.com, domain传入test.com
*/
func (manager *Manager) StaticPublishPlayURL(domain, nsId, streamId string, urlExpireSec int64, authEnable uint8) (*RouteRet, error) {

	expire := time.Now().Unix() + urlExpireSec
	path1 := fmt.Sprintf("/%s/%s?e=%d", nsId, streamId, expire)
	token := manager.mac.Sign([]byte(path1))
	path2 := fmt.Sprintf("/%s/%s", nsId, streamId)
	return &RouteRet{
		PublishUrl: fmt.Sprintf("rtmp://qvs-publish.%s%s&token=%s", domain, path1, token),
		PlayUrls: RoutePlayUrls{
			Rtmp: fmt.Sprintf("rtmp://qvs-live-rtmp.%s%s", domain, path2),
			Hls:  fmt.Sprintf("http://qvs-live-hls.%s%s", domain, path2),
			Flv:  fmt.Sprintf("http://qvs-live-hdl.%s%s", domain, path2),
		},
	}, nil
}

// 禁用流
func (manager *Manager) DisableStream(nsId string, streamId string) error {

	err := manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/streams/%s/disabled", nsId, streamId), nil, nil)
	return err
}

// 恢复流
func (manager *Manager) EnableStream(nsId string, streamId string) error {

	err := manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/streams/%s/enabled", nsId, streamId), nil, nil)
	return err
}

type StreamPublishHistory struct {
	StreamId string `json:"streamId"`
	NsId     string `json:"nsId"`
	Start    int64  `json:"start"`
	End      int64  `json:"end"`
	Duration int64  `json:"duration"`
	Snap     string `json:"snap"`
	Fmt      string `json:"format,omitempty"`
}

// 查询推流历史记录
func (manager *Manager) QueryStreamPubhistories(nsId string, streamId string, start, end int, line, offset int) ([]StreamPublishHistory, int64, error) {

	ret := struct {
		Items []StreamPublishHistory `json:"items"`
		Total int64                  `json:"total"`
	}{}

	query := url.Values{}
	setQuery(query, "start", start)
	setQuery(query, "end", end)
	setQuery(query, "offset", offset)
	setQuery(query, "line", line)

	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/streams/%s/pubhistories?%v", nsId, streamId, query.Encode()), nil)
	return ret.Items, ret.Total, err
}

// 查询截图列表
func (manager *Manager) StreamsSnapshots(nsId string, streamId string, start, end int, qtype int, limit int, marker string) ([]byte, error) {

	query := url.Values{}
	setQuery(query, "start", start)
	setQuery(query, "end", end)
	setQuery(query, "type", qtype)
	setQuery(query, "limit", limit)
	setQuery(query, "marker", marker)

	fmt.Println(manager.url("/namespaces/%s/streams/%s/snapshots?%v", nsId, streamId, query.Encode()))
	req, err := http.NewRequest("GET", manager.url("/namespaces/%s/streams/%s/snapshots?%v", nsId, streamId, query.Encode()), nil)
	if err != nil {
		return nil, err
	}
	resp, err := manager.client.Do(context.Background(), req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

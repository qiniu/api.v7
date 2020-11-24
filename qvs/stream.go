package qvs

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	DomainPublishRTMP string = "publishRtmp"
	DomainLiveRTMP    string = "liveRtmp"
	DomainLiveHLS     string = "liveHls"
	DomainLiveHDL     string = "liveHdl"
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

	// 以下字段只有在设备在线是才会出现
	UserCount      int    `json:"userCount"`
	ClientIp       string `json:"clientIp,omitempty"`
	AudioFrameRate int64  `json:"audioFrameRate,omitempty"`
	BitRate        int64  `json:"bitRate,omitempty"`
	VideoFrameRate int64  `json:"videoFrameRate,omitempty"`
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

type StaticLiveRoute struct {
	Domain       string `json:"domain"`       // 域名
	DomainType   string `json:"domainType"`   // 域名类型
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
	动态获取流地址API：推拉流IP地址计算最合适的设备端推拉流地址
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
	静态获取流地址API：根据domain生成推拉流地址
*/
func (manager *Manager) StaticPublishPlayURL(nsId, streamId string, route *StaticLiveRoute) (string, error) {

	var ret RouteRet
	err := manager.client.CallWithJson(context.Background(), &ret, "POST", manager.url("/namespaces/%s/streams/%s/domain", nsId, streamId), nil, route)
	if err != nil {
		return "", err
	}
	switch route.DomainType {
	case DomainPublishRTMP:
		return ret.PublishUrl, nil
	case DomainLiveRTMP:
		return ret.PlayUrls.Rtmp, nil
	case DomainLiveHLS:
		return ret.PlayUrls.Hls, nil
	case DomainLiveHDL:
		return ret.PlayUrls.Flv, nil
	}
	return "", nil
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

// 停止推流
func (manager *Manager) StopStream(nsId string, streamId string) error {

	err := manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/streams/%s/stop", nsId, streamId), nil, nil)
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

// 按需截图
func (manager *Manager) OndemandSnap(nsId, streamId string) error {
	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/streams/%s/snap", nsId, streamId), nil, nil)
}

// 删除截图
func (manager *Manager) DeleteSnapshots(nsId, streamId string, files []string) error {
	return manager.client.CallWithJson(context.Background(), nil, "DELETE", manager.url("/namespaces/%s/streams/%s/snapshots", nsId, streamId), nil, M{"files": files})
}

// 查询截图列表
func (manager *Manager) StreamsSnapshots(nsId string, streamId string, start, end int, qtype int, line int, marker string) ([]byte, error) {
	query := url.Values{}
	setQuery(query, "start", start)
	setQuery(query, "end", end)
	setQuery(query, "type", qtype)
	setQuery(query, "line", line)
	setQuery(query, "marker", marker)

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

type RecordHistory struct {
	Url      string `json:"url"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Duration int    `json:"duration"`
	Format   int    `json:"format"`
	Snap     string `json:"snap"`
	File     string `json:"file"`
}

// 查询视频流的录制记录
func (manager *Manager) QueryStreamRecordHistories(nsId string, streamId string, start, end int, marker string, line int, format string) ([]RecordHistory, string, error) {
	ret := struct {
		Items  []RecordHistory `json:"items"`
		Marker string          `json:"marker"`
	}{}

	query := url.Values{}
	setQuery(query, "start", start)
	setQuery(query, "end", end)
	setQuery(query, "marker", marker)
	setQuery(query, "line", line)
	setQuery(query, "format", format)

	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/streams/%s/recordhistories?%v", nsId, streamId, query.Encode()), nil)
	return ret.Items, ret.Marker, err
}

// 查询流封面
func (manager *Manager) QueryStreamCover(nsId string, streamId string) (string, error) {
	var ret = struct {
		Url string `json:"url"`
	}{}
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/streams/%s/cover", nsId, streamId), nil)
	return ret.Url, err
}

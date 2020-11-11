package qvs

import (
	"context"
	"net/url"
)

type Device struct {
	NamespaceId     string `json:"nsId"`
	Name            string `json:"name"`
	GBId            string `json:"gbId"`
	Type            int    `json:"type"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	PullIfRegister  bool   `json:"pullIfRegister"` //按需拉流
	Desc            string `json:"desc"`
	NamespaceName   string `json:"nsName"`
	State           string `json:"state"`
	Channels        int    `json:"channels"`
	Vendor          string `json:"vendor"`
	CreatedAt       int64  `json:"createdAt"`
	UpdatedAt       int64  `json:"updatedAt"`
	LastRegisterAt  int64  `json:"lastRegisterAt"`
	LastKeepaliveAt int64  `json:"lastKeepaliveAt"`
}

type QueryChannelsArgs struct {
	CmdArgs []string
	Prefix  string `json:"prefix"`
}

type Channel struct {
	GBID       string `json:"gbId"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Vendor     string `json:"vendor"`
	LastSyncAt int64  `json:"lastSyncAt"`
}

type DeviceChannels struct {
	OnlineCount  int       `json:"onlineCount"`
	OfflineCount int       `json:"offlineCount"`
	Total        int       `json:"total"`
	Items        []Channel `json:"items"`
}

type DeviceVideoItems struct {
	Items []deviceVideoItem `json:"items"`
}

type deviceVideoItem struct {
	PlayUrls RoutePlayUrls `json:"playUrls"`
	Start    int           `json:"start"`
	End      int           `json:"end"`
	Type     string        `json:"type"`
}

/*
   创建设备API
   参数device需要赋值字段:
       NamespaceId 必填
       Name
       GBId
       Username 必填
       Password 必填
       PullIfRegister
       Desc
       Type
*/
func (manager *Manager) AddDevice(device *Device) (*Device, error) {
	var ret Device
	err := manager.client.CallWithJson(context.Background(), &ret, "POST", manager.url("/namespaces/%s/devices", device.NamespaceId), nil, device)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
   删除设备API
*/
func (manager *Manager) DeleteDevice(nsId string, gbId string) error {
	return manager.client.Call(context.Background(), nil, "DELETE", manager.url("/namespaces/%s/devices/%s", nsId, gbId), nil)
}

/*
   查询设备API
*/
func (manager *Manager) QueryDevice(nsId string, gbId string) (*Device, error) {
	var ret Device
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/devices/%s", nsId, gbId), nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
   查询设备列表API
*/
func (manager *Manager) ListDevice(nsId string, offset, line int, prefix, state string, qType int) ([]Device, int64, error) {
	ret := struct {
		Items []Device `json:"items"`
		Total int64    `json:"total"`
	}{}

	query := url.Values{}
	setQuery(query, "offset", offset)
	setQuery(query, "line", line)
	setQuery(query, "state", state)
	setQuery(query, "prefix", prefix)
	setQuery(query, "qtype", qType)

	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/devices?%v", nsId, query.Encode()), nil)
	return ret.Items, ret.Total, err
}

/*
   更新设备API
*/
func (manager *Manager) UpdateDevice(nsId string, gbId string, ops []PatchOperation) (*Device, error) {
	req := M{"operations": ops}
	var ret Device
	err := manager.client.CallWithJson(context.Background(), &ret, "PATCH", manager.url("/namespaces/%s/devices/%s", nsId, gbId), nil, req)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
   启动设备拉流API
*/
func (manager *Manager) StartDevice(nsId string, gbId string, channels []string) error {
	req := M{"channels": channels}
	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/devices/%s/start", nsId, gbId), nil, req)
}

/*
   停止设备拉流API
*/
func (manager *Manager) StopDevice(nsId string, gbId string, channels []string) error {
	req := M{"channels": channels}
	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/devices/%s/stop", nsId, gbId), nil, req)
}

/*
   查询通道列表
*/
func (manager *Manager) ListChannels(nsId string, gbId string, prefix string) (*DeviceChannels, error) {
	var ret DeviceChannels
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/devices/%s/channels?prefix=%s", nsId, gbId, prefix), nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
   同步设备通道
*/
func (manager *Manager) FetchCatalog(nsId, gbId string) error {
	err := manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/devices/%s/catalog/fetch", nsId, gbId), nil, nil)
	if err != nil {
		return err
	}
	return nil
}

/*
   查询通道详情
*/
func (manager *Manager) QueryChannel(nsId, gbId, channelId string) (*Channel, error) {
	var ret Channel
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/devices/%s/channels/%s", nsId, gbId, channelId), nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
   删除通道
*/
func (manager *Manager) DeleteChannel(nsId, gbId, channelId string) error {
	err := manager.client.Call(context.Background(), nil, "DELETE", manager.url("/namespaces/%s/devices/%s/channels/%s", nsId, gbId, channelId), nil)
	if err != nil {
		return err
	}
	return nil
}

/*
   查询本地录像列表
   普通设备chId可以忽略, 置为空字符串即可
*/
func (manager *Manager) QueryGBRecordHistories(nsId, gbId, chId string, start, end int) (*DeviceVideoItems, error) {
	var ret DeviceVideoItems
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/devices/%s/recordhistories?start=%d&end=%d&chId=%s", nsId, gbId, start, end, chId), nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

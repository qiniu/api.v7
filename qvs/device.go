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

type Preset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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
	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/devices/%s/catalog/fetch", nsId, gbId), nil, nil)
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
	return manager.client.Call(context.Background(), nil, "DELETE", manager.url("/namespaces/%s/devices/%s/channels/%s", nsId, gbId, channelId), nil)
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

/*
    用于对摄像头进行 转动镜头，如水平、垂直、缩放等操作
	cmd: left(向左), right(向右), up(向上), down(向下), leftup(左上), rightup(右上), leftdown(左下), rightdown(右下), zoomin(放大), zoomout(缩小),stop(停止PTZ操作)
	speed: 调节速度(1~10, 默认位5)
*/
func (manager *Manager) DevicePtz(nsId, gbId, chId string, cmd string, speed int) error {
	req := M{"cmd": cmd, "speed": speed, "chId": chId}
	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/devices/%s/ptz", nsId, gbId), nil, req)
}

/*
	变焦控制
	cmd: focusnear(焦距变近), focusfar(焦距变远),stop(停止)
	speed: 调节速度(1~10, 默认位5)
*/
func (manager *Manager) DeviceFocus(nsId, gbId, chId string, cmd string, speed int) error {
	req := M{"cmd": cmd, "speed": speed, "chId": chId}
	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/devices/%s/focus", nsId, gbId), nil, req)
}

/*
	用于摄像头的光圈控制
	cmd: irisin(光圈变小), irisout(光圈变大),stop(停止)
	speed: 调节速度(1~10, 默认位5)
*/
func (manager *Manager) DeviceIris(nsId, gbId, chId string, cmd string, speed int) error {
	req := M{"cmd": cmd, "speed": speed, "chId": chId}
	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/devices/%s/iris", nsId, gbId), nil, req)
}

/*
	用于预置位的管理
	cmd: set(新增预置位), goto(设置),remove(删除)
	name: 预置位名称(cmd为set时有效,支持中文)
	presetId: 预置位ID
*/
func (manager *Manager) DevicePresets(nsId, gbId, chId string, cmd, name string, presetId int) (*Preset, error) {
	req := M{"cmd": cmd, "name": name, "presetId": presetId, "chId": chId}
	var ret Preset
	err := manager.client.CallWithJson(context.Background(), &ret, "POST", manager.url("/namespaces/%s/devices/%s/presets", nsId, gbId), nil, req)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
	用于获取预置位列表
	fetch: 是否强制从设备获取预置位
*/
func (manager *Manager) QueryDevicePresets(nsId, gbId, chId string, fetch bool) ([]Preset, error) {
	ret := struct {
		Items []Preset `json:"items"`
	}{}

	query := url.Values{}
	setQuery(query, "chId", chId)
	setQuery(query, "fetch", fetch)

	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/devices/%s/presets?%v", nsId, gbId, query.Encode()), nil)
	return ret.Items, err
}

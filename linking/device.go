package linking

import (
	"context"
	"encoding/base64"
	"net/url"
)

type M map[string]interface{}

type PatchOperation struct {
	Op    string      `json:"op"`    // 更该或删除某个属性，replace：更改，delete：删除
	Key   string      `json:"key"`   // 要修改或删除的属性
	Value interface{} `json:"value"` // 要修改或删除属性的值
}

type Channel struct {
	Channelid int    `json:"channelid"`
	Comment   string `json:"comment"`
}

type Device struct {
	Device   string `json:"device"`
	LoginAt  int64  `json:"loginAt,omitempty"`  // 查询条件 online 为 true 时才会出现该字段
	RemoteIp string `json:"remoteIp,omitempty"` // 查询条件 online 为 true 时才会出现该字段
	// 0 不录制
	// -1 永久
	// -2 继承app配置
	SegmentExpireDays int `json:"segmentExpireDays,has,omitempty"`

	// -1 继承app配置
	// 0 遵循设备端配置
	// 1 强制持续上传
	// 2 强制关闭上传
	UploadMode int `json:"uploadMode,has,omitempty"`

	State int `json:"state,omitempty"`

	ActivedAt int64 `json:"activedAt,omitempty"`
	CreatedAt int64 `json:"createdAt,omitempty"`
	UpdatedAt int64 `json:"updatedAt,omitempty"`

	// 1 免费使用
	// 0 正常收费
	LicenseMode int `json:"licenseMode,omitempty"`

	// meta data
	Meta []byte `json:"meta,omitempty"`

	// device type 0:normal type, 1:gateway
	Type int `json:"type"`
	// max channel of gateway [1,64]
	MaxChannel int       `json:"maxChannel,omitempty"`
	Channels   []Channel `json:"channels,omitempty"`
}

// 在指定的应用下添加新的设备
func (manager *Manager) AddDevice(appid string, dev *Device) (*Device, error) {

	var ret Device
	err := manager.client.CallWithJson(context.Background(), &ret, "POST", manager.url("/apps/%s/devices", appid), nil, dev)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// 查询指定设备的详细信息
func (manager *Manager) QueryDevice(appid, device string) (*Device, error) {

	device = base64.URLEncoding.EncodeToString([]byte(device))
	var ret Device
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/apps/%s/devices/%s", appid, device), nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// 更新设备配置信息的操作
func (manager *Manager) UpdateDevice(appid, device string, ops []PatchOperation) (*Device, error) {

	device = base64.URLEncoding.EncodeToString([]byte(device))
	req := M{"operations": ops}
	var ret Device
	err := manager.client.CallWithJson(context.Background(), &ret, "PATCH", manager.url("/apps/%s/devices/%s", appid, device), nil, req)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// 查询指定应用下所有设备的列表
func (manager *Manager) ListDevice(appid, prefix, marker string, limit int, online, status bool, deviceType int, batch string) ([]Device, string, error) {

	ret := struct {
		Items  []Device `json:"items"`
		Marker string   `json:"marker"`
	}{}

	query := url.Values{}
	setQuery(query, "online", online)
	setQuery(query, "status", status)
	if limit > 0 {
		setQuery(query, "limit", limit)
	}
	if prefix != "" {
		setQuery(query, "prefix", prefix)
	}
	if marker != "" {
		setQuery(query, "marker", marker)
	}
	setQuery(query, "type", deviceType)
	if batch != "" {
		setQuery(query, "batch", batch)
	}
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/apps/%s/devices?%v", appid, query.Encode()), nil)
	if err != nil {
		return nil, "", err
	}
	return ret.Items, ret.Marker, nil
}

// 删除指定的设备，删除后将不可恢复
func (manager *Manager) DeleteDevice(appid, device string) error {

	device = base64.URLEncoding.EncodeToString([]byte(device))
	return manager.client.Call(context.Background(), nil, "DELETE", manager.url("/apps/%s/devices/%s", appid, device), nil)
}

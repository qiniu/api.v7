package linking

import (
	"context"
	"encoding/base64"
)

// 设备密钥管理
type DeviceKey struct {
	AccessKey string `json:"accessKey"` // 设备的 accessKey
	SecretKey string `json:"secretKey"` // 设备的 secretkey
	State     int    `json:"state"`     // 密钥对状态，1表示被禁用,0表示已启用
	CreatedAt int64  `json:"createdAt"` // 创建时间
}

// 新增设备的密钥，每个设备最多有两对密钥
func (manager *Manager) AddDeviceKey(appid, device string) ([]DeviceKey, error) {

	device = base64.URLEncoding.EncodeToString([]byte(device))

	ret := struct {
		Keys []DeviceKey `json:"keys"`
	}{}
	err := manager.client.Call(context.Background(), &ret, "POST", manager.url("/apps/%s/devices/%s/keys", appid, device), nil)
	if err != nil {
		return nil, err
	}
	return ret.Keys, nil
}

// 查询指定设备的密钥
func (manager *Manager) QueryDeviceKey(appid, device string) ([]DeviceKey, error) {

	device = base64.URLEncoding.EncodeToString([]byte(device))

	ret := struct {
		Keys []DeviceKey `json:"keys"`
	}{}

	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/apps/%s/devices/%s/keys", appid, device), nil)
	if err != nil {
		return nil, err
	}
	return ret.Keys, nil
}

// 删除设备的密钥
func (manager *Manager) DeleteDeviceKey(appid, device, dak string) error {

	device = base64.URLEncoding.EncodeToString([]byte(device))

	return manager.client.Call(context.Background(), nil, "DELETE", manager.url("/apps/%s/devices/%s/keys/%s", appid, device, dak), nil)
}

// 禁用、启用设备的密钥
func (manager *Manager) UpdateDeviceKeyState(appid, device, dak string, state int) error {

	device = base64.URLEncoding.EncodeToString([]byte(device))

	req := M{"state": state}
	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/apps/%s/devices/%s/keys/%s/state", appid, device, dak), nil, req)
}

// 某个设备的密钥克隆给新的设备，不用重新对设备进行烧录新的密钥
func (manager *Manager) CloneDeviceKey(appid, fromdevice, todevice string, cleanSelfKeys, deleteDevice bool, deviceAccessKey string) ([]DeviceKey, error) {
	todevice = base64.URLEncoding.EncodeToString([]byte(todevice))
	req := M{"fromDevice": fromdevice, "cleanSelfKeys": cleanSelfKeys, "deleteDevice": deleteDevice, "deviceAccessKey": deviceAccessKey}
	ret := struct {
		Keys []DeviceKey `json:"cloneKeys"`
	}{}
	err := manager.client.CallWithJson(context.Background(), &ret, "POST", manager.url("/apps/%s/devices/%s/keys/clone", appid, todevice), nil, req)
	if err != nil {
		return nil, err
	}
	return ret.Keys, nil
}

func (manager *Manager) QueryAppidDeviceNameByAccessKey(dak string) (string, string, error) {
	var ret = struct {
		Appid  string `json:"appid"`
		Device string `json:"device"`
	}{}
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/keys/%s", dak), nil)
	if err != nil {
		return "", "", err
	}
	return ret.Appid, ret.Device, nil
}

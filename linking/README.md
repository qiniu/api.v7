 # Linking Cloud Server-Side Library for Go

## Features

- 设备管理
	- [x] 添加新的设备: AddDevice(appid string, dev *Device)
	- [x] 查询设备详细信息: QueryDevice(appid, device string)
	- [x] 查询所有设备的列表: ListDevice(appid, prefix, marker string, limit int, online bool)
	- [x] 查询设备的在线记录: ListDeviceHistoryactivity(appid, dev string, start, end int, marker string, limit int)
	- [x] 更新设备配置: UpdateDevice(appid, device string, ops []PatchOperation)
	- [x] 删除指定的设备: DeleteDevice(appid, device string)

- 设备密钥管理
	- [x] 新增设备密钥: AddDeviceKey(appid, device string)
	- [x] 查询设备密钥: QueryDeviceKey(appid, device string)
	- [x] 删除设备密钥: DeleteDeviceKey(appid, device, dak string)
	- [x] 禁用、启用设备的密钥：UpdateDeviceKeyState(appid, device, dak string, state int)
	- [x] 克隆设备: CloneDeviceKey(appid, fromdevice, todevice string, cleanSelfKeys, deleteDevice bool, deviceAccessKey string)
	- [x] 通过设备 acceskey 查询设备 appid 和 device name: QueryAppidDeviceNameByAccessKey(dak string)

- Vod
	- [x] 视频片段查询: Segments(appid, device string, start, end int, marker string, limit int)
	- [x] 视频片段进行收藏: Saveas(appid, device string, start, end int, fname, format string)

- Dtoken
	- [x] 视频回放/缩略图查询/倍速播放/延时直播/视频片段查询 Token 生成: VodToken(appid, device string, deadline int64)
	- [x] 在线记录查询/设备查询 Token 生成: StatusToken(appid, device string, deadline int64)

## Contents

- [Usage](#usage)
    - [Configuration](#configuration)
	- [Device](#device)
		- [Add Device](#add-device)
		- [Query Device](#query-device)
		- [List Device](#list-device)
		- [List Device history activity](#list-device-history-activity)
		- [Update Device](#update-device)
		- [Delete Device](#delete-device)

	- [Device Token](#deviceKey)
		- [Add Device Key](#add-device-key)
		- [Query Deivce Key](#query-deviceKey)
		- [Delete Device Key](#delete-deviceKey)
		- [Update Device Key State](#update-deviceKey))
		- [Clone Device Key](#clone-deviceKey)
		- [Query Appid Device By Device Key](#query-appid-device-by-deviceKey)
	- [Vod](#vod)
		- [Get Segments](#get-segments)
		- [Saveas](#saveas)

	- [Dtoken](#vod)
		- [Get Vod Token](#get-vod-token)
		- [Get Status Token](#get-status-token)

## Usage

### Configuration

```go
package main

import (
	// ...
	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/linking"
)

var (
	AccessKey = "<QINIU ACCESS KEY>" // 替换成自己 Qiniu 账号的 AccessKey.
	SecretKey = "<QINIU SECRET KEY>" // 替换成自己 Qiniu 账号的 SecretKey.
	Appid   = "<Appid>"    // App 必须事先创建.
)

func main() {
	// ...
	mac := auth.New(testAccessKey, testSecretKey)
	manager := linking.NewManager(mac, nil)
	// ...
}
```

### URL

#### Add Device

```go
manager.AddDevice(Appid, dev)
```

#### Query Device

```go
dev2, err := manager.QueryDevice(Appid, device)
```

#### List Device

```go
devices, marker, err = manager.ListDevice(Appid, "sdk-testListDevice", "", 1000, false)
```

#### List Device history activity

```go
segs, marker, err := manager.ListDeviceHistoryactivity(Appid, device, int(start), int(end), "", 1000)
```

#### Update Device

```go
ops := []PatchOperation{
		PatchOperation{Op: "replace", Key: "segmentExpireDays", Value: 30},
	}
dev3, err := manager.UpdateDevice(Appid, device, ops)
```

#### Delete Device
```go
manager.DeleteDevice(Appid, device)
```

### Device Key

#### Add Device Key

```go
manager.AddDeviceKey(Appid, device)
```

#### Query Deivce Key

```go
keys, err := manager.QueryDeviceKey(Appid, device)
```

#### Delete Device Key

```go
err := manager.DeleteDeviceKey(Appid, device, dak)
```

#### Update Device Key State

```go
err = manager.UpdateDeviceKeyState(Appid, device, dak, 1)
```

#### Clone Device Key

```go
keys, err := manager.CloneDeviceKey(Appid, device1, device2, false, false, dak1)
```

#### Query Appid Device By Device Key
```go
appid, device, err := manager.QueryAppidDeviceNameByAccessKey(dak)
```

### Vod

#### Get Segments

```go
end := time.Now().Unix()
start := time.Now().Add(-time.Hour).Unix()
segs, marker, err := manager.Segments(Appid, device, int(start), int(end), "", 1000)
```

#### Saveas

```go
end := time.Now().Unix()
start := time.Now().Add(-time.Hour).Unix()
saveasReply, _ := manager.Saveas(Appid, device, int(start), int(end), "testSaveas.mp4", "mp4")
```
### Dtoken

#### Get Vod Token
```go
token, err := manager.VodToken(Appid, device, time.Now().Add(time.Hour*5).Unix())
fmt.Println(token)
```
#### Get Status Token
```go
time.Sleep(time.Second)
token, err := manager.VodToken(Appid, device, time.Now().Add(time.Hour*5).Unix())
noError(t, err)
fmt.Println(token)
```


视频播放相关 api 编程模型说明：
如用户请求视频相关api(例如：[https://developer.qiniu.com/linking/api/5650/playback](https://developer.qiniu.com/linking/api/5650/playback))时进行业务服务器中转，会造成访问路径比较长, 因此建议服务端只进行dtoken的签算并提供给播放端，播放端直接请求视频播放相关的API，提升播放体验。
具体下token签算参考 [https://developer.qiniu.com/linking/api/5680/the-signature-in-the-url](https://developer.qiniu.com/linking/api/5680/the-signature-in-the-url)
![WeChatWorkScreenshot_c6c3e265-5ca3-48d9-8cc5-c40343506eeb](https://user-images.githubusercontent.com/34932312/63987548-3d000100-cb0b-11e9-971b-7aea84e07c67.png)

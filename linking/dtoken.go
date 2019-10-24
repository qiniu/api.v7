package linking

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Statement struct {
	Action string `json:"action"`
}

// 设备访问凭证
type DeviceAccessToken struct {
	Appid     string      `json:"appid"`     // appId
	Device    string      `json:"device"`    // device name
	DeadLine  int64       `json:"deadline"`  // 该token的有效期截止时间
	Random    int64       `json:"random"`    // 随机数，保证DEVICE ACCESS TOKEN全局唯一
	Statement []Statement `json:"statement"` // 针对那种功能进行授权
}

func (manager *Manager) deviceToken(policy *DeviceAccessToken) (string, error) {
	putPolicy, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}
	encodedDeviceAccessToken := base64.URLEncoding.EncodeToString(putPolicy)
	sign := manager.mac.Sign([]byte(encodedDeviceAccessToken))
	token := sign + ":" + encodedDeviceAccessToken
	return token, nil
}

// 视频回放/缩略图查询/倍速播放/延时直播/视频片段查询 Token
func (manager *Manager) VodToken(appid, device string, deadline int64) (string, error) {
	policy := &DeviceAccessToken{Appid: appid,
		Device:    device,
		DeadLine:  deadline,
		Random:    time.Now().UnixNano(),
		Statement: []Statement{Statement{Action: "linking:vod"}},
	}
	return manager.deviceToken(policy)
}

// 在线记录查询/设备查询 Token
func (manager *Manager) StatusToken(appid, device string, deadline int64) (string, error) {
	policy := &DeviceAccessToken{Appid: appid,
		Device:    device,
		DeadLine:  deadline,
		Random:    time.Now().UnixNano(),
		Statement: []Statement{Statement{Action: "linking:status"}},
	}
	return manager.deviceToken(policy)
}

func (manager *Manager) Token(appid, device string, deadline int64, actions []Statement) (string, error) {
	policy := &DeviceAccessToken{Appid: appid,
		Device:    device,
		DeadLine:  deadline,
		Random:    time.Now().UnixNano(),
		Statement: actions,
	}
	return manager.deviceToken(policy)
}

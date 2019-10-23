package linking

import (
	"context"
	"encoding/base64"
	"net/url"
)

//-----------------------------------------------------------------------------
// 历史记录
type DeviceHistoryItem struct {
	LoginAt      int64  `json:"loginAt"`
	LogoutAt     int64  `json:"logoutAt"`
	RemoteIp     string `json:"remoteIp,omitempty"`
	LogoutReason string `json:"logoutReason,omitempty"`
}

// 查询指定时间段内设备的在线记录
func (manager *Manager) ListDeviceHistoryactivity(appid, dev string, start, end int, marker string, limit int) ([]DeviceHistoryItem, string, error) {
	ret := struct {
		Items  []DeviceHistoryItem `json:"items"`
		Marker string              `json:"marker"`
	}{}
	dev = base64.URLEncoding.EncodeToString([]byte(dev))
	query := url.Values{}
	if limit > 0 {
		setQuery(query, "limit", limit)
	}
	if marker != "" {
		setQuery(query, "marker", marker)
	}
	setQuery(query, "start", start)
	setQuery(query, "end", end)
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/apps/%s/devices/%s/historyactivity?%v", appid, dev, query.Encode()), nil)
	if err != nil {
		return nil, "", err
	}
	return ret.Items, ret.Marker, nil
}

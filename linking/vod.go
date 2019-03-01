package linking

import (
	"context"
	"encoding/base64"
	"net/url"
)

type Segment struct {
	From  int    `json:"from"`
	To    int    `json:"to"`
	Frame string `json:"frame"`
}

// 视频片段查询
func (manager *Manager) Segments(appid, device string, start, end int, marker string, limit int) ([]Segment, string, error) {

	ret := struct {
		Items  []Segment `json:"items"`
		Marker string    `json:"marker"`
	}{}

	query := url.Values{}
	if limit > 0 {
		setQuery(query, "limit", limit)
	}
	if marker != "" {
		setQuery(query, "marker", marker)
	}
	if start > 0 {
		setQuery(query, "start", start)
	}
	if end > 0 {
		setQuery(query, "end", end)
	}

	device = base64.URLEncoding.EncodeToString([]byte(device))
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/apps/%s/devices/%s/vod/segments?%v", appid, device, query.Encode()), nil)
	if err != nil {
		return nil, "", err
	}
	return ret.Items, ret.Marker, nil
}

type SaveasReply struct {
	Fname       string `json:"fname"`
	PersistenId string `json:"persistentId,omitempty"`
	Duration    int    `json:"duration"` // ms
}

// 指定的视频片段进行收藏，保存在云存储上
func (manager *Manager) Saveas(appid, device string, start, end int, fname, format string) (*SaveasReply, error) {

	device = base64.URLEncoding.EncodeToString([]byte(device))

	req := M{
		"start": start,
		"end":   end,
	}
	if fname != "" {
		req["fname"] = fname
	}
	if format != "" {
		req["format"] = format
	}

	var reply *SaveasReply

	err := manager.client.CallWithJson(context.Background(), reply, "POST", manager.url("/apps/%s/devices/%s/vod/saveas", appid, device), nil, req)
	return reply, err
}

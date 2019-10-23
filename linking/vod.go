package linking

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

//-----------------------------------------------------------------------------
// mqtt rpc
type RpcRequest struct {
	Action   int             `json:"action"`
	Params   json.RawMessage `json:"params,omitempty"`
	Timeout  int             `json:"timeout,omitempty"`
	Response bool            `json:"response,omitempty"`
}
type DevResponse struct {
	ErrorCode int             `json:"errorCode,omitempty"`
	Error     string          `json:"error,omitempty"`
	Value     json.RawMessage `json:"value,omitempty"`
}
type RpcResponse struct {
	Id   string      `json:"id"`
	Resp DevResponse `json:"response"`
}

func (manager *Manager) RPC(appid, device string, req *RpcRequest) (*RpcResponse, error) {

	var ret RpcResponse
	device = base64.URLEncoding.EncodeToString([]byte(device))
	err := manager.client.CallWithJson(context.Background(), &ret, "POST", manager.url("/apps/%s/devices/%s/rpc", appid, device), nil, req)
	if err != nil {
		return nil, err
	}
	return &ret, nil
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

	var reply SaveasReply

	err := manager.client.CallWithJson(context.Background(), &reply, "POST", manager.url("/apps/%s/devices/%s/vod/saveas", appid, device), nil, req)
	return &reply, err
}

type LiveRequest struct {
	Appid      string `json:"appid"`
	DeviceName string `json:"deviceName"`
	PublishIP  string `json:"publishIP"`
	PlayIP     string `json:"playIP"`
}

type PlayUrls struct {
	Rtmp string `json:"rtmp"`
	Hls  string `json:"hls"`
	Flv  string `json:"flv"`
}

type LiveResponse struct {
	PublishUrl string   `json:"publishUrl"`
	PlayUrls   PlayUrls `json:"playUrls"`
}

// 指定的视频片段进行收藏，保存在云存储上
func (manager *Manager) StartLive(req *LiveRequest) (*LiveResponse, error) {

	var reply LiveResponse

	err := manager.client.CallWithJson(context.Background(), &reply, "POST", manager.url("/startlive"), nil, req)
	return &reply, err
}

// 查询活跃设备信息
type StatReq struct {
	Start  int    `form:"start"`
	End    int    `form:"end"`
	Group  string `form:"g"`
	Select string `form:"select"`
}

func (manager *Manager) Stat(req *StatReq) ([]M, error) {

	var ret []M
	query := url.Values{}
	setQuery(query, "start", req.Start)
	setQuery(query, "end", req.End)
	setQuery(query, "g", req.Group)
	setQuery(query, "select", req.Select)

	err := manager.client.Call(context.Background(), &ret, "GET", fmt.Sprintf("http://linking.qiniuapi.com/statd/device?%v", query.Encode()), nil)
	return ret, err
}

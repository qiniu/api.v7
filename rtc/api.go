package rtc

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/qiniu/api.v7/auth/qbox"
)

var (
	// RtcHost 为 Qiniu RTC Server API服务域名
	RtcHost = "rtc.qiniuapi.com"
)

// Manager 提供了 Qiniu RTC Server API 相关功能
type Manager struct {
	mac        *qbox.Mac
	httpClient *http.Client
}

// MergePublishRtmp  连麦合流转推 RTMP 的配置
// Enable: 布尔类型，用于开启和关闭所有房间的合流功能。
// AudioOnly: 布尔类型，可选，指定是否只合成音频。
// Height, Width: int，可选，指定合流输出的高和宽，默认为 640 x 480。
// OutputFps: int，可选，指定合流输出的帧率，默认为 25 fps 。
// OutputKbps: int，可选，指定合流输出的码率，默认为 1000 。
// URL: 合流后转推旁路直播的地址，可选，支持魔法变量配置按照连麦房间号生成不同的推流地址。如果是转推到七牛直播云，不建议使用该配置
// StreamTitle: 转推七牛直播云的流名，可选，支持魔法变量配置按照连麦房间号生成不同的流名。例如，配置 Hub 为 qn-zhibo ，配置 StreamTitle 为 $(roomName) ，则房间 meeting-001 的合流将会被转推到 rtmp://pili-publish.qn-zhibo.***.com/qn-zhibo/meeting-001地址。详细配置细则，请咨询七牛技术支持。
type MergePublishRtmp struct {
	Enable      bool   `json:"enable,omitempty"`
	AudioOnly   bool   `json:"audioOnly,omitempty"`
	Height      int    `json:"height,omitempty"`
	Width       int    `json:"width,omitempty"`
	OutputFps   int    `json:"fps,omitempty"`
	OutputKbps  int    `json:"kbps,omitempty"`
	URL         string `json:"url,omitempty"`
	StreamTitle string `json:"streamTitle,omitempty"`
}

// App 完整信息
// AppID: app 的唯一标识，创建的时候由系统生成。
type App struct {
	AppID string `json:"appId"`
	AppReq
	MergePublishRtmp MergePublishRtmp `json:"mergePublishRtmp,omitempty"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
}

// AppReq 创建 App 请求参数
// Title: app 的名称， 可选。
// Hub: 绑定的直播 hub，可选，用于合流后 rtmp 推流。
// MaxUsers: int 类型，可选，连麦房间支持的最大在线人数。
// NoAutoCloseRoom: bool 指针类型，可选，true 表示禁止自动关闭房间。
// NoAutoCreateRoom: bool 指针指型，可选，true 表示禁止自动创建房间。
// NoAutoKickUser: bool 类型，可选，禁止自动踢人。
type AppReq struct {
	Hub              string `json:"hub,omitempty"`
	Title            string `json:"title,omitempty"`
	MaxUsers         int    `json:"maxUsers,omitempty"`
	NoAutoCloseRoom  bool   `json:"noAutoCloseRoom,omitempty"`
	NoAutoCreateRoom bool   `json:"noAutoCreateRoom,omitempty"`
	NoAutoKickUser   bool   `json:"noAutoKickUser,omitempty"`
}

// MergePublishRtmpInfo 连麦合流转推 RTMP 的配置更改信息
type MergePublishRtmpInfo struct {
	Enable      *bool   `json:"enable,omitempty"`
	AudioOnly   *bool   `json:"audioOnly,omitempty"`
	Height      *int    `json:"height,omitempty"`
	Width       *int    `json:"width,omitempty"`
	OutputFps   *int    `json:"fps,omitempty"`
	OutputKbps  *int    `json:"kbps,omitempty"`
	URL         *string `json:"url,omitempty"`
	StreamTitle *string `json:"streamTitle,omitempty"`
}

// AppUpdateInfo 更改信息
// MergePublishRtmpInfo 连麦合流转推 RTMP 的配置更改信息
type AppUpdateInfo struct {
	Hub              *string               `json:"hub,omitempty"`
	Title            *string               `json:"title,omitempty"`
	MaxUsers         *int                  `json:"maxUsers,omitempty"`
	NoAutoCloseRoom  *bool                 `json:"noAutoCloseRoom,omitempty"`
	NoAutoCreateRoom *bool                 `json:"noAutoCreateRoom,omitempty"`
	NoAutoKickUser   *bool                 `json:"noAutoKickUser,omitempty"`
	MergePublishRtmp *MergePublishRtmpInfo `json:"mergePublishRtmp,omitempty"`
}

// User 连麦房间里的用户
type User struct {
	UserID string `json:"userId"`
}

// NewManager 用来构建一个新的 Manager
func NewManager(mac *qbox.Mac) *Manager {
	httpClient := http.DefaultClient
	return &Manager{mac: mac, httpClient: httpClient}
}

// CreateApp 新建实时音视频云
func (r *Manager) CreateApp(appReq AppReq) (App, ResInfo, error) {
	url := buildURL("/v3/apps")
	ret := App{}
	info := postReq(r.httpClient, r.mac, url, &appReq, &ret)
	return ret, *info, info.Err
}

// GetApp 根据 appID 获取 实时音视频云 信息
func (r *Manager) GetApp(appID string) (App, ResInfo, error) {
	url := buildURL("/v3/apps/" + appID)
	ret := App{}
	info := getReq(r.httpClient, r.mac, url, &ret)
	return ret, *info, info.Err
}

// DeleteApp 根据 appID 删除 实时音视频云
func (r *Manager) DeleteApp(appID string) (ResInfo, error) {
	url := buildURL("/v3/apps/" + appID)
	info := delReq(r.httpClient, r.mac, url, nil)
	return *info, info.Err
}

// UpdateApp 根据 appID, App 更改实时音视频云 信息
func (r *Manager) UpdateApp(appID string, appInfo AppUpdateInfo) (App, ResInfo, error) {
	url := buildURL("/v3/apps/" + appID)
	ret := App{}
	info := postReq(r.httpClient, r.mac, url, &appInfo, &ret)
	return ret, *info, info.Err
}

// ListUser 根据 appID, roomName 获取连麦房间里在线的用户
// appID: 连麦房间所属的 app 。
// roomName: 操作所查询的连麦房间。
func (r *Manager) ListUser(appID, roomName string) ([]User, ResInfo, error) {
	url := buildURL("/v3/apps/" + appID + "/rooms/" + roomName + "/users")
	users := struct {
		Users []User `json:"users"`
	}{}
	info := getReq(r.httpClient, r.mac, url, &users)
	return users.Users, *info, info.Err
}

// KickUser 根据 appID, roomName, UserID 剔除在线的用户
// appID: 连麦房间所属的 app 。
// roomName: 连麦房间。
// userID: 操作所剔除的用户。
func (r *Manager) KickUser(appID, roomName, userID string) (ResInfo, error) {
	url := buildURL("/v3/apps/" + appID + "/rooms/" + roomName + "/users/" + userID)
	info := delReq(r.httpClient, r.mac, url, nil)
	return *info, info.Err
}

// RoomQuery 房间查询响应结果
// IsEnd: bool 类型，分页查询是否已经查完所有房间。
// Offset: int 类型，下次分页查询使用的位移标记。
// Rooms: 当前活跃的房间名列表。
type RoomQuery struct {
	IsEnd  bool       `json:"end"`
	Offset int        `json:"offset"`
	Rooms  []RoomName `json:"rooms"`
}

// RoomName 房间名
type RoomName string

// ListActiveRoom 根据 appID, roomNamePrefix, offset, limit 查询当前活跃的房间
// appID: 连麦房间所属的 app 。
// roomNamePrefix: 所查询房间名的前缀索引，可以为空。
// offset: int 类型，分页查询的位移标记。
// limit: int 类型，此次查询的最大长度。
func (r *Manager) ListActiveRoom(appID, roomNamePrefix string, offset, limit int) (RoomQuery, ResInfo, error) {
	query := ""
	roomNamePrefix = strings.TrimSpace(roomNamePrefix)
	if len(roomNamePrefix) != 0 {
		query = "prefix=" + roomNamePrefix + "&"
	}
	query += fmt.Sprintf("offset=%v&limit=%v", offset, limit)
	url := buildURL("/v3/apps/" + appID + "/rooms?" + query)
	ret := RoomQuery{}
	info := getReq(r.httpClient, r.mac, url, &ret)
	return ret, *info, info.Err
}

// ListAllActiveRoom 根据 appID, roomNamePrefix 查询当前活跃的房间
// appID: 连麦房间所属的 app 。
// roomNamePrefix: 所查询房间名的前缀索引，可以为空。
func (r *Manager) ListAllActiveRoom(appID, roomNamePrefix string) ([]RoomName, ResInfo, error) {
	ns := []RoomName{}
	var err error = nil
	q := RoomQuery{}
	q.IsEnd = false
	q.Rooms = []RoomName{""}
	info := ResInfo{Code: 200}
	for offset := 0; err == nil && info.Code == 200 &&
		len(q.Rooms) > 0 && !q.IsEnd; offset += 100 {
		q, info, err = r.ListActiveRoom(appID, roomNamePrefix, offset, 100)
		if err != nil && info.Code != 401 {
			time.Sleep(100 * time.Millisecond)
			q, info, err = r.ListActiveRoom(appID, roomNamePrefix, offset, 100)
		}
		if err != nil && info.Code != 401 {
			time.Sleep(500 * time.Millisecond)
			q, info, err = r.ListActiveRoom(appID, roomNamePrefix, offset, 100)
		}
		if err == nil {
			ns = append(ns, q.Rooms...)
		}
	}
	return ns, info, err
}

// RoomAccess 房间管理凭证
// AppID: 房间所属帐号的 app 。
// RoomName: 房间名称，需满足规格 ^[a-zA-Z0-9_-]{3,64}$
// UserID: 请求加入房间的用户 ID，需满足规格 ^[a-zA-Z0-9_-]{3,50}$
// ExpireAt: int64 类型，鉴权的有效时间，传入以秒为单位的64位Unix绝对时间，token 将在该时间后失效。
// Permission: 该用户的房间管理权限，"admin" 或 "user"，默认为 "user" 。当权限角色为 "admin" 时，拥有将其他用户移除出房间等特权.
type RoomAccess struct {
	AppID      string `json:"appId"`
	RoomName   string `json:"roomName"`
	UserID     string `json:"userId"`
	ExpireAt   int64  `json:"expireAt"`
	Permission string `json:"permission"`
}

// RoomToken 生成房间管理鉴权，连麦用户终端通过房间管理鉴权获取七牛连麦服务。
func (r *Manager) RoomToken(roomAccess RoomAccess) (token string, err error) {
	roomAccessByte, err := json.Marshal(roomAccess)
	if err != nil {
		return
	}
	buf := make([]byte, base64.URLEncoding.EncodedLen(len(roomAccessByte)))
	base64.URLEncoding.Encode(buf, roomAccessByte)

	hmacsha1 := hmac.New(sha1.New, r.mac.SecretKey)
	hmacsha1.Write(buf)
	sign := hmacsha1.Sum(nil)

	encodedSign := base64.URLEncoding.EncodeToString(sign)
	token = r.mac.AccessKey + ":" + encodedSign + ":" + string(buf)
	return
}

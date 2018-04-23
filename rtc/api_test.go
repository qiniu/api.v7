package rtc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/qiniu/api.v7/auth/qbox"
)

var manager *Manager

func init() {
	accessKey := os.Getenv("QINIU_ACCESS_KEY")
	secretKey := os.Getenv("QINIU_SECRET_KEY")

	mac := qbox.NewMac(accessKey, secretKey)
	manager = NewManager(mac)
}

func TestApp(t *testing.T) {
	app := checkCreateApp(t)
	checkGetApp(t, app.AppID)
	checkAllActiveRoom(t, app.AppID)
	checkListUser(t, app.AppID)
	checkKickUser(t, app.AppID)
	checkUpdate(t, app.AppID)
	checkRoomToken(t, app.AppID)
	checkDel(t, app.AppID)
}

func checkApp(t *testing.T, app *App) {
	if app.Title != "gosdk-test" || !app.NoAutoCloseRoom || app.NoAutoCreateRoom || !app.NoAutoKickUser {
		t.Error(app)
	}
}

func checkInfo(t *testing.T, info *ResInfo, err error) {
	if err != nil || info.Code != 200 {
		t.Error(info)
	}
}

func checkCreateApp(t *testing.T) *App {
	appReq := AppReq{Hub: "hailong", Title: "gosdk-test", MaxUsers: 3, NoAutoCloseRoom: true, NoAutoCreateRoom: false, NoAutoKickUser: true}
	app, info, err := manager.CreateApp(appReq)
	checkApp(t, &app)
	checkInfo(t, &info, err)
	return &app
}

func checkGetApp(t *testing.T, appID string) {
	app, info, err := manager.GetApp(appID)
	checkApp(t, &app)
	checkInfo(t, &info, err)
}

func checkUpdate(t *testing.T, appID string) {
	appInfo := AppUpdateInfo{}
	mergePublishRtmp := MergePublishRtmpInfo{}
	mergePublishRtmp.Enable = &(&struct{ x bool }{true}).x

	AudioOnly := true
	mergePublishRtmp.AudioOnly = &AudioOnly
	StreamTitle := "ddddddddd"
	mergePublishRtmp.StreamTitle = &StreamTitle

	appInfo.NoAutoKickUser = &(&struct{ x bool }{false}).x
	appInfo.NoAutoCreateRoom = func(i bool) *bool { return &i }(true)
	appInfo.MergePublishRtmp = &mergePublishRtmp

	m := struct2JsonMap(appInfo)
	// should be nil , not ""
	if m["Title"] != nil {
		t.Errorf("m[\"Title\"] should be nil, but %v. and m is %v\n", m["Title"], m)
	}

	app, info, err := manager.UpdateApp(appID, appInfo)
	checkInfo(t, &info, err)
	if app.MaxUsers != 3 || !app.NoAutoCreateRoom || app.NoAutoKickUser || !app.MergePublishRtmp.AudioOnly || app.MergePublishRtmp.StreamTitle != StreamTitle || app.MergePublishRtmp.URL != "" {
		t.Error(app)
	}

	appInfo = AppUpdateInfo{}
	appInfo.MaxUsers = &(&struct{ x int }{4}).x
	app, info, err = manager.UpdateApp(appID, appInfo)
	checkInfo(t, &info, err)
	if app.MaxUsers != 4 || !app.NoAutoCreateRoom || app.NoAutoKickUser || !app.MergePublishRtmp.AudioOnly || app.MergePublishRtmp.StreamTitle != StreamTitle || app.MergePublishRtmp.URL != "" {
		t.Error(app)
	}
}

func checkAllActiveRoom(t *testing.T, appID string) {
	rooms, info, err := manager.ListAllActiveRoom(appID, "l")
	checkInfo(t, &info, err)
	if len(rooms) > 0 {
		t.Log(rooms)
	}
}

func checkListUser(t *testing.T, appID string) {
	users, info, err := manager.ListUser(appID, "roomName")
	checkInfo(t, &info, err)
	if len(users) > 0 {
		t.Log(users)
	}
}

func checkKickUser(t *testing.T, appID string) {
	info, err := manager.KickUser(appID, "roomName", "userID")
	if err != nil || (info.Code != 200 && info.Code/100 != 6) {
		t.Error(info)
	}
}

func checkRoomToken(t *testing.T, appID string) {
	roomName := "sdhfuexx"
	userID := "ghjkdfie"
	token, err := manager.RoomToken(RoomAccess{AppID: appID, RoomName: roomName, UserID: userID, ExpireAt: time.Now().Unix() + 3600})
	url := fmt.Sprintf("https://rtc.qiniuapi.com/v3/apps/%v/rooms/%v/auth?user=%v&token=%v", appID, roomName, userID, token)
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if res.ContentLength > 2*1024*1024 {
		t.Errorf("response is too long. url: %v, response headers: %v", req.URL.RequestURI(), res.Header)
	}
	resData, err := ioutil.ReadAll(res.Body)
	ret := string(resData)
	if res.StatusCode != 200 || strings.Index(ret, "accessToken") == -1 {
		t.Errorf("result is wrong: %v", ret)
	}
}

func checkDel(t *testing.T, appID string) {
	info, err := manager.DeleteApp(appID)
	if err != nil || (info.Code != 200 && info.Code/100 != 6) {
		t.Error(info)
	}
}

func struct2JsonMap(obj interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	j, _ := json.Marshal(obj)
	json.Unmarshal(j, &data)
	return data
}

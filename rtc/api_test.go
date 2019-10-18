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

	"github.com/qiniu/api.v7/v7/auth"
)

var manager *Manager

func init() {
	accessKey := os.Getenv("accessKey")
	secretKey := os.Getenv("secretKey")

	mac := auth.New(accessKey, secretKey)
	manager = NewManager(mac)
}

func TestApp(t *testing.T) {
	app := checkCreateApp(t)
	checkGetApp(t, app.AppID)
	rooms := checkAllActiveRooms(t, app.AppID)
	room := "roomName"
	if len(rooms) > 0 {
		room = string(rooms[0])
	}
	users := checkListUser(t, app.AppID, room)
	userID := "userID"
	if len(users) > 0 {
		userID = users[0].UserID
	}
	checkKickUser(t, app.AppID, room, userID)
	checkUpdate(t, app.AppID)
	checkRoomToken(t, app.AppID)
	checkDel(t, app.AppID)
}

func checkApp(t *testing.T, app *App) {
	if app.Title != "gosdk-test" || !app.NoAutoKickUser {
		t.Error(app)
	}
}

func checkInfo(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func checkCreateApp(t *testing.T) *App {
	appInitConf := AppInitConf{Hub: "hailong", Title: "gosdk-test", MaxUsers: 3, NoAutoKickUser: true}
	app, err := manager.CreateApp(appInitConf)
	checkApp(t, &app)
	checkInfo(t, err)
	return &app
}

func checkGetApp(t *testing.T, appID string) {
	app, err := manager.GetApp(appID)
	checkApp(t, &app)
	checkInfo(t, err)
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
	appInfo.MergePublishRtmp = &mergePublishRtmp

	m := struct2JsonMap(appInfo)
	// should be nil , not ""
	if m["Title"] != nil {
		t.Errorf("m[\"Title\"] should be nil, but %v. and m is %v\n", m["Title"], m)
	}

	app, err := manager.UpdateApp(appID, appInfo)
	checkInfo(t, err)
	if app.MaxUsers != 3 || app.NoAutoKickUser || !app.MergePublishRtmp.AudioOnly || app.MergePublishRtmp.StreamTitle != StreamTitle || app.MergePublishRtmp.URL != "" {
		t.Error(app)
	}

	appInfo = AppUpdateInfo{}
	appInfo.MaxUsers = &(&struct{ x int }{4}).x
	app, err = manager.UpdateApp(appID, appInfo)
	checkInfo(t, err)
	if app.MaxUsers != 4 || app.NoAutoKickUser || !app.MergePublishRtmp.AudioOnly || app.MergePublishRtmp.StreamTitle != StreamTitle || app.MergePublishRtmp.URL != "" {
		t.Error(app)
	}
}

func checkAllActiveRooms(t *testing.T, appID string) []RoomName {
	rooms, err := manager.ListAllActiveRooms(appID, "l")
	checkInfo(t, err)
	t.Logf("Rooms: %v", rooms[:min(10, len(rooms))])
	return rooms
}

func checkListUser(t *testing.T, appID, roomName string) []User {
	users, err := manager.ListUser(appID, roomName)
	checkInfo(t, err)
	t.Logf("Users: %v", users[:min(10, len(users))])
	return users
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func checkKickUser(t *testing.T, appID, roomName, userID string) {
	err := manager.KickUser(appID, roomName, userID)
	// 615   {"error":"room not active"}
	if err != nil {
		t.Log(err)
		if strings.Index(err.Error(), "room not active") == -1 {
			t.Error(err)
		}
	}
}

func checkRoomToken(t *testing.T, appID string) {
	roomName := "sdhfuexx"
	userID := "ghjkdfie"
	token, err := manager.GetRoomToken(RoomAccess{AppID: appID, RoomName: roomName, UserID: userID, ExpireAt: time.Now().Unix() + 3600})
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
	err := manager.DeleteApp(appID)
	if err != nil {
		t.Error(err)
	}
}

func struct2JsonMap(obj interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	j, _ := json.Marshal(obj)
	json.Unmarshal(j, &data)
	return data
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/rtc"
)

var (
	manager *rtc.Manager
)

func init() {
	accessKey := os.Getenv("QINIU_ACCESS_KEY")
	secretKey := os.Getenv("QINIU_SECRET_KEY")

	mac := auth.New(accessKey, secretKey)
	manager = rtc.NewManager(mac)
}

func getRoomToken(appId, roomName, userID string) (token string, err error) {
	token, err = manager.GetRoomToken(rtc.RoomAccess{AppID: appId, RoomName: roomName, UserID: userID, ExpireAt: time.Now().Unix() + 3600})
	return
}

func dosomethingbyRoomToken(token string, appId, roomName, userID string) (err error) {
	url := fmt.Sprintf("https://rtc.qiniuapi.com/v3/apps/%v/rooms/%v/auth?user=%v&token=%v", appId, "roomName", "userID", token)
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	fmt.Println(string(resData))
	jn, err := json.MarshalIndent(res.Header, "  ", "  ")
	fmt.Println(string(jn))
	return

}

func createApp() (rtc.App, error) {
	hubName := "hailong"
	appReq := rtc.AppInitConf{Hub: hubName, Title: "gosdk-test", MaxUsers: 3, NoAutoKickUser: true}
	app, err := manager.CreateApp(appReq)
	return app, err
}

func updateApp(appID string) (rtc.App, error) {
	appInfo := rtc.AppUpdateInfo{}
	mergePublishRtmp := rtc.MergePublishRtmpInfo{}
	mergePublishRtmp.Enable = &(&struct{ x bool }{true}).x

	appInfo.NoAutoKickUser = &(&struct{ x bool }{false}).x
	appInfo.MergePublishRtmp = &mergePublishRtmp

	app, err := manager.UpdateApp(appID, appInfo)
	return app, err
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func main() {
	fmt.Println("\napp := createApp()")
	app, err := createApp()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nroomToken(app.AppID)")
	token, err := getRoomToken(app.AppID, "roomName", "userID")
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		err = dosomethingbyRoomToken(token, app.AppID, "roomName", "userID")
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("\ngetApp(app.AppID)")
	app, err = manager.GetApp(app.AppID)
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		j, _ := json.MarshalIndent(app, "  ", "  ")
		fmt.Println(string(j))
	}

	fmt.Println("\nlistAllUser(app.AppID)")
	rooms, err := manager.ListAllActiveRooms(app.AppID, "l")
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		fmt.Printf("Rooms: %v\n", rooms[:min(10, len(rooms))])
	}

	fmt.Println("\nlistUser(app.AppID)")
	roomQuery, err := manager.ListActiveRooms(app.AppID, "l", 0, 1)
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		j, _ := json.MarshalIndent(roomQuery, "  ", "  ")
		fmt.Println(string(j))
	}

	fmt.Println("\nkickUser(app.AppID)")
	err = manager.KickUser(app.AppID, "roomName", "userID")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\nupdateApp(app.AppID)")
	app, err = updateApp(app.AppID)
	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		j, _ := json.MarshalIndent(app, "  ", "  ")
		fmt.Println(string(j))
	}

}

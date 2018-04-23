package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/rtc"
)

var (
	manager *rtc.Manager
)

func init() {
	accessKey := os.Getenv("QINIU_ACCESS_KEY")
	secretKey := os.Getenv("QINIU_SECRET_KEY")

	mac := qbox.NewMac(accessKey, secretKey)
	manager = rtc.NewManager(mac)
}

func roomToken(appId string) {
	token, err := manager.RoomToken(rtc.RoomAccess{AppID: appId, RoomName: "sdhfuexx", UserID: "ghjkdfie", ExpireAt: time.Now().Unix() + 3600})
	url := fmt.Sprintf("https://rtc.qiniuapi.com/v3/apps/%v/rooms/%v/auth?user=%v&token=%v", appId, "sdhfuexx", "ghjkdfie", token)
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	fmt.Println(req.URL.RequestURI())
	fmt.Println(err)
	res, err := http.DefaultClient.Do(req)
	fmt.Println(err)
	resData, err := ioutil.ReadAll(res.Body)
	fmt.Println(err)
	fmt.Println(string(resData))
	jn, err := json.MarshalIndent(res.Header, "  ", "  ")
	fmt.Println(err)
	fmt.Println(string(jn))
}

func createApp() *rtc.App {
	appReq := rtc.AppReq{Hub: "hailong", Title: "gosdk-test", MaxUsers: 3, NoAutoCloseRoom: true, NoAutoCreateRoom: true, NoAutoKickUser: true}
	app, info, err := manager.CreateApp(appReq)
	j, _ := json.MarshalIndent(app, "  ", "  ")
	fmt.Println(string(j))
	i, _ := json.MarshalIndent(info, "  ", "  ")
	fmt.Println(string(i))
	fmt.Println(err)
	fmt.Println()

	return &app
}

func delApp(appId string) {
	info, err := manager.DeleteApp(appId)
	fmt.Println(err)
	i, _ := json.MarshalIndent(info, "  ", "  ")
	fmt.Println(string(i))
	fmt.Println()
}

func getApp(appId string) {
	app, info, err := manager.GetApp(appId)
	fmt.Println(err)
	j, _ := json.MarshalIndent(app, "  ", "  ")
	fmt.Println(string(j))
	i, _ := json.MarshalIndent(info, "  ", "  ")
	fmt.Println(string(i))
	fmt.Println()
}

func updateApp(appId string) {
	appInfo := rtc.AppUpdateInfo{}
	mergePublishRtmp := rtc.MergePublishRtmpInfo{}
	mergePublishRtmp.Enable = &(&struct{ x bool }{true}).x

	appInfo.NoAutoKickUser = &(&struct{ x bool }{false}).x
	appInfo.NoAutoCreateRoom = &(&struct{ x bool }{true}).x
	appInfo.MergePublishRtmp = &mergePublishRtmp
	j, _ := json.MarshalIndent(appInfo, "  ", "  ")
	fmt.Println(string(j))

	app, info, err := manager.UpdateApp(appId, appInfo)
	fmt.Println(err)
	j, _ = json.MarshalIndent(app, "  ", "  ")
	fmt.Println(string(j))
	i, _ := json.MarshalIndent(info, "  ", "  ")
	fmt.Println(string(i))
	fmt.Printf("app.MaxUsers == 3 is %v\n", app.MaxUsers == 3)
	fmt.Printf("false app.NoAutoKickUser is %v\n", app.NoAutoKickUser)
	fmt.Printf("true app.NoAutoCreateRoom  is %v\n", app.NoAutoCreateRoom)
	fmt.Println()

}

func listUser(appId string) {
	rooms, info, err := manager.ListAllActiveRoom(appId, "l")
	fmt.Println(err)
	j, _ := json.MarshalIndent(rooms, "  ", "  ")
	fmt.Println(string(j))
	i, _ := json.MarshalIndent(info, "  ", "  ")
	fmt.Println(string(i))
	fmt.Println()
}

func listUser1(appId string) {
	rooms, info, err := manager.ListActiveRoom(appId, "l", 0, 1)
	fmt.Println(err)
	j, _ := json.MarshalIndent(rooms, "  ", "  ")
	fmt.Println(string(j))
	i, _ := json.MarshalIndent(info, "  ", "  ")
	fmt.Println(string(i))
	fmt.Println()
}

func kickUser(appId string) {
	info, err := manager.KickUser(appId, "lhurugt", "mvnrutirejur")
	fmt.Println(err)
	i, _ := json.MarshalIndent(info, "  ", "  ")
	fmt.Println(string(i))
	fmt.Println()
}

func main() {
	fmt.Println("\napp := createApp()")
	app := createApp()

	fmt.Println("\nroomToken(app.AppID)")
	roomToken(app.AppID)
	fmt.Println("\ngetApp(app.AppID)")
	getApp(app.AppID)
	fmt.Println("\nlistUser(app.AppID)")
	listUser(app.AppID)
	fmt.Println("\nlistUser1(app.AppID)")
	listUser1(app.AppID)
	fmt.Println("\nkickUser(app.AppID)")
	kickUser(app.AppID)
	fmt.Println("\nupdateApp(app.AppID)")
	updateApp(app.AppID)

	fmt.Println("\ndelApp(app.AppID)")
	delApp(app.AppID)
	fmt.Println("\ngetApp(app.AppID)")
	getApp(app.AppID)
	fmt.Println("\nlistUser1(app.AppID)")
	listUser(app.AppID)
	listUser1(app.AppID)
	fmt.Println("\nkickUser(app.AppID)")
	kickUser(app.AppID)

	updateApp(app.AppID)
}

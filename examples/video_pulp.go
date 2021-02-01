package main

import (
    "fmt"
    "github.com/qiniu/api.v7/auth/qbox"
    "github.com/qiniu/x/rpc.v7"
    "net/http"
)

type MediaAuditManager struct {
    client *rpc.Client
    mac    *qbox.Mac
}

func NewMediaAuditManager(mac *qbox.Mac) *MediaAuditManager {

    tp := &Transport{Mac: mac}
    cli := &http.Client{Transport: tp}

    return &MediaAuditManager{
        client: &rpc.Client{cli},
        mac:    mac,
    }
}

type Transport struct {
    Mac *qbox.Mac
}

func (c *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
    token, err := c.Mac.SignRequestV2(req)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Qiniu "+token)

    return http.DefaultTransport.RoundTrip(req)
}

type VideoNropParam struct {
    Data struct {
        URI string `json:"uri"`
    } `json:"data"`

    Ops []struct {
        Op string `json:"op"`
    } `json:"ops"`

    //Params struct {
    //  Async  bool `json:"async"`
    //  Vframe struct {
    //      Mode int `json:"mode"`,
    //  } `json:"vframe"`
    //} `json:"params"`
}

func (m *MediaAuditManager) VideoNrop(vid string, param interface{}) (ret interface{}, err error) {
    url1 := "https://argus.atlab.ai/v1/video/" + vid
    err = m.client.CallWithJson(nil, &ret, "POST", url1, param)
    return
}

func main() {
    mac := qbox.NewMac("ak", "sk")

    mgr := NewMediaAuditManager(mac)

    param := map[string]interface{}{
        "data":   map[string]string{"uri": "http://rwxf.qiniucdn.com/1.mp4"},
        "params": map[string]interface{}{"async": false},
        "ops":    []map[string]string{{"op": "pulp"}},
    }

    ret, err := mgr.VideoNrop("temp", param)

    fmt.Println("=====>", ret, err)
}

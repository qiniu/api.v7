package qvs

import (
	"context"
	"net/url"
)

type M map[string]interface{}

type PatchOperation struct {
	Op    string      `json:"op"`    // 更该或删除某个属性，replace：更改，delete：删除
	Key   string      `json:"key"`   // 要修改或删除的属性
	Value interface{} `json:"value"` // 要修改或删除属性的值
}

type DomainInfo struct {
	Domain string `json:"domain"`
	Type   string `json:"type"`
	CNAME  string `json:"cname"`
	State  int    `json:"state"`
}

type NameSpace struct {
	ID                     string       `json:"id"`
	Name                   string       `json:"name"`           // 空间名称(格式"^[a-zA-Z0-9_-]{1,100}$")
	Desc                   string       `json:"desc,omitempty"` // 空间描述
	AccessType             string       `json:"accessType"`     // 接入类型"gb28181"或者“rtmp”
	RTMPURLType            int          `json:"rtmpUrlType"`    // accessType为“rtmp”时，推拉流地址计算方式，1:static, 2:dynamic
	Domains                []string     `json:"domains"`        // 直播域名
	DomainDetails          []DomainInfo `json:"domainDetails"`
	Callback               string       `json:"callback,omitempty""`          // 后台服务器回调URL
	Disabled               bool         `json:"disabled"`                     // 流是否被启用, false:启用,true:禁用
	RecordTemplateId       string       `jons:"recordTemplateId,omitempty"`   // 录制模版id
	SnapShotTemplateId     string       `jons:"snapshotTemplateId,omitempty"` // 截图模版id
	RecordTemplateApplyAll bool         `json:"snapshotTemplateApplyAll"`     // 空间模版是否应用到全局
	SnapTemplateApplyAll   bool         `json:"snapshotTemplateApplyAll"`     // 截图模版是否应用到全局
	CreatedAt              int64        `json:"createdAt,omitempty"`          // 空间创建时间
	UpdatedAt              int64        `json:"updatedAt,omitempty"`          // 空间更新时间
}

/*
	创建空间API

	请求参数Body:
	name必填
	accessType必填
	rtmpUrlType当accessType为"rtmp"时必填
	domains当rtmpUrlType为1时必填
*/
func (manager *Manager) AddNamespace(ns *NameSpace) (*NameSpace, error) {

	var ret NameSpace
	err := manager.client.CallWithJson(context.Background(), &ret, "POST", manager.url("/namespaces"), nil, ns)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
	查询空间信息API
*/
func (manager *Manager) QueryNamespace(nsId string) (*NameSpace, error) {

	var ret NameSpace
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s", nsId), nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
	更新空间API

	可编辑参数: name/desc/callBack/recordTemplateId/snapshotTemplateId/recordTemplateApplyAll/snapshotTemplateApplyAll
*/
func (manager *Manager) UpdateNamespace(nsId string, ops []PatchOperation) (*NameSpace, error) {

	req := M{"operations": ops}
	var ret NameSpace
	err := manager.client.CallWithJson(context.Background(), &ret, "PATCH", manager.url("/namespaces/%s", nsId), nil, req)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
	删除空间API
*/
func (manager *Manager) DeleteNamespace(nsId string) error {

	return manager.client.Call(context.Background(), nil, "DELETE", manager.url("/namespaces/%s", nsId), nil)
}

/*
	列出空间API
*/
func (manager *Manager) ListNamespace(offset, line int, sortBy string) ([]NameSpace, int64, error) {

	ret := struct {
		Items []NameSpace `json:"items"`
		Total int64       `json:"total"`
	}{}

	query := url.Values{}
	setQuery(query, "offset", offset)
	setQuery(query, "line", line)
	setQuery(query, "sortBy", sortBy)

	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces?%v", query.Encode()), nil)
	return ret.Items, ret.Total, err
}

/*
	禁用空间API
*/
func (manager *Manager) DisableNamespace(nsId string) error {

	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/disabled", nsId), nil, nil)
}

/*
	启用空间API
*/
func (manager *Manager) EnableNamespace(nsId string) error {

	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/enabled", nsId), nil, nil)
}

/*
	添加域名API

	请求参数Body: 只需填入domain和type

	domainType支持四种类型 "publishRtmp":rtmp推流, "liveRtmp": rtmp播放, "liveHls": hls播放, "liveHdl": flv播放
*/
func (manager *Manager) AddDomain(nsId string, domainInfo *DomainInfo) error {

	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/domains", nsId), nil, domainInfo)
}

/*
	删除域名API
*/
func (manager *Manager) DeleteDomain(nsId string, domain string) error {

	return manager.client.Call(context.Background(), nil, "DELETE", manager.url("/namespaces/%s/domains/%s", nsId, domain), nil)
}

/*
	域名列表API
*/
func (manager *Manager) ListDomain(nsId string) ([]DomainInfo, error) {

	ret := struct {
		Items []DomainInfo `json:"items"`
	}{}

	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/domains", nsId), nil)
	return ret.Items, err
}

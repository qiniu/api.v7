package qvs

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
)

// 启动按需录制
func (manager *Manager) StartRecord(nsId, streamId string) error {
	err := manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/streams/%s/record/start", nsId, streamId), nil, nil)
	return err
}

// 停止按需录制
func (manager *Manager) StopRecord(nsId, streamId string) error {
	err := manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/streams/%s/record/stop", nsId, streamId), nil, nil)
	return err
}

// 删除录制片段
func (manager *Manager) DeleteStreamRecordHistories(nsId, streamId string, files []string) error {
	err := manager.client.CallWithJson(context.Background(), nil, "DELETE", manager.url("/namespaces/%s/streams/%s/recordhistories", nsId, streamId), nil, M{"files": files})
	return err
}

type saveasArgs struct {
	Fname     string `json:"fname"`
	Format   string `json:"format"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
	DeleteTs bool   `json:"deleteTs"` //sdk直接调用方式在不生成m3u8格式文件时是否删除对应的ts文件
	Pipeline  string `json:"pipeline"`
	NotifyUrl string `json:"notifyUrl"`
	DeleteAfterDays int  `json:"deleteAfterDays"`
}

type saveasReply struct {
	Fname       string `json:"fname"`
	PersistenId string `json:"persistentId,omitempty"`
	Bucket      string `json:"bucket"`
	Duration    int    `json:"duration"` // ms
}

// 录制视频片段合并
func (manager *Manager) RecordClipsSaveas(nsId, streamId string, arg *saveasArgs) (*saveasReply, error) {
	var ret saveasReply
	err := manager.client.CallWithJson(context.Background(), &ret, "POST", manager.url("/namespaces/%s/streams/%s/saveas", nsId, streamId), nil, arg)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

// 录制回放
func (manager *Manager) RecordsPlayback(nsId, streamId string, start, end int) (string, error) {
	var ret = struct {
		Url string `json:"url"`
	}{}
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/streams/%s/records/playback.m3u8?start=%d&end=%d", nsId, streamId, start, end), nil)
	if err != nil {
		return "", err
	}
	return ret.Url, nil
}

// 按需截图
func (manager *Manager) OndemandSnap(nsId, streamId string) error {
	err := manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/streams/%s/snap", nsId, streamId), nil, nil)
	return err
}

// 删除截图
func (manager *Manager) DeleteSnapshots(nsId, streamId string, files []string) error {
	err := manager.client.CallWithJson(context.Background(), nil, "DELETE", manager.url("/namespaces/%s/streams/%s/snapshots", nsId, streamId), nil, M{"files": files})
	return err
}

// 查询截图列表
func (manager *Manager) StreamsSnapshots(nsId string, streamId string, start, end int, qtype int, line int, marker string) ([]byte, error) {
	query := url.Values{}
	setQuery(query, "start", start)
	setQuery(query, "end", end)
	setQuery(query, "type", qtype)
	setQuery(query, "line", line)
	setQuery(query, "marker", marker)

	req, err := http.NewRequest("GET", manager.url("/namespaces/%s/streams/%s/snapshots?%v", nsId, streamId, query.Encode()), nil)
	if err != nil {
		return nil, err
	}
	resp, err := manager.client.Do(context.Background(), req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type RecordHistory struct {
	Url      string `json:"url"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Duration int    `json:"duration"`
	Format   int    `json:"format"`
	Snap     string `json:"snap"`
	File     string `json:"file"`
}

// 查询视频流的录制记录
func (manager *Manager) QueryStreamRecordHistories(nsId string, streamId string, start, end int, marker string, line int, format string) ([]RecordHistory, string, error) {
	ret := struct {
		Items  []RecordHistory `json:"items"`
		Marker string          `json:"marker"`
	}{}

	query := url.Values{}
	setQuery(query, "start", start)
	setQuery(query, "end", end)
	setQuery(query, "marker", marker)
	setQuery(query, "line", line)
	setQuery(query, "format", format)

	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/streams/%s/recordhistories?%v", nsId, streamId, query.Encode()), nil)
	return ret.Items, ret.Marker, err
}

// 查询流封面
func (manager *Manager) QueryStreamCover(nsId string, streamId string) (string, error) {
	var ret = struct {
		Url string `json:"url"`
	}{}
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/namespaces/%s/streams/%s/cover", nsId, streamId), nil)
	return ret.Url, err
}


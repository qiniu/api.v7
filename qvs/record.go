package qvs

import (
	"context"
)

// 启动按需录制
func (manager *Manager) StartRecord(nsId, streamId string) error {
	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/streams/%s/record/start", nsId, streamId), nil, nil)
}

// 停止按需录制
func (manager *Manager) StopRecord(nsId, streamId string) error {
	return manager.client.CallWithJson(context.Background(), nil, "POST", manager.url("/namespaces/%s/streams/%s/record/stop", nsId, streamId), nil, nil)
}

// 删除录制片段
func (manager *Manager) DeleteStreamRecordHistories(nsId, streamId string, files []string) error {
	return manager.client.CallWithJson(context.Background(), nil, "DELETE", manager.url("/namespaces/%s/streams/%s/recordhistories", nsId, streamId), nil, M{"files": files})
}

type saveasArgs struct {
	Fname           string `json:"fname"`
	Format          string `json:"format"`
	Start           int    `json:"start"`
	End             int    `json:"end"`
	DeleteTs        bool   `json:"deleteTs"` //sdk直接调用方式在不生成m3u8格式文件时是否删除对应的ts文件
	Pipeline        string `json:"pipeline"`
	NotifyUrl       string `json:"notifyUrl"`
	DeleteAfterDays int    `json:"deleteAfterDays"`
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

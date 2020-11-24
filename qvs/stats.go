package qvs

import (
	"context"
)

type FlowBadwidthData struct {
	Time []int64 `json:"time"`
	Data struct {
		Up   []int64 `json:"up"`
		Down []int64 `json:"down"`
	} `json:"data"`
}

/*
   查询流量数据
*/
func (manager *Manager) QueryFlow(nsId, streamId, tu string, start, end int) (*FlowBadwidthData, error) {
	var ret FlowBadwidthData
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/stats/flow?nsId=%s&streamId=%s&start=%d&end=%d&tu=%s", nsId, streamId, start, end, tu), nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
   查询带宽数据
*/
func (manager *Manager) QueryBandwidth(nsId, streamId, tu string, start, end int) (*FlowBadwidthData, error) {
	var ret FlowBadwidthData
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/stats/bandwidth?nsId=%s&stream=%s&start=%d&end=%d&tu=%s", nsId, streamId, start, end, tu), nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

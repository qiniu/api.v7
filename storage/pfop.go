package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/x/rpc.v7"
	"strings"
)

type OperationManager struct {
	client *rpc.Client
	mac    *qbox.Mac
	cfg    *Config
}

func NewOperationManager(mac *qbox.Mac, cfg *Config) *OperationManager {
	if cfg == nil {
		cfg = &Config{}
	}

	return &OperationManager{
		client: NewClient(mac, nil),
		mac:    mac,
		cfg:    cfg,
	}
}

// PfopResult pfop返回信息
type PfopRet struct {
	PersistentId string `json:"persistentId,omitempty"`
}

// FopRet 持久化云处理结果
type PrefopRet struct {
	Id          string `json:"id"`
	Code        int    `json:"code"`
	Desc        string `json:"desc"`
	InputBucket string `json:"inputBucket,omitempty"`
	InputKey    string `json:"inputKey,omitempty"`
	Pipeline    string `json:"pipeline,omitempty"`
	Reqid       string `json:"reqid,omitempty"`
	Items       []FopResult
}

func (this *PrefopRet) String() string {
	strData := fmt.Sprintf("Id: %s\r\nCode: %d\r\nDesc: %s\r\n", this.Id, this.Code, this.Desc)
	if this.InputBucket != "" {
		strData += fmt.Sprintln(fmt.Sprintf("InputBucket: %s", this.InputBucket))
	}
	if this.InputKey != "" {
		strData += fmt.Sprintln(fmt.Sprintf("InputKey: %s", this.InputKey))
	}
	if this.Pipeline != "" {
		strData += fmt.Sprintln(fmt.Sprintf("Pipeline: %s", this.Pipeline))
	}
	if this.Reqid != "" {
		strData += fmt.Sprintln(fmt.Sprintf("Reqid: %s", this.Reqid))
	}

	strData = fmt.Sprintln(strData)
	for _, item := range this.Items {
		strData += fmt.Sprintf("\tCmd:\t%s\r\n\tCode:\t%d\r\n\tDesc:\t%s\r\n", item.Cmd, item.Code, item.Desc)
		if item.Error != "" {
			strData += fmt.Sprintf("\tError:\t%s\r\n", item.Error)
		} else {
			if item.Hash != "" {
				strData += fmt.Sprintf("\tHash:\t%s\r\n", item.Hash)
			}
			if item.Key != "" {
				strData += fmt.Sprintf("\tKey:\t%s\r\n", item.Key)
			}
			if item.Keys != nil {
				if len(item.Keys) > 0 {
					strData += "\tKeys: {\r\n"
					for _, key := range item.Keys {
						strData += fmt.Sprintf("\t\t%s\r\n", key)
					}
					strData += "\t}\r\n"
				}
			}
		}
		strData += "\r\n"
	}
	return strData
}

// FopResult 云处理操作列表，包含每个云处理操作的状态信息
type FopResult struct {
	Cmd   string   `json:"cmd"`
	Code  int      `json:"code"`
	Desc  string   `json:"desc"`
	Error string   `json:"error,omitempty"`
	Hash  string   `json:"hash,omitempty"`
	Key   string   `json:"key,omitempty"`
	Keys  []string `json:"keys,omitempty"`
}

// Pfop 持久化数据处理
//
// @param bucket	资源空间
// @param key		源资源名
// @param fops		云处理操作列表，用`;``分隔，如:`avthumb/flv;saveas/xxx`，是将上传的视频文件转码成flv格式后存储为 bucket:key，
//                  其中 xxx 是 bucket:key 的URL安全的Base64编码结果。
// @param notifyURL	处理结果通知接收 URL，七牛将会向你设置的 URL 发起 Content-Type: application/json 的 POST 请求。
// @param pipeline	为空则表示使用公用队列，处理速度比较慢。建议指定私有队列，转码的时候使用独立的计算资源。
// @param force		强制执行数据处理。当服务端发现 fops 指定的数据处理结果已经存在，那就认为已经处理成功，避免重复处理浪费资源。
//                  本字段设为 `true`，则可强制执行数据处理并覆盖原结果。
//
func (m *OperationManager) Pfop(bucket, key, fops, pipeline, notifyURL string, force bool) (persistentId string, err error) {
	pfopParams := map[string][]string{
		"bucket": []string{bucket},
		"key":    []string{key},
		"fops":   []string{fops},
	}

	if pipeline != "" {
		pfopParams["pipeline"] = []string{pipeline}
	}

	if notifyURL != "" {
		pfopParams["notifyURL"] = []string{notifyURL}
	}

	if force {
		pfopParams["force"] = []string{"1"}
	}
	var ret PfopRet
	ctx := context.TODO()
	reqHost, reqErr := m.ApiHost(bucket)
	if reqErr != nil {
		err = reqErr
		return
	}
	reqUrl := fmt.Sprintf("%s/pfop/", reqHost)
	err = m.client.CallWithForm(ctx, &ret, "POST", reqUrl, pfopParams)
	if err != nil {
		return
	}

	persistentId = ret.PersistentId
	return
}

// Prefop 持久化处理状态查询
func (m *OperationManager) Prefop(persistentId string) (ret PrefopRet, err error) {
	ctx := context.TODO()
	reqHost, reqErr := m.prefopApiHost(persistentId)
	if reqErr != nil {
		err = reqErr
		return
	}

	reqUrl := fmt.Sprintf("%s/status/get/prefop?id=%s", reqHost, persistentId)
	err = m.client.Call(ctx, &ret, "GET", reqUrl)
	return
}

// 获取资源管理域名
func (m *OperationManager) ApiHost(bucket string) (apiHost string, err error) {
	zone, zoneErr := GetZone(m.mac.AccessKey, bucket)
	if zoneErr != nil {
		err = zoneErr
		return
	}

	if m.cfg.UseHttps {
		apiHost = fmt.Sprintf("https://%s", zone.ApiHost)
	} else {
		apiHost = fmt.Sprintf("http://%s", zone.ApiHost)
	}
	return
}

func (m *OperationManager) prefopApiHost(persistentId string) (apiHost string, err error) {
	if strings.Contains(persistentId, "z1.") {
		apiHost = Zone_z1.ApiHost
	} else if strings.Contains(persistentId, "z2.") {
		apiHost = Zone_z2.ApiHost
	} else if strings.Contains(persistentId, "na0.") {
		apiHost = Zone_na0.ApiHost
	} else if strings.Contains(persistentId, "z0.") {
		apiHost = DefaultApiHost
	} else {
		err = errors.New("invalid persistent id")
	}

	if m.cfg.UseHttps {
		apiHost = fmt.Sprintf("https://%s", apiHost)
	} else {
		apiHost = fmt.Sprintf("http://%s", apiHost)
	}

	return

}

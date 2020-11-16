# QVS Cloud Server-Side Library for Go

## Features

- 空间管理
    - [x] 创建空间: AddNamespace(ns *NameSpace)
	- [x] 删除空间: DeleteNamespace(nsId string)
	- [x] 更新空间: UpdateNamespace(nsId string, ops []PatchOperation)
	- [x] 查询空间信息: QueryNamespace(nsId string)
	- [x] 获取空间列表: ListNamespace(offset, line int, sortBy string)
	- [x] 禁用空间: DisableNamespace(nsId string)
    - [x] 启用空间: EnableNamespace(nsId string)
    
- 流管理
    - [x] 创建流: AddStream(nsId string, stream *Stream)
    - [x] 删除流: DeleteStream(nsId string, streamId string)
    - [x] 查询流信息: QueryStream(nsId string, streamId string)
    - [x] 更新流: UpdateStream(nsId string, streamId string, ops []PatchOperation)
    - [x] 获取流列表: ListStream(nsId string, offset, line int, prefix, sortBy string, qType int)
    - [x] 获取流地址
        - [x] 动态模式: DynamicPublishPlayURL(nsId string, streamId string, route *DynamicLiveRoute)
        - [x] 静态模式: StaticPublishPlayURL(nsId string, streamId string, route *StaticLiveRoute)
    - [x] 禁用流: DisableStream(nsId string, streamId string)  
    - [x] 启用流: EnableStream(nsId string, streamId string) 
    - [x] 查询推流记录: QueryStreamPubhistories(nsId string, streamId string, start, end int, line, offset int)  

- 模板管理
    - [x] 创建模板: AddTemplate(tmpl *Template)
    - [x] 删除模板: DeleteTemplate(templId string)
    - [x] 更新模板: UpdateTemplate(templId string, ops []PatchOperation)
    - [x] 查询模板信息: QueryTemplate(templId string)
    - [x] 获取模板列表: ListTemplate(offset, line int, sortBy string, templateType int, match string)

- 录制管理相关接口
    - [x] 查询录制记录: QueryStreamRecordHistories(nsId string, streamId string, start, end int, marker string, line int)
    - [x] 获取截图列表: StreamsSnapshots(nsId string, streamId string, start, end int, qtype int, limit int, marker string)
    - [x] 获取直播封面截图: QueryStreamCover(nsId string, streamId string)


## Usage

```go
package main

import (
	// ...
	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/qvs"
)

var (
	AccessKey = "<QINIU ACCESS KEY>" // 替换成自己 Qiniu 账号的 AccessKey.
	SecretKey = "<QINIU SECRET KEY>" // 替换成自己 Qiniu 账号的 SecretKey.
)

func main() {
	// ...
	mac := auth.New(AccessKey, SecretKey)
	manager := qvs.NewManager(mac, nil)
	// ...
}
```


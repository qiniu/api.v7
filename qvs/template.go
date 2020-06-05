package qvs

import (
	"context"
	"net/url"
)

type Template struct {
	ID               string `json:"id"`
	Name             string `json:"name"`             // 模版名称，格式为 4 ~ 100个字符，可包含小写字母、数字、中划线、汉字)
	Desc             string `json:"desc,omitempty"`   // 模版描述
	Bucket           string `json:"bucket"`           // 模版对应的对象存储的bucket
	DeleteAfterDays  int    `json:"deleteAfterDays"`  // 存储过期时间,默认永久不过期
	TemplateType     int    `json:"templateType"`     // 模板类型,取值：0（录制模版）, 1（截图模版）
	FileType         int    `json:"fileType"`         // 文件存储类型,取值：0（普通存储）,1（低频存储）
	RecordType       int    `json:"recordType"`       // 录制模式, 0（不录制）,1（实时录制）, 2（按需录制）
	RecordFileFormat int    `json:"recordFileFormat"` // 录制文件存储格式,取值：0（ts格式存储）

	//record/ts/${namespaceId}/${streamId}/${startMs}-${endMs}.ts
	TSFileNameTemplate string `json:"tsFileNameTemplate"`
	//record/snap/${namespaceId}/${streamId}/${startMs}.jpg // 录制封面
	RecordSnapFileNameFmt string `json:"recordSnapFileNameTemplate"`
	RecordInterval        int    `json:"recordInterval"` //录制文件长度

	M3u8FileNameTemplate string `json:"m3u8FileNameTemplate,omitempty"` // m3u8文件命名格式

	JpgOverwriteStatus bool `json:"jpgOverwriteStatus"` // 开启覆盖式截图(一般用于流封面)
	JpgSequenceStatus  bool `json:"jpgSequenceStatus"`  // 开启序列式截图
	JpgOnDemandStatus  bool `json:"jpgOnDemandStatus"`  // 开启按需截图

	// 覆盖式截图文件命名格式:snapshot/jpg/${namespaceId}/${streamId}/${streamId}.jpg
	JpgOverwriteFileNameTemplate string `json:"jpgOverwriteFileNameTemplate"`
	// 序列式截图文件命名格式:snapshot/jpg/${namespaceId}/${streamId}/${startMs}.jpg
	JpgSequenceFileNameTemplate string `json:"jpgSequenceFileNameTemplate"`
	// 按需式截图文件命名格式:snapshot/jpg/${namespaceId}/${streamId}/ondemand-${startMs}.jpg
	JpgOnDemandFileNameTemplate string `json:"jpgOnDemandFileNameTemplate"`
	SnapInterval                int    `json:"snapInterval"` // 序列式截图时间间隔

	CreatedAt int64  `json:"createdAt,omitempty"` // 模板创建时间
	UpdatedAt int64  `json:"updatedAt,omitempty"` // 模板更新时间
	Zone      string `json:"zone"`                // zone为服务区域配置项，可选项为z0, z1, z2,默认为z0. z0表示华东, z1表示华北、z2表示华南
}

/*
	创建模板API
*/
func (manager *Manager) AddTemplate(tmpl *Template) (*Template, error) {

	var ret Template
	err := manager.client.CallWithJson(context.Background(), &ret, "POST", manager.url("/templates"), nil, tmpl)
	if err != nil {
		return nil, err
	}
	return &ret, err
}

/*
	查询模板信息API
*/
func (manager *Manager) QueryTemplate(templId string) (*Template, error) {

	var ret Template
	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/templates/%s", templId), nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
	修改模板API
*/
func (manager *Manager) UpdateTemplate(templId string, ops []PatchOperation) (*Template, error) {

	req := M{"operations": ops}
	var ret Template
	err := manager.client.CallWithJson(context.Background(), &ret, "PATCH", manager.url("/templates/%s", templId), nil, req)
	if err != nil {
		return nil, err
	}
	return &ret, err
}

/*
	删除模板API
*/
func (manager *Manager) DeleteTemplate(templId string) error {

	return manager.client.Call(context.Background(), nil, "DELETE", manager.url("/templates/%s", templId), nil)
}

/*
	获取模版列表API
*/
func (manager *Manager) ListTemplate(offset, line int, sortBy string, templateType int, match string) ([]Template, int64, error) {

	ret := struct {
		Items []Template `json:"items"`
		Total int64      `json:"total"`
	}{}

	query := url.Values{}
	setQuery(query, "offset", offset)
	setQuery(query, "line", line)
	setQuery(query, "sortBy", sortBy)
	setQuery(query, "match", match)
	setQuery(query, "templateType", templateType)

	err := manager.client.Call(context.Background(), &ret, "GET", manager.url("/templates?%v", query.Encode()), nil)
	return ret.Items, ret.Total, err
}

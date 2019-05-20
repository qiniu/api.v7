package sms

import (
	"fmt"
	"net/url"
)

// Template 模板
type Template struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         TemplateType `json:"type"`
	Template     string       `json:"template"`
	Description  string       `json:"description"`
	AuditStatus  AuditStatus  `json:"audit_status"`
	RejectReason string       `json:"reject_reason"`

	UpdatedAt uint64 `json:"updated_at"`
	CreatedAt uint64 `json:"created_at"`
}

// TemplatePagination 模板分页
type TemplatePagination struct {
	Page     int        `json:"page"`      // 页码，默认为 1
	PageSize int        `json:"page_size"` // 分页大小，默认为 20
	Total    int        `json:"total"`     // 总记录条数
	Items    []Template `json:"items"`     // 模板
}

// TemplateRequest 创建模板请求参数
type TemplateRequest struct {
	UID         uint32       `json:"uid"`
	Name        string       `json:"name"`
	Type        TemplateType `json:"type"`
	Template    string       `json:"template"`
	Description string       `json:"description"`
}

// TemplateResponse 模板响应
type TemplateResponse struct {
	TemplateID string `json:"template_id"`
}

// CreateTemplate 创建模板
func (m *Manager) CreateTemplate(args TemplateRequest) (ret TemplateResponse, err error) {
	url := fmt.Sprintf("%s%s", Host, "/v1/template")
	err = m.client.CallWithJSON(&ret, url, args)
	return
}

// UpdateTemplate 更新模板
func (m *Manager) UpdateTemplate(id string, args TemplateRequest) (err error) {
	url := fmt.Sprintf("%s%s/%s", Host, "/v1/template", id)
	_, err = m.client.PutWithJSON(url, args)
	return
}

// QueryTemplateRequest 查询模板参数
type QueryTemplateRequest struct {
	AuditStatus AuditStatus `json:"audit_status"` // 审核状态
	Page        int         `json:"page"`         // 页码，默认为 1
	PageSize    int         `json:"page_size"`    // 分页大小，默认为 20
}

// QueryTemplate 查询模板
func (m *Manager) QueryTemplate(args QueryTemplateRequest) (pagination TemplatePagination, err error) {
	values := url.Values{}

	if args.AuditStatus.IsValid() {
		values.Set("audit_status", args.AuditStatus.String())
	}

	if args.Page > 0 {
		values.Set("page", fmt.Sprintf("%d", args.Page))
	}

	if args.PageSize > 0 {
		values.Set("page_size", fmt.Sprintf("%d", args.PageSize))
	}

	url := fmt.Sprintf("%s%s?%s", Host, "/v1/template", values.Encode())
	err = m.client.GetCall(&pagination, url)
	return
}

// DeleteTemplate 删除模板
func (m *Manager) DeleteTemplate(id string) (err error) {
	url := fmt.Sprintf("%s%s/%s", Host, "/v1/template", id)
	_, err = m.client.Delete(url)
	return
}

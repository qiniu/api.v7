package sms

import (
	"fmt"
	"net/url"
)

// Signature 签名
type Signature struct {
	ID           string           `json:"id"`
	Source       SignatureSrcType `json:"source"`
	Signature    string           `json:"signature"`
	AuditStatus  AuditStatus      `json:"audit_status"`
	RejectReason string           `json:"reject_reason,omitempty"`
	Description  string           `json:"description"`

	UpdatedAt uint64 `json:"updated_at"`
	CreatedAt uint64 `json:"created_at"`
}

// SignaturePagination 签名分页
type SignaturePagination struct {
	Page     int         `json:"page"`      // 页码，默认为 1
	PageSize int         `json:"page_size"` // 分页大小，默认为 20
	Total    int         `json:"total"`     // 总记录条数
	Items    []Signature `json:"items"`     // 签名
}

// CreateSignatureRequest 创建签名请求参数
type CreateSignatureRequest struct {
	Signature   string           `json:"signature"`
	Source      SignatureSrcType `json:"source"`
	Pic         string           `json:"pic"`
	Description string           `json:"decription"`
}

// CreateSignatureResponse 签名响应
type CreateSignatureResponse struct {
	SignatureID string `json:"signature_id"`
}

// CreateSignature 创建签名
func (m *Manager) CreateSignature(args CreateSignatureRequest) (ret CreateSignatureResponse, err error) {
	url := fmt.Sprintf("%s%s", Host, "/v1/signature")
	err = m.client.CallWithJSON(&ret, url, args)
	return
}

// UpdateSignature 更新签名
func (m *Manager) UpdateSignature(args CreateSignatureRequest) (ret CreateSignatureResponse, err error) {
	// url := fmt.Sprintf("%s%s", Host, "/v1/signature")
	// err = m.client.PutCallWithJSON()
	return
}

// QuerySignatureRequest 查询签名参数
type QuerySignatureRequest struct {
	AuditStatus AuditStatus `json:"audit_status"` // 审核状态
	Page        int         `json:"page"`         // 页码，默认为 1
	PageSize    int         `json:"page_size"`    // 分页大小，默认为 20
}

// QuerySignature 查询签名
func (m *Manager) QuerySignature(args QuerySignatureRequest) (pagination SignaturePagination, err error) {
	value := url.Values{}

	if args.AuditStatus.IsValid() {
		value.Set("audit_status", args.AuditStatus.String())
	}

	if args.Page > 0 {
		value.Set("page", fmt.Sprintf("%d", args.Page))
	}

	if args.PageSize > 0 {
		value.Set("page_size", fmt.Sprintf("%d", args.PageSize))
	}

	_url := fmt.Sprintf("%s%s?%s", Host, "/v1/signature", value.Encode())
	err = m.client.GetCall(&pagination, _url)
	return
}

// DeleteSignature 删除签名
func (m *Manager) DeleteSignature(args CreateSignatureRequest) (ret CreateSignatureResponse, err error) {

	return
}

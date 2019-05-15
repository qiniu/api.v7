package sms

import (
	"fmt"
	"net/http"

	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/sms/client"
	"github.com/qiniu/api.v7/sms/rpc"
)

var (
	// Host 为 Qiniu SMS Server API 服务域名
	Host = "https://sms.qiniuapi.com"
)

// SignatureSrcType 签名类型
type SignatureSrcType string

const (
	// EnterprisesAndInstitutions 企事业单位的全称或简称
	EnterprisesAndInstitutions SignatureSrcType = "enterprises_and_institutions"

	// Website 工信部备案网站的全称或简称
	Website SignatureSrcType = "website"

	// APP APP应用的全称或简称
	APP SignatureSrcType = "app"

	// PublicNumberOrSmallProgram 公众号或小程序的全称或简称
	PublicNumberOrSmallProgram SignatureSrcType = "public_number_or_small_program"

	// StoreName 电商平台店铺名的全称或简称
	StoreName SignatureSrcType = "store_name"

	// TradeName 商标名的全称或简称
	TradeName SignatureSrcType = "trade_name"
)

// Manager 提供了 Qiniu SMS Server API 相关功能
type Manager struct {
	mac    *auth.Credentials
	client rpc.Client
}

// Error 统一错误
type Error struct {
	Code      string
	Message   string
	RequestID string
}

// NewManager 用来构建一个新的 Manager
func NewManager(mac *auth.Credentials) (manager *Manager) {

	manager = &Manager{}

	mac1 := &client.Mac{
		AccessKey: mac.AccessKey,
		SecretKey: []byte(mac.SecretKey),
	}

	transport := client.NewTransport(mac1, nil)
	manager.client = rpc.Client{Client: &http.Client{Transport: transport}}

	return
}

// CreateSignatureRequest 创建签名请求参数
type CreateSignatureRequest struct {
	Signature   string           `json:"signature"`
	Source      SignatureSrcType `json:"source"`
	Pic         string           `json:"pic"`
	Description string           `json:"decription"`
}

// SignatureRet 签名响应
type SignatureRet struct {
	SignatureID string `json:"signature_id"`
}

// CreateSignature 创建签名
func (m *Manager) CreateSignature(log rpc.Logger, args CreateSignatureRequest) (ret SignatureRet, err error) {
	url := fmt.Sprintf("%s%s", Host, "/v1/signature")
	m.client.CallWithJSON(log, &ret, url, args)
	return
}

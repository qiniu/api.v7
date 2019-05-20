package sms

import "fmt"

// MessagesRequest 短信消息
type MessagesRequest struct {
	SignatureID string                 `json:"signature_id"`
	TemplateID  string                 `json:"template_id"`
	Mobiles     []string               `json:"mobiles"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// MessagesResponse 发送短信响应
type MessagesResponse struct {
	JobID string `json:"job_id"`
}

// SendMessage 发送短信
func (m *Manager) SendMessage(args MessagesRequest) (ret MessagesResponse, err error) {
	url := fmt.Sprintf("%s%s", Host, "/v1/message")
	err = m.client.CallWithJSON(&ret, url, args)
	return
}

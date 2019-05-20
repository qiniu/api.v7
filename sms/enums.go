package sms

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

// IsValid 验证
func (s SignatureSrcType) IsValid() bool {
	switch s {
	case EnterprisesAndInstitutions, Website, APP, PublicNumberOrSmallProgram, StoreName, TradeName:
		return true
	}
	return false
}

func (s SignatureSrcType) String() string {
	return string(s)
}

// TemplateType 模版类型
type TemplateType string

const (
	// NotificationType 通知类短信
	NotificationType TemplateType = "notification"

	// VerificationType 验证码短信
	VerificationType TemplateType = "verification"

	// MarketingType 营销类短信
	MarketingType TemplateType = "marketing"

	// VoiceType 语音短信
	VoiceType TemplateType = "voice"
)

// IsValid 是否合法Template类型
func (t TemplateType) IsValid() bool {
	switch t {
	case NotificationType, VerificationType, MarketingType, VoiceType:
		return true
	}
	return false
}

// String to string
func (t TemplateType) String() string {
	return string(t)
}

// AuditStatus 审核状态
type AuditStatus string

const (
	// AuditStatusPassed 通过
	AuditStatusPassed AuditStatus = "passed"

	// AuditStatusReject 未通过
	AuditStatusReject AuditStatus = "rejected"

	// AuditStatusReviewing 审核中
	AuditStatusReviewing AuditStatus = "reviewing"
)

// IsValid 验证
func (a AuditStatus) IsValid() bool {
	switch a {
	case AuditStatusPassed, AuditStatusReject, AuditStatusReviewing:
		return true
	}
	return false
}

func (a AuditStatus) String() string {
	return string(a)
}

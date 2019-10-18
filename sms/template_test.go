package sms_test

import (
	"testing"

	"github.com/qiniu/api.v7/v7/sms"
)

func TestTemplate(t *testing.T) {

	// CreateTemplate
	args := sms.TemplateRequest{
		Name:     "Test",
		Type:     sms.VerificationType,
		Template: "您的验证码是 ${{code}}， 5分钟内有效",
	}

	ret, err := manager.CreateTemplate(args)

	if err != nil {
		t.Fatalf("CreateTemplate() error: %v\n", err)
	}

	if len(ret.TemplateID) == 0 {
		t.Fatal("CreateTemplate() error: The template ID cannot be empty")
	}

	// QueryTemplate
	query := sms.QueryTemplateRequest{}

	pagination, err := manager.QueryTemplate(query)

	if err != nil {
		t.Fatalf("QueryTemplate() error: %v\n", err)
	}

	if len(pagination.Items) == 0 {
		t.Fatal("QueryTemplate() error: templates cannot be empty")
	}

	if pagination.Total == 0 {
		t.Fatal("QueryTemplate() error: total cannot be 0")
	}

	// UpdateTemplate
	update := sms.TemplateRequest{
		Name: "test",
	}

	err = manager.UpdateTemplate(ret.TemplateID, update)
	if err != nil {
		t.Fatalf("UpdateTemplate() error: %v\n", err)
	}

	// DeleteTemplate
	err = manager.DeleteTemplate(ret.TemplateID)
	if err != nil {
		t.Fatalf("DeleteTemplate() error: %v\n", err)
	}
}

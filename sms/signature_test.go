package sms_test

import (
	"testing"

	"github.com/qiniu/api.v7/v7/sms"
)

func TestSignature(t *testing.T) {

	// CreateSignature
	args := sms.SignatureRequest{
		Signature: "Test",
		Source:    sms.Website,
	}

	ret, err := manager.CreateSignature(args)

	if err != nil {
		t.Fatalf("CreateSignature() error: %v\n", err)
	}

	if len(ret.SignatureID) == 0 {
		t.Fatal("CreateSignature() error: The signature ID cannot be empty")
	}

	// QuerySignature
	query := sms.QuerySignatureRequest{}

	pagination, err := manager.QuerySignature(query)

	if err != nil {
		t.Fatalf("QuerySignature() error: %v\n", err)
	}

	if len(pagination.Items) == 0 {
		t.Fatal("QuerySignature() error: signatures cannot be empty")
	}

	if pagination.Total == 0 {
		t.Fatal("QuerySignature() error: total cannot be 0")
	}

	// UpdateSignature
	update := sms.SignatureRequest{
		Signature: "test",
	}

	err = manager.UpdateSignature(ret.SignatureID, update)
	if err != nil {
		t.Fatalf("UpdateSignature() error: %v\n", err)
	}

	// DeleteSignature
	err = manager.DeleteSignature(ret.SignatureID)
	if err != nil {
		t.Fatalf("DeleteSignature() error: %v\n", err)
	}

}

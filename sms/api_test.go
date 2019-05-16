package sms_test

import (
	"os"
	"testing"

	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/sms"
)

var manager *sms.Manager

func init() {
	accessKey := os.Getenv("accessKey")
	secretKey := os.Getenv("secretKey")

	mac := auth.New(accessKey, secretKey)
	manager = sms.NewManager(mac)
}

func TestCreateSignature(t *testing.T) {
	args := sms.CreateSignatureRequest{
		Signature: "QVM-ZW",
		Source:    sms.Website,
	}

	ret, err := manager.CreateSignature(args)

	if err != nil {
		t.Fatalf("CreateSignature() error: %v\n", err)
	}

	if len(ret.SignatureID) == 0 {
		t.Fatal("CreateSignature() error: The signature ID cannot be empty")
	}
}

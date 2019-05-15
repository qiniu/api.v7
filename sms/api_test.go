package sms_test

import (
	"os"
	"testing"

	"github.com/qiniu/api.v7/auth"
	"github.com/qiniu/api.v7/sms"
)

type Logger struct {
}

func (log Logger) ReqID() string {
	return "req_id"
}

func (log Logger) Xput(logs []string) {

	return
}

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

	ret, err := manager.CreateSignature(Logger{}, args)

	if err != nil {
		t.Fatalf("CreateSignature() error: %v\n", err)
	}

	if len(ret.SignatureID) == 0 {
		t.Fatal("CreateSignature() error: The signature ID cannot be empty")
	}
}

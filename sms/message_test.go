package sms_test

// import (
// 	"testing"

// 	"github.com/qiniu/api.v7/v7/sms"
// )

// func TestMessage(t *testing.T) {

// 	// SendMessage
// 	args := sms.MessagesRequest{
// 		SignatureID: "",
// 		TemplateID:  "",
// 		Mobiles:     []string{""},
// 		Parameters: map[string]interface{}{
// 			"code": 123456,
// 		},
// 	}

// 	ret, err := manager.SendMessage(args)

// 	if err != nil {
// 		t.Fatalf("SendMessage() error: %v\n", err)
// 	}

// 	if len(ret.JobID) == 0 {
// 		t.Fatal("SendMessage() error: The job id cannot be empty")
// 	}
// }

package qbox

import (
	"net/http"
	"testing"
)

var mac *Mac

func init() {
	mac = NewMac("ak", "sk")
}

func TestMacSign(t *testing.T) {
	testStrs := []struct {
		Data   string
		Signed string
	}{
		{Data: "hello", Signed: "ak:MzQzMzdjNzBjZDJiYzI4YjMxODQ3MjdhNDAwNzA4ZWUyNmE1YWY0OA=="},
		{Data: "world", Signed: "ak:YzE5ZmFjYzM4ZDVhY2ExZGNmMTQzOTkwMDNlMGY3YTNiNzgxMjQ4Ng=="},
		{Data: "-test", Signed: "ak:YTA5ZTdkYjE5NmFjODk2NGFhMmY1YTNiYmEwNjZjZTRlMjI3MTBhZQ=="},
		{Data: "ba#a-", Signed: "ak:YjZhMWNiZjE1ZDgxNmNkMjM0NzU1MGQ3YjJmYjllNjY5ZmY2NDI3Mg=="},
	}

	for _, b := range testStrs {
		got := mac.Sign([]byte(b.Data))
		if got != b.Signed {
			t.Errorf("Sign %q, want=%q, got=%q\n", b.Data, b.Signed, got)
		}
	}
}

package qbox

import (
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
		{Data: "hello", Signed: "ak:NDN8cM0rwosxhHJ6QAcI7ialr0g="},
		{Data: "world", Signed: "ak:wZ-sw41ayh3PFDmQA-D3o7eBJIY="},
		{Data: "-test", Signed: "ak:oJ59sZasiWSqL1o7ugZs5OInEK4="},
		{Data: "ba#a-", Signed: "ak:tqHL8V2BbNI0dVDXsvueZp_2QnI="},
	}
	for _, b := range testStrs {
		got := mac.Sign([]byte(b.Data))
		if got != b.Signed {
			t.Errorf("Sign %q, want=%q, got=%q\n", b.Data, b.Signed, got)
		}
	}
}

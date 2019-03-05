package log

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	b := bytes.Buffer{}
	logger := New(&b, InfoPrefix, log.LstdFlags, LogInfo)

	logger.Info("hello world")

	splits := strings.Split(b.String(), " ")
	if splits[0] != strings.Trim(InfoPrefix, " ") {
		t.Errorf("got prefix: %q, want: %q\n", splits[0], InfoPrefix)
	}
}

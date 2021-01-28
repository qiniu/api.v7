package deprecated

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

var inited = false

func init() {
	if value := os.Getenv("SUPPRESS_DEPRECATION_WARNING"); value != "" {
		inited = true
	}
}

func PrintDeprecationWarning() {
	if !inited {
		inited = true

		isStderrTerminal := false
		if fileInfo, err := os.Stderr.Stat(); err == nil {
			if fileInfo.Mode()&os.ModeCharDevice != 0 {
				isStderrTerminal = true
			}
		}
		content := "[DEPRECATION WARNING] github.com/qiniu/api.v7 is deprecated. And the repository is moved to github.com/qiniu/api.v7 now"
		if isStderrTerminal {
			fmt.Fprintf(os.Stderr, "%s\n", color.Warn.Render(content))
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", content)
		}
	}
}

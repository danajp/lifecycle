package style

import (
	"fmt"

	"github.com/fatih/color"
)

var Symbol = func(format string, a ...interface{}) string {
	if color.NoColor {
		format = fmt.Sprintf("'%s'", format)
	}
	return Key(format, a...)
}

var Key = color.HiBlueString

var Warn = color.New(color.FgYellow, color.Bold).SprintfFunc()

var Error = color.New(color.FgRed, color.Bold).SprintfFunc()
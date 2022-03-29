package xgo

import (
	"regexp"
	"runtime"
)

var LineBreak = "\n"

// Variable regexp pattern: %(variable)s
var VarPattern = regexp.MustCompile(`%\(([^)]+)\)s`)

func init() {
	if runtime.GOOS == "windows" {
		LineBreak = "\r\n"
	}
}

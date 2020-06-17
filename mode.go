package xgo

import (
	"io"
	"os"
)

const (
	// DebugMode indicates gin mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates gin mode is release.
	//ReleaseMode = "release"
	// TestMode indicates gin mode is test.
	//TestMode = "test"
)

const (
	debugCode = iota // 0
	//releaseCode			// 1
	//testCode			// 2
)

var xGoMode = debugCode
var modeName = DebugMode

var DefaultWriter io.Writer = os.Stdout

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

// DefaultWriter is the default io.Writer used by xgo for debug output and
// middleware output like Logger() or Recovery().
// Note that both Logger and Recovery provides custom ways to configure their
// output io.Writer.
// To support coloring in Windows use:
// 		import "github.com/mattn/go-colorable"
// 		xgo.DefaultWriter = colorable.NewColorableStdout()
// 可用于将日志输出到文件
var DefaultWriter io.Writer = os.Stdout

// DefaultErrorWriter is the default io.Writer used by Gin to debug errors
var DefaultErrorWriter io.Writer = os.Stderr

//go:build !windows && !plan9 && !appengine && !wasm
// +build !windows,!plan9,!appengine,!wasm

package term

import (
	"golang.org/x/sys/unix"
)

// 获得终端的列数.
func getTerminalColumns() int {
	ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
	if err != nil {
		return 80
	}
	return int(ws.Col)
}

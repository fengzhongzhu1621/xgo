//go:build !appengine && !js && windows
// +build !appengine,!js,windows

package term

import (
	"io"
	"os"

	"golang.org/x/sys/windows"
)

func CheckIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		handle := windows.Handle(v.Fd())
		var mode uint32
		if err := windows.GetConsoleMode(handle, &mode); err != nil {
			return false
		}
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		if err := windows.SetConsoleMode(handle, mode); err != nil {
			return false
		}
		return true
	}
	return false
}
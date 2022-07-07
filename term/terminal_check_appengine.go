//go:build appengine
// +build appengine

package term

import (
	"io"
)

func CheckIfTerminal(w io.Writer) bool {
	return true
}

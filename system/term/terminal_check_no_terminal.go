//go:build js || nacl || plan9
// +build js nacl plan9

package term

import (
	"io"
)

func CheckIfTerminal(w io.Writer) bool {
	return false
}

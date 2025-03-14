//go:build js
// +build js

package term

func isTerminal(fd int) bool {
	return false
}

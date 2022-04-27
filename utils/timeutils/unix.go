//go:build freebsd || openbsd || netbsd || dragonfly || darwin
// +build freebsd openbsd netbsd dragonfly darwin

package timeutils

import (
	"time"

	"golang.org/x/sys/unix"
)

// DurationToTimespec prepares a timeout value.
func DurationToTimespec(d time.Duration) unix.Timespec {
	return unix.NsecToTimespec(d.Nanoseconds())
}

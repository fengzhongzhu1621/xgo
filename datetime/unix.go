//go:build freebsd || openbsd || netbsd || dragonfly || darwin
// +build freebsd openbsd netbsd dragonfly darwin

package datetime

import (
	"time"

	"golang.org/x/sys/unix"
)

// DurationToTimespec prepares a timeout value.
func DurationToTimespec(d time.Duration) unix.Timespec {
	return unix.NsecToTimespec(d.Nanoseconds())
}

// MaxTime from http://stackoverflow.com/questions/25065055#32620397
// This is a long time in the future. It's an imaginary number that is
// used for b-tree ordering.
var MaxTime = time.Unix(1<<63-62135596801, 999999999)

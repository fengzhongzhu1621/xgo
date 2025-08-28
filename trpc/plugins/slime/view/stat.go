package view

import "time"

// Stat defines a stat each retry/hedging must have.
type IStat interface {
	Cost() time.Duration
	Attempts() []IAttempt
	Throttled() bool
	InflightN() int
	Error() error
}

// Attempt defines the stat of each retry/hedging attempt.
type IAttempt interface {
	Start() time.Time
	End() time.Time
	Error() error
	Inflight() bool
	NoMoreAttempt() bool
}

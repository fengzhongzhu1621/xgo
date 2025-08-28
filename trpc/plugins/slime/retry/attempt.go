package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/fengzhongzhu1621/xgo/trpc/plugins/slime/pushback"
	"trpc.group/trpc-go/trpc-go/codec"
)

const (
	timeFormat = "15:04:05.000"
)

// attempt preserves the info the each attempt.
type attempt struct {
	*impl

	ctx context.Context
	rsp interface{}
	err error

	attempt       int
	start, end    time.Time
	pushbackDelay *time.Duration
}

// SyncStart starts the attempt synchronously.
func (a *attempt) SyncStart() {
	a.start = time.Now()
	a.err = a.handler(a.ctx, a.req, a.rsp)
	a.end = time.Now()

	a.pushbackDelay = pushback.FromMsg(codec.Message(a.ctx))
}

// OnReturn just report to throttle.
func (a *attempt) OnReturn() {
	a.ackThrottle()
}

// ackThrottle ack the throttle with success or failure.
func (a *attempt) ackThrottle() {
	if !a.noRetry() {
		if a.err == nil {
			a.throttle.OnSuccess()
			return
		}
		if !a.isRetryableErr() {
			return
		}
	}
	a.throttle.OnFailure()
}

func (a *attempt) isRetryableErr() bool {
	return a.impl.isRetryableErr(a.err)
}

func (a *attempt) noRetry() bool {
	return a.pushbackDelay != nil && *a.pushbackDelay < 0
}

// Start implements view.Attempt.
func (a *attempt) Start() time.Time {
	return a.start
}

// End implements view.Attempt.
func (a *attempt) End() time.Time {
	return a.end
}

// Error implements view.Attempt.
func (a *attempt) Error() error {
	return a.err
}

// Inflight implements view.Attempt.
func (a *attempt) Inflight() bool {
	return false
}

// NoMoreAttempt implements view.NoMoreAttempt.
func (a *attempt) NoMoreAttempt() bool {
	if a.pushbackDelay == nil {
		return false
	}
	return *a.pushbackDelay < 0
}

// String implements fmt.Stringer.
func (a *attempt) String() string {
	if a.pushbackDelay == nil {
		return fmt.Sprintf("%dth attempt, start: %v, end: %v, pushbackDelay: nil, err: %v",
			a.attempt, a.start.Format(timeFormat), a.end.Format(timeFormat), a.err)
	}
	if *a.pushbackDelay < 0 {
		return fmt.Sprintf("%dth attempt, start: %v, end: %v, pushbackDelay: no_retry, err: %v",
			a.attempt, a.start.Format(timeFormat), a.end.Format(timeFormat), a.err)
	}
	return fmt.Sprintf("%dth attempt, start: %v, end: %v, pushbackDelay: %v, err: %v",
		a.attempt, a.start.Format(timeFormat), a.end.Format(timeFormat), a.pushbackDelay, a.err)
}

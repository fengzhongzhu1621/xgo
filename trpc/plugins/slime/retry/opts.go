package retry

import (
	"errors"
	"fmt"
	"time"

	"github.com/fengzhongzhu1621/xgo/collections/backoff"
	lazylog "github.com/fengzhongzhu1621/xgo/logging/lazy_log"
	"github.com/fengzhongzhu1621/xgo/trpc/plugins/slime/view"
	"github.com/fengzhongzhu1621/xgo/trpc/plugins/slime/view/metrics"
	tlog "trpc.group/trpc-go/trpc-go/log"
)

// Opt is the option function to modify Retry.
type Opt func(*Retry) error

// WithRetryableErr allows user to register an additional function to check retryable errors.
func WithRetryableErr(retryableErr func(error) bool) Opt {
	return func(r *Retry) error {
		if retryableErr == nil {
			return errors.New("need a non-nil retryableErr")
		}

		r.retryableErr = retryableErr
		return nil
	}
}

// WithRspToErr allows user to register an additional function to convert rsp body errors.
func WithRspToErr(rspToErr func(interface{}) error) Opt {
	return func(r *Retry) error {
		if rspToErr == nil {
			return errors.New("need a non-nil rspToErr")
		}
		r.rspToErr = func(rsp interface{}) (err error) {
			defer func() {
				if rc := recover(); rc != nil {
					err = fmt.Errorf("retry rspToErr paniced: %v", rc)
				}
			}()
			return rspToErr(rsp)
		}
		return nil
	}
}

// WithExpBackoff set backoff strategy as exponential backoff.
func WithExpBackoff(
	initial, maximum time.Duration,
	multiplier int,
) Opt {
	return func(r *Retry) error {
		if _, ok := r.bf.(*backoff.CustomizedBackoff); ok {
			tlog.Trace("omit exponentialBackoff, since a customizedBackoff has already been set")
			return nil
		}

		bf, err := backoff.NewExponentialBackoff(initial, maximum, multiplier)
		if err != nil {
			return fmt.Errorf("failed to create new exponentialBackoff, err: %w", err)
		}

		r.bf = bf
		return nil
	}
}

// WithLinearBackoff set backoff strategy as linear backoff.
func WithLinearBackoff(bfs ...time.Duration) Opt {
	return func(r *Retry) error {
		switch r.bf.(type) {
		case *backoff.CustomizedBackoff:
			tlog.Trace("omit linearBackoff, since a customizedBackoff has already been set")
			return nil
		case *backoff.ExponentialBackoff:
			tlog.Trace("omit linearBackoff, since an exponentialBackoff has already been set")
			return nil
		default:
		}

		bf, err := backoff.NewLinearBackoffs(bfs...)
		if err != nil {
			return fmt.Errorf("failed to create new linearBackoff, err: %w", err)
		}

		r.bf = bf
		return nil
	}
}

// WithBackoff set a user defined backoff function.
func WithBackoff(bf func(attempt int) time.Duration) Opt {
	return func(r *Retry) error {
		bf, err := backoff.NewCustomizedBackoff(bf)
		if err != nil {
			return fmt.Errorf("failed to create new customizedBackoff, err: %w", err)
		}

		r.bf = bf
		return nil
	}
}

// WithSkipVisitedNodes set whether to skip visited nodes in next retry request.
//
// The behavior depends on selector implementation.
// If skip is true, selector **must** always not return a visited node.
// If skip is false, selectors of each hedging request act absolutely independently.
// Without this Opt, as the default behavior, selector **should** try its best to return a non-visited node.
// If all nodes has been visited, it **may** returns a node as its wish.
func WithSkipVisitedNodes(skip bool) Opt {
	return func(r *Retry) error {
		r.skipVisitedNodes = &skip
		return nil
	}
}

// WithConditionalLog set a conditional log for retry policy.
// Only requests which meet the condition will be displayed.
func WithConditionalLog(l lazylog.Logger, condition func(stat view.IStat) bool) Opt {
	return func(r *Retry) error {
		r.logCondition = condition
		r.newLazyLog = func() ILazyLogger {
			return lazylog.NewLazyLog(l)
		}
		return nil
	}
}

// WithConditionalCtxLog set a conditional log for retry policy.
// Only requests which meet the condition will be displayed.
func WithConditionalCtxLog(l lazylog.CtxLogger, condition func(stat view.IStat) bool) Opt {
	return func(r *Retry) error {
		r.logCondition = condition
		r.newLazyLog = func() ILazyLogger {
			return lazylog.NewLazyCtxLog(l)
		}
		return nil
	}
}

// WithEmitter set the emitter for retry policy.
func WithEmitter(emitter metrics.IEmitter) Opt {
	return func(r *Retry) error {
		r.reporter = metrics.NewReport(emitter)
		return nil
	}
}

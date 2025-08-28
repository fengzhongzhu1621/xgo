package retry

import (
	"context"
	"time"

	"github.com/fengzhongzhu1621/xgo/collections/flowctrl/throttle"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/naming/bannednodes"
)

// ThrottledRetry defines a retry policy with throttle.
//
// A Retry should not be bound to a throttle. Instead, user may bind one throttle to some Retrys.
// This is why we introduce a new struct instead of adding a new field to Retry.
type ThrottledRetry struct {
	*Retry
	throttle throttle.Throttler
}

// NewThrottledRetry create a new ThrottledRetry.
func NewThrottledRetry(
	maxAttempts int,
	ecs []int,
	throttle throttle.Throttler,
	opts ...Opt,
) (*ThrottledRetry, error) {
	// 执行父类的构造函数
	r, err := New(maxAttempts, ecs, opts...)
	if err != nil {
		return nil, err
	}

	return r.NewThrottledRetry(throttle), nil
}

// Invoke invokes handler f with Retry policy.
func (r *ThrottledRetry) Invoke(ctx context.Context, req, rsp interface{}, f filter.ClientHandleFunc) error {
	// 注入不可变标记
	ctx = client.WithOptionsImmutable(ctx)
	// mandatory 标志允许区分临时禁止（如熔断）和永久禁止（如黑名单）
	if r.skipVisitedNodes == nil {
		ctx = bannednodes.NewCtx(ctx, false)
	} else if *r.skipVisitedNodes {
		ctx = bannednodes.NewCtx(ctx, true)
	}

	l := r.newLazyLog()

	impl := r.newImpl(ctx, req, rsp, f, l)
	impl.Start()

	if r.logCondition(impl) {
		l.Printf(impl.String())
		l.FlushCtx(ctx)
	}
	r.reporter.Report(ctx, impl)

	return impl.err
}

// newImpl create an impl from ThrottledRetry.
func (r *ThrottledRetry) newImpl(
	ctx context.Context,
	req, rsp interface{},
	handler filter.ClientHandleFunc,
	log ILogger,
) *impl {
	return &impl{
		ThrottledRetry: r,
		ctx:            ctx,
		req:            req,
		rsp:            rsp,
		handler:        handler,
		timer:          time.NewTimer(0), // start first attempt at once
		log:            log,
	}
}

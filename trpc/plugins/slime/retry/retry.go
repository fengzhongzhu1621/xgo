package retry

import (
	"context"
	"errors"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/collections/flowctrl/throttle"
	lazylog "github.com/fengzhongzhu1621/xgo/logging/lazy_log"
	"github.com/fengzhongzhu1621/xgo/trpc/plugins/slime/view"
	"github.com/fengzhongzhu1621/xgo/trpc/plugins/slime/view/metrics"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/filter"
)

const (
	// MaximumAttempts maximum attempts.
	MaximumAttempts = 5 // 最大允许重试次数（安全阈值，防止无限重试）
)

// iBackoff is used to implement backoff priorities.
// customizedBackoff has the highest priority.
// exponentialBackoff comes second.
// linearBackoff has the lowest priority.

// Retry retry policy.
type Retry struct {
	maxAttempts      int                     // 最大重试次数（不超过MaximumAttempts）
	bf               IBackoff                // 退避策略接口（支持指数/线性/自定义退避）
	retryableECs     map[int]struct{}        // 可重试的错误码集合（HTTP状态码或业务错误码）
	retryableErr     func(error) bool        // 自定义错误判断函数（动态决定是否重试）
	rspToErr         func(interface{}) error // 将响应体转换为错误的函数
	skipVisitedNodes *bool                   // 是否跳过已访问节点（用于分布式场景）

	logCondition func(view.IStat) bool // 日志记录条件函数
	newLazyLog   func() ILazyLogger    // 惰性日志生成函数
	reporter     IReporter             // 监控上报接口（如Prometheus）
}

// New create a Retry policy.
// An error will be returned if provided args cannot build a valid Retry.
func New(maxAttempts int, ecs []int, opts ...Opt) (*Retry, error) {
	if maxAttempts <= 0 {
		return nil, errors.New("maxAttempts must be positive")
	}

	if maxAttempts > MaximumAttempts {
		maxAttempts = MaximumAttempts // 强制限制最大尝试次数
	}

	// 初始化默认值
	r := Retry{
		maxAttempts:  maxAttempts,
		retryableECs: make(map[int]struct{}),                           // 初始化错误码集合
		rspToErr:     func(interface{}) error { return nil },           // 默认响应转换函数
		logCondition: func(view.IStat) bool { return false },           // 默认不记录日志
		newLazyLog:   func() ILazyLogger { return &lazylog.NoopLog{} }, // 默认空日志
		reporter:     &metrics.Noop{},                                  // 默认空监控上报
	}

	// 配置错误码集合
	for _, ec := range ecs {
		r.retryableECs[ec] = struct{}{}
	}

	// 应用可选参数
	for _, opt := range opts {
		if err := opt(&r); err != nil {
			return nil, fmt.Errorf("failed to apply Retry Opt(s), err: %w", err)
		}
	}

	// 校验必填参数
	if r.bf == nil {
		return nil, errors.New("backoff is uninitialized")
	}

	if len(r.retryableECs) == 0 && r.retryableErr == nil {
		return nil, errors.New("one of retryableECs or retryableErr must be provided")
	}

	if r.retryableErr == nil {
		r.retryableErr = func(err error) bool { return false }
	}

	return &r, nil
}

// isRetryableErr checks whether the error is retryable.
func (r *Retry) isRetryableErr(err error) bool {
	if _, ok := r.retryableECs[int(errs.Code(err))]; ok {
		return true
	}

	return r.retryableErr(err)
}

// NewThrottledRetry create a new ThrottledRetry from receiver Retry.
// 由父类创建一个子类
func (r *Retry) NewThrottledRetry(throttle throttle.Throttler) *ThrottledRetry {
	return &ThrottledRetry{Retry: r, throttle: throttle}
}

// Invoke calls Invoke of ThrottledRetry with a Noop throttle.
// 执行子类的默认操作
func (r *Retry) Invoke(ctx context.Context, req, rsp interface{}, f filter.ClientHandleFunc) error {
	// 默认不支持限频
	return r.NewThrottledRetry(throttle.NewNoop()).
		Invoke(ctx, req, rsp, f)
}

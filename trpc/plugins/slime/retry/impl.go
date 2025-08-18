package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/fengzhongzhu1621/xgo/buildin/reflectutils"
	copyutils "github.com/fengzhongzhu1621/xgo/copier"
	"github.com/fengzhongzhu1621/xgo/trpc/plugins/slime/view"
	"github.com/fengzhongzhu1621/xgo/trpc/utils/cpmsg"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/filter"
)

// TimeoutErr request deadline exceeded.
var TimeoutErr = errs.NewFrameError(errs.RetClientTimeout, "request timeout")

// impl contains some useful fields to implement retry request.
// 实现了 IStat 接口
type impl struct {
	*ThrottledRetry                         // 继承限流重试基类（提供限流逻辑）
	ctx             context.Context         // 请求上下文（支持超时/取消）
	req             interface{}             // 原始请求对象
	rsp             interface{}             // 最终响应对象
	err             error                   // 最终错误
	handler         filter.ClientHandleFunc // TRPC客户端处理函数

	cost      time.Duration // 总耗时
	throttled bool          // 是否被限流
	frozen    bool          // 是否终止重试（如达到最大尝试次数）
	attempts  []*attempt    // 所有尝试记录
	timer     *time.Timer   // 控制重试间隔的定时器
	log       ILogger       // 日志接口
}

// newAttempt creates a new attempt.
//
// The msg and rsp in impl are copied to attempt.
// newAttempt freeze impl if all attempts has been drained or throttle check is failed.
func (impl *impl) newAttempt() (*attempt, error) {
	// 复制上下文和消息
	ctx, msg := codec.WithNewMessage(impl.ctx)
	if err := cpmsg.CopyMsg(msg, codec.Message(impl.ctx)); err != nil {
		return nil, fmt.Errorf("failed to create new attempt: %w", err)
	}

	a := attempt{
		impl:    impl,
		ctx:     ctx,
		rsp:     reflectutils.New(impl.rsp), // 深拷贝响应对象
		attempt: len(impl.attempts) + 1,     // 当前尝试序号
	}
	impl.attempts = append(impl.attempts, &a) // 记录尝试

	impl.log.Printf("start %dth attempt", a.attempt)

	if len(impl.attempts) == impl.maxAttempts || !impl.throttle.Allow() {
		if len(impl.attempts) == impl.maxAttempts {
			impl.log.Printf("freeze hedging for no more attempts")
		} else {
			impl.throttled = true
			impl.log.Printf("freeze hedging for throttle")
		}

		impl.frozen = true
	}

	return &a, nil
}

// Start start the main loop of retry.
func (impl *impl) Start() {
	start := time.Now()
	defer func() {
		impl.cost = time.Since(start)
	}()

	for {
		select {
		case <-impl.ctx.Done():
			impl.log.Printf("retry finished for timeout error")
			impl.err = TimeoutErr
			return
		case <-impl.timer.C:
			a, err := impl.newAttempt()
			if err != nil {
				impl.err = err
				return
			}
			a.SyncStart()
			if impl.onReturn(a) {
				impl.log.Printf("%dth attempt is return to client", a.attempt)
				return
			}
		}
	}
}

// onReturn process the returned attempt.
//
// It returns a boolean indicate whether should the attempt terminate main loop of impl.
func (impl *impl) onReturn(a *attempt) (final bool) {
	impl.log.Printf("%dth attempt has returned", a.attempt)
	a.OnReturn()

	defer func() {
		if final {
			if err := cpmsg.CopyMsg(codec.Message(impl.ctx), codec.Message(a.ctx)); err != nil {
				impl.err = fmt.Errorf("failed to copy back msg: %w, attempt err: %s", err, a.err)
			} else {
				impl.err = a.err
			}
		}
		codec.PutBackMessage(codec.Message(a.ctx))
	}()
	if a.err == nil {
		a.err = impl.rspToErr(a.rsp)
	}
	if a.err == nil {
		a.err = copyutils.ShallowCopy(impl.rsp, a.rsp)
		return true
	}
	if !impl.isRetryableErr(a.err) {
		return true
	}

	if a.pushbackDelay == nil {
		impl.scheduleNext(impl.bf.Backoff(a.attempt))
	} else {
		impl.log.Printf("server issues a pushback delay: %v", *a.pushbackDelay)
		impl.scheduleNext(*a.pushbackDelay)
	}

	return impl.frozen
}

// scheduleNext schedules next retry request.
func (impl *impl) scheduleNext(delay time.Duration) {
	if impl.frozen {
		return
	}

	if delay < 0 {
		impl.timer.Stop()
		impl.frozen = true
		impl.log.Printf("freeze retry for negative delay")
		return
	}

	if !impl.timer.Stop() {
		select {
		case <-impl.timer.C:
		default:
		}
	}
	impl.timer.Reset(delay)
}

// Cost implements view.Stat.
func (impl *impl) Cost() time.Duration {
	return impl.cost
}

// Attempts implements view.Stat.
func (impl *impl) Attempts() []view.IAttempt {
	attempts := make([]view.IAttempt, 0, len(impl.attempts))
	for _, att := range impl.attempts {
		attempts = append(attempts, att)
	}
	return attempts
}

// Throttled implements view.Stat.
func (impl *impl) Throttled() bool {
	return impl.throttled
}

// InflightN implements view.Stat.
func (impl *impl) InflightN() int {
	return 0
}

// Error implements view.Stat.
func (impl *impl) Error() error {
	return impl.err
}

// String implements fmt.Stringer.
func (impl *impl) String() string {
	var s string
	s += fmt.Sprintf("totalAttempts: %d, throttled: %t, finalErr: %v\n",
		len(impl.attempts), impl.throttled, impl.err)
	for _, a := range impl.attempts {
		s += "\t" + a.String() + "\n"
	}
	return s[:len(s)-1]
}

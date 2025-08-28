package mock

import (
	"context"
	"math/rand"
	"time"

	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/filter"
)

// ClientFilter set client request mock interceptor.
func ClientFilter(opts ...Option) filter.ClientFilter {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	rand.Seed(time.Now().Unix())

	return func(ctx context.Context, req, rsp interface{}, handler filter.ClientHandleFunc) error {
		msg := trpc.Message(ctx)

		for _, mock := range o.mocks {
			if mock.Method != "" && mock.Method != msg.ClientRPCName() {
				continue
			}

			if mock.Percent == 0 || rand.Intn(100) >= mock.Percent {
				// Triggered by percentage. For example, if 20%, the random number 0-99, only 0-19 will trigger.
				continue
			}

			if mock.Timeout {
				// 模拟框架超时
				<-ctx.Done()
				return errs.NewFrameError(errs.RetClientTimeout, "mock filter: timeout")
			}
			if mock.Delay > 0 {
				// 模拟自定义超时
				select {
				case <-ctx.Done():
					return errs.NewFrameError(
						errs.RetClientTimeout,
						"mock filter: timeout during delay mock",
					)
				case <-time.After(mock.delay):
				}
			}

			// 模拟返回自定义错误码和错误信息
			if mock.Retcode > 0 {
				return errs.New(mock.Retcode, mock.Retmsg)
			}

			// 模拟返回成功的body
			if mock.Body != "" {
				if err := codec.Unmarshal(mock.Serialization, mock.data, rsp); err != nil {
					return errs.NewFrameError(
						errs.RetClientDecodeFail,
						"mock filter Unmarshal: "+err.Error(),
					)
				}
				return nil
			}
		}

		return handler(ctx, req, rsp)
	}
}

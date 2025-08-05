package degrade

import (
	"context"
	"math/rand"

	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/filter"
)

// Filter 声明熔断组件的 filter 来充当拦截器
func Filter(
	ctx context.Context, req interface{}, handler filter.ServerHandleFunc,
) (interface{}, error) {
	// 触发熔断，保留部分流量
	if isDegrade {
		randNum := rand.Intn(100)
		if randNum >= cfg.DegradeRate {
			return nil, errs.New(systemDegradeErrNo, errDegardeReturn)
		}
	}

	// 限制最大并发数
	if enableConcurrency() {
		select {
		case sema <- struct{}{}: // 未达到最大并发请求数，占用一个并发槽位
			// 这个 defer 属于当前 case 分支，作用域是整个 Filter 函数
			// 不管 Filter 函数后续是正常返回、还是发生了错误返回（比如 handler 执行出错），
			// 这个 defer 都会在 Filter 函数退出之前被执行，从而调用 <-sema 释放信号量。
			// 不会 只作用于 case 分支内部，因为 defer 是 Go 的函数级关键字，
			// 一旦注册，它就绑定到了 当前函数（Filter）的执行栈帧 上，会在函数返回前执行。
			defer func() {
				<-sema // 释放一个并发槽位
			}()

		default: // 达到最大并发请求数，直接丢弃请求
			return nil, errs.New(systemDegradeErrNo, errDegardeReturn)
		}
	}

	return handler(ctx, req)
}

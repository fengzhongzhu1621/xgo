package hystrix

import (
	"context"
	"fmt"
	"runtime"

	"github.com/afex/hystrix-go/hystrix"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/log"
)

const (
	filterName = "hystrix"
)

var (
	// WildcardKey is a wild card symbol.
	WildcardKey = "*"
	// ExcludeKey excludes prefix key.
	ExcludeKey = "_"

	panicBufLen = 1024
)

// recoveryHandler prints call stack information.
var recoveryHandler = func(ctx context.Context, e interface{}) error {
	buf := make([]byte, panicBufLen)
	buf = buf[:runtime.Stack(buf, false)]
	log.ErrorContextf(ctx, "[Hystrix-Panic] %v\n%s\n", e, buf)
	return fmt.Errorf("%v", e)
}

// ServerFilter fuses server request.
func ServerFilter() filter.ServerFilter {
	return func(ctx context.Context, req interface{}, handler filter.ServerHandleFunc) (interface{}, error) {
		// Get routing and configuration.
		cmd := trpc.Message(ctx).ServerRPCName()

		// 根据路由获取熔断配置
		if _, ok := cfg[cmd]; !ok {
			// Configuration does not exist.
			// Whether there is wild card symbol.
			// *表示通配配置
			if _, ok := cfg[WildcardKey]; !ok {
				// 没有直接返回
				return handler(ctx, req)
			}

			// Whether to exclude parts when opening wildcards.
			// _前缀表示排除接口
			if _, ok := cfg[ExcludeKey+cmd]; ok {
				return handler(ctx, req)
			}
			cmd = WildcardKey
		}

		var rsp interface{}
		// 执行熔断逻辑
		// 第一个参数 cmd 是熔断命令名称
		// 第二个参数是执行函数，在其中调用实际的 RPC 处理函数 handler
		return rsp, hystrix.Do(cmd, func() (err error) {
			defer func() {
				// 使用 defer 捕获 panic 并调用恢复处理器
				if errPanic := recover(); errPanic != nil {
					err = recoveryHandler(ctx, errPanic)
				}
			}()
			rsp, err = handler(ctx, req)
			return err
		}, nil)
	}
}

// ClientFilter fuses client request.
func ClientFilter() filter.ClientFilter {
	return func(ctx context.Context, req, rsp interface{}, handler filter.ClientHandleFunc) error {
		// Get routing and configuration.
		cmd := trpc.Message(ctx).ClientRPCName()
		if _, ok := cfg[cmd]; !ok {
			// Configuration does not exist.
			// Whether there is wild card symbol.
			if _, ok := cfg[WildcardKey]; !ok {
				// 没有直接返回
				return handler(ctx, req, rsp)
			}
			//  Whether to exclude parts when opening wildcards.
			if _, ok := cfg[ExcludeKey+cmd]; ok {
				return handler(ctx, req, rsp)
			}
			cmd = WildcardKey
		}
		return hystrix.Do(cmd, func() (err error) {
			defer func() {
				if errPanic := recover(); errPanic != nil {
					err = recoveryHandler(ctx, errPanic)
				}
			}()
			return handler(ctx, req, rsp)
		}, nil)
	}
}

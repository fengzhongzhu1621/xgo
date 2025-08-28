package recovery

import (
	"context"
	"fmt"
	"runtime"

	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/metrics"
)

// PanicBufLen is the size of the buffer for storing the panic call stack log. The default value as below.
// 存储 panic 调用栈的缓冲区大小（默认 4KB），限制调用栈日志的长度，避免内存浪费。
var PanicBufLen = 4096

type options struct {
	rh RecoveryHandler
}

// RecoveryHandler is the Recovery handle function.
// Deprecated: Use recovery.Handler instead.
type RecoveryHandler = Handler //nolint:revive

// Handler is the Recovery handle function.
// ctx：上下文，用于传递请求链路的元数据（如 TraceID）。
// err：panic 捕获到的值（类型为 any）。
type Handler func(ctx context.Context, err any) error

var defaultOptions = &options{
	rh: defaultRecoveryHandler, // Recovery 处理函数，默认为 defaultRecoveryHandler
}

var defaultRecoveryHandler = func(ctx context.Context, e any) error {
	buf := make([]byte, PanicBufLen)                                  // 分配缓冲区
	buf = buf[:runtime.Stack(buf, false)]                             // 获取当前协程的调用栈
	log.ErrorContextf(ctx, "[PANIC]%v\n%s\n", e, buf)                 // 记录错误日志（带上下文）
	metrics.IncrCounter("trpc.PanicNum", 1)                           // 监控计数
	return errs.NewFrameError(errs.RetServerSystemErr, fmt.Sprint(e)) // 返回系统错误
}

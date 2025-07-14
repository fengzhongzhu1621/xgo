package entity

import (
	"context"
	"runtime"
	"time"

	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

// 存储 panic 调用栈的缓冲区大小（默认 4KB），限制调用栈日志的长度，避免内存浪费。
var PanicBufLen = 4096

// Do 创建异步任务
func Do(fn func(context.Context), timeout time.Duration) {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(trpc.BackgroundContext(), timeout)

	go func() {
		defer cancel() // 主动取消，冗余保护，确保goroutine退出时取消context

		// panic捕获与日志记录
		defer func() {
			if e := recover(); e != nil {
				buf := make([]byte, PanicBufLen)                  // 分配缓冲区
				buf = buf[:runtime.Stack(buf, false)]             // 获取当前协程的调用栈
				log.ErrorContextf(ctx, "[PANIC]%v\n%s\n", e, buf) // 记录错误日志（带上下文）
			}
		}()

		// 记录请求开始时间
		begin := time.Now()

		// 任务执行日志
		log.DebugContextf(ctx, "start an async task")

		// 执行实际业务逻辑
		fn(ctx)

		log.DebugContextf(ctx, "async task finished, cost: %s", time.Since(begin))
	}()
}

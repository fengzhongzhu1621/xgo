package debuglog

import (
	"context"
	"time"

	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/filter"
)

// ServerFilter is the server-side filter.
func ServerFilter(opts ...Option) filter.ServerFilter {
	// 获得默认配置
	o := getFilterOptions(opts...)

	// 设置追加的日志格式器
	nilLogFormat := getLogFormat(debugLevel, o.enableColor, "server request:%s, cost:%s, from:%s%s")
	errLogFormat := getLogFormat(
		errorLevel,
		o.enableColor,
		"server request:%s, cost:%s, from:%s, err:%s%s",
	)
	deadlineLogFormat := getLogFormat(errorLevel, o.enableColor,
		"server request:%s, cost:%s, from:%s, err:%s, total timeout:%s%s")

	return func(ctx context.Context, req any, handler filter.ServerHandleFunc) (rsp any, err error) {
		// 记录请求开始时间
		begin := time.Now()

		// 执行实际业务逻辑
		rsp, err = handler(ctx, req)

		// 从上下文中提取 RPC 元数据（如方法名、远程地址）
		msg := trpc.Message(ctx)

		// 检查是否满足日志记录条件（如方法名或错误码过滤），只打印默写规则下的调试日志
		if !o.passed(msg.ServerRPCName(), int(errs.Code(err))) {
			return rsp, err
		}

		// 记录请求结束时间
		end := time.Now()

		// 获取客户端地址
		var addr string
		if msg.RemoteAddr() != nil {
			addr = msg.RemoteAddr().String()
		}

		if err == nil { // 成功请求
			// 调用日志函数记录成功日志
			o.nilLogLevelFunc(
				ctx,
				nilLogFormat,             // 日志格式
				msg.ServerRPCName(),      // 服务端方法名
				end.Sub(begin),           // 请求耗时
				addr,                     // 客户端地址
				o.logFunc(ctx, req, rsp), // 追加额外信息
			)
		} else { // 错误请求
			// 返回该Context被取消的截止时间，用于检查是否有超时设置
			deadline, ok := ctx.Deadline()
			if ok { // 记录超时错误
				o.errLogLevelFunc(
					ctx,
					deadlineLogFormat,        // 日志格式
					msg.ServerRPCName(),      // 服务端方法名
					end.Sub(begin),           // 请求耗时
					addr,                     // 客户端地址
					err.Error(),              // 错误原因
					deadline.Sub(begin),      // 取消时的耗时
					o.logFunc(ctx, req, rsp), // 追加额外信息
				)
			} else { // 记录普通错误
				o.errLogLevelFunc(
					ctx,
					errLogFormat,             // 日志格式
					msg.ServerRPCName(),      // 服务端方法名
					end.Sub(begin),           // 请求耗时
					addr,                     // 客户端地址
					err.Error(),              // 错误原因
					o.logFunc(ctx, req, rsp), // 追加额外信息
				)
			}
		}
		return rsp, err
	}
}

// ClientFilter is the client-side filter.
func ClientFilter(opts ...Option) filter.ClientFilter {
	// 获得默认配置
	o := getFilterOptions(opts...)

	// 设置追加的日志格式器
	nilLogFormat := getLogFormat(debugLevel, o.enableColor, "client request:%s, cost:%s, to:%s%s")
	errLogFormat := getLogFormat(
		errorLevel,
		o.enableColor,
		"client request:%s, cost:%s, to:%s, err:%s%s",
	)

	return func(ctx context.Context, req, rsp any, handler filter.ClientHandleFunc) (err error) {
		msg := trpc.Message(ctx)

		// 记录开始时间
		begin := time.Now()

		// 执行远程调用
		err = handler(ctx, req, rsp)

		// 检查日志条件
		if !o.passed(msg.ClientRPCName(), int(errs.Code(err))) {
			return err
		}

		// 记录结束时间
		end := time.Now()

		// 获取服务端地址
		var addr string
		if msg.RemoteAddr() != nil {
			addr = msg.RemoteAddr().String()
		}

		if err == nil { // 成功请求
			o.nilLogLevelFunc(
				ctx,
				nilLogFormat,
				msg.ClientRPCName(),
				end.Sub(begin),
				addr,
				o.logFunc(ctx, req, rsp),
			)
		} else { // 错误请求
			o.errLogLevelFunc(
				ctx, errLogFormat, msg.ClientRPCName(), end.Sub(begin), addr, err.Error(), o.logFunc(ctx, req, rsp),
			)
		}

		return err
	}
}

package filter

import (
	"context"
)

// ServerChain 服务端过滤器链，用于串联多个服务端过滤器
type ServerChain []ServerFilter

// Filter 按顺序调用过滤器链中的每个服务端过滤器
// 参数:
//
//	ctx: 上下文对象，用于传递请求相关信息（如超时控制、链路追踪等）
//	req: 请求对象，包含客户端发送的原始数据
//	next: 下一个处理函数，通常是实际的业务处理逻辑
//
// 返回:
//
//	interface{}: 处理后的响应对象
//	error: 处理过程中的错误（如过滤器拦截、业务逻辑错误等）
func (c ServerChain) Filter(ctx context.Context, req interface{}, next ServerHandleFunc) (interface{}, error) {
	// 创建包装后的处理函数，将next函数包装成标准格式
	// 这个闭包函数负责在执行实际业务逻辑前后添加额外处理
	nextF := func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		// 创建新的链路追踪上下文（这里可以添加链路追踪逻辑）
		// 调用实际的业务处理函数处理请求
		rsp, err = next(ctx, req)
		return rsp, err
	}

	// 从后往前遍历过滤器链，构建嵌套的处理函数
	// 这种反向构建方式确保过滤器按照声明顺序执行（先进先出）
	for i := len(c) - 1; i >= 0; i-- {
		// 保存当前循环的变量（避免闭包捕获问题）
		// curHandleFunc: 当前已构建的处理链
		// curFilter: 当前要添加的过滤器
		// _: 当前索引（这里忽略未使用）
		curHandleFunc, curFilter, _ := nextF, c[i], i

		// 重新定义nextF，将当前过滤器加入到调用链中
		// 每次循环都会创建一个新的闭包，形成嵌套调用结构
		nextF = func(ctx context.Context, req interface{}) (interface{}, error) {
			// 调用当前过滤器，并传入下一个处理函数
			// 这样当过滤器执行时，可以通过调用curHandleFunc来触发后续链式调用
			rsp, err := curFilter(ctx, req, curHandleFunc)
			return rsp, err
		}
	}

	// 执行构建好的过滤器链，启动整个处理流程
	// 从这里开始，请求会依次通过所有过滤器的预处理、业务处理、后处理阶段[5](@ref)
	return nextF(ctx, req)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// NoopServerFilter 是 ServerFilter 的空操作实现
// 通常用于测试或作为默认过滤器
func NoopServerFilter(ctx context.Context, req interface{}, next ServerHandleFunc) (rsp interface{}, err error) {
	return next(ctx, req)
}

package filter

import "context"

// EmptyChain is an empty chain.
var EmptyChain = ClientChain{}

// ClientChain chains client side filters.
type ClientChain []ClientFilter

// Filter invokes every client side filters in the chain.
func (c ClientChain) Filter(ctx context.Context, req, rsp interface{}, next ClientHandleFunc) error {
	var nextF ClientHandleFunc

	nextF = func(ctx context.Context, req, rsp interface{}) error {
		err := next(ctx, req, rsp)
		return err
	}

	// 递归调用，从后往前
	// 每次调用都将nextF作为参数传递给当前过滤器，这样当前过滤器可以调用nextF来触发后续过滤器的执行
	for i := len(c) - 1; i >= 0; i-- {
		var (
			curHandleFunc ClientHandleFunc
			curFilter     ClientFilter
		)
		curHandleFunc, curFilter, _ = nextF, c[i], i
		nextF = func(ctx context.Context, req, rsp interface{}) error {
			return curFilter(ctx, req, rsp, curHandleFunc)
		}
	}

	return nextF(ctx, req, rsp)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// NoopClientFilter is the noop implementation of ClientFilter.
func NoopClientFilter(ctx context.Context, req, rsp interface{}, next ClientHandleFunc) error {
	return next(ctx, req, rsp)
}

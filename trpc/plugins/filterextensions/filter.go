package filterextensions

import (
	"context"

	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/filter"
)

func newClientIntercept(
	serviceFilters serviceMethodClientFilters,
) filter.ClientFilter {
	return func(ctx context.Context, req, rsp interface{}, handler filter.ClientHandleFunc) error {
		msg := trpc.Message(ctx)
		// 获取被调用的服务名称
		service := msg.CalleeServiceName()
		// 获取被调用的方法名称
		method := msg.CalleeMethod()

		if methodFilters, ok := serviceFilters[service]; ok {
			if filters, ok := methodFilters[method]; ok {
				// 获得被调方法关联的过滤器，添加到过滤器链中
				return filter.ClientChain(filters).Filter(ctx, req, rsp, handler)
			}
		}
		return handler(ctx, req, rsp)
	}
}

func newServerIntercept(
	serviceFilters serviceMethodServerFilters,
) filter.ServerFilter {
	return func(ctx context.Context, req interface{}, handler filter.ServerHandleFunc) (interface{}, error) {
		msg := trpc.Message(ctx)
		// 获取被调用的服务名称
		service := msg.CalleeServiceName()
		// 获取被调用的方法名称
		method := msg.CalleeMethod()
		if methodFilters, ok := serviceFilters[service]; ok {
			if filters, ok := methodFilters[method]; ok {
				return filter.ServerChain(filters).Filter(ctx, req, handler)
			}
		}
		return handler(ctx, req)
	}
}

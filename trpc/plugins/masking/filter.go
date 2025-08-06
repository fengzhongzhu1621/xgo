package masking

import (
	"context"

	"trpc.group/trpc-go/trpc-go/filter"
)

// ServerFilter 服务端RPC调用自动校验req输入参数
func ServerFilter() filter.ServerFilter {
	return func(ctx context.Context, req interface{}, handler filter.ServerHandleFunc) (rsp interface{}, err error) {
		rsp, err = handler(ctx, req)
		if err != nil {
			return nil, err
		}
		DeepCheck(rsp)
		return rsp, nil
	}
}

package recovery

import (
	"context"

	"trpc.group/trpc-go/trpc-go/filter"
)

// ServerFilter adds the recovery filter to the server.
func ServerFilter(opts ...Option) filter.ServerFilter {
	o := defaultOptions
	for _, opt := range opts {
		opt(o)
	}
	return func(ctx context.Context, req interface{}, handler filter.ServerHandleFunc) (rsp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = o.rh(ctx, r)
			}
		}()

		return handler(ctx, req)
	}
}

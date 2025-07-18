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
	return func(ctx context.Context, req any, handler filter.ServerHandleFunc) (rsp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = o.rh(ctx, r)
			}
		}()

		return handler(ctx, req)
	}
}

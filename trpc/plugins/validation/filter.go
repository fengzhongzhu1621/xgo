package validation

import (
	"context"

	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/http"
	"trpc.group/trpc-go/trpc-go/log"
)

// ServerFilter automatically validates the req input parameters during server-side RPC invocation.
// Deprecated: Use ServerFilterWithOptions instead.
func ServerFilter(opts ...Option) filter.ServerFilter {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	return ServerFilterWithOptions(o)
}

// ClientFilter automatically validates the rsp response parameters during client-side RPC invocation.
// Deprecated: Use ClientFilterWithOptions instead.
func ClientFilter(opts ...Option) filter.ClientFilter {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	return ClientFilterWithOptions(o)
}

// ServerFilterWithOptions automatically validates the req input parameters during server-side RPC invocation.
// 根据配置创建一个服务端过滤器
func ServerFilterWithOptions(o options) filter.ServerFilter {
	return func(ctx context.Context, req interface{}, handler filter.ServerHandleFunc) (interface{}, error) {
		// The request structure has not been validated by Validator.
		valid, ok := req.(Validator)
		if !ok {
			return handler(ctx, req)
		}

		// Verification passed.
		// 请求验证
		err := valid.Validate()
		if err == nil {
			return handler(ctx, req)
		}

		// Record logs as needed when the verification fails.
		errMsg := err.Error()
		if o.EnableErrorLog {
			// 获得请求头
			if head, ok := ctx.Value(http.ContextKeyHeader).(*http.Header); ok {
				reqPath := head.Request.URL.Path
				reqRawQuery := head.Request.URL.RawQuery
				reqUserAgent := head.Request.Header.Get("User-Agent")
				reqReferer := head.Request.Header.Get("Referer")
				log.WithContext(ctx,
					log.Field{Key: "request_content", Value: req},
					log.Field{Key: "request_path", Value: reqPath},
					log.Field{Key: "request_query", Value: reqRawQuery},
					log.Field{Key: "request_useragent", Value: reqUserAgent},
					log.Field{Key: "request_referer", Value: reqReferer},
				).Errorf("validation request error: %s", errMsg)
			} else {
				log.WithContext(ctx, log.Field{Key: "request_content", Value: req}).
					Errorf("validation request error: %s", errMsg)
			}
		}

		return nil, errs.New(o.ServerValidateErrCode, errMsg)
	}
}

// ClientFilterWithOptions automatically validates the rsp response parameters during client-side RPC invocation.
func ClientFilterWithOptions(o options) filter.ClientFilter {
	return func(ctx context.Context, req, rsp interface{}, handler filter.ClientHandleFunc) error {
		// rsp does not need to be validated if An error occurred when calling downstream.
		if err := handler(ctx, req, rsp); err != nil {
			return err
		}

		// 验证服务端的响应
		// The response structure has not been validated by Validator.
		valid, ok := rsp.(Validator)
		if !ok {
			return nil
		}

		// Verification passed.
		err := valid.Validate()
		if err == nil {
			return nil
		}

		// Record logs as needed when the verification fails.
		if o.EnableErrorLog {
			log.WithContext(ctx, log.Field{Key: "response_content", Value: rsp}).
				Errorf("validation response error: %v", err)
		}

		return errs.New(o.ClientValidateErrCode, err.Error())
	}
}

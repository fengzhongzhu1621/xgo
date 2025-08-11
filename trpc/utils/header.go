package utils

import (
	"context"

	"trpc.group/trpc-go/trpc-go/http"
)

// GetHTTPHeaderFromCtx 从请求上下文获取 Header 信息
func GetHTTPHeaderFromCtx(ctx context.Context) *http.Header {
	head, ok := ctx.Value(http.ContextKeyHeader).(*http.Header)
	if ok {
		return head
	}

	return nil
}

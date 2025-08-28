package auth

import (
	"context"
	"strings"

	trpcHttp "trpc.group/trpc-go/trpc-go/http"
)

// GetBearerToken 默认获取 token 的函数, 用户可自定义实现
func GetBearerToken(ctx context.Context, req any) (string, error) {
	head := trpcHttp.Head(ctx)
	token := head.Request.Header.Get("Authorization")
	return strings.TrimPrefix(token, "Bearer "), nil
}

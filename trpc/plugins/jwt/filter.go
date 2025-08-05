package jwt

import (
	"context"
	"fmt"

	"github.com/fengzhongzhu1621/xgo/trpc/utils/auth"
	"github.com/mitchellh/mapstructure"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/filter"
	trpcHttp "trpc.group/trpc-go/trpc-go/http"
)

// ServerFilter 设置服务端增加 jwt 验证
func ServerFilter(opts ...Option) filter.ServerFilter {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return func(ctx context.Context, req interface{}, handler filter.ServerHandleFunc) (interface{}, error) {
		head := trpcHttp.Head(ctx)
		// 非http请求
		if head == nil {
			return handler(ctx, req)
		}

		// 是否跳过OA验证(path白名单)
		r := head.Request
		if r.URL != nil && o.isInExcludePath(r.URL.Path) {
			return handler(ctx, req)
		}

		// 从auth头获取token
		token, err := auth.GetBearerToken(ctx, req)
		if err != nil {
			return nil, err
		}

		// 验证token
		customInfo, err := DefaultSigner.Verify(token)
		if err != nil {
			return nil, errs.NewFrameError(errs.RetServerAuthFail, err.Error())
		}

		// 将认证信息存储到ctx中用于业务侧使用
		innerCtx := context.WithValue(ctx, AuthJwtCtxKey, customInfo)
		return handler(innerCtx, req)
	}
}

// GetCustomInfo 获取用户信息 (参数 ptr 为struct的指针对象)
func GetCustomInfo(ctx context.Context, ptr interface{}) error {
	if data, ok := ctx.Value(AuthJwtCtxKey).(map[string]interface{}); ok {
		return mapstructure.Decode(data, ptr)
	}
	return fmt.Errorf("fail to find ctx value! key=(%v)", AuthJwtCtxKey)
}

// Package jwt 身份认证
package jwt

import (
	"github.com/fengzhongzhu1621/xgo/network/nethttp/auth/jwt"
)

// ContextKey 定义类型
type ContextKey string

const (
	// AuthJwtCtxKey context key
	AuthJwtCtxKey = ContextKey("AuthJwtCtxKey")
)

// DefaultSigner 默认的 signer
var DefaultSigner jwt.Signer

// SetDefaultSigner 设置默认 signer
func SetDefaultSigner(s jwt.Signer) {
	if s != nil {
		DefaultSigner = s
	}
}

// options 插件配置
type options struct {
	ExcludePathSet map[string]bool // path 白名单
}

// isInExcludePath 是否在白名单列表中
func (o *options) isInExcludePath(path string) bool {
	_, ok := o.ExcludePathSet[path]
	return ok
}

// Option 设置参数选项
type Option func(*options)

// WithExcludePathSet 设置 path 白名单
func WithExcludePathSet(set map[string]bool) Option {
	return func(o *options) {
		o.ExcludePathSet = set
	}
}

package common

import (
	"context"
	"time"
)

type Cache interface {
	// 初始化缓存
	Parse() error
	// 初始缓存数据
	Scan() (map[string][]string, error)
	// 从缓存中获取指定键的值
	Get(key string) string
	// 返回缓存中指定键的所有值
	GetMany(key string) []string
	// 将一个值设置到缓存中指定键的位置
	Set(key string, value string)
	// 将多个值设置到缓存中指定键的位置
	SetMany(key string, value []string)
}

// RandomExtraExpirationDurationFunc is the type of the function generate extra expiration duration
// 生成随机过期时间偏移量的函数，用于打散过期时间
type RandomExtraExpirationDurationFunc func() time.Duration

// RetrieveFunc is the type of the retrieve function.
// it retrieves the value from database, redis, apis, etc.
// 禁用缓存时根据 key 获取 value 的函数
type RetrieveFunc func(ctx context.Context, key Key) (interface{}, error)

package cacheimpls

import (
	"time"

	"github.com/fengzhongzhu1621/xgo/cache/cleaner"
	"github.com/fengzhongzhu1621/xgo/cache/memory"
	"github.com/fengzhongzhu1621/xgo/cache/redis"
)

var (
	LocalAPIGatewayJWTClientIDCache memory.Cache

	TestRedisCache   *redis.Cache
	TestCacheCleaner *cleaner.CacheCleaner
)

func InitCaches(disabled bool) {
	LocalAPIGatewayJWTClientIDCache = memory.NewCache(
		"local_apigw_jwt_client_id",
		disabled,
		retrieveAPIGatewayJWTClientID,
		30*time.Second,
		nil,
	)

	TestRedisCache = redis.NewCache(
		"test",
		30*time.Minute,
	)

	TestCacheCleaner = cleaner.NewCacheCleaner("TestCacheCleaner", testCacheDeleter{})
	go TestCacheCleaner.Run()
}

func init() {
	InitCaches(false)
}

package cacheimpls

import (
	"time"

	"github.com/fengzhongzhu1621/xgo/cache/cleaner"
	"github.com/fengzhongzhu1621/xgo/cache/memory"
	"github.com/fengzhongzhu1621/xgo/cache/redis"
	gocache "github.com/patrickmn/go-cache"
)

// CacheLayer ...
const CacheLayer = "Cache"

// LocalAppCodeAppSecretCache ...
var (
	LocalAppCodeAppSecretCache *gocache.Cache
	LocalAuthAppAccessKeyCache *gocache.Cache

	LocalAPIGatewayJWTClientIDCache memory.Cache

	// for unittest
	TestRedisCache   *redis.Cache
	TestCacheCleaner *cleaner.CacheCleaner
)

// InitCaches
// Cache should only know about get/retrieve data
// ! DO NOT CARE ABOUT WHAT THE DATA WILL BE USED FOR
func InitCaches(disabled bool) {

	LocalAppCodeAppSecretCache = gocache.New(12*time.Hour, 5*time.Minute)

	// auth app_code/app_secret cache
	LocalAuthAppAccessKeyCache = gocache.New(12*time.Hour, 5*time.Minute)

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

package cache

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
	LocalAppCodeAppSecretCache      *gocache.Cache
	LocalAuthAppAccessKeyCache      *gocache.Cache
	LocalSubjectCache               memory.Cache
	LocalSubjectRoleCache           memory.Cache
	LocalSystemClientsCache         memory.Cache
	LocalRemoteResourceListCache    memory.Cache
	LocalSubjectPKCache             memory.Cache
	LocalSubjectDepartmentCache     memory.Cache
	LocalAPIGatewayJWTClientIDCache memory.Cache
	LocalActionCache                memory.Cache // for iam engine
	LocalUnmarshaledExpressionCache *gocache.Cache
	LocalGroupSystemAuthTypeCache   *gocache.Cache
	LocalActionDetailCache          memory.Cache
	LocalSubjectBlackListCache      memory.Cache
	LocalResourceTypePKCache        memory.Cache
	LocalThinResourceTypeCache      memory.Cache

	RemoteResourceCache     *redis.Cache
	ResourceTypeCache       *redis.Cache
	SubjectDepartmentCache  *redis.Cache
	SubjectPKCache          *redis.Cache
	SubjectSystemGroupCache *redis.Cache

	SystemCache         *redis.Cache
	ActionPKCache       *redis.Cache
	ActionDetailCache   *redis.Cache
	ActionListCache     *redis.Cache
	ResourceTypePKCache *redis.Cache

	PolicyCache                  *redis.Cache
	GroupResourcePolicyCache     *redis.Cache
	ExpressionCache              *redis.Cache
	TemporaryPolicyCache         *redis.Cache
	GroupSystemAuthTypeCache     *redis.Cache
	GroupActionResourceCache     *redis.Cache
	SubjectActionExpressionCache *redis.Cache

	LocalPolicyCache          *gocache.Cache
	LocalExpressionCache      *gocache.Cache
	LocalTemporaryPolicyCache *gocache.Cache
	ChangeListCache           *redis.Cache

	ActionCacheCleaner            *cleaner.CacheCleaner
	ActionListCacheCleaner        *cleaner.CacheCleaner
	ResourceTypeCacheCleaner      *cleaner.CacheCleaner
	SubjectDepartmentCacheCleaner *cleaner.CacheCleaner
	SystemCacheCleaner            *cleaner.CacheCleaner
)

// InitCaches
// Cache should only know about get/retrieve data
// ! DO NOT CARE ABOUT WHAT THE DATA WILL BE USED FOR
func InitCaches(disabled bool) {

	LocalAppCodeAppSecretCache = gocache.New(12*time.Hour, 5*time.Minute)

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

	UserAccessCache = redis.NewCache(
		"useraccess",
		120*time.Second,
	)

	// cleaner
	TestCacheCleaner = cleaner.NewCacheCleaner("TestCacheCleaner", testCacheDeleter{})
	go TestCacheCleaner.Run()

}

// InitCacheForContainer Init Cache for Container
func InitCacheForContainer() {
	GoCacheClient = gocache.New(time.Minute*5, 0)
}

// GetCacheForContainer Get Cache For Container
func GetCacheForContainer() *gocache.Cache {
	return GoCacheClient
}

func init() {
	InitCaches(false)
}

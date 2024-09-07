package memory

import (
	"context"
	"time"

	"github.com/fengzhongzhu1621/xgo/cache"
)

// RetrieveFunc is the type of the retrieve function.
// it retrieves the value from database, redis, apis, etc.
type RetrieveFunc func(ctx context.Context, key cache.Key) (interface{}, error)

// Cache is the interface for the cache.
type Cache interface {
	Get(ctx context.Context, key cache.Key) (interface{}, error)
	Set(ctx context.Context, key cache.Key, data interface{})

	GetString(ctx context.Context, key cache.Key) (string, error)
	GetBool(ctx context.Context, key cache.Key) (bool, error)
	GetInt(ctx context.Context, key cache.Key) (int, error)
	GetInt8(ctx context.Context, key cache.Key) (int8, error)
	GetInt16(ctx context.Context, key cache.Key) (int16, error)
	GetInt32(ctx context.Context, key cache.Key) (int32, error)
	GetInt64(ctx context.Context, key cache.Key) (int64, error)
	GetUint(ctx context.Context, key cache.Key) (uint, error)
	GetUint8(ctx context.Context, key cache.Key) (uint8, error)
	GetUint16(ctx context.Context, key cache.Key) (uint16, error)
	GetUint32(ctx context.Context, key cache.Key) (uint32, error)
	GetUint64(ctx context.Context, key cache.Key) (uint64, error)
	GetFloat32(ctx context.Context, key cache.Key) (float32, error)
	GetFloat64(ctx context.Context, key cache.Key) (float64, error)
	GetTime(ctx context.Context, key cache.Key) (time.Time, error)

	Delete(ctx context.Context, key cache.Key) error
	Exists(ctx context.Context, key cache.Key) bool

	DirectGet(ctx context.Context, key cache.Key) (interface{}, bool)

	Disabled() bool
}

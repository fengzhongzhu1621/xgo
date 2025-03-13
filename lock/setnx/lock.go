package setnx

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
)

type lock struct {
	cache      *redis.Client
	key        string
	needUnlock bool // 是否需要释放key
	isFirst    bool
}

type Locker interface {
	Lock(key string, expire time.Duration) (locked bool, err error)
	Unlock() error
}

func NewLocker(cache *redis.Client) Locker {
	return &lock{
		isFirst:    false,
		cache:      cache,
		key:        "",
		needUnlock: false,
	}
}

func (l *lock) Lock(key string, expire time.Duration) (locked bool, err error) {
	if l.isFirst {
		return false, fmt.Errorf("repeat lock")
	}
	l.isFirst = true
	l.key = key

	uuid := xid.New().String()
	locked, err = l.cache.SetNX(context.Background(), l.key, uuid, expire).Result()
	if locked {
		l.needUnlock = true
	}
	return locked, err
}

func (l *lock) Unlock() error {
	if !l.needUnlock {
		return nil
	}
	return l.cache.Del(context.Background(), l.key).Err()
}

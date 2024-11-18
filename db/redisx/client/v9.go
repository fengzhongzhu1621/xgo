package client

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var (
	initOnce sync.Once
	rdsV9    *redis.Client
)

func InitRedisV9Client(ctx context.Context, cfg *RedisConfig) {
	opts, err := redis.ParseURL(cfg.DSN())
	if err != nil {
		log.Fatalf("redis parse url error: %s", err.Error())
	}

	opts.DialTimeout = time.Duration(dialTimeout) * time.Second
	opts.ReadTimeout = time.Duration(readTimeout) * time.Second
	opts.WriteTimeout = time.Duration(writeTimeout) * time.Second
	opts.ConnMaxIdleTime = time.Duration(idleTimeout) * time.Second
	opts.PoolSize = poolSizeMultiple * runtime.NumCPU()
	opts.MinIdleConns = minIdleConnectionMultiple * runtime.NumCPU()

	initOnce.Do(func() {
		rdsV9 = redis.NewClient(opts)
		if _, err = rds.Ping(ctx).Result(); err != nil {
			log.Fatalf("redis connect error: %s", err.Error())
		} else {
			log.Infof("redis: %s:%d/%d connected", cfg.Host, cfg.Port, cfg.DB)
		}
	})
}

func GetDefaultRedisV9Client() *redis.Client {
	return rdsV9
}

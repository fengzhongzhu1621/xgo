package client

import (
	"context"
	"runtime"
	"sync"
	"time"

	redisV9 "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var (
	// 确保Redis客户端只初始化一次
	initOnce sync.Once
	// 确保Redis客户端只初始化一次
	rdsV9 *redisV9.Client
)

// InitRedisV9Client 初始化Redis客户端
func InitRedisV9Client(ctx context.Context, cfg *RedisConfig) {
	// 解析配置中的DSN字符串
	opts, err := redisV9.ParseURL(cfg.DSN())
	if err != nil {
		log.Fatalf("redis parse url error: %s", err.Error())
	}

	// 设置连接到Redis服务器的超时时间
	opts.DialTimeout = time.Duration(dialTimeout) * time.Second
	// 设置从Redis服务器读取数据的超时时间
	opts.ReadTimeout = time.Duration(readTimeout) * time.Second
	// 设置向Redis服务器写入数据的超时时间
	opts.WriteTimeout = time.Duration(writeTimeout) * time.Second
	// 设置连接在空闲状态下的最大时间，超过这个时间连接将被关闭
	opts.ConnMaxIdleTime = time.Duration(idleTimeout) * time.Second
	// 设置连接池的大小，即最多可以同时存在的连接数，乘以当前机器的CPU核心数
	opts.PoolSize = poolSizeMultiple * runtime.NumCPU()
	// 设置连接池中最少保持的空闲连接数
	opts.MinIdleConns = minIdleConnectionMultiple * runtime.NumCPU()

	initOnce.Do(func() {
		rdsV9 = redisV9.NewClient(opts)
		if _, err = rds.Ping(ctx).Result(); err != nil {
			log.Fatalf("redis connect error: %s", err.Error())
		} else {
			log.Infof(
				"redis connected %s [db=%d, dialTimeout=%s, readTimeout=%s, writeTimeout=%s, poolSize=%d, minIdleConns=%d, idleTimeout=%d]",
				cfg.Host,
				cfg.DB,
				opts.DialTimeout,
				opts.ReadTimeout,
				opts.WriteTimeout,
				opts.PoolSize,
				opts.MinIdleConns,
				opts.MinIdleConns,
			)
		}
	})
}

func GetDefaultRedisV9Client() *redisV9.Client {
	return rdsV9
}

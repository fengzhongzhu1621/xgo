package client

import (
	"context"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

const (
	// 尝试连接超时 单位：s
	dialTimeout = 2
	// 读超时 单位：s
	readTimeout = 1
	// 写超时 单位：s
	writeTimeout = 1
	// 闲置超时 单位: s
	idleTimeout = 3 * 60
	// 连接池大小 * 核
	poolSizeMultiple = 20
	// 最小空闲连接数 * 核
	minIdleConnectionMultiple = 10
)

var (
	rds *redis.Client
	mq  *redis.Client
)

var (
	redisClientInitOnce   sync.Once
	mqRedisClientInitOnce sync.Once
)

func newStandaloneClient(redisConfig *Redis) *redis.Client {
	opt := &redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	}

	// set default options
	opt.DialTimeout = time.Duration(dialTimeout) * time.Second
	opt.ReadTimeout = time.Duration(readTimeout) * time.Second
	opt.WriteTimeout = time.Duration(writeTimeout) * time.Second
	opt.PoolSize = poolSizeMultiple * runtime.NumCPU()
	opt.MinIdleConns = minIdleConnectionMultiple * runtime.NumCPU()
	opt.IdleTimeout = time.Duration(idleTimeout) * time.Second

	// set custom options, from config.yaml
	if redisConfig.DialTimeout > 0 {
		opt.DialTimeout = time.Duration(redisConfig.DialTimeout) * time.Second
	}
	if redisConfig.ReadTimeout > 0 {
		opt.ReadTimeout = time.Duration(redisConfig.ReadTimeout) * time.Second
	}
	if redisConfig.WriteTimeout > 0 {
		opt.WriteTimeout = time.Duration(redisConfig.WriteTimeout) * time.Second
	}

	if redisConfig.PoolSize > 0 {
		opt.PoolSize = redisConfig.PoolSize
	}
	if redisConfig.MinIdleConns > 0 {
		opt.MinIdleConns = redisConfig.MinIdleConns
	}

	log.Infof(
		"connect to redis: "+
			"%s [db=%d, dialTimeout=%s, readTimeout=%s, writeTimeout=%s, poolSize=%d, minIdleConns=%d, idleTimeout=%s]",
		opt.Addr,
		opt.DB,
		opt.DialTimeout,
		opt.ReadTimeout,
		opt.WriteTimeout,
		opt.PoolSize,
		opt.MinIdleConns,
		opt.IdleTimeout,
	)

	return redis.NewClient(opt)
}

func newSentinelClient(redisConfig *Redis) *redis.Client {
	sentinelAddrs := strings.Split(redisConfig.SentinelAddr, ",")
	opt := &redis.FailoverOptions{
		MasterName:    redisConfig.MasterName,
		SentinelAddrs: sentinelAddrs,
		DB:            redisConfig.DB,
		Password:      redisConfig.Password,
	}

	if redisConfig.SentinelPassword != "" {
		opt.SentinelPassword = redisConfig.SentinelPassword
	}

	// set default options
	opt.DialTimeout = 2 * time.Second
	opt.ReadTimeout = 1 * time.Second
	opt.WriteTimeout = 1 * time.Second
	opt.PoolSize = 20 * runtime.NumCPU()
	opt.MinIdleConns = 10 * runtime.NumCPU()
	opt.IdleTimeout = 3 * time.Minute

	// set custom options, from config.yaml
	if redisConfig.DialTimeout > 0 {
		opt.DialTimeout = time.Duration(redisConfig.DialTimeout) * time.Second
	}
	if redisConfig.ReadTimeout > 0 {
		opt.ReadTimeout = time.Duration(redisConfig.ReadTimeout) * time.Second
	}
	if redisConfig.WriteTimeout > 0 {
		opt.WriteTimeout = time.Duration(redisConfig.WriteTimeout) * time.Second
	}

	if redisConfig.PoolSize > 0 {
		opt.PoolSize = redisConfig.PoolSize
	}
	if redisConfig.MinIdleConns > 0 {
		opt.MinIdleConns = redisConfig.MinIdleConns
	}

	return redis.NewFailoverClient(opt)
}

func initRedisClient(redisConfig *Redis) (cli *redis.Client) {
	switch redisConfig.Type {
	case ModeStandalone:
		cli = newStandaloneClient(redisConfig)
	case ModeSentinel:
		cli = newSentinelClient(redisConfig)
	default:
		panic("init redis client fail, invalid redis.id, should be `standalone` or `sentinel`")
	}

	_, err := cli.Ping(context.TODO()).Result()
	if err != nil {
		log.WithError(err).Error("connect to redis fail")
		panic(err)
	}
	return cli
}

// InitRedisClient ...
func InitRedisClient(redisConfig *Redis) {
	if rds == nil {
		redisClientInitOnce.Do(func() {
			rds = initRedisClient(redisConfig)
		})
	}
}

// InitMQRedisClient ...
func InitMQRedisClient(redisConfig *Redis) {
	if mq == nil {
		mqRedisClientInitOnce.Do(func() {
			mq = initRedisClient(redisConfig)
		})
	}
}

// GetDefaultRedisClient 获取默认的Redis实例
func GetDefaultRedisClient() *redis.Client {
	return rds
}

// GetDefaultMQRedisClient 获取默认的MQ Redis实例
func GetDefaultMQRedisClient() *redis.Client {
	return mq
}

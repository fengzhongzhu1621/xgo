package client

import (
	"context"
	"errors"
	"strings"

	"github.com/go-redis/redis/v8"
)

// TestConnection 测试 redis 连接
func TestConnection(redisConfig *Redis) error {
	var rds *redis.Client
	switch redisConfig.Type {
	case ModeStandalone:
		opt := &redis.Options{
			Addr:     redisConfig.Addr,
			Password: redisConfig.Password,
			DB:       redisConfig.DB,
			// just for live test
			PoolSize: 1,
		}

		rds = redis.NewClient(opt)
	case ModeSentinel:
		sentinelAddrs := strings.Split(redisConfig.SentinelAddr, ",")
		opt := &redis.FailoverOptions{
			MasterName:    redisConfig.MasterName,
			SentinelAddrs: sentinelAddrs,
			DB:            redisConfig.DB,
			Password:      redisConfig.Password,
			PoolSize:      1,
		}

		if redisConfig.SentinelPassword != "" {
			opt.SentinelPassword = redisConfig.SentinelPassword
		}

		rds = redis.NewFailoverClient(opt)
	default:
		return errors.New("invalid redis ID, should be `standalone` or `sentinel`")
	}

	defer rds.Close()

	// 测试连接
	_, err := rds.Ping(context.Background()).Result()

	return err
}

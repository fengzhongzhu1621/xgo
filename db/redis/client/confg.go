package client

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

// Config define redis config
type Config struct {
	Address          string
	Password         string
	Database         string
	MasterName       string
	SentinelPassword string
	// for datacollection, notify if the snapshot redis is in use
	Enable       string
	MaxOpenConns int
}

// NewFromConfig returns new redis client from config
func NewFromConfig(cfg Config) (*redis.Client, error) {
	dbNum, err := strconv.Atoi(cfg.Database)
	if nil != err {
		return nil, err
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 3000
	}

	var client *redis.Client
	if cfg.MasterName == "" {
		option := &redis.Options{
			Addr:     cfg.Address,
			Password: cfg.Password,
			DB:       dbNum,
			PoolSize: cfg.MaxOpenConns,
		}
		client = redis.NewClient(option)
	} else {
		hosts := strings.Split(cfg.Address, ",")
		option := &redis.FailoverOptions{
			MasterName:       cfg.MasterName,
			SentinelAddrs:    hosts,
			Password:         cfg.Password,
			DB:               dbNum,
			PoolSize:         cfg.MaxOpenConns,
			SentinelPassword: cfg.SentinelPassword,
		}
		client = redis.NewFailoverClient(option)
	}

	err = client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return client, err
}

package redis

import "github.com/go-redis/redis/v8"

// IsNilErr returns whether err is nil error
func IsNilErr(err error) bool {
	return redis.Nil == err
}

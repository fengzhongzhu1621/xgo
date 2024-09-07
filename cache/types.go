package cache

import "time"

type Cache interface {
	Parse() error
	Scan() (map[string][]string, error)
	Get(key string) string
	GetMany(key string) []string
	Set(key string, value string)
	SetMany(key string, value []string)
}

// RandomExtraExpirationDurationFunc is the type of the function generate extra expiration duration
type RandomExtraExpirationDurationFunc func() time.Duration

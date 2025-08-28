package backend

import "time"

// Backend is the interface that wraps the basic cache backend operations.
type Backend interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string) error
}

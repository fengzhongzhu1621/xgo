package config

import (
	"context"
	"errors"
	"sync"
)

var (
	globalKV KV = &noopKV{}
	lock        = sync.RWMutex{}
)

// ErrConfigNotSupport is not supported config error
var ErrConfigNotSupport = errors.New("xgo/config: not support")

var (
	configMap = make(map[string]KVConfig)
)

var _ KV = (*noopKV)(nil)

// KV defines a kv storage for config center.
type KV interface {
	// Put puts or updates config value by key.
	Put(ctx context.Context, key, val string, opts ...Option) error

	// Get returns config value by key.
	Get(ctx context.Context, key string, opts ...Option) (IResponse, error)

	// Del deletes config value by key.
	Del(ctx context.Context, key string, opts ...Option) error
}

// KVConfig defines a kv config interface.
type KVConfig interface {
	KV
	Watcher
	Name() string
}

// Watcher defines the interface of config center watch event.
type Watcher interface {
	// Watch watches the config key change event.
	Watch(ctx context.Context, key string, opts ...Option) (<-chan IResponse, error)
}

// GlobalKV returns an instance of kv config center.
func GlobalKV() KV {
	return globalKV
}

// SetGlobalKV sets the instance of kv config center.
func SetGlobalKV(kv KV) {
	globalKV = kv
}

// Register registers a kv config by its name.
func Register(c KVConfig) {
	lock.Lock()
	configMap[c.Name()] = c
	lock.Unlock()
}

// Get returns a kv config by name.
func Get(name string) KVConfig {
	lock.RLock()
	c := configMap[name]
	lock.RUnlock()
	return c
}

// noopKV is an empty implementation of KV interface.
type noopKV struct{}

// Put does nothing but returns nil.
func (kv *noopKV) Put(ctx context.Context, key, val string, opts ...Option) error {
	return nil
}

// Get returns not supported error.
func (kv *noopKV) Get(ctx context.Context, key string, opts ...Option) (IResponse, error) {
	return nil, ErrConfigNotSupport
}

// Del does nothing but returns nil.
func (kv *noopKV) Del(ctx context.Context, key string, opts ...Option) error {
	return nil
}

package config

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/fengzhongzhu1621/xgo/crypto/unmarshaler"
)

var (
	globalKV KV = &noopKV{}
	lock        = sync.RWMutex{}
)

// ErrConfigNotSupport is not supported config error
var ErrConfigNotSupport = errors.New("xgo/config: not support")

// KV defines a kv storage for config center.
type KV interface {
	// Put puts or updates config value by key.
	Put(ctx context.Context, key, val string, opts ...Option) error

	// Get returns config value by key.
	Get(ctx context.Context, key string, opts ...Option) (IResponse, error)

	// Del deletes config value by key.
	Del(ctx context.Context, key string, opts ...Option) error
}

// GlobalKV returns an instance of kv config center.
func GlobalKV() KV {
	return globalKV
}

// SetGlobalKV sets the instance of kv config center.
func SetGlobalKV(kv KV) {
	globalKV = kv
}

// GetString returns string value get from
// kv storage by key.
func GetString(key string) (string, error) {
	val, err := globalKV.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	return val.Value(), nil
}

// GetStringWithDefault returns string value get by key.
// If anything wrong, returns default value specified by input param def.
func GetStringWithDefault(key, def string) string {
	val, err := globalKV.Get(context.Background(), key)
	if err != nil {
		return def
	}
	return val.Value()
}

// GetInt returns int value get by key.
func GetInt(key string) (int, error) {
	val, err := globalKV.Get(context.Background(), key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val.Value())
}

// GetIntWithDefault returns int value get by key.
// If anything wrong, returns default value specified by input param def.
func GetIntWithDefault(key string, def int) int {
	val, err := globalKV.Get(context.Background(), key)
	if err != nil {
		return def
	}
	i, err := strconv.Atoi(val.Value())
	if err != nil {
		return def
	}
	return i
}

// GetWithUnmarshal gets the specific encoding data
// by key. the encoding type is defined by unmarshalName parameter.
func GetWithUnmarshal(key string, val interface{}, unmarshalName string) error {
	v, err := globalKV.Get(context.Background(), key)
	if err != nil {
		return err
	}
	return unmarshaler.GetUnmarshaler(unmarshalName).Unmarshal([]byte(v.Value()), val)
}

// GetWithUnmarshalProvider gets the specific encoding data by key
// the encoding type is defined by unmarshalName parameter
// the provider name is defined by provider parameter.
func GetWithUnmarshalProvider(key string, val interface{}, unmarshalName string, provider string) error {
	p := Get(provider)
	if p == nil {
		return fmt.Errorf("xgo/config: failed to get %s", provider)
	}
	v, err := p.Get(context.Background(), key)
	if err != nil {
		return err
	}
	return unmarshaler.GetUnmarshaler(unmarshalName).Unmarshal([]byte(v.Value()), val)
}

// GetJSON gets json data by key. The value will unmarshal into val parameter.
func GetJSON(key string, val interface{}) error {
	return GetWithUnmarshal(key, val, "json")
}

// GetJSONWithProvider gets json data by key. The value will unmarshal into val parameter
// the provider name is defined by provider parameter.
func GetJSONWithProvider(key string, val interface{}, provider string) error {
	return GetWithUnmarshalProvider(key, val, "json", provider)
}

// GetYAML gets yaml data by key. The value will unmarshal into val parameter.
func GetYAML(key string, val interface{}) error {
	return GetWithUnmarshal(key, val, "yaml")
}

// GetYAMLWithProvider gets yaml data by key. The value will unmarshal into val parameter
// the provider name is defined by provider parameter.
func GetYAMLWithProvider(key string, val interface{}, provider string) error {
	return GetWithUnmarshalProvider(key, val, "yaml", provider)
}

// GetTOML gets toml data by key. The value will unmarshal into val parameter.
func GetTOML(key string, val interface{}) error {
	return GetWithUnmarshal(key, val, "toml")
}

// GetTOMLWithProvider gets toml data by key. The value will unmarshal into val parameter
// the provider name is defined by provider parameter.
func GetTOMLWithProvider(key string, val interface{}, provider string) error {
	return GetWithUnmarshalProvider(key, val, "toml", provider)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////
// KVConfig defines a kv config interface.
var _ KV = (KVConfig)(nil)

var (
	configMap = make(map[string]KVConfig)
)

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

// ////////////////////////////////////////////////////////////////////////////////////////////////
var _ KV = (*noopKV)(nil)

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

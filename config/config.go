package config

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fengzhongzhu1621/xgo/crypto/unmarshaler"
)

// IConfig defines the common config interface. We can
// implement different config center by this interface.
type IConfig interface {
	// Load loads config.
	Load() error

	// Reload reloads config.
	Reload()

	// Get returns config by key.
	Get(string, interface{}) interface{}

	// Unmarshal deserializes the config into input param.
	Unmarshal(interface{}) error

	// IsSet returns if the config specified by key exists.
	IsSet(string) bool

	// GetInt returns int value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetInt(string, int) int

	// GetInt32 returns int32 value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetInt32(string, int32) int32

	// GetInt64 returns int64 value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetInt64(string, int64) int64

	// GetUint returns uint value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetUint(string, uint) uint

	// GetUint32 returns uint32 value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetUint32(string, uint32) uint32

	// GetUint64 returns uint64 value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetUint64(string, uint64) uint64

	// GetFloat32 returns float32 value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetFloat32(string, float32) float32

	// GetFloat64 returns float64 value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetFloat64(string, float64) float64

	// GetString returns string value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetString(string, string) string

	// GetBool returns bool value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetBool(string, bool) bool

	// Bytes returns config data as bytes.
	Bytes() []byte
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

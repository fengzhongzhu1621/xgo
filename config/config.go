package config

// IConfig defines the common config interface. We can
// implement different config center by this interface.
// Iconfig 接口为业务代码提供了获取配置项的标准接口，每种数据类型都有一个独立的接口，接口支持返回 default 值。
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

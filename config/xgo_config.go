package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/fengzhongzhu1621/xgo/crypto/unmarshaler"
	"github.com/fengzhongzhu1621/xgo/env"
	"github.com/fengzhongzhu1621/xgo/logging"
	"trpc.group/trpc-go/tnet/log"
)

var _ IConfig = (*TrpcConfig)(nil)

// TrpcConfig is used to parse yaml config file for trpc.
type TrpcConfig struct {
	id  string       // config identity
	msg WatchMessage // new to init message for notify only copy

	p         IDataProvider      // config provider
	path      string             // config name
	decoder   unmarshaler.ICodec // config codec
	expandEnv bool               // status for whether replace the variables in the configuration with environment variables

	// because function is not support comparable in singleton, so the following options work only for the first load
	watch     bool
	watchHook func(message WatchMessage)

	mutex sync.RWMutex
	value *entity // store config value
}

// WatchMessage change message
type WatchMessage struct {
	Provider  string // provider name
	Path      string // config path
	ExpandEnv bool   // expend env status
	Codec     string // codec
	Watch     bool   // status for start watch
	Value     []byte // config content diff ?
	Error     error  // load error message, success is empty string
}

// defaultNotifyChange default hook for notify config changed
var defaultWatchHook = func(message WatchMessage) {}

// SetDefaultWatchHook set default hook notify when config changed
func SetDefaultWatchHook(f func(message WatchMessage)) {
	defaultWatchHook = f
}

// entity 配置实体，包含配置内容的原始数据和解析后的数据
type entity struct {
	// 配置原始内容
	raw []byte // current binary data
	// 解析后的配置对象
	data interface{} // unmarshal type to use point type, save latest no error data
}

func newEntity() *entity {
	return &entity{
		data: make(map[string]interface{}),
	}
}

// newTrpcConfig create a new config instance
func newTrpcConfig(path string, opts ...LoadOption) (*TrpcConfig, error) {
	c := &TrpcConfig{
		path:      path,                         // 配置文件路径
		p:         GetProvider("file"),          // 获得 FileProvider 对象
		decoder:   unmarshaler.GetCodec("yaml"), // 获得 yaml 编解码器对象
		expandEnv: true,                         // 是否展开环境变量
		watch:     true,                         // 是否开启监听
		watchHook: func(message WatchMessage) {
			defaultWatchHook(message)
		},
	}
	for _, o := range opts {
		o(c)
	}

	if c.p == nil {
		return nil, ErrProviderNotExist
	}
	if c.decoder == nil {
		return nil, ErrCodecNotExist
	}

	c.msg.Provider = c.p.Name()
	c.msg.Path = c.path
	c.msg.Codec = c.decoder.Name()
	c.msg.ExpandEnv = c.expandEnv
	c.msg.Watch = c.watch

	// since reflect.String() cannot uniquely identify a type, this id is used as a preliminary judgment basis
	const idFormat = "provider:%s path:%s codec:%s env:%t watch:%t"
	c.id = fmt.Sprintf(idFormat, c.p.Name(), c.path, c.decoder.Name(), c.expandEnv, c.watch)
	return c, nil
}

// doWatch 监听回调函数，类型是 ProviderCallback
// 当配置文件内容发生变化后，触发此回调
// data 监听的文件内容
func (c *TrpcConfig) doWatch(data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.set(data)
}

func (c *TrpcConfig) set(data []byte) error {
	if c.expandEnv {
		data = env.ExpandEnv(data)
	}

	e := newEntity()
	e.raw = data                              // 配置原始内容
	err := c.decoder.Unmarshal(data, &e.data) // 解析后的配置对象
	if err != nil {
		return fmt.Errorf("trpc/config: failed to parse:%w, id:%s", err, c.id)
	}
	c.value = e
	return nil
}

func (c *TrpcConfig) notify(data []byte, err error) {
	m := c.msg

	m.Value = data
	if err != nil {
		m.Error = err
	}

	c.watchHook(m)
}

// Load loads config.
func (c *TrpcConfig) Load() error {
	if c.p == nil {
		return ErrProviderNotExist
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	data, err := c.p.Read(c.path)
	if err != nil {
		return fmt.Errorf("trpc/config failed to load error: %w config id: %s", err, c.id)
	}

	return c.set(data)
}

// Reload reloads config.
func (c *TrpcConfig) Reload() {
	if err := c.Load(); err != nil {
		logging.Tracef("trpc/config: failed to reload %s: %v", c.id, err)
	}
}

// init return config entity error when entity is empty and load run loads config once
func (c *TrpcConfig) init() error {
	c.mutex.RLock()
	if c.value != nil {
		c.mutex.RUnlock()
		return nil
	}
	c.mutex.RUnlock()

	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.value != nil {
		return nil
	}

	// 读取配置
	data, err := c.p.Read(c.path)
	if err != nil {
		return fmt.Errorf("trpc/config failed to load error: %w config id: %s", err, c.id)
	}
	// 解析配置内容为配置对象
	return c.set(data)
}

func (c *TrpcConfig) get() *entity {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.value != nil {
		return c.value
	}
	return newEntity()
}

// Get returns config value by key. If key is absent will return the default value.
func (c *TrpcConfig) Get(key string, defaultValue interface{}) interface{} {
	if v, ok := c.search(key); ok {
		return v
	}
	return defaultValue
}

// Unmarshal deserializes the config into input param.
func (c *TrpcConfig) Unmarshal(out interface{}) error {
	return c.decoder.Unmarshal(c.get().raw, out)
}

// Bytes returns original config data as bytes.
func (c *TrpcConfig) Bytes() []byte {
	return c.get().raw
}

// GetInt returns int value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *TrpcConfig) GetInt(key string, defaultValue int) int {
	return c.findWithDefaultValue(key, defaultValue).(int)
}

// GetInt32 returns int32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *TrpcConfig) GetInt32(key string, defaultValue int32) int32 {
	return c.findWithDefaultValue(key, defaultValue).(int32)
}

// GetInt64 returns int64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *TrpcConfig) GetInt64(key string, defaultValue int64) int64 {
	return c.findWithDefaultValue(key, defaultValue).(int64)
}

// GetUint returns uint value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *TrpcConfig) GetUint(key string, defaultValue uint) uint {
	return c.findWithDefaultValue(key, defaultValue).(uint)
}

// GetUint32 returns uint32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *TrpcConfig) GetUint32(key string, defaultValue uint32) uint32 {
	return c.findWithDefaultValue(key, defaultValue).(uint32)
}

// GetUint64 returns uint64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *TrpcConfig) GetUint64(key string, defaultValue uint64) uint64 {
	return c.findWithDefaultValue(key, defaultValue).(uint64)
}

// GetFloat64 returns float64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *TrpcConfig) GetFloat64(key string, defaultValue float64) float64 {
	return c.findWithDefaultValue(key, defaultValue).(float64)
}

// GetFloat32 returns float32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *TrpcConfig) GetFloat32(key string, defaultValue float32) float32 {
	return c.findWithDefaultValue(key, defaultValue).(float32)
}

// GetBool returns bool value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *TrpcConfig) GetBool(key string, defaultValue bool) bool {
	return c.findWithDefaultValue(key, defaultValue).(bool)
}

// GetString returns string value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *TrpcConfig) GetString(key string, defaultValue string) string {
	return c.findWithDefaultValue(key, defaultValue).(string)
}

// IsSet returns if the config specified by key exists.
func (c *TrpcConfig) IsSet(key string) bool {
	_, ok := c.search(key)
	return ok
}

// findWithDefaultValue ensures that the type of `value` is same as `defaultValue`
func (c *TrpcConfig) findWithDefaultValue(key string, defaultValue interface{}) (value interface{}) {
	v, ok := c.search(key)
	if !ok {
		return defaultValue
	}

	var err error
	switch defaultValue.(type) {
	case bool:
		v, err = cast.ToBoolE(v)
	case string:
		v, err = cast.ToStringE(v)
	case int:
		v, err = cast.ToIntE(v)
	case int32:
		v, err = cast.ToInt32E(v)
	case int64:
		v, err = cast.ToInt64E(v)
	case uint:
		v, err = cast.ToUintE(v)
	case uint32:
		v, err = cast.ToUint32E(v)
	case uint64:
		v, err = cast.ToUint64E(v)
	case float64:
		v, err = cast.ToFloat64E(v)
	case float32:
		v, err = cast.ToFloat32E(v)
	default:
	}

	if err != nil {
		return defaultValue
	}
	return v
}

func (c *TrpcConfig) search(key string) (interface{}, bool) {
	e := c.get()

	unmarshalledData, ok := e.data.(map[string]interface{})
	if !ok {
		return nil, false
	}

	subkeys := strings.Split(key, ".")
	value, err := search(unmarshalledData, subkeys)
	if err != nil {
		log.Debugf("trpc config: search key %s failed: %+v", key, err)
		return value, false
	}

	return value, true
}

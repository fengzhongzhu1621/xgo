package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fengzhongzhu1621/xgo/cast"
	"github.com/fengzhongzhu1621/xgo/config/entity"
	"github.com/fengzhongzhu1621/xgo/config/hooks"
	"github.com/fengzhongzhu1621/xgo/config/provider"
	"github.com/fengzhongzhu1621/xgo/crypto/unmarshaler"
	"github.com/fengzhongzhu1621/xgo/env"
	"github.com/fengzhongzhu1621/xgo/logging"
)

var _ IConfig = (*XGoConfig)(nil)

// XGoConfig is used to parse yaml config file for xgo.
type XGoConfig struct {
	id  string             // config identity
	msg hooks.WatchMessage // new to init message for notify only copy

	p         provider.IDataProvider // config provider
	path      string                 // config name
	decoder   unmarshaler.ICodec     // config codec
	expandEnv bool                   // status for whether replace the variables in the configuration with environment variables

	// because function is not support comparable in singleton, so the following options work only for the first load
	watch     bool
	watchHook hooks.WatchMessageHookFunc

	mutex sync.RWMutex
	value *entity.Entity // store config value
}

// newXGoConfig create a new config instance
func newXGoConfig(path string, opts ...LoadOption) (*XGoConfig, error) {
	c := &XGoConfig{
		p:         provider.GetProvider("file"), // 获得 FileProvider 对象
		path:      path,                         // 配置文件路径
		decoder:   unmarshaler.GetCodec("yaml"), // 获得 yaml 编解码器对象
		expandEnv: true,                         // 是否展开环境变量

		watch: true, // 是否开启监听
		watchHook: func(message hooks.WatchMessage) {
			defaultWatchHook := hooks.GetWatchMessageHook()
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

// doWatch 解析并记录配置内容
// 当配置文件内容发生变化后，触发此回调
// raw 监听的文件内容
func (c *XGoConfig) doWatch(raw []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.set(raw)
}

// set 解析并记录配置内容到 entity 对象中
func (c *XGoConfig) set(raw []byte) error {
	var (
		data = make(map[string]interface{})
	)

	// 解析环境变量
	if c.expandEnv {
		raw = env.ExpandEnv(raw)
	}

	// 记录配置内容到 Entity 中
	e := entity.NewEntity()
	e.SetRaw(raw)                         // 配置原始内容
	err := c.decoder.Unmarshal(raw, data) // 解析后的配置对象
	if err != nil {
		return fmt.Errorf("xgo/config: failed to parse:%w, id:%s", err, c.id)
	}
	e.SetData(data)

	c.value = e
	return nil
}

// notify 执行 WatchMessageHook
func (c *XGoConfig) notify(raw []byte, err error) {
	m := c.msg

	m.Value = raw
	if err != nil {
		m.Error = err
	}

	// 执行 watch message hook，传递WatchMessage对象，包含配置的原始内容
	c.watchHook(m)
}

// Load loads config.
func (c *XGoConfig) Load() error {
	if c.p == nil {
		return ErrProviderNotExist
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 读取配置内容
	data, err := c.p.Read(c.path)
	if err != nil {
		return fmt.Errorf("xgo/config failed to load error: %w config id: %s", err, c.id)
	}

	// 解析并记录配置内容到 entity 对象中
	return c.set(data)
}

// Reload reloads config.
func (c *XGoConfig) Reload() {
	if err := c.Load(); err != nil {
		logging.Tracef("xgo/config: failed to reload %s: %v", c.id, err)
	}
}

// init return config entity error when entity is empty and load run loads config once
func (c *XGoConfig) init() error {
	// 防止重复初始化
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

	// 读取配置内容
	data, err := c.p.Read(c.path)
	if err != nil {
		return fmt.Errorf("xgo/config failed to load error: %w config id: %s", err, c.id)
	}

	// 解析并记录配置内容到 entity 对象中
	return c.set(data)
}

// get 获取 entity 对象
func (c *XGoConfig) get() *entity.Entity {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.value != nil {
		return c.value
	}
	return entity.NewEntity()
}

// Get returns config value by key. If key is absent will return the default value.
func (c *XGoConfig) Get(key string, defaultValue interface{}) interface{} {
	// 根据 keys 查询配置中的 value 值
	if v, ok := c.search(key); ok {
		return v
	}

	return defaultValue
}

// Unmarshal deserializes the config into input param.
func (c *XGoConfig) Unmarshal(out interface{}) error {
	return c.decoder.Unmarshal(c.get().GetRaw(), out)
}

// Bytes returns original config data as bytes.
func (c *XGoConfig) Bytes() []byte {
	return c.get().GetRaw()
}

// GetInt returns int value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *XGoConfig) GetInt(key string, defaultValue int) int {
	return c.findWithDefaultValue(key, defaultValue).(int)
}

// GetInt32 returns int32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *XGoConfig) GetInt32(key string, defaultValue int32) int32 {
	return c.findWithDefaultValue(key, defaultValue).(int32)
}

// GetInt64 returns int64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *XGoConfig) GetInt64(key string, defaultValue int64) int64 {
	return c.findWithDefaultValue(key, defaultValue).(int64)
}

// GetUint returns uint value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *XGoConfig) GetUint(key string, defaultValue uint) uint {
	return c.findWithDefaultValue(key, defaultValue).(uint)
}

// GetUint32 returns uint32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *XGoConfig) GetUint32(key string, defaultValue uint32) uint32 {
	return c.findWithDefaultValue(key, defaultValue).(uint32)
}

// GetUint64 returns uint64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *XGoConfig) GetUint64(key string, defaultValue uint64) uint64 {
	return c.findWithDefaultValue(key, defaultValue).(uint64)
}

// GetFloat64 returns float64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *XGoConfig) GetFloat64(key string, defaultValue float64) float64 {
	return c.findWithDefaultValue(key, defaultValue).(float64)
}

// GetFloat32 returns float32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *XGoConfig) GetFloat32(key string, defaultValue float32) float32 {
	return c.findWithDefaultValue(key, defaultValue).(float32)
}

// GetBool returns bool value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *XGoConfig) GetBool(key string, defaultValue bool) bool {
	return c.findWithDefaultValue(key, defaultValue).(bool)
}

// GetString returns string value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *XGoConfig) GetString(key string, defaultValue string) string {
	return c.findWithDefaultValue(key, defaultValue).(string)
}

// IsSet returns if the config specified by key exists.
func (c *XGoConfig) IsSet(key string) bool {
	_, ok := c.search(key)
	return ok
}

// findWithDefaultValue ensures that the type of `value` is same as `defaultValue`
func (c *XGoConfig) findWithDefaultValue(key string, defaultValue interface{}) (value interface{}) {
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

// search 根据 keys 查询配置中的 value 值
func (c *XGoConfig) search(key string) (interface{}, bool) {
	// 获取 entity 对象
	e := c.get()

	// 获得解析后的配置内容
	unmarshalledData, ok := e.GetData().(map[string]interface{})
	if !ok {
		return nil, false
	}

	// 根据 keys 查询配置 keys 的 value
	subkeys := strings.Split(key, ".")
	value, err := searchByKeys(unmarshalledData, subkeys)
	if err != nil {
		logging.Debugf("xgo config: search key %s failed: %+v", key, err)
		return value, false
	}

	return value, true
}

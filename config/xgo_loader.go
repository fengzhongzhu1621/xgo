package config

import "sync"

// LoadOption defines the option function for loading configuration.
type LoadOption func(*TrpcConfig)

// TrpcConfigLoader is a config loader for trpc.
type TrpcConfigLoader struct {
	watchers sync.Map
}

// ////////////////////////////////////////////////////////////////////////////////////////////
// DefaultConfigLoader is the default config loader.
var DefaultConfigLoader = newTrpcConfigLoad()

func newTrpcConfigLoad() *TrpcConfigLoader {
	return &TrpcConfigLoader{}
}

// Load returns the config specified by input parameter.
func Load(path string, opts ...LoadOption) (IConfig, error) {
	return DefaultConfigLoader.Load(path, opts...)
}

// Reload reloads config data.
func Reload(path string, opts ...LoadOption) error {
	return DefaultConfigLoader.Reload(path, opts...)
}

// ////////////////////////////////////////////////////////////////////////////////////////////
// Load returns the config specified by input parameter.
func (loader *TrpcConfigLoader) Load(path string, opts ...LoadOption) (IConfig, error) {
	// 创建配置对象
	c, err := newTrpcConfig(path, opts...)
	if err != nil {
		return nil, err
	}

	// 创建一个 watcher 对象，与 provider 对象绑定
	w := &watcher{}
	// 是一个原子操作，用于根据指定的 key 是否存在来执行读取或存储操作。
	// 如果 key 存在，则返回其对应的现有值；
	// 如果 key 不存在，则存储给定的 key-value 对并返回这个新值。
	// 返回值：
	// actual：实际获取到的值。如果 key 存在，返回已存储的值；如果 key 不存在，返回传入的 value。
	// loaded：一个布尔值，指示 key 是否存在于 map 中。存在为 true，不存在为 false。
	i, loaded := loader.watchers.LoadOrStore(c.p, w)
	if !loaded {
		// 监听配置变更，执行回调函数 w.watch，更新s.items 和 更新对应的 TrpcConfig 配置
		c.p.Watch(w.watch)
	} else {
		w = i.(*watcher)
	}

	// 获取或存储set配置对象
	s := w.getOrCreate(c.path)

	// 获取或存储 TrpcConfig 配置对象
	c = s.getOrStore(c)

	// 初始化 TrpcConfig 配置
	if err = c.init(); err != nil {
		return nil, err
	}

	return c, nil
}

// Reload reloads config data.
func (loader *TrpcConfigLoader) Reload(path string, opts ...LoadOption) error {
	c, err := newTrpcConfig(path, opts...)
	if err != nil {
		return err
	}

	// 获得 provider 对象绑定的 watcher
	v, ok := loader.watchers.Load(c.p)
	if !ok {
		return ErrConfigNotExist
	}
	w := v.(*watcher)

	// 获得 set 管理对象
	s := w.get(path)
	if s == nil {
		return ErrConfigNotExist
	}

	// 获得 set 管理的 TrpcConfig 对象
	oc := s.get(c.id)
	if oc == nil {
		return ErrConfigNotExist
	}

	// 重新读取配置
	return oc.Load()
}

// watch manage one data provider
type watcher struct {
	sets sync.Map // *set
}

// get config item by path
func (w *watcher) get(path string) *set {
	if i, ok := w.sets.Load(path); ok {
		return i.(*set)
	}
	return nil
}

// getOrCreate get config item by path if not exist and create and return
func (w *watcher) getOrCreate(path string) *set {
	i, _ := w.sets.LoadOrStore(path, &set{})
	return i.(*set)
}

// watch 监听回调函数，类型是 ProviderCallback
// 当配置文件内容发生变化后，触发此回调
// path 监听的文件路径
// data 监听的文件内容
func (w *watcher) watch(path string, data []byte) {
	if v := w.get(path); v != nil {
		v.watch(data)
	}
}

// set manages configs with same provider and name with different type
// used config.id as unique identifier
type set struct {
	path  string
	mutex sync.RWMutex
	items []*TrpcConfig
}

// get data
func (s *set) get(id string) *TrpcConfig {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, v := range s.items {
		if v.id == id {
			return v
		}
	}
	return nil
}

func (s *set) getOrStore(tc *TrpcConfig) *TrpcConfig {
	if v := s.get(tc.id); v != nil {
		return v
	}

	s.mutex.Lock()
	for _, item := range s.items {
		if item.id == tc.id {
			s.mutex.Unlock()
			return item
		}
	}
	// not found and add
	s.items = append(s.items, tc)
	s.mutex.Unlock()
	return tc
}

// watch data change, delete no watch model config and update watch model config and target notify
// data 监听的文件内容
func (s *set) watch(data []byte) {
	var items []*TrpcConfig
	var del []*TrpcConfig
	s.mutex.Lock()
	for _, v := range s.items {
		if v.watch {
			items = append(items, v)
		} else {
			del = append(del, v)
		}
	}
	s.items = items
	s.mutex.Unlock()

	for _, item := range items {
		// 保存配置的原始内容
		err := item.doWatch(data)
		item.notify(data, err)
	}

	for _, item := range del {
		item.notify(data, nil)
	}
}

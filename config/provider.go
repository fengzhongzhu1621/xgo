package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/fengzhongzhu1621/xgo/logging"
	"github.com/fsnotify/fsnotify"
	"trpc.group/trpc-go/tnet/log"
)

// providerMap 存储所有注册的数据提供者，以名称作为键
var providerMap = make(map[string]IDataProvider)

// ProviderCallback 是配置提供者处理配置变更的回调函数类型
// 参数1: 配置文件路径
// 参数2: 配置文件内容的字节数组
type ProviderCallback func(string, []byte)

// IDataProvider 定义了通用的数据提供者接口
// 我们可以实现这个接口来定义不同的数据提供者（如文件、TConf、ETCD、configmap等）
// 并通过编解码器将配置数据解析为标准格式（如json、toml、yaml等）
type IDataProvider interface {
	// Name 返回数据提供者的名称
	Name() string

	// Read 读取指定路径的文件，返回其内容的字节数组
	Read(string) ([]byte, error)

	// Watch 监听配置变更，变更将通过回调函数处理
	Watch(ProviderCallback)
}

// RegisterProvider 根据名称注册一个数据提供者
func RegisterProvider(p IDataProvider) {
	providerMap[p.Name()] = p
}

// GetProvider 根据名称获取数据提供者
func GetProvider(name string) IDataProvider {
	return providerMap[name]
}

// 编译时检查 FileProvider 是否实现了 IDataProvider 接口
var _ IDataProvider = (*FileProvider)(nil)

// FileProvider 是一个从文件系统获取配置的配置提供者
type FileProvider struct {
	disabledWatcher bool                  // 是否禁用文件监听器
	watcher         *fsnotify.Watcher     // 文件系统监听器
	cb              chan ProviderCallback // 回调函数通道
	cache           map[string]string     // 文件路径缓存，key为清理后的路径，value为原始路径
	modTime         map[string]int64      // 文件修改时间缓存
	mu              sync.RWMutex          // 读写锁，保护并发访问
}

// Name 返回文件提供者的名称
func (*FileProvider) Name() string {
	return "file"
}

// newFileProvider 创建一个新的文件提供者实例
func newFileProvider() *FileProvider {
	fp := &FileProvider{
		cb:              make(chan ProviderCallback), // 初始化回调函数通道
		disabledWatcher: true,                        // 默认禁用监听器
		cache:           make(map[string]string),     // 初始化路径缓存
		modTime:         make(map[string]int64),      // 初始化修改时间缓存
	}
	// 尝试创建文件系统监听器
	watcher, err := fsnotify.NewWatcher()
	if err == nil {
		fp.disabledWatcher = false // 启用监听器
		fp.watcher = watcher
		go fp.run() // 启动监听协程
		return fp
	}
	log.Debugf("fsnotify.NewWatcher err: %+v", err)
	return fp
}

// Read 读取指定路径的文件，返回其内容的字节数组
func (fp *FileProvider) Read(path string) ([]byte, error) {
	// 如果监听器未被禁用，则添加文件目录到监听列表
	if !fp.disabledWatcher {
		if err := fp.watcher.Add(filepath.Dir(path)); err != nil {
			return nil, err
		}
		// 缓存文件路径映射
		fp.mu.Lock()
		fp.cache[filepath.Clean(path)] = path
		fp.mu.Unlock()
	}
	// 读取文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		logging.Tracef("Failed to read file %v", err)
		return nil, err
	}
	return data, nil
}

// Watch 监听配置变更，变更将通过回调函数处理
func (fp *FileProvider) Watch(cb ProviderCallback) {
	// 只有在监听器未被禁用时才添加回调函数
	if !fp.disabledWatcher {
		fp.cb <- cb
	}
}

// run 启动文件监听循环，处理回调函数注册和文件变更事件
func (fp *FileProvider) run() {
	fn := make([]ProviderCallback, 0) // 存储所有注册的回调函数
	for {
		select {
		case i := <-fp.cb:
			// 接收到新的回调函数，添加到列表中
			fn = append(fn, i)
		case e := <-fp.watcher.Events:
			// 接收到文件系统事件，检查是否为修改事件
			if t, ok := fp.isModified(e); ok {
				fp.trigger(e, t, fn) // 触发所有回调函数
			}
		}
	}
}

// isModified 检查文件是否被修改
// 返回修改时间和是否修改的布尔值
func (fp *FileProvider) isModified(e fsnotify.Event) (int64, bool) {
	// 只处理写入事件
	if e.Op&fsnotify.Write != fsnotify.Write {
		return 0, false
	}
	fp.mu.RLock()
	defer fp.mu.RUnlock()

	// 检查文件是否在缓存中（即是否被监听）
	if _, ok := fp.cache[filepath.Clean(e.Name)]; !ok {
		return 0, false
	}
	// 获取文件信息
	fi, err := os.Stat(e.Name)
	if err != nil {
		return 0, false
	}
	// 检查修改时间是否更新
	if fi.ModTime().Unix() > fp.modTime[e.Name] {
		return fi.ModTime().Unix(), true
	}
	return 0, false
}

// trigger 触发所有注册的回调函数
func (fp *FileProvider) trigger(e fsnotify.Event, t int64, fn []ProviderCallback) {
	// 读取修改后的文件内容
	data, err := os.ReadFile(e.Name)
	if err != nil {
		return
	}
	// 更新缓存中的修改时间
	fp.mu.Lock()
	path := fp.cache[filepath.Clean(e.Name)] // 获取原始路径
	fp.modTime[e.Name] = t                   // 更新修改时间
	fp.mu.Unlock()

	// 并发调用所有回调函数
	for _, f := range fn {
		go f(path, data)
	}
}

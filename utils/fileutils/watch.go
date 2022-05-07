package fileutils

import (
	"log"
	"path/filepath"
	"sync"

	jww "xgo/log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
)

var SupportedExts = []string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "tfvars", "dotenv", "env", "ini"}
var SupportedRemoteProviders = []string{"etcd", "consul", "firestore"}

type IWatch interface {
	// 根据配置文件路径，读取配置
	ReadInConfig() error
	// 监听配置文件，如果有变化则调用 ReadInConfig() 重新加载配置
	WatchConfig()
}

type Watcher struct {
	IWatch

	// The filesystem to read config from.
	fs afero.Fs

	configName  string   // 配置文件名
	configPaths []string // 在指定路径下搜索配置文件
	configType  string   // 配置文件类型

	configFile string // 搜索到的配置文件的完整路径

	logger jww.Logger

	onConfigChange func(fsnotify.Event)
}

func (v *Watcher) getConfigFile() (string, error) {
	if v.configFile == "" {
		cf, err := FindConfigFile(v.fs, v.configPaths, v.configName, SupportedExts, v.configType)
		if err != nil {
			return "", err
		}
		v.configFile = cf
	}
	return v.configFile, nil
}

func (v *Watcher) getConfigType() string {
	if v.configType != "" {
		return v.configType
	}

	cf, err := v.getConfigFile()
	if err != nil {
		return ""
	}

	// 获得文件扩展名，包括 .
	ext := filepath.Ext(cf)
	// 去掉 .
	if len(ext) > 1 {
		return ext[1:]
	}

	return ""
}

// WatchConfig 监听配置文件的变化
func (v *Watcher) WatchConfig() {
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		// 创建文件监听器
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		// 获得配置文件
		// we have to watch the entire directory to pick up renames/atomic saves in a cross-platform way
		filename, err := v.getConfigFile()
		if err != nil {
			log.Printf("error: %v\n", err)
			initWG.Done()
			return
		}

		// 清理路径中的多余字符，比如 /// 或 ../ 或 ./
		configFile := filepath.Clean(filename)
		// 获得配置文件所在的目录
		configDir, _ := filepath.Split(configFile)
		// 会将所有路径的符号链接都解析出来。除此之外，它返回的路径，是直接可访问的。
		realConfigFile, _ := filepath.EvalSymlinks(filename)

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok { // 'Events' channel is closed
						eventsWG.Done()
						return
					}
					currentConfigFile, _ := filepath.EvalSymlinks(filename)
					// we only care about the config file with the following cases:
					// 1 - if the config file was modified or created
					// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
					const WriteOrCreateMask = fsnotify.Write | fsnotify.Create
					if (filepath.Clean(event.Name) == configFile &&
						event.Op&WriteOrCreateMask != 0) ||
						(currentConfigFile != "" && currentConfigFile != realConfigFile) {
						// 创建或更新了配置文件
						realConfigFile = currentConfigFile
						// 读取配置文件
						err := v.ReadInConfig()
						if err != nil {
							log.Printf("error reading config file: %v\n", err)
						}
						// 执行配置文件创建或修改后的自定义事件处理方法
						if v.onConfigChange != nil {
							v.onConfigChange(event)
						}
					} else if filepath.Clean(event.Name) == configFile &&
						event.Op&fsnotify.Remove != 0 {
						// 删除了配置文件
						eventsWG.Done()
						return
					}

				case err, ok := <-watcher.Errors:
					if ok { // 'Errors' channel is not closed
						log.Printf("watcher error: %v\n", err)
					}
					eventsWG.Done()
					return
				}
			}
		}()
		// 监听配置文件所在的目录
		watcher.Add(configDir)
		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
}

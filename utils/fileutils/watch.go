package fileutils

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"xgo/encoding"
	"xgo/encoding/dotenv"
	"xgo/encoding/hcl"
	"xgo/encoding/ini"
	"xgo/encoding/javaproperties"
	"xgo/encoding/json"
	"xgo/encoding/toml"
	"xgo/encoding/yaml"
	jww "xgo/log"
	"xgo/utils"
	"xgo/utils/stringutils"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
)

var SupportedExts = []string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "tfvars", "dotenv", "env", "ini"}
var SupportedRemoteProviders = []string{"etcd", "consul", "firestore"}

type Watcher struct {
	// The filesystem to read config from.
	fs afero.Fs

	// Delimiter that separates a list of keys
	// used to access a nested value in one go
	keyDelim string

	// Specific commands for ini parsing
	iniLoadOptions ini.LoadOptions

	envPrefix     string
	allowEmptyEnv bool

	configPaths []string // 在指定路径下搜索配置文件
	configName  string   // 配置文件名
	configType  string   // 配置文件后缀 (不带 .)

	configFile string                 // 搜索到的配置文件的完整路径
	config     map[string]interface{} // 从configFile中读取的配置

	logger jww.Logger

	onConfigChange func(fsnotify.Event) // 接受了文件创建和修改实践后自定义处理方法

	encoderRegistry *encoding.EncoderRegistry
	decoderRegistry *encoding.DecoderRegistry
}

// Option configures Viper using the functional options paradigm popularized by Rob Pike and Dave Cheney.
// If you're unfamiliar with this style,
// see https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html and
// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis.
type Option interface {
	apply(v *Watcher)
}

type optionFunc func(v *Watcher)

func (fn optionFunc) apply(v *Watcher) {
	fn(v)
}

// KeyDelimiter sets the delimiter used for determining key parts.
// By default it's value is ".".
func KeyDelimiter(d string) Option {
	return optionFunc(func(v *Watcher) {
		v.keyDelim = d
	})
}

// IniLoadOptions sets the load options for ini parsing.
func IniLoadOptions(in ini.LoadOptions) Option {
	return optionFunc(func(v *Watcher) {
		v.iniLoadOptions = in
	})
}

func (v *Watcher) OnConfigChange(run func(in fsnotify.Event)) {
	v.onConfigChange = run
}

func (v *Watcher) resetEncoding() {
	encoderRegistry := encoding.NewEncoderRegistry()
	decoderRegistry := encoding.NewDecoderRegistry()

	{
		codec := yaml.Codec{}

		encoderRegistry.RegisterEncoder("yaml", codec)
		decoderRegistry.RegisterDecoder("yaml", codec)

		encoderRegistry.RegisterEncoder("yml", codec)
		decoderRegistry.RegisterDecoder("yml", codec)
	}

	{
		codec := json.Codec{}

		encoderRegistry.RegisterEncoder("json", codec)
		decoderRegistry.RegisterDecoder("json", codec)
	}

	{
		codec := toml.Codec{}

		encoderRegistry.RegisterEncoder("toml", codec)
		decoderRegistry.RegisterDecoder("toml", codec)
	}

	{
		codec := hcl.Codec{}

		encoderRegistry.RegisterEncoder("hcl", codec)
		decoderRegistry.RegisterDecoder("hcl", codec)

		encoderRegistry.RegisterEncoder("tfvars", codec)
		decoderRegistry.RegisterDecoder("tfvars", codec)
	}

	{
		codec := ini.Codec{
			KeyDelimiter: v.keyDelim,
			LoadOptions:  v.iniLoadOptions,
		}

		encoderRegistry.RegisterEncoder("ini", codec)
		decoderRegistry.RegisterDecoder("ini", codec)
	}

	{
		codec := &javaproperties.Codec{
			KeyDelimiter: v.keyDelim,
		}

		encoderRegistry.RegisterEncoder("properties", codec)
		decoderRegistry.RegisterDecoder("properties", codec)

		encoderRegistry.RegisterEncoder("props", codec)
		decoderRegistry.RegisterDecoder("props", codec)

		encoderRegistry.RegisterEncoder("prop", codec)
		decoderRegistry.RegisterDecoder("prop", codec)
	}

	{
		codec := &dotenv.Codec{}

		encoderRegistry.RegisterEncoder("dotenv", codec)
		decoderRegistry.RegisterDecoder("dotenv", codec)

		encoderRegistry.RegisterEncoder("env", codec)
		decoderRegistry.RegisterDecoder("env", codec)
	}

	v.encoderRegistry = encoderRegistry
	v.decoderRegistry = decoderRegistry
}

// getConfigFile 在指定路径下查找配置文件的路径
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
		v.configFile = ext[1:]
		return v.configFile
	}

	return ""
}

// ReadInConfig 读取配置文件，返回反序列化后的内容
func (v *Watcher) ReadInConfig() error {
	// 获得配置文件路径
	v.logger.Info("attempting to read in config file")
	filename, err := v.getConfigFile()
	if err != nil {
		return err
	}

	// 判断文件后缀是否有效
	configType := v.getConfigType()
	if !stringutils.StringInSlice(configType, SupportedExts) {
		return UnsupportedConfigError(v.getConfigType())
	}

	// 读取配置文件
	v.logger.Debug("reading file", "file", filename)
	file, err := afero.ReadFile(v.fs, filename)
	if err != nil {
		return err
	}

	config := make(map[string]interface{})

	err = v.unmarshalReader(bytes.NewReader(file), config)
	if err != nil {
		return err
	}

	v.config = config
	return nil
}

func (v *Watcher) SetConfigFile(in string) {
	if in != "" {
		v.configFile = in
	}
}

func (v *Watcher) SetEnvPrefix(in string) {
	if in != "" {
		v.envPrefix = in
	}
}

func (v *Watcher) mergeWithEnvPrefix(in string) string {
	if v.envPrefix != "" {
		return strings.ToUpper(v.envPrefix + "_" + in)
	}

	return strings.ToUpper(in)
}

func (v *Watcher) AllowEmptyEnv(allowEmptyEnv bool) {
	v.allowEmptyEnv = allowEmptyEnv
}

// getEnv is a wrapper around os.Getenv which replaces characters in the original
// key. This allows env vars which have different keys than the config object
// keys.
func (v *Watcher) getEnv(key string) (string, bool) {
	val, ok := os.LookupEnv(key)

	return val, ok && (v.allowEmptyEnv || val != "")
}

// ConfigFileUsed returns the file used to populate the config registry.
func (v *Watcher) ConfigFileUsed() string {
	return v.configFile
}

// AddConfigPath adds a path for Viper to search for the config file in.
// Can be called multiple times to define multiple search paths.
// 添加搜索路径 .
func (v *Watcher) AddConfigPath(in string) {
	if in != "" {
		// 获得绝对路径
		absin := AbsPathify(in)

		v.logger.Info("adding path to search paths", "path", absin)
		if !stringutils.StringInSlice(absin, v.configPaths) {
			v.configPaths = append(v.configPaths, absin)
		}
	}
}

func (v *Watcher) unmarshalReader(in io.Reader, c map[string]interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(in)

	switch format := strings.ToLower(v.getConfigType()); format {
	case "yaml", "yml", "json", "toml", "hcl", "tfvars", "ini", "properties", "props", "prop", "dotenv", "env":
		// 解码文件内容，结果放到c中
		err := v.decoderRegistry.Decode(format, buf.Bytes(), c)
		if err != nil {
			return ConfigParseError{err}
		}
	}

	// 将字典的key转换为小写，返回原字典
	utils.InsensitiviseMap(c)

	return nil
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

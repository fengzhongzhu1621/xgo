package fileutils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"xgo/cast"
	"xgo/encoding"
	"xgo/encoding/dotenv"
	"xgo/encoding/hcl"
	"xgo/encoding/ini"
	"xgo/encoding/javaproperties"
	"xgo/encoding/json"
	"xgo/encoding/mapstructure"
	"xgo/encoding/toml"
	"xgo/encoding/yaml"
	jww "xgo/log"
	"xgo/remote"
	"xgo/utils"
	"xgo/utils/stringutils"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
)

var SupportedExts = []string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "tfvars", "dotenv", "env", "ini"}
var SupportedRemoteProviders = []string{"etcd", "consul", "firestore"}

type defaultRemoteProvider struct {
	provider      string
	endpoint      string
	path          string
	secretKeyring string
}

func (rp defaultRemoteProvider) Provider() string {
	return rp.provider
}

func (rp defaultRemoteProvider) Endpoint() string {
	return rp.endpoint
}

func (rp defaultRemoteProvider) Path() string {
	return rp.path
}

func (rp defaultRemoteProvider) SecretKeyring() string {
	return rp.secretKeyring
}

type Watcher struct {
	logger jww.Logger

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

	configFile        string                 // 搜索到的配置文件的完整路径
	config            map[string]interface{} // 从configFile中读取的配置，即本地文件配置文件
	configPermissions os.FileMode            // 配置文件模式

	onConfigChange func(fsnotify.Event) // 接受了文件创建和修改实践后自定义处理方法

	override map[string]interface{} // 保存设置的值

	// A set of remote providers to search for the configuration
	remoteProviders []*defaultRemoteProvider
	kvstore         map[string]interface{} // 远程配置

	automaticEnvApplied bool // 是否从环境变量获取value

	// 远程配置服务器客户端
	RemoteConfig remote.IRemoteConfigFactory

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

func (v *Watcher) OnConfigChange(run func(in fsnotify.Event)) {
	v.onConfigChange = run
}

// AutomaticEnv makes Viper check if environment variables match any of the existing keys
// (config, default or flags). If matching env vars are found, they are loaded into Viper.
func (v *Watcher) AutomaticEnv() {
	v.automaticEnvApplied = true
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

// getConfigFile 在指定路径下查找配置文件的路径.
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

// ReadInConfig 读取配置文件，返回反序列化后的内容.
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

	// 将配置文件内容解码为结构体对象
	config := make(map[string]interface{})
	err = v.unmarshalReader(bytes.NewReader(file), config)
	if err != nil {
		return err
	}

	v.config = config
	return nil
}

// ReadConfig will read a configuration file, setting existing keys to nil if the
// key does not exist in the file.
func (v *Watcher) ReadConfig(in io.Reader) error {
	v.config = make(map[string]interface{})
	return v.unmarshalReader(in, v.config)
}

// MergeInConfig merges a new configuration with an existing config.
func (v *Watcher) MergeInConfig() error {
	v.logger.Info("attempting to merge in config file")
	filename, err := v.getConfigFile()
	if err != nil {
		return err
	}

	if !stringutils.StringInSlice(v.getConfigType(), SupportedExts) {
		return UnsupportedConfigError(v.getConfigType())
	}

	file, err := afero.ReadFile(v.fs, filename)
	if err != nil {
		return err
	}

	return v.MergeConfig(bytes.NewReader(file))
}

// MergeConfig merges a new configuration with an existing config.
func (v *Watcher) MergeConfig(in io.Reader) error {
	cfg := make(map[string]interface{})
	if err := v.unmarshalReader(in, cfg); err != nil {
		return err
	}
	return v.MergeConfigMap(cfg)
}

// MergeConfigMap merges the configuration from the map given with an existing config.
// Note that the map given may be modified.
func (v *Watcher) MergeConfigMap(cfg map[string]interface{}) error {
	if v.config == nil {
		v.config = make(map[string]interface{})
	}
	utils.InsensitiviseMap(cfg)
	utils.MergeMaps(cfg, v.config, nil)
	return nil
}

// Marshal a map into Writer.
func (v *Watcher) marshalWriter(f afero.File, configType string) error {
	c := v.AllSettings()
	switch configType {
	case "yaml", "yml", "json", "toml", "hcl", "tfvars", "ini", "prop", "props", "properties", "dotenv", "env":
		b, err := v.encoderRegistry.Encode(configType, c)
		if err != nil {
			return ConfigMarshalError{err}
		}

		_, err = f.WriteString(string(b))
		if err != nil {
			return ConfigMarshalError{err}
		}
	}
	return nil
}

func (v *Watcher) ReadRemoteConfig() error {
	return v.getKeyValueConfig()
}

func (v *Watcher) WatchRemoteConfig() error {
	return v.watchKeyValueConfig()
}

func (v *Watcher) WatchRemoteConfigOnChannel() error {
	return v.watchKeyValueConfigOnChannel()
}

// getKeyValueConfig 从多个远程配置服务器，根据path获得远程配置，并转换为struct对象
// Retrieve the first found remote configuration.
func (v *Watcher) getKeyValueConfig() error {
	if v.RemoteConfig == nil {
		return RemoteConfigError("Enable the remote features by doing a blank import of the viper/remote package: '_ github.com/spf13/viper/remote'")
	}

	for _, rp := range v.remoteProviders {
		// 根据path获得远程配置，并转换为struct对象
		val, err := v.getRemoteConfig(rp)
		if err != nil {
			v.logger.Error(fmt.Errorf("get remote config: %w", err).Error())

			continue
		}

		v.kvstore = val

		return nil
	}
	return RemoteConfigError("No Files Found")
}

// 从单个远程服务器，根据path获得远程配置，并转换为struct对象.
func (v *Watcher) getRemoteConfig(provider remote.IRemoteProvider) (map[string]interface{}, error) {
	// 获得远程配置
	reader, err := v.RemoteConfig.Get(provider)
	if err != nil {
		return nil, err
	}
	// 将io.Reader转换为struct对象
	err = v.unmarshalReader(reader, v.kvstore)
	return v.kvstore, err
}

// Retrieve the first found remote configuration.
// 监听并获取远程配置服务器的值.
func (v *Watcher) watchKeyValueConfigOnChannel() error {
	for _, rp := range v.remoteProviders {
		// 监听远程配置服务器
		respc, _ := v.RemoteConfig.WatchChannel(rp)
		// Todo: Add quit channel
		go func(rc <-chan *remote.RemoteResponse) {
			for {
				b := <-rc
				reader := bytes.NewReader(b.Value)
				v.unmarshalReader(reader, v.kvstore)
			}
		}(respc)
		return nil
	}
	return RemoteConfigError("No Files Found")
}

// Retrieve the first found remote configuration.
func (v *Watcher) watchKeyValueConfig() error {
	for _, rp := range v.remoteProviders {
		val, err := v.watchRemoteConfig(rp)
		if err != nil {
			continue
		}
		v.kvstore = val
		return nil
	}
	return RemoteConfigError("No Files Found")
}

func (v *Watcher) watchRemoteConfig(provider remote.IRemoteProvider) (map[string]interface{}, error) {
	reader, err := v.RemoteConfig.Watch(provider)
	if err != nil {
		return nil, err
	}
	err = v.unmarshalReader(reader, v.kvstore)
	return v.kvstore, err
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

func (v *Watcher) SetFs(fs afero.Fs) {
	v.fs = fs
}

func (v *Watcher) SetConfigName(in string) {
	if in != "" {
		v.configName = in
		v.configFile = ""
	}
}

func (v *Watcher) SetConfigType(in string) {
	if in != "" {
		v.configType = in
	}
}

func (v *Watcher) SetConfigPermissions(perm os.FileMode) {
	v.configPermissions = perm.Perm()
}

// IniLoadOptions sets the load options for ini parsing.
func IniLoadOptions(in ini.LoadOptions) Option {
	return optionFunc(func(v *Watcher) {
		v.iniLoadOptions = in
	})
}

// mergeWithEnvPrefix 将环境变量的key加上前缀并大写.
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

// 判断远程配置对象是否存在.
func (v *Watcher) providerPathExists(p *defaultRemoteProvider) bool {
	for _, y := range v.remoteProviders {
		if reflect.DeepEqual(y, p) {
			return true
		}
	}
	return false
}

// AddRemoteProvider adds a remote configuration source.
// Remote Providers are searched in the order they are added.
// provider is a string value: "etcd", "consul" or "firestore" are currently supported.
// endpoint is the url.  etcd requires http://ip:port  consul requires ip:port
// path is the path in the k/v store to retrieve configuration
// To retrieve a config file called myapp.json from /configs/myapp.json
// you should set path to /configs and set config name (SetConfigName()) to "myapp"
// 添加远程配置对象 .
func (v *Watcher) AddRemoteProvider(providerName, endpoint, path string) error {
	// 判断是否支持远程key/value存储的类型
	if !stringutils.StringInSlice(providerName, SupportedRemoteProviders) {
		return UnsupportedRemoteProviderError(providerName)
	}
	if providerName != "" && endpoint != "" {
		v.logger.Info("adding remote provider", "provider", providerName, "endpoint", endpoint)
		// 新增一个远程配置对象
		rp := &defaultRemoteProvider{
			endpoint: endpoint,
			provider: providerName,
			path:     path,
		}
		if !v.providerPathExists(rp) {
			v.remoteProviders = append(v.remoteProviders, rp)
		}
	}
	return nil
}

// AddSecureRemoteProvider adds a remote configuration source.
// Secure Remote Providers are searched in the order they are added.
// provider is a string value: "etcd", "consul" or "firestore" are currently supported.
// endpoint is the url.  etcd requires http://ip:port  consul requires ip:port
// secretkeyring is the filepath to your openpgp secret keyring.  e.g. /etc/secrets/myring.gpg
// path is the path in the k/v store to retrieve configuration
// To retrieve a config file called myapp.json from /configs/myapp.json
// you should set path to /configs and set config name (SetConfigName()) to
// "myapp"
// Secure Remote Providers are implemented with github.com/bketelsen/crypt.
func (v *Watcher) AddSecureRemoteProvider(provider, endpoint, path, secretkeyring string) error {
	if !stringutils.StringInSlice(provider, SupportedRemoteProviders) {
		return UnsupportedRemoteProviderError(provider)
	}
	if provider != "" && endpoint != "" {
		v.logger.Info("adding remote provider", "provider", provider, "endpoint", endpoint)

		rp := &defaultRemoteProvider{
			endpoint:      endpoint,
			provider:      provider,
			path:          path,
			secretKeyring: secretkeyring,
		}
		if !v.providerPathExists(rp) {
			v.remoteProviders = append(v.remoteProviders, rp)
		}
	}
	return nil
}

// isPathShadowedInAutoEnv makes sure the given path is not shadowed somewhere
// in the environment, when automatic env is on.
// e.g., if "foo.bar" has a value in the environment, it “shadows”
//       "foo.bar.baz" in a lower-priority map
// 子路径是否是环境变量的key.
func (v *Watcher) isPathShadowedInAutoEnv(path []string) string {
	var parentKey string
	for i := 1; i < len(path); i++ {
		parentKey = strings.Join(path[0:i], v.keyDelim)
		// 将环境变量的key加上前缀并大写，然后获得环境变量的值
		if _, ok := v.getEnv(v.mergeWithEnvPrefix(parentKey)); ok {
			return parentKey
		}
	}
	return ""
}

// Given a key, find the value.
//
// Viper will check to see if an alias exists first.
// Viper will then check in the following order:
// flag, env, config file, key/value store.
// Lastly, if no value was found and flagDefault is true, and if the key
// corresponds to a flag, the flag's default value is returned.
//
// Note: this assumes a lower-cased key given.
// 根据key依次从环境变量，本地配置，远程配置获取值.
func (v *Watcher) find(lcaseKey string) interface{} {
	var (
		val    interface{}
		path   = strings.Split(lcaseKey, v.keyDelim)
		nested = len(path) > 1
	)

	// Set() override first
	val = utils.SearchMap(v.override, path)
	if val != nil {
		return val
	}
	if nested && utils.IsPathShadowedInDeepMap(path, v.override, v.keyDelim) != "" {
		return nil
	}

	// Env override next
	if v.automaticEnvApplied {
		// even if it hasn't been registered, if automaticEnv is used,
		// check any Get request
		if val, ok := v.getEnv(v.mergeWithEnvPrefix(lcaseKey)); ok {
			return val
		}
		if nested && v.isPathShadowedInAutoEnv(path) != "" {
			return nil
		}
	}

	// 从本地配置文件内容中获取
	val = utils.SearchIndexableWithPathPrefixes(v.config, path, v.keyDelim)
	if val != nil {
		return val
	}
	if nested && utils.IsPathShadowedInDeepMap(path, v.config, v.keyDelim) != "" {
		return nil
	}

	// K/V store next
	val = utils.SearchMap(v.kvstore, path)
	if val != nil {
		return val
	}
	if nested && utils.IsPathShadowedInDeepMap(path, v.kvstore, v.keyDelim) != "" {
		return nil
	}

	return nil
}

// Get can retrieve any value given the key to use.
// Get is case-insensitive for a key.
// Get has the behavior of returning the value associated with the first
// place from where it is set. Viper will check in the following order:
// override, flag, env, config file, key/value store, default
//
// Get returns an interface. For a specific value use one of the Get____ methods.
func (v *Watcher) Get(key string) interface{} {
	lcaseKey := strings.ToLower(key)
	// 根据key依次从环境变量，本地配置，远程配置获取值
	val := v.find(lcaseKey)
	if val == nil {
		return nil
	}

	return val
}

// GetString returns the value associated with the key as a string.
func (v *Watcher) GetString(key string) string {
	return cast.ToString(v.Get(key))
}

// GetBool returns the value associated with the key as a boolean.
func (v *Watcher) GetBool(key string) bool {
	return cast.ToBool(v.Get(key))
}

// GetInt returns the value associated with the key as an integer.
func (v *Watcher) GetInt(key string) int {
	return cast.ToInt(v.Get(key))
}

// GetInt32 returns the value associated with the key as an integer.
func (v *Watcher) GetInt32(key string) int32 {
	return cast.ToInt32(v.Get(key))
}

// GetInt64 returns the value associated with the key as an integer.
func (v *Watcher) GetInt64(key string) int64 {
	return cast.ToInt64(v.Get(key))
}

// GetUint returns the value associated with the key as an unsigned integer.
func (v *Watcher) GetUint(key string) uint {
	return cast.ToUint(v.Get(key))
}

// GetUint32 returns the value associated with the key as an unsigned integer.
func (v *Watcher) GetUint32(key string) uint32 {
	return cast.ToUint32(v.Get(key))
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (v *Watcher) GetUint64(key string) uint64 {
	return cast.ToUint64(v.Get(key))
}

// GetFloat64 returns the value associated with the key as a float64.
func (v *Watcher) GetFloat64(key string) float64 {
	return cast.ToFloat64(v.Get(key))
}

// GetTime returns the value associated with the key as time.
func (v *Watcher) GetTime(key string) time.Time {
	return cast.ToTime(v.Get(key))
}

// GetDuration returns the value associated with the key as a duration.
func (v *Watcher) GetDuration(key string) time.Duration {
	return cast.ToDuration(v.Get(key))
}

// GetIntSlice returns the value associated with the key as a slice of int values.
func (v *Watcher) GetIntSlice(key string) []int {
	return cast.ToIntSlice(v.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (v *Watcher) GetStringSlice(key string) []string {
	return cast.ToStringSlice(v.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (v *Watcher) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(v.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (v *Watcher) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(v.Get(key))
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (v *Watcher) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(v.Get(key))
}

// GetSizeInBytes returns the size of the value associated with the given key
// in bytes.
func (v *Watcher) GetSizeInBytes(key string) uint {
	sizeStr := cast.ToString(v.Get(key))
	return utils.ParseSizeInBytes(sizeStr)
}

// IsSet checks to see if the key has been set in any of the data locations.
// IsSet is case-insensitive for a key.
// 判断key是否在环境变量，本地配置，远程配置中存在.
func (v *Watcher) IsSet(key string) bool {
	lcaseKey := strings.ToLower(key)
	val := v.find(lcaseKey)
	return val != nil
}

// InConfig checks to see if the given key (or an alias) is in the config file.
// 支持根据字典的key和数组的索引进行搜索.
func (v *Watcher) InConfig(key string) bool {
	lcaseKey := strings.ToLower(key)
	path := strings.Split(lcaseKey, v.keyDelim)

	return utils.SearchIndexableWithPathPrefixes(v.config, path, v.keyDelim) != nil
}

// Set sets the value for the key in the override register.
// Set is case-insensitive for a key.
// Will be used instead of values obtained via
// flags, config file, ENV, default, or key/value store.
// 将字典的key转换为小写，返回新的字典.
func (v *Watcher) Set(key string, value interface{}) {
	value = utils.ToCaseInsensitiveValue(value)

	path := strings.Split(key, v.keyDelim)
	lastKey := strings.ToLower(path[len(path)-1])
	// 根据路径构造字典
	deepestMap := utils.DeepSearch(v.override, path[0:len(path)-1])

	// set innermost value
	deepestMap[lastKey] = value
}

// AllKeys 获得所有需要查询的key.
func (v *Watcher) AllKeys() []string {
	m := map[string]interface{}{}
	m = utils.FlattenAndMergeMap(m, v.override, "", v.keyDelim)
	m = utils.FlattenAndMergeMap(m, v.config, "", v.keyDelim)
	m = utils.FlattenAndMergeMap(m, v.kvstore, "", v.keyDelim)

	// convert set of paths to list
	a := make([]string, 0, len(m))
	for x := range m {
		a = append(a, x)
	}
	return a
}

// AllSettings 获得所有需要查询的key的值.
func (v *Watcher) AllSettings() map[string]interface{} {
	m := map[string]interface{}{}
	// start from the list of keys, and construct the map one value at a time
	for _, k := range v.AllKeys() {
		value := v.Get(k)
		if value == nil {
			// should not happen, since AllKeys() returns only keys holding a value,
			// check just in case anything changes
			continue
		}
		path := strings.Split(k, v.keyDelim)
		lastKey := strings.ToLower(path[len(path)-1])
		deepestMap := utils.DeepSearch(m, path[0:len(path)-1])
		// set innermost value
		deepestMap[lastKey] = value
	}
	return m
}

// UnmarshalKey takes a single key and unmarshals it into a Struct.
// 搜索key的值，转换为golang对象rawVal .
func (v *Watcher) UnmarshalKey(key string, rawVal interface{}, opts ...mapstructure.DecoderConfigOption) error {
	return mapstructure.Decode(v.Get(key), mapstructure.DefaultDecoderConfig(rawVal, opts...))
}

// Unmarshal unmarshals the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
func (v *Watcher) Unmarshal(rawVal interface{}, opts ...mapstructure.DecoderConfigOption) error {
	return mapstructure.Decode(v.AllSettings(), mapstructure.DefaultDecoderConfig(rawVal, opts...))
}

func (v *Watcher) UnmarshalExact(rawVal interface{}, opts ...mapstructure.DecoderConfigOption) error {
	config := mapstructure.DefaultDecoderConfig(rawVal, opts...)
	config.ErrorUnused = true

	return mapstructure.Decode(v.AllSettings(), config)
}

// 将io.Reader转换为struct对象.
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

// WatchConfig 监听配置文件的变化.
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

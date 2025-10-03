package config

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/fengzhongzhu1621/xgo/logging/zaplogger"
	"github.com/spf13/viper"
)

var loggerInitOnce sync.Once

var (
	cfgFile      string
	globalConfig *Config
)

func GetGlobalConfig() *Config {
	return globalConfig
}

// SetCfgFile 设置开发配置文件
func SetCfgFile() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	cfgFile := path.Join(dir, "config.yaml")
	viper.SetConfigFile(cfgFile)

	return cfgFile
}

// SetCfgFile 设置开发配置文件
func GetAdminCfgFilePath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	adminCfgFilePath := path.Join(dir, "admin.yaml")

	return adminCfgFilePath
}

// GetCfgFile 获取配置文件
func GetCfgFile() string {
	return cfgFile
}

func LoadConfig() {
	var err error

	// 设置默认配置文件
	SetCfgFile()

	// 读取并解析配置文件
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Error reading config file: %v", err))
	}

	// 将配置文件转换为全局配置对象
	globalConfig, err = Load(viper.GetViper())
	if err != nil {
		panic(fmt.Sprintf("Could not load configurations from file, error: %v", err))
	}
}

// InitLogger 初始化日志记录器，只能执行一次
func InitLogger(cache bool) {
	globalConfig := GetGlobalConfig()
	logger := globalConfig.Logger

	// 设置系统日志记录器
	zaplogger.InitSystemLogger(&logger.System)

	loggerInitOnce.Do(func() {
		// 设置 web 服务器日志记录器
		appLogger := zaplogger.NewZapJSONLogger(&logger.Web, cache)
		zaplogger.SetDbLogger(&zaplogger.DBLogger{appLogger})
	})
}

func init() {
	InitLogger(false)
}

func init() {
	LoadConfig()
}

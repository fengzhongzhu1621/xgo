package config

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

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

func init() {
	LoadConfig()
}

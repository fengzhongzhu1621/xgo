package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var globalConfig *Config

func GetGlobalConfig() *Config {
	return globalConfig
}

func init() {
	var err error
	globalConfig, err = Load(viper.GetViper())
	if err != nil {
		panic(fmt.Sprintf("Could not load configurations from file, error: %v", err))
	}
}

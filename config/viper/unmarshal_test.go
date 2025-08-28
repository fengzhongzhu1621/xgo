package viper

import (
	"fmt"
	"log"
	"testing"

	"github.com/spf13/viper"
)

type configStructs struct {
	MySQL struct {
		Username string `toml:"username"`
		Password string `toml:"password"`
	} `toml:"mysql"`
}

// TestUnmarshal 将配置文件解析到结构体
func TestUnmarshal(t *testing.T) {
	viper.SetConfigName("config_demo")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	//  读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal  error  config  file:  %w", err))
	}

	// 将配置文件解析到结构体
	configs := new(configStructs)
	err = viper.Unmarshal(configs)
	if err != nil {
		log.Fatal(fmt.Errorf("fatal  error  config  file:  %w", err))
	}

	fmt.Println("MySQL-Username:", configs.MySQL.Username)
}

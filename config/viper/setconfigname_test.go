package viper

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestSetConfigName(t *testing.T) {
	// 设置配置文件名称（不包括文件扩展名）
	viper.SetConfigName("config_demo")

	// 设置配置文件类型
	viper.SetConfigType("toml")

	// 添加配置文件搜索路径，在当前目录搜索
	viper.AddConfigPath(".")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal  error  config  file:  %w", err))
	}

	//  获取配置项
	message := viper.GetString("message")
	fmt.Println("Message:", message)
}

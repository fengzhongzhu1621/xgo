package viper

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetConfigFile(t *testing.T) {
	// 使用 SetConfigFile 指定配置文件的路径
	viper.SetConfigFile("config.yaml")

	// 读取并解析配置文件
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	// 获取配置文件中的值
	appName := viper.GetString("app.name")
	appVersion := viper.GetString("app.version")

	assert.Equal(t, "MyApp", appName)
	assert.Equal(t, "My1.0.0App", appVersion)
}

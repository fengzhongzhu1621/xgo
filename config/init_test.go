package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestSetCfgFile 测试读取开发配置文件的内容
func TestSetCfgFile(t *testing.T) {
	// 设置开发配置文件
	SetCfgFile()

	appName := viper.GetString("app.name")
	assert.Equal(t, "MyApp", appName)
}

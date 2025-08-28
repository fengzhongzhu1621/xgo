package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestSetCfgFile 测试读取开发配置文件的内容
func TestSetCfgFile(t *testing.T) {
	// 设置开发配置文件
	cfgFilePath := SetCfgFile()
	assert.Equal(t, true, strings.HasSuffix(cfgFilePath, "config/config.yaml"))

	appName := viper.GetString("app.name")
	assert.Equal(t, "MyApp", appName)
}

func TestGetAdminCfgFilePath(t *testing.T) {
	adminFilePath := GetAdminCfgFilePath()
	assert.Equal(t, true, strings.HasSuffix(adminFilePath, "config/admin.yaml"))
}

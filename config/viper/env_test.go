package viper

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAutomaticEnv(t *testing.T) {
	// 启用自动环境变量支持
	// 允许 Viper 自动从环境变量中读取配置。当调用此函数时，Viper 会查找与其已知键匹配的环境变量，并将这些环境变量的值作为配置。
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("app.name", "MyApp")
	viper.SetDefault("app.version", "1.0.0")

	// 从环境变量中读取配置
	appName := viper.GetString("app.name")
	appVersion := viper.GetString("app.version")

	assert.Equal(t, "MyApp", appName)
	assert.Equal(t, "1.0.0", appVersion)

	// export APP_NAME="MyCustomApp"
	// export APP_VERSION="2.0.0"
	// 运行上述代码将输出：

	// App Name: MyCustomApp
	// App Version: 2.0.0
}

func TestSetEnvPrefix(t *testing.T) {
	// 设置环境变量前缀
	// 用于设置环境变量的前缀。当设置了前缀后，Viper 会在读取环境变量时自动添加这个前缀，并且会将环境变量名中的大写字母转换为小写字母，同时用下划线替换点号
	viper.SetEnvPrefix("myapp")

	// 启用自动环境变量支持
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("name", "MyApp")
	viper.SetDefault("version", "1.0.0")

	// 从环境变量中读取配置
	appName := viper.GetString("name")
	appVersion := viper.GetString("version")

	fmt.Println("App Name:", appName)
	fmt.Println("App Version:", appVersion)

	// export MYAPP_NAME="MyCustomApp"
	// export MYAPP_VERSION="2.0.0"
	// 运行上述代码将输出：

	// App Name: MyCustomApp
	// App Version: 2.0.0
}

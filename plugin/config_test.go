package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert/yaml"
	"github.com/stretchr/testify/require"
)

type config struct {
	Plugins Config
}

func TestConfig_Setup(t *testing.T) {
	const configInfoNotRegister = `
plugins:
  mock_type:
    mock_not_register:
      address: localhost:8000
`
	cfg := config{}
	err := yaml.Unmarshal([]byte(configInfoNotRegister), &cfg)
	assert.Nil(t, err)

	// 因为插件没有注册，所以加载插件失败
	_, err = cfg.Plugins.SetupClosables()
	assert.NotNil(t, err)

	const configInfo = `
plugins:
  mock_type:
    mock_name:
      address: localhost:8000
`
	// 注册插件
	Register(pluginName, &mockPlugin{})
	cfg = config{}
	err = yaml.Unmarshal([]byte(configInfo), &cfg)
	assert.Nil(t, err)

	clo, err := cfg.Plugins.SetupClosables()
	assert.Nil(t, err)
	require.Nil(t, clo())
}

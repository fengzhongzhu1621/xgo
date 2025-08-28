package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v3"
	"trpc.group/trpc-go/trpc-go"
)

const configInfo = `
server:
  filter:
    - jwt
plugins:
  auth:
    jwt:
      secret: q7wt3n1t
      expired: 7200  # 2 hours
      issuer: tencent
      exclude_paths:
        - /v1/login
`

// TestPlugin_PluginType PluginType单元测试
func TestPlugin_PluginType(t *testing.T) {
	p := &pluginImp{}
	assert.Equal(t, pluginType, p.Type())
}

// TestPlugin_Setup Setup单元测试
func TestPlugin_Setup(t *testing.T) {
	cfg := trpc.Config{}
	err := yaml.Unmarshal([]byte(configInfo), &cfg)
	assert.Nil(t, err)

	assert.Equal(t, len(cfg.Server.Filter), 1)

	JWTCfg := cfg.Plugins[pluginType][pluginName]
	p := &pluginImp{}
	err = p.Setup(pluginName, &JWTCfg)
	assert.Nil(t, err)
}

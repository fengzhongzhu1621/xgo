package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v3"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const confNewFormat = `
server:
 filter:
  - validation
plugins:
  auth:
    validation:
      enable_error_log: true
      server_validate_err_code: 100101
      client_validate_err_code: 100102
`

// 测试插件类型
func TestValidationPlugin_Type(t *testing.T) {
	p := &ValidationPlugin{}
	assert.Equal(t, pluginType, p.Type())
}

// readConf
func readConf(conf string) plugin.Decoder {
	// 解析插件的yaml配置
	cfg := trpc.Config{}
	if err := yaml.Unmarshal([]byte(conf), &cfg); err != nil {
		return nil
	}
	// 返回 yaml.Node
	validCfg := cfg.Plugins[pluginType][pluginName]
	return &validCfg
}

func TestValidationPlugin_Setup(t *testing.T) {
	type args struct {
		name      string
		configDec plugin.Decoder
	}
	tests := []struct {
		name    string
		p       *ValidationPlugin
		args    args
		wantErr bool
	}{
		{
			name:    "test succ with confNewFormat",
			p:       &ValidationPlugin{},
			args:    args{name: pluginName, configDec: readConf(confNewFormat)},
			wantErr: false,
		},
		{
			name:    "test succ with configDec nil",
			p:       &ValidationPlugin{},
			args:    args{name: pluginName, configDec: nil},
			wantErr: false,
		},
		{
			name:    "test err configDec decode error",
			p:       &ValidationPlugin{},
			args:    args{name: pluginName, configDec: &plugin.YamlNodeDecoder{}},
			wantErr: true, // 插件配置解析错误
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ValidationPlugin{}
			err := p.Setup(tt.args.name, tt.args.configDec)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

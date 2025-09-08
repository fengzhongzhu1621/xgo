package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	pluginType        = "mock_type"
	pluginName        = "mock_name"
	pluginFailName    = "mock_fail_name"
	pluginTimeoutName = "mock_timeout_name"
	pluginDependName  = "mock_depend_name"
)

type mockPlugin struct{}

func (p *mockPlugin) Type() string {
	return pluginType
}

func (p *mockPlugin) Setup(name string, decoder IDecoder) error {
	return nil
}

func TestGet(t *testing.T) {
	Register(pluginName, &mockPlugin{})
	// test duplicate registration
	Register(pluginName, &mockPlugin{})
	p := Get(pluginType, pluginName)
	assert.NotNil(t, p)

	pNo := Get("notexist", pluginName)
	assert.Nil(t, pNo)
}

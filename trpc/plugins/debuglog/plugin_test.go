package debuglog

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/plugin"
)

const configInfo = `
plugins:
  tracing:
    debuglog:
      log_type: simple
      err_log_level: error
      nil_log_level: info
      server_log_type: prettyjson
      client_log_type: json
      enable_color: true
      exclude:
        - method: /trpc.app.server.service/method
        - retcode: 51
`

type testReq struct {
	A int
	B string
}

type testRsp struct {
	C int
	D string
}

// TestFilter_PluginType is the unit test for the Log interceptor plugin type.
func TestPlugin_PluginType(t *testing.T) {
	p := &Plugin{}
	assert.Equal(t, pluginType, p.Type())
}

// TestPlugin_Setup is the unit test for setting the attributes of the Log interceptor plugin.
func TestPlugin_Setup(t *testing.T) {
	cfg := trpc.Config{}
	err := yaml.Unmarshal([]byte(configInfo), &cfg)
	assert.Nil(t, err)

	conf, ok := cfg.Plugins[pluginType][pluginName]
	if !ok {
		assert.Nil(t, conf)
	}

	p := &Plugin{}
	err = p.Setup(pluginName, &plugin.YamlNodeDecoder{Node: &conf})
	assert.Nil(t, err)
}

func TestFilter_Filter(t *testing.T) {
	rsp := testRsp{
		C: 456,
		D: "456",
	}
	testHandleFunc1 := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &rsp, nil
	}
	testHandleFunc2 := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errs.New(errs.RetServerSystemErr, "system error")
	}

	testHandleFunc3 := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errs.New(errs.RetServerSystemErr, "system error")
	}

	testHandleFunc4 := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errs.New(errs.RetServerDecodeFail, "system decode error")
	}

	// 构造 message
	ctx := trpc.BackgroundContext()
	msg := trpc.Message(ctx)
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:6379")
	msg.WithRemoteAddr(addr)
	msg.WithServerRPCName("/trpc.app.server.service/methodA")

	req := testReq{
		A: 123,
		B: "123",
	}

	ex := newTestRule("/trpc.app.server.service/methodA", int(errs.RetClientCanceled))
	in := newTestRule("/trpc.app.server.service/methodA", int(errs.RetServerSystemErr))
	sf := ServerFilter(
		WithExclude(ex), WithInclude(in),
		WithInclude(newTestRule("/trpc.app.server.service/methodA", int(errs.RetOK))))

	_, err := sf(ctx, req, testHandleFunc1)
	assert.NotNil(t, sf)
	assert.Nil(t, err)

	ret, err := ServerFilter(WithExclude(ex))(ctx, req, testHandleFunc1)
	assert.NotNil(t, ret)
	assert.Nil(t, err)

	_, err = sf(ctx, req, testHandleFunc2)
	assert.NotNil(t, err)

	deadLineCtx, deadLineCancel := context.WithDeadline(ctx, time.Now().Add(time.Second*1))
	defer deadLineCancel()
	_, err = sf(deadLineCtx, req, testHandleFunc3)
	assert.NotNil(t, err)

	_, err = sf(deadLineCtx, req, testHandleFunc4)
	assert.NotNil(t, err)

	testClientHandleFunc1 := func(ctx context.Context, req, rsp interface{}) error {
		return nil
	}
	testClientHandleFunc2 := func(ctx context.Context, req, rsp interface{}) error {
		return errs.New(errs.RetClientConnectFail, "connect fail")
	}
	cf := ClientFilter()
	assert.NotNil(t, cf)
	assert.Nil(t, cf(ctx, req, rsp, testClientHandleFunc1))
	assert.NotNil(t, cf(ctx, req, rsp, testClientHandleFunc2))
}

func newTestRule(method string, retcode int) *RuleItem {
	return &RuleItem{Method: &method, Retcode: &retcode}
}

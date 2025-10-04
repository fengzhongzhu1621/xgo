package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/fengzhongzhu1621/xgo/crypto/unmarshaler"
	"github.com/stretchr/testify/require"
)

var _ unmarshaler.ICodec = (*dstMustBeMapCodec)(nil)

type dstMustBeMapCodec struct{}

func (c dstMustBeMapCodec) Name() string {
	return "map"
}

func (c dstMustBeMapCodec) Unmarshal(bts []byte, dst interface{}) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr ||
		rv.Elem().Kind() != reflect.Interface ||
		rv.Elem().Elem().Kind() != reflect.Map ||
		rv.Elem().Elem().Type().Key().Kind() != reflect.String ||
		rv.Elem().Elem().Type().Elem().Kind() != reflect.Interface {
		return errors.New("the dst of codec.Unmarshal must be a map")
	}
	return nil
}

func NewEnvProvider(name string, data []byte) *EnvProvider {
	return &EnvProvider{
		name: name,
		data: data,
	}
}

var _ IDataProvider = (*EnvProvider)(nil)

type EnvProvider struct {
	name string
	data []byte
}

func (ep *EnvProvider) Name() string {
	return ep.name
}

func (ep *EnvProvider) Read(string) ([]byte, error) {
	return ep.data, nil
}

func (ep *EnvProvider) Watch(cb ProviderCallback) {
	cb("", ep.data)
}

var _ IDataProvider = (*manualTriggerWatchProvider)(nil)

type manualTriggerWatchProvider struct {
	values    sync.Map
	callbacks []ProviderCallback
}

func (m *manualTriggerWatchProvider) Name() string {
	return "manual_trigger_watch_provider"
}

func (m *manualTriggerWatchProvider) Read(s string) ([]byte, error) {
	if v, ok := m.values.Load(s); ok {
		return v.([]byte), nil
	}
	return nil, fmt.Errorf("not found config")
}

func (m *manualTriggerWatchProvider) Watch(callback ProviderCallback) {
	m.callbacks = append(m.callbacks, callback)
}

func (m *manualTriggerWatchProvider) Set(key string, v []byte) {
	m.values.Store(key, v)
	for _, callback := range m.callbacks {
		callback(key, v)
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////
func TestCodecUnmarshalDstMustBeMap(t *testing.T) {
	filePath := t.TempDir() + "/conf.map"
	require.Nil(t, os.WriteFile(filePath, []byte{}, 0600))

	unmarshaler.RegisterCodec(dstMustBeMapCodec{})
	_, err := DefaultConfigLoader.Load(filePath, WithCodec(dstMustBeMapCodec{}.Name()))
	require.Nil(t, err)
}

func TestEnvExpanded(t *testing.T) {
	RegisterProvider(NewEnvProvider(t.Name(), []byte(`
password: ${pwd}
`)))

	t.Setenv("pwd", t.Name())

	cfg, err := DefaultConfigLoader.Load(
		t.Name(),
		WithProvider(t.Name()),
		WithExpandEnv())
	require.Nil(t, err)

	require.Equal(t, t.Name(), cfg.GetString("password", ""))
	require.Contains(t, string(cfg.Bytes()), fmt.Sprintf("password: %s", t.Name()))
}

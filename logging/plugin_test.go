package logging

import (
	"testing"

	"github.com/fengzhongzhu1621/xgo/datetime"
	"github.com/fengzhongzhu1621/xgo/logging/config"
	"github.com/gookit/goutil/testutil/assert"
)

type fakeDecoder struct{}

func (c *fakeDecoder) Decode(conf interface{}) error {
	return nil
}

func TestWriterFactory(t *testing.T) {
	f1 := &ConsoleWriterFactory{}
	assert.Equal(t, "log", f1.Type())

	// empty decoder
	err := f1.Setup("default", nil)
	assert.NotNil(t, err)

	f2 := &FileWriterFactory{}
	assert.Equal(t, "log", f2.Type())
	// empty decoder
	err = f2.Setup("default", nil)
	assert.NotNil(t, err)

	f3 := &ConsoleWriterFactory{}
	assert.Equal(t, "log", f3.Type())
	err = f3.Setup("default", &fakeDecoder{})
	assert.NotNil(t, err)

	f4 := &FileWriterFactory{}
	assert.Equal(t, "log", f4.Type())
	err = f4.Setup("default", &fakeDecoder{})
	assert.NotNil(t, err)
}

func TestFileWriterFactory_Setup(t *testing.T) {
	var fileCfg = []config.LogOutputConfig{
		{
			Writer: "file", // writer 插件名称
			WriteConfig: config.LogWriteConfig{
				Filename:   "trpc_time.log",
				MaxAge:     7,
				MaxBackups: 10,
				MaxSize:    100,
				TimeUnit:   datetime.Day,
				LogPath:    "log",
			},
		},
	}
	logger := NewZapLog(fileCfg)
	assert.NotNil(t, logger)
}

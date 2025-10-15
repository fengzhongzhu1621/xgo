package logging

import (
	"errors"
	"testing"

	"github.com/fengzhongzhu1621/xgo/logging/config"
	"github.com/fengzhongzhu1621/xgo/logging/output"
	"github.com/fengzhongzhu1621/xgo/plugin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const observewriter = "observewriter"

type observeWriter struct {
	core zapcore.Core
}

func (f *observeWriter) Type() string { return "log" }

func (f *observeWriter) Setup(name string, dec plugin.IDecoder) error {
	if dec == nil {
		return errors.New("empty decoder")
	}
	decoder, ok := dec.(*Decoder)
	if !ok {
		return errors.New("invalid decoder")
	}
	decoder.Core = f.core
	decoder.ZapLevel = zap.NewAtomicLevel()
	return nil
}

func TestWithFields(t *testing.T) {
	// register Writer.
	// use zap observer to support test.
	core, ob := observer.New(zap.InfoLevel)
	output.RegisterWriter(observewriter, &observeWriter{core: core})

	// config is configuration.
	cfg := []config.LogOutputConfig{
		{
			Writer: observewriter,
		},
	}

	// create a zap logger.
	zl := NewZapLog(cfg)
	assert.NotNil(t, zl)

	// test With.
	field := Field{Key: "abc", Value: int32(123)}
	logger := zl.With(field)
	assert.NotNil(t, logger)
	SetLogger(logger)
	Warn("with fields warning")
	assert.Equal(t, 1, ob.Len())
	entry := ob.All()[0]
	assert.Equal(t, zap.WarnLevel, entry.Level)
	assert.Equal(t, "with fields warning", entry.Message)
	assert.Equal(t, []zapcore.Field{{Key: "abc", Type: zapcore.Int32Type, Integer: 123}}, entry.Context)
}

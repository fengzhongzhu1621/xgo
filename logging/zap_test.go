package logging

import (
	"bytes"
	"testing"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/logging/level"
	"github.com/fengzhongzhu1621/xgo/logging/zaplogger"
	"github.com/gookit/goutil/testutil/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewZapLog(t *testing.T) {
	logger := NewZapLog(defaultConfig)
	assert.NotNil(t, logger)

	logger.SetLevel("0", level.LevelInfo)
	lvl := logger.GetLevel("0")
	assert.Equal(t, lvl, level.LevelInfo)

	l := logger.With(Field{Key: "test", Value: "a"})
	l.SetLevel("output", level.LevelDebug)
	assert.Equal(t, level.LevelDebug, l.GetLevel("output"))
}

func TestZapLogWithLevel(t *testing.T) {
	logger := NewZapLog(defaultConfig)
	assert.NotNil(t, logger)

	l := logger.With(Field{Key: "test", Value: "a"})
	l.SetLevel("0", level.LevelFatal)
	assert.Equal(t, level.LevelFatal, l.GetLevel("0"))

	l = l.With(Field{Key: "key1", Value: "val1"})
	l.SetLevel("0", level.LevelError)
	assert.Equal(t, level.LevelError, l.GetLevel("0"))
}

// NewZapBufLogger return a buffer logger
func NewZapBufLogger(buf *bytes.Buffer, skip int) ILogger {
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewCore(encoder, zapcore.AddSync(buf), zapcore.DebugLevel)
	return &zapLog{
		levels: []zap.AtomicLevel{},
		logger: zap.New(
			core,
			zap.AddCallerSkip(skip),
			zap.AddCaller(),
		),
	}
}

// NewZapFatalLogger return a fatal hook logger
func NewZapFatalLogger(h zapcore.CheckWriteHook) ILogger {
	core, _ := zaplogger.NewConsoleCore(&config.LogOutputConfig{
		Writer:    "console",
		Level:     "debug",
		Formatter: "console",
	})
	return &zapLog{
		levels: []zap.AtomicLevel{},
		logger: zap.New(
			core,
			zap.AddCallerSkip(1),
			zap.AddCaller(),
			zap.WithFatalHook(h),
		)}
}

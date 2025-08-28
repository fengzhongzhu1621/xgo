package logging

import (
	"os"
	"testing"

	"github.com/fengzhongzhu1621/xgo/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZapCoreEncoderConfig(t *testing.T) {
	// 创建 EncoderConfig 实例，并设置自定义的 NameEncoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 创建编码器
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 创建一个 WriteSyncer，例如输出到标准输出
	writeSyncer := zapcore.AddSync(os.Stdout)

	// 创建一个 Core 实例
	core := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)

	// 使用自定义的 Core 创建一个 Logger 实例
	logger := zap.New(core)

	// 使用 logger 输出日志
	logger.Info("这是一条信息级别的日志")
	// {"level":"info","time":"2024-10-03T10:27:32.757+0800","msg":"这是一条信息级别的日志"}
}

func TestNewZapJSONLogger(t *testing.T) {
	logConfig := config.LogConfig{
		Writer:   "os",
		Level:    "info",
		Settings: map[string]string{"name": "stdout"},
	}
	logger := NewZapJSONLogger(&logConfig, false)
	logger.Info("这是一条信息级别的日志")
	logger.Sync()
	// {"level":"info","time":"2024-10-03T11:34:14.395+0800","msg":"这是一条信息级别的日志"}
}

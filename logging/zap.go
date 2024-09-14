package logging

import (
	"fmt"
	"time"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/logging/zapx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapJSONLogger 设置  json 格式的日志输出记录器
func NewZapJSONLogger(cfg *config.LogConfig) *zap.Logger {
	writer, err := GetWriter(cfg.Writer, cfg.Settings)
	if err != nil {
		panic(err)
	}
	w := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(writer),
		Size:          256 * 1024, // 256 kB
		FlushInterval: 30 * time.Second,
	}

	// 从配置文件获得日志级别
	l, err := zapx.ParseZapLogLevel(cfg.Level)
	if err != nil {
		fmt.Println("api logger settings level invalid, will use level: info")
		l = zap.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		l,
	)

	return zap.New(core)
}

package logging

import (
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
	// zapcore.Core 是 zap 库的核心结构，它负责管理日志的编码、输出和级别过滤。
	// AddSync 函数将一个 io.Writer 接口添加到 zapcore.Core 的同步器列表中，允许将自定义的输出目标（如文件、网络连接等）添加到日志系统的同步器列表中。
	//
	// zapcore.BufferedWriteSyncer 实现了 zapcore.WriteSyncer 接口
	// 提供了一种缓冲写入的方式，可以在内存中缓存一定数量的日志消息，然后在达到一定阈值时一次性将它们写入底层的 WriteSyncer
	// Size: 指定了缓冲区的大小（以字节为单位）。当缓冲区满时，日志消息将被一次性写入文件，可以根据需要调整缓冲区的大小，以平衡内存使用和写入性能
	w := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(writer),
		Size:          256 * 1024,
		FlushInterval: 30 * time.Second,
	}

	// 将字符串转换为 Zap 的日志级别
	l, err := zapx.ParseZapLogLevel(cfg.Level)
	if err != nil {
		l = zap.InfoLevel
	}

	// 用于配置日志编码器的行为。编码器负责将日志记录转换为字节流，以便将其写入输出目标（如文件、控制台等）
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",                         // 时间戳的键名，默认为 "time"
		LevelKey:       "level",                        // 日志级别的键名，默认为 "level"
		NameKey:        "logger",                       // 日志记录器名称的键名，默认为 "logger"
		CallerKey:      "linenum",                      // 调用者信息的键名，默认为 "caller"
		MessageKey:     "msg",                          // 日志消息的键名，默认为 "msg"
		StacktraceKey:  "stacktrace",                   // 堆栈跟踪信息的键名，默认为 "stacktrace"
		LineEnding:     zapcore.DefaultLineEnding,      // 行结束符，默认为 zapcore.DefaultLineEnding
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 日志级别编码器 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // 时间戳编码器 ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 持续时间编码器
		EncodeCaller:   zapcore.FullCallerEncoder,      // 调用者信息编码器
		EncodeName:     zapcore.FullNameEncoder,        // 想要自定义日志记录器名称的编码方式，你可以实现自己的 NameEncoder 函数
	}

	// 创建一个 Core 实例
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		l,
	)

	// 使用自定义的 Core 创建一个 Logger 实例
	return zap.New(core)
}

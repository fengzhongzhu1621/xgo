package logging

import (
	"log"
	"testing"

	"gopkg.in/natefinch/lumberjack.v2"
)

func TestLogger(t *testing.T) {
	// lumberjack.Logger 是一个 Go 语言的日志库，它提供了简单的日志滚动和压缩功能。lumberjack.Logger 结构体用于配置和使用这个日志库。
	// 配置 lumberjack.Logger
	logger := &lumberjack.Logger{
		Filename:   "logs/myapp.log", // 日志文件的位置
		MaxSize:    10,               // 每个日志文件的最大尺寸（以MB为单位）
		MaxBackups: 3,                // 保留的最大日志文件数量
		MaxAge:     28,               // 保留的最大日志文件天数
		Compress:   true,             // 是否压缩旧的日志文件
		LocalTime:  true,             // 是否在日志文件名中使用本地时间而不是 UTC 时间。默认为 false
	}

	// 设置 log 包的输出到 lumberjack.Logger
	log.SetOutput(logger)

	// 使用 log 包记录日志
	log.Println("This is a test log message.")

	// 2024/09/30 18:24:58 This is a test log message.
	// 2024/10/03 10:21:17 This is a test log message.
}

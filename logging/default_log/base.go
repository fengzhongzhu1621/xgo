package default_go

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"logur.dev/logur"

	"github.com/fengzhongzhu1621/xgo/logging"
)

// SetLogger sets a logger to the logging package.
func SetLogger(name string, logger *log.Logger, level logur.Level) {
	logging.SetLogger(name, New(logger, level))
}

// EnsureDefaultLogger will replace the default logger by log default logger.
func EnsureDefaultLogger(level logur.Level) {
	SetLogger(logging.DefaultLoggerName, log.New(os.Stderr, log.Prefix(), log.Flags()), level)
}

// Logger is a Logur adapter for TEMPLATE.
type Logger struct {
	level  logur.Level
	logger *log.Logger
}

// New returns a new Logur logger.
// If logger is nil, a default instance is created.
func New(logger *log.Logger, level logur.Level) *Logger {
	return &Logger{
		level:  level,
		logger: logger,
	}
}

func (l *Logger) log(level logur.Level, prefix string, msg string, fields ...map[string]interface{}) {
	// 检查传入的日志级别是否低于 Logger 实例设置的最低日志级别。如果是，则直接返回，不进行任何日志记录。
	if level < l.level {
		return
	}

	// 计算了日志消息中需要包含的部分总数。
	// 默认情况下，至少有 prefix 和 msg 两部分。如果有额外的字段（fields 不为空），
	// 则将第一个映射中的键值对数量加到 fieldCount 上。
	fieldCount := 2
	if len(fields) != 0 {
		fieldCount += len(fields[0])
	}

	// 创建一个字符串切片 parts，其容量为 fieldCount，用于存储日志消息的各个部分。
	// 然后将 prefix 和 msg 添加到 parts 中
	parts := make([]string, 0, fieldCount)
	parts = append(parts, prefix, msg)

	// 如果提供了额外的字段（fields 不为空），则遍历第一个映射中的所有键值对，
	// 并使用 fmt.Sprintf 将它们格式化为 key=value 的形式，然后将这些字符串添加到 parts 中
	if len(fields) != 0 {
		for key, value := range fields[0] {
			parts = append(parts, fmt.Sprintf("%v=%v", key, value))
		}
	}

	// 使用 strings.Join 将 parts 中的所有部分用空格连接起来，形成一个完整的日志消息字符串，
	// 然后调用底层 log.Logger 的 Print 方法来输出日志。
	l.logger.Print(strings.Join(parts, " "))
}

// Trace implements the Logur Logger interface.
func (l *Logger) Trace(msg string, fields ...map[string]interface{}) {
	l.log(logur.Trace, "TRACE", msg, fields...)
}

// Debug implements the Logur Logger interface.
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	l.log(logur.Debug, "DEBUG", msg, fields...)
}

// Info implements the Logur Logger interface.
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	l.log(logur.Info, "INFO", msg, fields...)
}

// Warn implements the Logur Logger interface.
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	l.log(logur.Warn, "WARN", msg, fields...)
}

// Error implements the Logur Logger interface.
func (l *Logger) Error(msg string, fields ...map[string]interface{}) {
	l.log(logur.Error, "ERROR", msg, fields...)
}

// TraceContext 在这个实现中，context.Context 参数被忽略了，
// 因为 log 包本身不支持上下文。如果需要上下文感知的日志记录，你可能需要扩展这个适配器或使用其他支持上下文的日志库。
func (l *Logger) TraceContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Trace(msg, fields...)
}

func (l *Logger) DebugContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Debug(msg, fields...)
}

func (l *Logger) InfoContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Info(msg, fields...)
}

func (l *Logger) WarnContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Warn(msg, fields...)
}

func (l *Logger) ErrorContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Error(msg, fields...)
}

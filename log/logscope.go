package log

import (
	"sync"
	"time"
)

const (
	LogDebug = iota
	LogTrace
	LogAccess
	LogInfo
	LogError
	LogFatal

	LogNone    = 126
	LogUnknown = 127
)


// 定义日志输出接口，需要使用方自定义输出函数的实现
type LogFunc func(logScope *LogScope, level LogLevel, levelName string, category string, msg string, args ...interface{})

type LogScope struct {
	Parent     *LogScope   // Parent scope
	MinLevel   LogLevel    // Minimum log level to write out.
	Mutex      *sync.Mutex // Should be shared across all LogScopes that write to the same destination.
	Associated []AssocPair // Additional data associated with scope
	LogFunc    LogFunc
}

type AssocPair struct {
	Key   string
	Value string
}

/**
 * 添加日志元素
 */
func (l *LogScope) Associate(key string, value string) {
	l.Associated = append(l.Associated, AssocPair{key, value})
}

func (l *LogScope) Debug(category string, msg string, args ...interface{}) {
	l.LogFunc(l, LogDebug, "DEBUG", category, msg, args...)
}

func (l *LogScope) Trace(category string, msg string, args ...interface{}) {
	l.LogFunc(l, LogTrace, "TRACE", category, msg, args...)
}

func (l *LogScope) Access(category string, msg string, args ...interface{}) {
	l.LogFunc(l, LogAccess, "ACCESS", category, msg, args...)
}

func (l *LogScope) Info(category string, msg string, args ...interface{}) {
	l.LogFunc(l, LogInfo, "INFO", category, msg, args...)
}

func (l *LogScope) Error(category string, msg string, args ...interface{}) {
	l.LogFunc(l, LogError, "ERROR", category, msg, args...)
}

func (l *LogScope) Fatal(category string, msg string, args ...interface{}) {
	l.LogFunc(l, LogFatal, "FATAL", category, msg, args...)
}

func (l *LogScope) NewLevel(logFunc LogFunc) *LogScope {
	return &LogScope{
		Parent:     l,
		MinLevel:   l.MinLevel,
		Mutex:      l.Mutex,
		Associated: make([]AssocPair, 0),
		LogFunc:    logFunc}
}

func RootLogScope(minLevel LogLevel, logFunc LogFunc) *LogScope {
	return &LogScope{
		Parent:     nil,
		MinLevel:   minLevel,
		Mutex:      &sync.Mutex{},
		Associated: make([]AssocPair, 0),
		LogFunc:    logFunc}
}

func Timestamp() string {
	return time.Now().Format(time.RFC1123Z)
}

func LevelFromString(s string) LogLevel {
	switch s {
	case "debug":
		return LogDebug
	case "trace":
		return LogTrace
	case "access":
		return LogAccess
	case "info":
		return LogInfo
	case "error":
		return LogError
	case "fatal":
		return LogFatal
	case "none":
		return LogNone
	default:
		return LogUnknown
	}
}

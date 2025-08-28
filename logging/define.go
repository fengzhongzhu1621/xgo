package logging

import (
	"sync"

	"logur.dev/logur"
)

// Looger is a structured logger interface.
// 定义了一个结构化日志记录器接口 Logger，它是一个 logur.LoggerFacade 类型的别名
type LogurLogger = logur.LoggerFacade

// DefaultLoggerName is the name of the default logger.
// 定义了一个默认的日志记录器名称 DefaultLoggerName
const DefaultLoggerName = ""

// 日志记录器的插件仓库
var loggerAlias, loggers *sync.Map

// tryGetLogger 尝试从 loggers 映射中加载对应的日志记录器实例。
func tryGetLogger(name string) (LogurLogger, bool) {
	// 根据给定的名称从 loggers 这个全局映射中获取一个日志记录器实例
	v, ok := loggers.Load(name)
	if !ok {
		return nil, false
	}
	// 转换为 LogurLogger 类型
	logger, ok := v.(LogurLogger)
	if !ok {
		return nil, false
	}

	return logger, true
}

// GetLogger returns a named logger. If the logger does not exist, it will return the default logger.
func GetLogger(name string) LogurLogger {
	// 尝试从 loggers 映射中加载对应的日志记录器实例
	logger, ok := tryGetLogger(name)
	if ok {
		return logger
	}

	// 尝试查找并获得别名
	// search for alias when name is not found
	realName, ok := loggerAlias.Load(name)
	if !ok {
		// 如果 realName 不存在，则直接返回默认的日志记录器
		logger, _ = tryGetLogger(DefaultLoggerName)
		return logger
	}

	// 尝试根据这个别名获取日志记录器
	logger, ok = tryGetLogger(realName.(string))
	if ok {
		return logger
	}

	// 如果 realName 不存在，则直接返回默认的日志记录器
	logger, _ = tryGetLogger(DefaultLoggerName)
	return logger
}

// SetLogger sets a named logger.
func SetLogger(name string, logger LogurLogger) {
	loggers.Store(name, logger)
}

// SetAlias is used to define a logger alias.
func SetAlias(name string, aliases ...string) {
	for _, alias := range aliases {
		loggerAlias.Store(alias, name)
	}
}

func init() {
	loggerAlias = &sync.Map{}
	loggers = &sync.Map{}

	// make sure the default logger is created
	SetLogger(DefaultLoggerName, logur.NoopLogger{})
}

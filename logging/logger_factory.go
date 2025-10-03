package logging

import (
	"sync"

	"github.com/fengzhongzhu1621/xgo/logging/level"
)

const (
	defaultLoggerName = "default"
)

var (
	// DefaultLogger the default Logger. The initial output is console. When frame start, it is
	// over write by configuration.
	DefaultLogger ILogger
	// DefaultLogFactory is the default log loader. Users may replace it with their own
	// implementation.
	DefaultLogFactory = &Factory{}

	mu        sync.RWMutex
	loggerMap = make(map[string]ILogger)
)

// GetDefaultLogger gets the default Logger.
// To configure it, set key in configuration file to default.
// The console output is the default value.
func GetDefaultLogger() ILogger {
	mu.RLock()
	l := DefaultLogger
	mu.RUnlock()
	return l
}

// SetLogger sets the default Logger.
func SetLogger(logger ILogger) {
	mu.Lock()
	DefaultLogger = logger
	mu.Unlock()
}

// Get returns the Logger implementation by log name.
// log.Debug use DefaultLogger to print logs. You may also use log.Get("name").Debug.
func Get(name string) ILogger {
	mu.RLock()
	l := loggerMap[name]
	mu.RUnlock()
	return l
}

// Sync syncs all registered loggers.
func Sync() {
	mu.RLock()
	defer mu.RUnlock()
	for _, logger := range loggerMap {
		_ = logger.Sync()
	}
}

// SetLevel sets log level for different output which may be "0", "1" or "2".
func SetLevel(output string, level level.Level) {
	GetDefaultLogger().SetLevel(output, level)
}

// GetLevel gets log level for different output.
func GetLevel(output string) level.Level {
	return GetDefaultLogger().GetLevel(output)
}

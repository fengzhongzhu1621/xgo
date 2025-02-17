package logrus

import (
	"github.com/sirupsen/logrus"
	adapter "logur.dev/adapter/logrus"

	"github.com/fengzhongzhu1621/xgo/logging"
)

// SetLogger sets a logger to the logging package.
func SetLogger(name string, logger *logrus.Logger) {
	logging.SetLogger(name, adapter.New(logger))
}

// SetLoggerFromEntry sets a logger entry to the logging package.
func SetLoggerFromEntry(name string, entry *logrus.Entry) {
	logging.SetLogger(name, adapter.NewFromEntry(entry))
}

// EnsureDefaultLogger will replace the default logger by logrus.StandardLogger.
// 启用默认日志实现（输出到标准错误）
func EnsureDefaultLogger() {
	SetLogger(logging.DefaultLoggerName, logrus.StandardLogger())
}

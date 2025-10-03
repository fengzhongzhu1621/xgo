package zaplogger

import (
	"github.com/fengzhongzhu1621/xgo/logging/config"
	"github.com/fengzhongzhu1621/xgo/logging/output"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

var (
	dbLogger  *DBLogger
	appLogger *zap.Logger
)

func InitSystemLogger(cfg *config.LogConfig) {
	writer, err := output.GetWriter(cfg.Writer, cfg.Settings)
	if err != nil {
		panic(err)
	}
	// 日志输出到stdout
	logrus.SetOutput(writer)
	// 不输出颜色
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})

	// 设置日志级别
	l, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		l = logrus.InfoLevel
	}
	logrus.SetLevel(l)
}

// GetSystemLogger 获得预定义的日志记录器实例
func GetSystemLogger() *logrus.Logger {
	// 一个全局标准日志记录器，可以直接使用它来记录日志，而无需自己创建一个新的日志记录器实例。
	return logrus.StandardLogger()
}

// GetAppLogger 获得 web 应用的日志记录器
func GetAppLogger() *zap.Logger {
	return appLogger
}

func GetDbLogger() *DBLogger {
	return dbLogger
}

func SetDbLogger(logger *DBLogger) {
	dbLogger = logger
}

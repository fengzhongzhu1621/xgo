package logging

import (
	"sync"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

var loggerInitOnce sync.Once

var appLogger *zap.Logger

// GetSystemLogger 获得预定义的日志记录器实例
func GetSystemLogger() *logrus.Logger {
	// 一个全局标准日志记录器，可以直接使用它来记录日志，而无需自己创建一个新的日志记录器实例。
	return logrus.StandardLogger()
}

// GetAppLogger 获得 web 应用的日志记录器
func GetAppLogger() *zap.Logger {
	if appLogger == nil {
		appLogger, _ = zap.NewProduction()
		// 刷新所有缓冲的日志条目
		defer appLogger.Sync()
	}
	return appLogger
}

func initSystemLogger(cfg *config.LogConfig) {
	writer, err := GetWriter(cfg.Writer, cfg.Settings)
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

// InitLogger 初始化日志记录器，只能执行一次
func InitLogger() {
	globalConfig := config.GetGlobalConfig()
	logger := globalConfig.Logger

	// 设置系统日志记录器
	initSystemLogger(&logger.System)

	// TOOD InitSLogger

	loggerInitOnce.Do(func() {
		// 设置 api 服务器日志记录器
		// apiLogger = NewZapJSONLogger(&logger.API)
		// 设置 web 服务器日志记录器
		appLogger = NewZapJSONLogger(&logger.Web)
	})

}

func init() {
	InitLogger()
}

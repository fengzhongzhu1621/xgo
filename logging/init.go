package logging

import (
	"fmt"
	"sync"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

var loggerInitOnce sync.Once

var webLogger *zap.Logger

// GetWebLogger 获得 web 应用的日志记录器
func GetWebLogger() *zap.Logger {
	if webLogger == nil {
		// 创建一个新的zap.Logger实例
		webLogger, _ = zap.NewProduction()
		defer webLogger.Sync()
	}
	return webLogger
}

// InitLogger 初始化日志记录器，只能执行一次
func InitLogger(logger *config.Logger) {
	// 设置系统日志记录器
	initSystemLogger(&logger.System)

	loggerInitOnce.Do(func() {
		// 设置 web 服务器日志记录器
		webLogger = NewZapJSONLogger(&logger.Web)
	})
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
		fmt.Println("system logger settings level invalid, will use level: info")
		l = logrus.InfoLevel
	}
	logrus.SetLevel(l)
}

func init() {
	var globalConfig = config.GetGlobalConfig()
	var logger = globalConfig.Logger
	InitLogger(&logger)
}

package logrus

import (
	"testing"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func TestLumberjack(t *testing.T) {
	log := logrus.New()

	// 使用 lumberjack 实现日志文件轮转
	log.SetOutput(&lumberjack.Logger{
		Filename:   "app.log", // 日志文件路径
		MaxSize:    100,       // 单个日志文件的最大大小（MB）
		MaxBackups: 3,         // 最多保留的日志文件数量
		MaxAge:     28,        // 日志文件的最大保存天数
		Compress:   true,      // 是否压缩备份文件
	})

	log.Info("This is an info log")
}

package logrus

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSetFormatter(t *testing.T) {
	// 创建一个默认的 logger
	var log = logrus.New()

	log.SetLevel(logrus.WarnLevel) // 设置日志级别为 Warn
	log.SetFormatter(&logrus.JSONFormatter{})

	log.Info("This is an info log")
	log.Warn("This is a warning log")
	log.Error("This is an error log")
}

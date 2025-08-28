package logrus

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNew(t *testing.T) {
	// 创建一个默认的 logger
	log := logrus.New()

	log.SetLevel(logrus.WarnLevel) // 设置日志级别为 Warn

	log.Info("This is an info log")
	log.Warn("This is a warning log")
	log.Error("This is an error log")
}

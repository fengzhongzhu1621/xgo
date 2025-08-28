package logrus

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSetOutputToStdout(t *testing.T) {
	log := logrus.New()

	// 创建文件输出
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	// 同时将日志输出到文件和控制台
	log.SetOutput(os.Stdout)
	log.SetOutput(file)

	log.Info("This is an info log")
	log.Warn("This is a warning log")
}

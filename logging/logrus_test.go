package logging

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

type AppHook struct {
	AppName string
}

func (h *AppHook) Levels() []Level {
	return AllLevels
}

func (h *AppHook) Fire(entry *Entry) error {
	entry.Data["app"] = h.AppName
	return nil
}

func TestLogrusInfo(t *testing.T) {
	LogrusSetLevel(TraceLevel)
	LogrusSetReportCaller(true)
	LogrusInfo("info msg")
	LogrusWithFields(Fields{
		"name": "dj",
		"age":  18,
	}).Info("info msg")

	writer1 := &bytes.Buffer{}
	writer2 := os.Stdout
	writer3, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	LogrusSetOutput(io.MultiWriter(writer1, writer2, writer3))
	LogrusInfo("info msg")

	LogrusSetFormatter(&JSONFormatter{})
	LogrusInfo("info msg")

	h := &AppHook{AppName: "awesome-web"}
	LogrusAddHook(h)
	LogrusInfo("info msg")
}

func TestWithFields(t *testing.T) {
	// 创建一个新的 logrus 实例
	logger := logrus.New()

	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})

	// 记录不同级别的日志
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// 使用 WithFields 记录结构化日志
	logger.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	// 使用 WithField 记录带有单个字段的结构化日志
	logger.WithField("omg", true).Warn("The ice breaks!")
}

func TestParseLevel(t *testing.T) {
	// 解析字符串为日志级别
	level, err := logrus.ParseLevel("debug")
	if err != nil {
		fmt.Println("Error parsing level:", err)
		return
	}

	// 设置日志级别
	logrus.SetLevel(level)

	// 记录不同级别的日志
	logrus.Debug("This is a debug message")
	logrus.Info("This is an info message")
	logrus.Warn("This is a warning message")
	logrus.Error("This is an error message")
}

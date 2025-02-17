package logrus

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

type MyHook struct{}

func (hook *MyHook) Levels() []logrus.Level {
	return logrus.AllLevels // 钩子作用于所有日志级别
}

func (hook *MyHook) Fire(entry *logrus.Entry) error {
	// 在每次日志输出时执行的逻辑
	fmt.Println("My custom hook triggered!")
	fmt.Println("hook output: " + entry.Message)
	return nil
}

func TestAddHook(t *testing.T) {
	log := logrus.New()
	log.AddHook(&MyHook{}) // 注册自定义 hook

	log.Info("This is an info log")   // 会触发钩子
	log.Error("This is an error log") // 会触发钩子
}

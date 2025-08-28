package logging

import (
	"testing"
)

func TestGetSystemLogger(t *testing.T) {
	logger := GetSystemLogger()
	logger.Info("这是一条信息级别的日志")
	// time="2024-10-03T11:16:10+08:00" level=info msg="这是一条信息级别的日志"
}

func TestGetAppLogger(t *testing.T) {
	webLogger := GetAppLogger()
	webLogger.Info("这是一条信息级别的日志")

	webLogger.Sync()
	// {"level":"info","time":"2024-10-03T11:35:03.025+0800","msg":"这是一条信息级别的日志"}
}

package logging

import (
	"log/slog"

	"github.com/go-resty/resty/v2"
)

// RestyLogger 用于实现 resty.Logger
type RestyLogger struct{}

func NewRestyLogger() *RestyLogger {
	return &RestyLogger{}
}

func (l *RestyLogger) Errorf(format string, v ...interface{}) {
	SLog(slog.LevelError, format, v...)
}

func (l *RestyLogger) Warnf(format string, v ...interface{}) {
	SLog(slog.LevelWarn, format, v...)
}

func (l *RestyLogger) Debugf(format string, v ...interface{}) {
	SLog(slog.LevelDebug, format, v...)
}

var _ resty.Logger = (*RestyLogger)(nil)

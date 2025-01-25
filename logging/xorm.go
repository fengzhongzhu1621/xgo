package logging

import (
	"fmt"

	"go.uber.org/zap"
	xormLog "xorm.io/xorm/log"
)

// DBLogger 实现了 zap 需要的接口
type DBLogger struct {
	logger *zap.Logger
}

var _ xormLog.Logger = (*DBLogger)(nil)

func (il *DBLogger) Sync() {
	_ = il.logger.Sync()
}
func (il *DBLogger) Debug(v ...interface{}) {
	il.logger.Debug(fmt.Sprint(v...))
}
func (il *DBLogger) Debugf(format string, v ...interface{}) {
	il.logger.Debug(fmt.Sprintf(format, v...))
}

func (il *DBLogger) Info(v ...interface{}) {
	il.logger.Info(fmt.Sprint(v...))
}

func (il *DBLogger) Infof(format string, v ...interface{}) {
	il.logger.Info(fmt.Sprintf(format, v...))
}

func (il *DBLogger) Warn(v ...interface{}) {
	il.logger.Warn(fmt.Sprint(v...))
}

func (il *DBLogger) Warnf(format string, v ...interface{}) {
	il.logger.Warn(fmt.Sprintf(format, v...))
}

func (il *DBLogger) Error(v ...interface{}) {
	il.logger.Error(fmt.Sprint(v...))
}

func (il *DBLogger) Errorf(format string, v ...interface{}) {
	il.logger.Error(fmt.Sprintf(format, v...))
}

// Level 返回对应的 xorm 日志级别
func (il *DBLogger) Level() xormLog.LogLevel {
	switch il.logger.Core().Enabled(zap.DebugLevel) {
	case true:
		return xormLog.LOG_DEBUG
	case false:
		return xormLog.LOG_ERR
	default:
		return xormLog.LOG_INFO
	}
}

// SetLevel implement ILogger 设置日志级别
func (il *DBLogger) SetLevel(_ xormLog.LogLevel) {

}

// ShowSQL implement ILogger
func (il *DBLogger) ShowSQL(_ ...bool) {

}

// IsShowSQL implement ILogger 显示 SQL 语句
func (il *DBLogger) IsShowSQL() bool {
	return true
}

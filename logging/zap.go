package logging

import (
	"fmt"
	"strconv"

	"github.com/fengzhongzhu1621/xgo/logging/config"
	"github.com/fengzhongzhu1621/xgo/logging/level"
	"github.com/fengzhongzhu1621/xgo/logging/output"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultConfig = []config.LogOutputConfig{
	{
		Writer:    "console",
		Level:     "debug",
		Formatter: "console",
	},
}

// NewZapLog creates a trpc default Logger from zap whose caller skip is set to 2.
func NewZapLog(c config.LogOutputConfigs) ILogger {
	return NewZapLogWithCallerSkip(c, 2)
}

// NewZapLogWithCallerSkip creates a trpc default Logger from zap.
func NewZapLogWithCallerSkip(cfg config.LogOutputConfigs, callerSkip int) ILogger {
	var (
		cores  []zapcore.Core
		levels []zap.AtomicLevel
	)
	for _, c := range cfg {
		// 根据插件名称从插件仓库获取 writer 插件
		writer := output.GetLogWriter(c.Writer)
		if writer == nil {
			panic("log: writer core: " + c.Writer + " no registered")
		}

		decoder := &Decoder{OutputConfig: &c}
		if err := writer.Setup(c.Writer, decoder); err != nil {
			panic("log: writer core: " + c.Writer + " setup fail: " + err.Error())
		}
		cores = append(cores, decoder.Core)
		levels = append(levels, decoder.ZapLevel)
	}
	return &zapLog{
		levels: levels,
		logger: zap.New(
			zapcore.NewTee(cores...),
			zap.AddCallerSkip(callerSkip),
			zap.AddCaller(),
		),
	}
}

// zapLog is a Logger implementation based on zaplogger.
type zapLog struct {
	levels []zap.AtomicLevel
	logger *zap.Logger
}

func (l *zapLog) WithOptions(opts ...Option) ILogger {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return &zapLog{
		levels: l.levels,
		logger: l.logger.WithOptions(zap.AddCallerSkip(o.skip)),
	}
}

// With add user defined fields to Logger. Fields support multiple values.
func (l *zapLog) With(fields ...Field) ILogger {
	zapFields := make([]zap.Field, len(fields))
	for i := range fields {
		zapFields[i] = zap.Any(fields[i].Key, fields[i].Value)
	}

	return &zapLog{
		levels: l.levels,
		logger: l.logger.With(zapFields...)}
}

// Trace logs to TRACE log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Trace(args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(getLogMsg(args...))
	}
}

// Tracef logs to TRACE log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Tracef(format string, args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(getLogMsgf(format, args...))
	}
}

// Debug logs to DEBUG log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Debug(args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(getLogMsg(args...))
	}
}

// Debugf logs to DEBUG log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Debugf(format string, args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(getLogMsgf(format, args...))
	}
}

// Info logs to INFO log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Info(args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.InfoLevel) {
		l.logger.Info(getLogMsg(args...))
	}
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Infof(format string, args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.InfoLevel) {
		l.logger.Info(getLogMsgf(format, args...))
	}
}

// Warn logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Warn(args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.WarnLevel) {
		l.logger.Warn(getLogMsg(args...))
	}
}

// Warnf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Warnf(format string, args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.WarnLevel) {
		l.logger.Warn(getLogMsgf(format, args...))
	}
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Error(args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.ErrorLevel) {
		l.logger.Error(getLogMsg(args...))
	}
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Errorf(format string, args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.ErrorLevel) {
		l.logger.Error(getLogMsgf(format, args...))
	}
}

// Fatal logs to FATAL log. Arguments are handled in the manner of fmt.Println.
func (l *zapLog) Fatal(args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.FatalLevel) {
		l.logger.Fatal(getLogMsg(args...))
	}
}

// Fatalf logs to FATAL log. Arguments are handled in the manner of fmt.Printf.
func (l *zapLog) Fatalf(format string, args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.FatalLevel) {
		l.logger.Fatal(getLogMsgf(format, args...))
	}
}

// Sync calls the zap logger's Sync method, and flushes any buffered log entries.
// Applications should take care to call Sync before exiting.
func (l *zapLog) Sync() error {
	return l.logger.Sync()
}

// SetLevel sets output log level.
func (l *zapLog) SetLevel(output string, _level level.Level) {
	i, e := strconv.Atoi(output)
	if e != nil {
		return
	}
	if i < 0 || i >= len(l.levels) {
		return
	}
	l.levels[i].SetLevel(level.LevelToZapLevel[_level])
}

// GetLevel gets output log level.
func (l *zapLog) GetLevel(output string) level.Level {
	i, e := strconv.Atoi(output)
	if e != nil {
		return level.LevelDebug
	}
	if i < 0 || i >= len(l.levels) {
		return level.LevelDebug
	}
	return level.ZapLevelToLevel[l.levels[i].Level()]
}

func getLogMsg(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	msg = msg[:len(msg)-1]
	//report.LogWriteSize.IncrBy(float64(len(msg)))
	return msg
}

func getLogMsgf(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	// report.LogWriteSize.IncrBy(float64(len(msg)))
	return msg
}

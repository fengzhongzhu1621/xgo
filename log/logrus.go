package log

import (
	"context"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"xgo"

	"xgo/pool"
	"xgo/utils"
)

type exitFunc func(int)

// LogFunction For big messages, it can be more efficient to pass a function
// and only call it if the log level is actually enables rather than
// generating the log message and then checking if the level is enabled
type LogFunction func() []interface{}

// StdLogger is what your logrus-enabled library should take, that way
// it'll accept a stdlib logger and a logrus logger. There's no standard
// interface, this is the closest we get, unfortunately.
type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

// The FieldLogger interface generalizes the Entry and Logger types
type FieldLogger interface {
	WithField(key string, value interface{}) *Entry
	WithFields(fields Fields) *Entry
	WithError(err error) *Entry

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	// IsDebugEnabled() bool
	// IsInfoEnabled() bool
	// IsWarnEnabled() bool
	// IsErrorEnabled() bool
	// IsFatalEnabled() bool
	// IsPanicEnabled() bool
}

// Ext1FieldLogger (the first extension to FieldLogger) is superfluous, it is
// here for consistancy. Do not use. Use Logger or Entry instead.
type Ext1FieldLogger interface {
	FieldLogger
	Tracef(format string, args ...interface{})
	Trace(args ...interface{})
	Traceln(args ...interface{})
}

// var _ LogrusLogger = (*Ext1FieldLogger)(null)

type LogrusLogger struct {
	// The logs are `io.Copy`'d to this in a mutex. It's common to set this to a
	// file, or leave it default which is `os.Stderr`. You can also set this to
	// something more adventurous, such as logging to Kafka.
	Out io.Writer
	// Hooks for the logger instance. These allow firing events based on logging
	// levels and log entries. For example, to send errors to an error tracking
	// service, log to StatsD or dump the core on fatal errors.
	Hooks LevelHooks
	// All log entries pass through the formatter before logged to Out. The
	// included formatters are `TextFormatter` and `JSONFormatter` for which
	// TextFormatter is the default. In development (when a TTY is attached) it
	// logs with colors, but to a file it wouldn't. You can easily implement your
	// own that implements the `Formatter` interface, see the `README` or included
	// formatters for examples.
	Formatter Formatter

	// Flag for whether to log caller info (off by default)
	ReportCaller bool

	// The logging level the logger should log at. This is typically (and defaults
	// to) `logrus.Info`, which allows Info(), Warn(), Error() and Fatal() to be
	// logged.
	Level Level
	// Used to sync writing to the log. Locking is enabled by Default
	mu utils.MutexWrap
	// Reusable empty entry
	entryPool sync.Pool
	// Function to exit the application, defaults to `os.Exit()`
	ExitFunc exitFunc
	// The buffer pool used to format the log. If it is nil, the default global
	// buffer pool will be used.
	BufferPool pool.BufferPool
}

// Creates a new logger. Configuration should be set by changing `Formatter`,
// `Out` and `Hooks` directly on the default logger instance.
// It's recommended to make this a global instance called `log`.
func New() *LogrusLogger {
	return &LogrusLogger{
		Out:          os.Stderr,          // 日志默认输出到标准输出
		Formatter:    new(TextFormatter), // 日志输出的格式
		Hooks:        make(LevelHooks),   // 日志的hooks，可以有多个
		Level:        InfoLevel,          // 默认日志级别
		ExitFunc:     os.Exit,            // 日志的默认退出行为，是停止当前进程
		ReportCaller: false,              // 是否输出堆栈信息
	}
}

// 新建一条空日志
func (logger *LogrusLogger) newEntry() *Entry {
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger)
}

func (logger *LogrusLogger) releaseEntry(entry *Entry) {
	entry.Data = map[string]interface{}{}
	logger.entryPool.Put(entry)
}

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func (logger *LogrusLogger) WithField(key string, value interface{}) *Entry {
	// 新建一条空日志
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	// 给空日志填充内容，生成一条新的日志
	return entry.WithField(key, value)
}

// WithFields Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *LogrusLogger) WithFields(fields Fields) *Entry {
	// 新建一条空日志
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	// 给空日志填充内容，生成一条新的日志
	return entry.WithFields(fields)
}

// WithError Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
// 给空日志填充一个key为error的错误
func (logger *LogrusLogger) WithError(err error) *Entry {
	// 新建一条空日志
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithError(err)
}

// WithContext Add a context to the log entry.
func (logger *LogrusLogger) WithContext(ctx context.Context) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithContext(ctx)
}

// Overrides the time of the log entry.
func (logger *LogrusLogger) WithTime(t time.Time) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithTime(t)
}

// level 获得日志级别，注意并发原子性
func (logger *LogrusLogger) level() Level {
	return Level(atomic.LoadUint32((*uint32)(&logger.Level)))
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (logger *LogrusLogger) IsLevelEnabled(level Level) bool {
	return logger.level() >= level
}

// Logf 打印日志
func (logger *LogrusLogger) Logf(level Level, format string, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.Logf(level, format, args...)
		logger.releaseEntry(entry)
	}
}

// Log will log a message at the level given as parameter.
// Warning: using Log at Panic or Fatal level will not respectively Panic nor Exit.
// For this behaviour Logger.Panic or Logger.Fatal should be used instead.
func (logger *LogrusLogger) Log(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.Log(level, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *LogrusLogger) Info(args ...interface{}) {
	logger.Log(InfoLevel, args...)
}

func (logger *LogrusLogger) Fatal(args ...interface{}) {
	logger.Log(FatalLevel, args...)
	logger.Exit(1)
}

// SetLevel sets the logger level.
func (logger *LogrusLogger) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&logger.Level), uint32(level))
}

func (logger *LogrusLogger) Panic(args ...interface{}) {
	logger.Log(PanicLevel, args...)
}

func (logger *LogrusLogger) Exit(code int) {
	xgo.RunHandlers()
	if logger.ExitFunc == nil {
		logger.ExitFunc = os.Exit
	}
	logger.ExitFunc(code)
}

func (logger *LogrusLogger) SetReportCaller(reportCaller bool) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.ReportCaller = reportCaller
}

// SetOutput sets the logger output.
func (logger *LogrusLogger) SetOutput(output io.Writer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Out = output
}

// SetFormatter sets the logger formatter.
func (logger *LogrusLogger) SetFormatter(formatter Formatter) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Formatter = formatter
}

// AddHook adds a hook to the logger hooks.
func (logger *LogrusLogger) AddHook(hook Hook) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Hooks.Add(hook)
}

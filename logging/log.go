package logging

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/fengzhongzhu1621/xgo/logging/level"
	loglevel "github.com/fengzhongzhu1621/xgo/logging/level"
)

const DATETIME_DEFAULT_FORMAT string = "2006-01-02 15:04:05"

type LoggerConfig struct {
	flag      int // properties
	Formatter ILogFormatter
	level     loglevel.LogLevel
	mu        sync.Mutex
	buf       []byte
	out       io.Writer
}

type ILogFormatter interface {
	Format(loggerConfig *LoggerConfig, t time.Time, level loglevel.LogLevel, message string) error
}

type DefaultFormatter struct{}

func (formatter *DefaultFormatter) Format(
	loggerConfig *LoggerConfig,
	t time.Time,
	level level.LogLevel,
	message string,
) error {
	buf := &loggerConfig.buf
	*buf = append(*buf, t.Format(DATETIME_DEFAULT_FORMAT)...)
	// year, month, day := t.Date()
	// bytesconv.Itoa(buf, year, 4)
	// *buf = append(*buf, '-')
	// bytesconv.Itoa(buf, int(month), 2)
	// *buf = append(*buf, '-')
	// bytesconv.Itoa(buf, day, 2)
	*buf = append(*buf, '|')
	*buf = append(*buf, loglevel.LogLevelToString(level)...)
	*buf = append(*buf, '|')
	*buf = append(*buf, message...)
	if len(message) == 0 || message[len(message)-1] != '\n' {
		*buf = append(*buf, '\n')
	}
	return nil
}

var loggerConfig = NewLogger(os.Stderr)

func NewLogger(w io.Writer) *LoggerConfig {
	return &LoggerConfig{out: os.Stderr, level: level.LOG_INFO, Formatter: &DefaultFormatter{}}
}

func (loggerConfig *LoggerConfig) SetFlags(flag int) {
	loggerConfig.mu.Lock()
	defer loggerConfig.mu.Unlock()
	loggerConfig.flag = flag
}

func (loggerConfig *LoggerConfig) SetFormatter(formatter ILogFormatter) {
	loggerConfig.mu.Lock()
	defer loggerConfig.mu.Unlock()
	loggerConfig.Formatter = formatter
}

func (loggerConfig *LoggerConfig) SetOutputWriter(out io.Writer) {
	loggerConfig.mu.Lock()
	defer loggerConfig.mu.Unlock()
	loggerConfig.out = out
}

func (loggerConfig *LoggerConfig) SetLevel(level level.LogLevel) {
	loggerConfig.mu.Lock()
	defer loggerConfig.mu.Unlock()
	loggerConfig.level = level
}

func (loggerConfig *LoggerConfig) Output(calldepth int, level level.LogLevel, message string) error {
	now := time.Now() // get this early.
	loggerConfig.mu.Lock()
	defer loggerConfig.mu.Unlock()
	loggerConfig.buf = loggerConfig.buf[:0]
	loggerConfig.Formatter.Format(loggerConfig, now, level, message)
	_, err := loggerConfig.out.Write(loggerConfig.buf)
	return err
}

func (loggerConfig *LoggerConfig) log(level level.LogLevel, message string) error {
	if loggerConfig.level > level {
		return nil
	}
	return loggerConfig.Output(2, level, message)
}

func (loggerConfig *LoggerConfig) Info(message string) error {
	return loggerConfig.log(level.LOG_INFO, message)
}

func Info(message string) error {
	return loggerConfig.Info(message)
}

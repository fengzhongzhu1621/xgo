package log

import (
	"io"
	"os"
	"sync"
	"time"
)

const DATETIME_DEFAULT_FORMAT string = "2006-01-02 15:04:05"

type (
	LogLevel int
)

const (
	LOG_EMERG   = LogLevel(0) //  system is unusable
	LOG_ALERT   = LogLevel(1) //  action must be taken immediately
	LOG_CRIT    = LogLevel(2) //  critical conditions
	LOG_ERR     = LogLevel(3) //  error conditions
	LOG_WARNING = LogLevel(4) //  warning conditions
	LOG_NOTICE  = LogLevel(5) //  normal but significant condition
	LOG_INFO    = LogLevel(6) //  informational
	LOG_DEBUG   = LogLevel(7) //  debug-level messages
)

func logLevelToString(t LogLevel) string {
	switch t {
	case LOG_ERR:
		return "ERROR"
	case LOG_WARNING:
		return "WARNING"
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_INFO:
		return "INFO"
	}
	return "unknown"
}

type LoggerConfig struct {
	flag      int // properties
	Formatter ILogFormatter
	level     LogLevel
	mu        sync.Mutex
	buf       []byte
	out       io.Writer
}

type ILogFormatter interface {
	Format(logger *LoggerConfig, t time.Time, level LogLevel, message string) error
}

type DefaultFormatter struct {
}

func (formatter *DefaultFormatter) Format(logger *LoggerConfig, t time.Time, level LogLevel, message string) error {
	buf := &logger.buf
	*buf = append(*buf, t.Format(DATETIME_DEFAULT_FORMAT)...)
	//year, month, day := t.Date()
	//bytesconv.Itoa(buf, year, 4)
	//*buf = append(*buf, '-')
	//bytesconv.Itoa(buf, int(month), 2)
	//*buf = append(*buf, '-')
	//bytesconv.Itoa(buf, day, 2)
	*buf = append(*buf, '|')
	*buf = append(*buf, logLevelToString(level)...)
	*buf = append(*buf, '|')
	*buf = append(*buf, message...)
	if len(message) == 0 || message[len(message)-1] != '\n' {
		*buf = append(*buf, '\n')
	}
	return nil
}

var logger = NewLogger(os.Stderr)

func NewLogger(w io.Writer) *LoggerConfig {
	return &LoggerConfig{out: os.Stderr, level: LOG_INFO, Formatter: &DefaultFormatter{}}
}

func (logger *LoggerConfig) SetFlags(flag int) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.flag = flag
}

func (logger *LoggerConfig) SetFormatter(formatter ILogFormatter) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.Formatter = formatter
}

func (logger *LoggerConfig) SetOutputWriter(out io.Writer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.out = out
}

func (logger *LoggerConfig) SetLevel(level LogLevel) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.level = level
}

func (logger *LoggerConfig) Output(calldepth int, level LogLevel, message string) error {
	now := time.Now() // get this early.
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.buf = logger.buf[:0]
	logger.Formatter.Format(logger, now, level, message)
	_, err := logger.out.Write(logger.buf)
	return err
}

func (logger *LoggerConfig) log(level LogLevel, message string) error {
	if logger.level > level {
		return nil
	}
	return logger.Output(2, level, message)
}

func (logger *LoggerConfig) Info(message string) error {
	return logger.log(LOG_INFO, message)
}

func Info(message string) error {
	return logger.Info(message)
}

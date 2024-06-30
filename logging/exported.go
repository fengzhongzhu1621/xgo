package logging

import (
	"context"
	"io"
	"time"
)

var (
	// std is the name of the standard logger in stdlib `log`
	std = New()
)

func LogrusStandardLogger() *LogrusLogger {
	return std
}

// SetLevel sets the standard logger level.
func LogrusSetLevel(level Level) {
	std.SetLevel(level)
}

// IsLevelEnabled checks if the log level of the standard logger is greater than the level param
func LogrusIsLevelEnabled(level Level) bool {
	return std.IsLevelEnabled(level)
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func LogrusWithError(err error) *Entry {
	return std.WithField(ErrorKey, err)
}

// WithContext creates an entry from the standard logger and adds a context to it.
func LogrusWithContext(ctx context.Context) *Entry {
	return std.WithContext(ctx)
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func LogrusWithField(key string, value interface{}) *Entry {
	return std.WithField(key, value)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func LogrusWithFields(fields Fields) *Entry {
	return std.WithFields(fields)
}

// WithTime creates an entry from the standard logger and overrides the time of
// logs generated with it.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func LogrusWithTime(t time.Time) *Entry {
	return std.WithTime(t)
}

// Info logs a message at level Info on the standard logger.
func LogrusInfo(args ...interface{}) {
	std.Info(args...)
}

// SetReportCaller sets whether the standard logger will include the calling
// method as a field.
func LogrusSetReportCaller(include bool) {
	std.SetReportCaller(include)
}

// SetOutput sets the standard logger output.
func LogrusSetOutput(out io.Writer) {
	std.SetOutput(out)
}

// SetFormatter sets the standard logger formatter.
func LogrusSetFormatter(formatter Formatter) {
	std.SetFormatter(formatter)
}

// AddHook adds a hook to the standard logger hooks.
func LogrusAddHook(hook Hook) {
	std.AddHook(hook)
}

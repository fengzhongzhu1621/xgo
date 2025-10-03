package logging

import (
	"io"

	"github.com/fengzhongzhu1621/xgo/logging/level"
)

// Option modifies the options of optionLogger.
type Option func(*options)

type options struct {
	skip int
}

// WithAdditionalCallerSkip adds additional caller skip.
func WithAdditionalCallerSkip(skip int) Option {
	return func(o *options) {
		o.skip = skip
	}
}

// LoggerOptions is the log options.
type LoggerOptions struct {
	LogLevel level.Level
	Pattern  string
	Writer   io.Writer
}

// LoggerOption modifies the LoggerOptions.
type LoggerOption func(*LoggerOptions)

// IOptionLogger defines logger with additional options.
type IOptionLogger interface {
	WithOptions(opts ...Option) ILogger
}

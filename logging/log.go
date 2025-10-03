package logging

import (
	"context"
	"errors"

	"github.com/fengzhongzhu1621/xgo/codec"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// With adds user defined fields to Logger. Field support multiple values.
func With(fields ...Field) ILogger {
	if ol, ok := GetDefaultLogger().(IOptionLogger); ok {
		return ol.WithOptions(WithAdditionalCallerSkip(-1)).With(fields...)
	}
	return GetDefaultLogger().With(fields...)
}

// WithContext add user defined fields to the Logger of context. Fields support multiple values.
func WithContext(ctx context.Context, fields ...Field) ILogger {
	logger, ok := codec.Message(ctx).Logger().(ILogger)
	if !ok {
		// 降级到默认logger
		return With(fields...)
	}
	if ol, ok := logger.(IOptionLogger); ok {
		return ol.WithOptions(WithAdditionalCallerSkip(-1)).With(fields...)
	}
	return logger.With(fields...)
}

// WithContextFields sets some user defined data to logs, such as uid, imei, etc.
// Fields must be paired.
// If ctx has already set a Msg, this function returns that ctx, otherwise, it returns a new one.
func WithContextFields(ctx context.Context, fields ...string) context.Context {
	tagCapacity := len(fields) / 2
	tags := make([]Field, 0, tagCapacity)
	for i := 0; i < tagCapacity; i++ {
		tags = append(tags, Field{
			Key:   fields[2*i],
			Value: fields[2*i+1],
		})
	}

	// 获取消息
	ctx, msg := codec.EnsureMessage(ctx)
	logger, ok := msg.Logger().(ILogger)
	if ok && logger != nil {
		logger = logger.With(tags...)
	} else {
		logger = GetDefaultLogger().With(tags...)
	}

	// 将 logger 加入到 msg 中
	msg.WithLogger(logger)
	return ctx
}

// RedirectStdLog redirects std log to trpc logger as log level INFO.
// After redirection, log flag is zero, the prefix is empty.
// The returned function may be used to recover log flag and prefix, and redirect output to
// os.Stderr.
func RedirectStdLog(logger ILogger) (func(), error) {
	return RedirectStdLogAt(logger, zap.InfoLevel)
}

// RedirectStdLogAt redirects std log to trpc logger with a specific level.
// After redirection, log flag is zero, the prefix is empty.
// The returned function may be used to recover log flag and prefix, and redirect output to
// os.Stderr.
func RedirectStdLogAt(logger ILogger, level zapcore.Level) (func(), error) {
	if l, ok := logger.(*zapLog); ok {
		return zap.RedirectStdLogAt(l.logger, level)
	}

	return nil, errors.New("log: only supports redirecting std logs to trpc zap logger")
}

// Trace logs to TRACE log. Arguments are handled in the manner of fmt.Println.
func Trace(args ...interface{}) {
	if traceEnabled {
		GetDefaultLogger().Trace(args...)
	}
}

// Tracef logs to TRACE log. Arguments are handled in the manner of fmt.Printf.
func Tracef(format string, args ...interface{}) {
	if traceEnabled {
		GetDefaultLogger().Tracef(format, args...)
	}
}

// TraceContext logs to TRACE log. Arguments are handled in the manner of fmt.Println.
func TraceContext(ctx context.Context, args ...interface{}) {
	if !traceEnabled {
		return
	}
	if l, ok := codec.Message(ctx).Logger().(ILogger); ok {
		l.Trace(args...)
		return
	}
	GetDefaultLogger().Trace(args...)
}

// TraceContextf logs to TRACE log. Arguments are handled in the manner of fmt.Printf.
func TraceContextf(ctx context.Context, format string, args ...interface{}) {
	if !traceEnabled {
		return
	}
	if l, ok := codec.Message(ctx).Logger().(ILogger); ok {
		l.Tracef(format, args...)
		return
	}
	GetDefaultLogger().Tracef(format, args...)
}

// Debug logs to DEBUG log. Arguments are handled in the manner of fmt.Println.
func Debug(args ...interface{}) {
	GetDefaultLogger().Debug(args...)
}

// Debugf logs to DEBUG log. Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, args ...interface{}) {
	GetDefaultLogger().Debugf(format, args...)
}

// DebugContext logs to DEBUG log. Arguments are handled in the manner of fmt.Println.
func DebugContext(ctx context.Context, args ...interface{}) {
	if l, ok := codec.Message(ctx).Logger().(ILogger); ok {
		l.Debug(args...)
		return
	}
	GetDefaultLogger().Debug(args...)
}

// DebugContextf logs to DEBUG log. Arguments are handled in the manner of fmt.Printf.
func DebugContextf(ctx context.Context, format string, args ...interface{}) {
	if l, ok := codec.Message(ctx).Logger().(ILogger); ok {
		l.Debugf(format, args...)
		return
	}
	GetDefaultLogger().Debugf(format, args...)
}

// Info logs to INFO log. Arguments are handled in the manner of fmt.Println.
func Info(args ...interface{}) {
	GetDefaultLogger().Info(args...)
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func Infof(format string, args ...interface{}) {
	GetDefaultLogger().Infof(format, args...)
}

// InfoContext logs to INFO log. Arguments are handled in the manner of fmt.Println.
func InfoContext(ctx context.Context, args ...interface{}) {
	if l, ok := codec.Message(ctx).Logger().(ILogger); ok {
		l.Info(args...)
		return
	}
	GetDefaultLogger().Info(args...)
}

// InfoContextf logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func InfoContextf(ctx context.Context, format string, args ...interface{}) {
	if l, ok := codec.Message(ctx).Logger().(ILogger); ok {
		l.Infof(format, args...)
		return
	}
	GetDefaultLogger().Infof(format, args...)
}

// Warn logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func Warn(args ...interface{}) {
	GetDefaultLogger().Warn(args...)
}

// Warnf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func Warnf(format string, args ...interface{}) {
	GetDefaultLogger().Warnf(format, args...)
}

// WarnContext logs to WARNING log. Arguments are handled in the manner of fmt.Println.
func WarnContext(ctx context.Context, args ...interface{}) {
	if l, ok := codec.Message(ctx).Logger().(ILogger); ok {
		l.Warn(args...)
		return
	}
	GetDefaultLogger().Warn(args...)
}

// WarnContextf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func WarnContextf(ctx context.Context, format string, args ...interface{}) {
	if l, ok := codec.Message(ctx).Logger().(ILogger); ok {
		l.Warnf(format, args...)
		return
	}
	GetDefaultLogger().Warnf(format, args...)

}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func Error(args ...interface{}) {
	GetDefaultLogger().Error(args...)
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, args ...interface{}) {
	GetDefaultLogger().Errorf(format, args...)
}

// ErrorContext logs to ERROR log. Arguments are handled in the manner of fmt.Println.
func ErrorContext(ctx context.Context, args ...interface{}) {
	if l, ok := codec.Message(ctx).Logger().(ILogger); ok {
		l.Error(args...)
		return
	}
	GetDefaultLogger().Error(args...)
}

// ErrorContextf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func ErrorContextf(ctx context.Context, format string, args ...interface{}) {
	if l, ok := codec.Message(ctx).Logger().(ILogger); ok {
		l.Errorf(format, args...)
		return
	}
	GetDefaultLogger().Errorf(format, args...)
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Println.
// All Fatal logs will exit by calling os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func Fatal(args ...interface{}) {
	GetDefaultLogger().Fatal(args...)
}

// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func Fatalf(format string, args ...interface{}) {
	GetDefaultLogger().Fatalf(format, args...)
}

// FatalContext logs to ERROR log. Arguments are handled in the manner of fmt.Println.
// All Fatal logs will exit by calling os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func FatalContext(ctx context.Context, args ...interface{}) {
	if l, ok := codec.Message(ctx).Logger().(ILogger); ok {
		l.Fatal(args...)
		return
	}
	GetDefaultLogger().Fatal(args...)
}

// FatalContextf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func FatalContextf(ctx context.Context, format string, args ...interface{}) {
	if l, ok := codec.Message(ctx).Logger().(ILogger); ok {
		l.Fatalf(format, args...)
		return
	}
	GetDefaultLogger().Fatalf(format, args...)
}

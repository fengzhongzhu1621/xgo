package logging

import (
	"fmt"

	jww "github.com/spf13/jwalterweatherman"
)

// Logger is a unified interface for various logging use cases and practices, including:
//   - leveled logging
//   - structured logging
type Logger interface {
	// Trace logs a Trace event.
	//
	// Even more fine-grained information than Debug events.
	// Loggers not supporting this level should fall back to Debug.
	Trace(msg string, keyvals ...interface{})

	// Debug logs a Debug event.
	//
	// A verbose series of information events.
	// They are useful when debugging the system.
	Debug(msg string, keyvals ...interface{})

	// Info logs an Info event.
	//
	// General information about what's happening inside the system.
	Info(msg string, keyvals ...interface{})

	// Warn logs a Warn(ing) event.
	//
	// Non-critical events that should be looked at.
	Warn(msg string, keyvals ...interface{})

	// Error logs an Error event.
	//
	// Critical events that require immediate attention.
	// Loggers commonly provide Fatal and Panic levels above Error level,
	// but exiting and panicing is out of scope for a logging library.
	Error(msg string, keyvals ...interface{})
}

type JwwLogger struct{}

func (JwwLogger) Trace(msg string, keyvals ...interface{}) {
	jww.TRACE.Print(jwwLogMessage(msg, keyvals...)) //nolint:govet
}

func (JwwLogger) Debug(msg string, keyvals ...interface{}) {
	jww.DEBUG.Print(jwwLogMessage(msg, keyvals...)) //nolint:govet
}

func (JwwLogger) Info(msg string, keyvals ...interface{}) {
	jww.INFO.Print(jwwLogMessage(msg, keyvals...)) //nolint:govet
}

func (JwwLogger) Warn(msg string, keyvals ...interface{}) {
	jww.WARN.Print("%s", jwwLogMessage(msg, keyvals...)) //vet:ignore
}

func (JwwLogger) Error(msg string, keyvals ...interface{}) {
	jww.ERROR.Print("%s", jwwLogMessage(msg, keyvals...)) //vet:ignore
}

// 格式化消息，加上参数值.
func jwwLogMessage(msg string, keyvals ...interface{}) string {
	out := msg

	// 参数为计数个，补一个nil，变成偶数数量
	if len(keyvals) > 0 && len(keyvals)%2 == 1 {
		keyvals = append(keyvals, nil)
	}

	// 打印日志时，同时打印参数值
	for i := 0; i <= len(keyvals)-2; i += 2 {
		out = fmt.Sprintf("%s %v=%v", out, keyvals[i], keyvals[i+1])
	}

	return out
}

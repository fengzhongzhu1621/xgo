package debuglog

import (
	"context"

	"trpc.group/trpc-go/trpc-go/log"
)

// LogLevelFunc specifies the log level.
type LogLevelFunc func(ctx context.Context, format string, args ...any)

// logLevel is log level.
type logLevel = string

var (
	traceLevel   logLevel = "trace"
	debugLevel   logLevel = "debug"
	warningLevel logLevel = "warning"
	infoLevel    logLevel = "info"
	errorLevel   logLevel = "error"
	fatalLevel   logLevel = "fatal"
)

// LogContextfFuncs is a map of methods for logging at different levels.
var LogContextfFuncs = map[string]func(ctx context.Context, format string, args ...any){
	traceLevel:   log.TraceContextf,
	debugLevel:   log.DebugContextf,
	warningLevel: log.WarnContextf,
	infoLevel:    log.InfoContextf,
	errorLevel:   log.ErrorContextf,
	fatalLevel:   log.FatalContextf,
}

// getLogLevelFunc gets the log print method for the corresponding log level.
func getLogLevelFunc(level string, defaultLevel string) LogLevelFunc {
	logFunc, ok := LogContextfFuncs[level]
	if !ok {
		logFunc = LogContextfFuncs[defaultLevel]
	}
	return logFunc
}

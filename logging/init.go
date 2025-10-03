package logging

import (
	"github.com/fengzhongzhu1621/xgo/logging/output"
	"github.com/fengzhongzhu1621/xgo/plugin"
)

func init() {
	output.RegisterWriter(output.OutputConsole, DefaultConsoleWriterFactory)
	output.RegisterWriter(output.OutputFile, DefaultFileWriterFactory)
	Register(defaultLoggerName, NewZapLog(defaultConfig))
	plugin.Register(defaultLoggerName, DefaultLogFactory)
}

// Register registers Logger. It supports multiple Logger implementation.
func Register(name string, logger ILogger) {
	mu.Lock()
	defer mu.Unlock()
	if logger == nil {
		panic("log: Register logger is nil")
	}
	if _, dup := loggerMap[name]; dup && name != defaultLoggerName {
		panic("log: Register called twiced for logger name " + name)
	}
	loggerMap[name] = logger
	if name == defaultLoggerName {
		DefaultLogger = logger
	}
}

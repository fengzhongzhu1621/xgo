package logging

import (
	"github.com/fengzhongzhu1621/xgo/logging/output"
	"github.com/fengzhongzhu1621/xgo/plugin"
)

func init() {
	// 注册 log writer 插件
	output.RegisterWriter(output.OutputConsole, DefaultConsoleWriterFactory)
	output.RegisterWriter(output.OutputFile, DefaultFileWriterFactory)
	// 创建默认 logger 实例
	Register(defaultLoggerName, NewZapLog(defaultConfig))
	// 注册 log default 插件
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

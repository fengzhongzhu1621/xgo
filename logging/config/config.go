package config

import (
	"github.com/fengzhongzhu1621/xgo/datetime"
	yaml "gopkg.in/yaml.v3"
)

// Config is the log config. Each log may have multiple outputs.
type LogOutputConfigs []LogOutputConfig

// OutputConfig is the output config, includes console, file and remote.
type LogOutputConfig struct {
	// Writer is the output of log, such as console or file.
	Writer      string         `yaml:"writer"`
	WriteConfig LogWriteConfig `yaml:"writer_config"`

	// Formatter is the format of log, such as console or json.
	Formatter    string          `yaml:"formatter"`
	FormatConfig LogFormatConfig `yaml:"formatter_config"`

	// RemoteConfig is the remote config. It's defined by business and should be registered by
	// third-party modules.
	RemoteConfig yaml.Node `yaml:"remote_config"`

	// Level controls the log level, like debug, info or error.
	Level string `yaml:"level"`

	// CallerSkip controls the nesting depth of log function.
	CallerSkip int `yaml:"caller_skip"`

	// EnableColor determines if the output is colored. The default value is false.
	EnableColor bool `yaml:"enable_color"`
}

// LogWriteConfig is the local file config.
type LogWriteConfig struct {
	// LogPath is the log path like /usr/local/trpc/log/.
	LogPath string `yaml:"log_path"`
	// Filename is the file name like trpc.log.
	Filename string `yaml:"filename"`
	// WriteMode is the log write mod. 1: sync, 2: async, 3: fast(maybe dropped), default as 3.
	WriteMode int `yaml:"write_mode"`
	// RollType is the log rolling type. Split files by size/time, default by size.
	RollType string `yaml:"roll_type"`
	// MaxAge is the max expire times(day).
	MaxAge int `yaml:"max_age"`
	// MaxBackups is the max backup files.
	MaxBackups int `yaml:"max_backups"`
	// Compress defines whether log should be compressed.
	Compress bool `yaml:"compress"`
	// MaxSize is the max size of log file(MB).
	MaxSize int `yaml:"max_size"`

	// TimeUnit splits files by time unit, like year/month/hour/minute, default day.
	// It takes effect only when split by time.
	TimeUnit datetime.TimeUnit `yaml:"time_unit"`
}

// FormatConfig is the log format config.
type LogFormatConfig struct {
	// TimeFmt is the time format of log output, default as "2006-01-02 15:04:05.000" on empty.
	TimeFmt string `yaml:"time_fmt"`

	// TimeKey is the time key of log output, default as "T".
	TimeKey string `yaml:"time_key"`
	// LevelKey is the level key of log output, default as "L".
	LevelKey string `yaml:"level_key"`
	// NameKey is the name key of log output, default as "N".
	NameKey string `yaml:"name_key"`
	// CallerKey is the caller key of log output, default as "C".
	CallerKey string `yaml:"caller_key"`
	// FunctionKey is the function key of log output, default as "", which means not to print
	// function name.
	FunctionKey string `yaml:"function_key"`
	// MessageKey is the message key of log output, default as "M".
	MessageKey string `yaml:"message_key"`
	// StackTraceKey is the stack trace key of log output, default as "S".
	StacktraceKey string `yaml:"stacktrace_key"`
}

type LogConfig struct {
	Level    string
	Writer   string
	Settings map[string]string // 日志详细配置
}

type Logger struct {
	System LogConfig
	API    LogConfig
	Web    LogConfig
}

type ServiceLogConfig struct {
	Level string
	Dir   string
}

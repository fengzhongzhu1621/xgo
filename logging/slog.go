package logging

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

var sloggers map[string]*slog.Logger

type Options struct {
	Level        string
	HandlerName  string
	WriterType   string
	WriterConfig map[string]string
}

func GetSLogger(name string) *slog.Logger {
	if logger, ok := sloggers[name]; ok {
		return logger
	}

	// 找不到获得默认日志处理器
	return slog.Default()
}

// toSlogLevel 将字符串转换为 slog 的日志级别
func toSlogLevel(level string) (slog.Level, error) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN", "WARNING":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	}

	// 默认是 info 级别
	return slog.LevelInfo, fmt.Errorf("[%s] level not supported", level)
}

func newLogger(opts *Options) (*slog.Logger, error) {
	w, err := NewSlogWriter(opts.WriterType, opts.WriterConfig)
	if err != nil {
		return nil, err
	}
	level, err := toSlogLevel(opts.Level)
	if err != nil {
		return nil, err
	}
	handlerOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) != 0 {
				return a
			}
			// 参考 python 的 logging 字段的设计
			switch a.Key {
			case slog.MessageKey:
				a.Key = "message"
			case slog.LevelKey:
				a.Key = "levelname"
			case slog.SourceKey:
				a.Key = "pathname"
			}

			return a
		},
	}

	// 支持两种日志输出方式
	switch opts.HandlerName {
	case "text":
		return slog.New(slog.NewTextHandler(w, handlerOpts)), nil
	case "json":
		return slog.New(slog.NewJSONHandler(w, handlerOpts)), nil
	}

	return nil, fmt.Errorf("[%s] handler not supported", opts.HandlerName)
}

func SLogDebug(format string, vars ...any) {
	SLog(slog.LevelDebug, format, vars...)
}

func SLogInfo(format string, vars ...any) {
	SLog(slog.LevelInfo, format, vars...)
}

func SLogWarn(format string, vars ...any) {
	SLog(slog.LevelWarn, format, vars...)
}

func SLogError(format string, vars ...any) {
	SLog(slog.LevelError, format, vars...)
}

func SLogFatal(format string, vars ...any) {
	// slog 没有 Fatal 级别
	logger := log.New(os.Stderr, "", log.LstdFlags)
	logger.Fatalf(format, vars...)
}

func SLog(level slog.Level, format string, vars ...any) {
	// 创建一个默认的上下文 ctx
	ctx := context.Background()
	// 获取 slog 的默认日志记录器 logger
	logger := slog.Default()
	// 判断当前的日志级别是否启用，过滤不需要记录的日志
	if !logger.Enabled(ctx, level) {
		return
	}

	// pcs 是一个长度为 1 的 uintptr 数组，用于存储获取到的程序计数器值。
	var pcs [1]uintptr
	// 使用 runtime.Callers 函数获取调用 SLog 函数的上一层函数的程序计数器（PC）值。这里的 3 表示向上追溯三层调用栈，
	// 这是因为 runtime.Callers 在 SLog 内部被调用了一次，SLog 被调用了一次，再加上我们要获取的是调用 SLog 的函数，所以总共是三层。
	runtime.Callers(3, pcs[:])

	// 创建一个新的日志记录 r，包含时间戳、日志级别、格式化后的消息和程序计数器值。
	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, vars...), pcs[0])
	// 将日志记录传递给日志处理器的 Handle 方法进行处理
	_ = logger.Handler().Handle(ctx, r)
}

func InitSLogger(name string, opts *Options) (err error) {
	// 尝试获取一个已初始化的 slogger
	if sloggers == nil {
		sloggers = make(map[string]*slog.Logger)
	}
	if _, ok := sloggers[name]; ok {
		return nil
	}

	// 没有则创建并缓存
	if sloggers[name], err = newLogger(opts); err != nil {
		return err
	}
	// 设置默认的 slogger
	if name == "default" {
		slog.SetDefault(sloggers[name])
	}

	return nil
}

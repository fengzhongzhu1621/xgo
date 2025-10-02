package slog

import (
	"log/slog"
	"os"
	"testing"
)

func TestSlogDefault(t *testing.T) {
	// 使用默认日志记录器,此模式日志输出格式为text文本型,且不会输出Debug级别日志
	logger := slog.Default()

	logger.Debug("This is a debug message")
	// 2024/12/24 00:25:14 INFO This is a info message a=b
	logger.Info("This is a info message", "a", "b")
	// 2024/12/24 00:25:14 WARN This is a warn message c=d
	logger.Warn("This is a warn message", "c", "d")
	// 2024/12/24 00:25:14 ERROR This is a error message e=f
	logger.Error("This is a error message", "e", "f")
}

func TestNewSlogNew(t *testing.T) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler).With("a", "b", "c", "d")

	// {"time":"2024-12-24T00:25:01.845953+08:00","level":"DEBUG","msg":"This is a debug message","a":"b","c":"d"}
	logger.Debug("This is a debug message")
	// {"time":"2024-12-24T00:26:03.516023+08:00","level":"INFO","msg":"This is a info message","a":"b","c":"d","a":"b"}
	logger.Info("This is a info message", "a", "b")
	// {"time":"2024-12-24T00:26:03.516033+08:00","level":"ERROR","msg":"This is a error message","a":"b","c":"d","e":"f"}
	logger.Error("This is a error message", "e", "f")
}

func TestSlogReplaceAttr(t *testing.T) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 替换关键字
			if a.Key == "a" {
				return slog.Attr{Key: a.Key, Value: slog.StringValue("a+")}
			}
			return a
		},
	})

	logger := slog.New(handler).With("a", "b", "c", "d")
	// {"time":"2024-12-24T00:28:42.86368+08:00","level":"DEBUG","msg":"This is a debug message","a":"a+","c":"d"}
	logger.Debug("This is a debug message")
	// {"time":"2024-12-24T00:28:42.864144+08:00","level":"INFO","msg":"This is a info message","a":"a+","c":"d","a":"a+"}
	logger.Info("This is a info message", "a", "b")
}

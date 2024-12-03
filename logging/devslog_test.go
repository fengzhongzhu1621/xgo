package logging

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/golang-cz/devslog"
	"github.com/google/uuid"
)

func LoggerWithTraceID(ctx context.Context, logger *slog.Logger) *slog.Logger {
	traceID := uuid.New().String()
	return logger.With(
		slog.String("trace_id", traceID),
	)
}

func TestNewHandler(t *testing.T) {
	logger := slog.New(devslog.NewHandler(os.Stdout, nil))
	logger.Info("这是信息日志")
	logger.Debug("这是调试日志")
	logger.Warn("这是警告日志")
	logger.Error("这是错误日志")
}

func TestLoggerWithTraceID(t *testing.T) {
	logger := slog.New(devslog.NewHandler(os.Stdout, &devslog.Options{
		NoColor: true,
	}))

	ctx := context.Background()
	contextLogger := LoggerWithTraceID(ctx, logger)

	contextLogger.Info("trace",
		slog.String("A", "a"),
		slog.String("B", "b"),
	)
}

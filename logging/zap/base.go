package zap

import (
	"go.uber.org/zap"
	adapter "logur.dev/adapter/zap"

	"github.com/fengzhongzhu1621/xgo/logging"
)

// SetLogger sets a logger to the logging package.
func SetLogger(name string, logger *zap.Logger) {
	logging.SetLogger(name, adapter.New(logger))
}

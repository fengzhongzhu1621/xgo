package report

import "github.com/fengzhongzhu1621/xgo/opentelemetry/metrics"

var (
	// -----------------------------log----------------------------- //
	// log is dropped because the queue is full.
	LogQueueDropNum = metrics.Counter("LogQueueDropNum")
)

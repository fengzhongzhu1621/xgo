package retry

import (
	"context"
	"time"

	"github.com/fengzhongzhu1621/xgo/trpc/plugins/slime/view"
)

type IBackoff interface {
	Backoff(attempt int) time.Duration
}

type ILogger interface {
	Printf(string, ...interface{})
}

type ILazyLogger interface {
	ILogger
	FlushCtx(context.Context)
}

type IReporter interface {
	Report(context.Context, view.IStat)
}

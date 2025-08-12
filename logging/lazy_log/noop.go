package lazylog

import "context"

// NoopLog empty implementation.
type NoopLog struct{}

// Printf empty implementation.
func (NoopLog) Printf(string, ...interface{}) {}

// Flush empty implementation.
func (NoopLog) Flush() {}

// FlushCtx empty implementation.
func (NoopLog) FlushCtx(context.Context) {}

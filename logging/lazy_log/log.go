package lazylog

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	timeFormat = "15:04:05.000"
)

// Logger is a simple interface to print an string.
type Logger interface {
	Println(string)
}

// CtxLogger accepts an additional context.
type CtxLogger interface {
	Println(context.Context, string)
}

// LazyLog buffers messages and flush them by Logger.
//
// LazyLog is not concurrent safe.
type LazyLog struct {
	log CtxLogger
	buf []string
}

// NewLazyLog create a new LazyLog.
func NewLazyLog(log Logger) *LazyLog {
	return &LazyLog{log: &noopCtxLog{log: log}, buf: []string{"[lazy log]"}}
}

// NewLazyCtxLog creates a new LazyLog.
func NewLazyCtxLog(log CtxLogger) *LazyLog {
	return &LazyLog{log: log, buf: []string{"[lazy log]"}}
}

// Printf provides a format printer.
func (l *LazyLog) Printf(format string, a ...interface{}) {
	l.buf = append(l.buf, time.Now().Format(timeFormat)+"]\t"+fmt.Sprintf(format, a...))
}

// Flush flushes messages in buffer. Messages are separated by a new line
// and flushed with a single call of Logger.Println.
func (l *LazyLog) Flush() {
	l.FlushCtx(context.Background())
}

// FlushCtx flushes messages in buffer. Messages are separated by a new line
// and flushed with a single call of Logger.Println.
func (l *LazyLog) FlushCtx(ctx context.Context) {
	l.log.Println(ctx, strings.Join(l.buf, "\n"))
	l.buf = l.buf[:0]
}

// noopCtxLog 将 Logger转换为 CtxLogger
type noopCtxLog struct {
	log Logger
}

func (l *noopCtxLog) Println(_ context.Context, s string) {
	l.log.Println(s)
}

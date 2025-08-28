package lazylog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type bufLog struct {
	buf string
}

func (l *bufLog) Println(s string) {
	l.buf += s
}

type ctxBufLog struct {
	buf string
}

func (l *ctxBufLog) Println(ctx context.Context, s string) {
	l.buf += s
}

func TestLazyLog(t *testing.T) {
	bufLog := bufLog{}
	l := NewLazyLog(&bufLog)
	l.Printf("aaa")
	l.Printf("%s", "bbb")
	require.Equal(t, "", bufLog.buf)

	l.Flush()
	require.Contains(t, bufLog.buf, "aaa")
	require.Contains(t, bufLog.buf, "bbb")

	l.Printf("ccc")
	l.Flush()
	require.Contains(t, bufLog.buf, "ccc")
}

func TestLazyCtxLog(t *testing.T) {
	ctxBufLog := ctxBufLog{}
	l := NewLazyCtxLog(&ctxBufLog)
	l.Printf("aaa")
	require.Equal(t, "", ctxBufLog.buf)
	l.FlushCtx(context.Background())
	require.Contains(t, ctxBufLog.buf, "aaa")
}

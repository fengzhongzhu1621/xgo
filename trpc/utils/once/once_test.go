package once

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

func TestInvoke(t *testing.T) {
	o := New()

	var attempt atomic.Int32
	handler := func(ctx context.Context, req, rsp interface{}) error {
		attempt.Inc()

		// 设置响应结果
		*rsp.(*int) = 1

		return nil
	}

	var rsp int
	err := o.Invoke(context.Background(), nil, &rsp, handler)
	require.Nil(t, err)
	require.Equal(t, rsp, 1)
	require.Equal(t, int32(1), attempt.Load())
}

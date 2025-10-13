package filter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerFilterChain(t *testing.T) {
	ctx := context.Background()
	req := "req"
	businessAction := func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		return nil, nil
	}
	// 注册 NoopServerFilter 过滤器
	sc := ServerChain{NoopServerFilter}

	// 执行过滤器
	// Filter() -> NoopServerFilter begin -> NoopServerFilter.next() -> businessAction() -> NoopServerFilter.end
	_, err := sc.Filter(ctx, req, businessAction)
	require.Nil(t, err)
}

package filter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNamedFilter(t *testing.T) {
	const filterName = "filterName"
	Register(filterName, NoopServerFilter, NoopClientFilter)
	require.NotNil(t, GetClient(filterName))
	require.NotNil(t, GetServer(filterName))

	ctx := context.Background()
	cc := ClientChain{NoopClientFilter}
	require.Nil(t, cc.Filter(ctx, nil, nil,
		func(ctx context.Context, req, rsp interface{}) error { return nil }))

	sc := ServerChain{NoopServerFilter}
	_, err := sc.Filter(ctx, nil,
		func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil })
	require.Nil(t, err)
}

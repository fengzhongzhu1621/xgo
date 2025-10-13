package filter

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestClientFilterChain(t *testing.T) {
	ctx := context.Background()
	req, rsp := "req", "rsp"

	cc := ClientChain{NoopClientFilter}
	require.Nil(t, cc.Filter(ctx, req, rsp,
		func(ctx context.Context, req, rsp interface{}) error {
			return nil
		}))
}

func TestChainConcurrentHandle(t *testing.T) {
	const concurrentN = 4
	var calledTimes [concurrentN]int32

	filterFunc1 := func(ctx context.Context, req interface{}, rsp interface{}, f ClientHandleFunc) error {
		atomic.AddInt32(&calledTimes[0], 1)
		return f(ctx, req, rsp)
	}
	filterFunc2 := func(ctx context.Context, req interface{}, rsp interface{}, f ClientHandleFunc) error {
		atomic.AddInt32(&calledTimes[1], 1)
		var eg errgroup.Group
		for i := 0; i < concurrentN; i++ {
			eg.Go(func() error {
				return f(ctx, req, rsp)
			})
		}
		return eg.Wait()
	}
	filterFunc3 := func(ctx context.Context, req interface{}, rsp interface{}, f ClientHandleFunc) (err error) {
		atomic.AddInt32(&calledTimes[2], 1)
		return f(ctx, req, rsp)
	}

	filterFunc4 := func(ctx context.Context, req interface{}, rsp interface{}, f ClientHandleFunc) (err error) {
		atomic.AddInt32(&calledTimes[3], 1)
		return f(ctx, req, rsp)
	}

	// 注册 4 个过滤器
	cc := ClientChain{filterFunc1, filterFunc2, filterFunc3, filterFunc4}

	// 执行过滤器
	// Filter() -> filterFunc1 begin -> filterFunc1.f() ->
	//             filterFunc2 begin -> filterFunc2.f() ->
	//             filterFunc3 begin -> filterFunc3.f() ->
	//             filterFunc4 begin -> filterFunc4.f() ->
	//             businessAction() ->
	//             filterFunc4.end ->
	//             filterFunc3.end ->
	//             filterFunc2.end ->
	//             filterFunc1.end
	businessAction := func(ctx context.Context, req, rsp interface{}) (err error) {
		return nil
	}
	require.Nil(t, cc.Filter(context.Background(), nil, nil, businessAction))

	require.Equal(t, int32(1), atomic.LoadInt32(&calledTimes[0]))
	require.Equal(t, int32(1), atomic.LoadInt32(&calledTimes[1]))
	require.Equal(t, int32(concurrentN), atomic.LoadInt32(&calledTimes[2]))
	require.Equal(t, int32(concurrentN), atomic.LoadInt32(&calledTimes[3]))
}

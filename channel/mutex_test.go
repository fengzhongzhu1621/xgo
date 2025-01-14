package channel

import (
	"sync"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/tests"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// Synchronize 	将底层回调封装在互斥对象中。它接收一个可选的互斥锁
func TestSynchronize(t *testing.T) {
	t.Parallel()
	tests.TestWithTimeout(t, 100*time.Millisecond)
	is := assert.New(t)

	// check that callbacks are not executed concurrently
	{
		start := time.Now()

		wg := sync.WaitGroup{}
		wg.Add(10)

		s := lo.Synchronize()

		for i := 0; i < 10; i++ {
			go s.Do(func() {
				time.Sleep(5 * time.Millisecond)
				wg.Done()
			})
		}

		wg.Wait()

		duration := time.Since(start)

		is.Greater(duration, 50*time.Millisecond)
		is.Less(duration, 60*time.Millisecond)
	}

	// check locker is locked
	{
		mu := &sync.Mutex{}
		s := lo.Synchronize(mu)

		s.Do(func() {
			// 测试不能重复加锁
			is.False(mu.TryLock())
		})
		is.True(mu.TryLock())

		Try0(func() {
			mu.Unlock()
		})
	}

	// check we don't accept multiple arguments
	{
		is.PanicsWithValue("unexpected arguments", func() {
			mu := &sync.Mutex{}
			lo.Synchronize(mu, mu, mu)
		})
	}
}

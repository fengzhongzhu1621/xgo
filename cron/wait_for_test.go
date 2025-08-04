package cron

import (
	"context"
	"testing"
	"time"

	"github.com/fengzhongzhu1621/xgo/tests"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// TestWaitFor 等待定时任务执行结束
func TestWaitFor(t *testing.T) {
	t.Parallel()

	testTimeout := 100 * time.Millisecond
	longTimeout := 2 * testTimeout
	shortTimeout := 4 * time.Millisecond

	t.Run("exist condition works", func(t *testing.T) {
		t.Parallel()

		tests.TestWithTimeout(t, testTimeout)
		is := assert.New(t)

		// 超时也会终止定时任务
		// i 表示执行次数
		sum := 0
		laterTrue := func(i int) bool {
			// 根据执行次数，决定是否结束定时任务，返回 true 表示结束定时任务
			sum += 1
			return i >= 5
		}
		// laterTrue 终止条件
		// longTimeout 超时时间
		// time.Microsecond 心跳间隔
		//
		// iter 执行次数
		// duration 执行耗时
		// ok 超时返回false
		iter, duration, ok := lo.WaitFor(laterTrue, longTimeout, time.Millisecond)

		// 执行次数
		is.Equal(6, iter, "unexpected iteration count")
		is.Equal(6, sum, "unexpected iteration count")
		// 执行耗时
		is.InEpsilon(6*time.Millisecond, duration, float64(500*time.Microsecond))
		// 因为符合终止条件，终止了定时任务
		is.True(ok)
	})

	t.Run("counter is incremented", func(t *testing.T) {
		t.Parallel()

		tests.TestWithTimeout(t, testTimeout)
		is := assert.New(t)

		counter := 0
		alwaysFalse := func(i int) bool {
			is.Equal(counter, i)
			counter++
			return false
		}

		// 超时退出
		iter, duration, ok := lo.WaitFor(alwaysFalse, shortTimeout, 1050*time.Microsecond)
		is.Equal(counter, iter, "unexpected iteration count")
		is.InEpsilon(10*time.Millisecond, duration, float64(500*time.Microsecond))
		is.False(ok)
	})

	alwaysTrue := func(_ int) bool { return true }
	alwaysFalse := func(_ int) bool { return false }

	t.Run("short timeout works", func(t *testing.T) {
		t.Parallel()

		tests.TestWithTimeout(t, testTimeout)
		is := assert.New(t)

		iter, duration, ok := lo.WaitFor(alwaysFalse, shortTimeout, 10*time.Millisecond)
		is.Equal(0, iter, "unexpected iteration count")
		is.InEpsilon(10*time.Millisecond, duration, float64(500*time.Microsecond))
		is.False(ok)
	})

	t.Run("timeout works", func(t *testing.T) {
		t.Parallel()

		tests.TestWithTimeout(t, testTimeout)
		is := assert.New(t)

		shortTimeout := 4 * time.Millisecond
		iter, duration, ok := lo.WaitFor(alwaysFalse, shortTimeout, 10*time.Millisecond)
		is.Equal(0, iter, "unexpected iteration count")
		is.InEpsilon(10*time.Millisecond, duration, float64(500*time.Microsecond))
		is.False(ok)
	})

	t.Run("exist on first condition", func(t *testing.T) {
		t.Parallel()

		tests.TestWithTimeout(t, testTimeout)
		is := assert.New(t)

		iter, duration, ok := lo.WaitFor(alwaysTrue, 10*time.Millisecond, time.Millisecond)
		is.Equal(1, iter, "unexpected iteration count")
		is.InEpsilon(time.Millisecond, duration, float64(5*time.Microsecond))
		is.True(ok)
	})
}

func TestWaitForWithContext(t *testing.T) {
	t.Parallel()

	testTimeout := 100 * time.Millisecond
	longTimeout := 2 * testTimeout
	shortTimeout := 4 * time.Millisecond

	t.Run("exist condition works", func(t *testing.T) {
		t.Parallel()

		tests.TestWithTimeout(t, testTimeout)
		is := assert.New(t)

		laterTrue := func(_ context.Context, i int) bool {
			return i >= 5
		}

		iter, duration, ok := lo.WaitForWithContext(
			context.Background(),
			laterTrue,
			longTimeout,
			time.Millisecond,
		)
		is.Equal(6, iter, "unexpected iteration count")
		is.InEpsilon(6*time.Millisecond, duration, float64(500*time.Microsecond))
		is.True(ok)
	})

	t.Run("counter is incremented", func(t *testing.T) {
		t.Parallel()

		tests.TestWithTimeout(t, testTimeout)
		is := assert.New(t)

		counter := 0
		alwaysFalse := func(_ context.Context, i int) bool {
			is.Equal(counter, i)
			counter++
			return false
		}

		iter, duration, ok := lo.WaitForWithContext(
			context.Background(),
			alwaysFalse,
			shortTimeout,
			1050*time.Microsecond,
		)
		is.Equal(counter, iter, "unexpected iteration count")
		is.InEpsilon(10*time.Millisecond, duration, float64(500*time.Microsecond))
		is.False(ok)
	})

	alwaysTrue := func(_ context.Context, _ int) bool { return true }
	alwaysFalse := func(_ context.Context, _ int) bool { return false }

	t.Run("short timeout works", func(t *testing.T) {
		t.Parallel()

		tests.TestWithTimeout(t, testTimeout)
		is := assert.New(t)

		iter, duration, ok := lo.WaitForWithContext(
			context.Background(),
			alwaysFalse,
			shortTimeout,
			10*time.Millisecond,
		)
		is.Equal(0, iter, "unexpected iteration count")
		is.InEpsilon(10*time.Millisecond, duration, float64(500*time.Microsecond))
		is.False(ok)
	})

	t.Run("timeout works", func(t *testing.T) {
		t.Parallel()

		tests.TestWithTimeout(t, testTimeout)
		is := assert.New(t)

		shortTimeout := 4 * time.Millisecond
		iter, duration, ok := lo.WaitForWithContext(
			context.Background(),
			alwaysFalse,
			shortTimeout,
			10*time.Millisecond,
		)
		is.Equal(0, iter, "unexpected iteration count")
		is.InEpsilon(10*time.Millisecond, duration, float64(500*time.Microsecond))
		is.False(ok)
	})

	t.Run("exist on first condition", func(t *testing.T) {
		t.Parallel()

		tests.TestWithTimeout(t, testTimeout)
		is := assert.New(t)

		iter, duration, ok := lo.WaitForWithContext(
			context.Background(),
			alwaysTrue,
			10*time.Millisecond,
			time.Millisecond,
		)
		is.Equal(1, iter, "unexpected iteration count")
		is.InEpsilon(time.Millisecond, duration, float64(5*time.Microsecond))
		is.True(ok)
	})

	t.Run("context cancellation stops everything", func(t *testing.T) {
		t.Parallel()

		tests.TestWithTimeout(t, testTimeout)
		is := assert.New(t)

		expiringCtx, clean := context.WithTimeout(context.Background(), 8*time.Millisecond)
		t.Cleanup(func() {
			clean()
		})

		iter, duration, ok := lo.WaitForWithContext(
			expiringCtx,
			alwaysFalse,
			100*time.Millisecond,
			3*time.Millisecond,
		)
		is.Equal(2, iter, "unexpected iteration count")
		is.InEpsilon(10*time.Millisecond, duration, float64(500*time.Microsecond))
		is.False(ok)
	})

	t.Run("canceled context stops everything", func(t *testing.T) {
		t.Parallel()

		tests.TestWithTimeout(t, testTimeout)
		is := assert.New(t)

		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		iter, duration, ok := lo.WaitForWithContext(
			canceledCtx,
			alwaysFalse,
			100*time.Millisecond,
			1050*time.Microsecond,
		)
		is.Equal(0, iter, "unexpected iteration count")
		is.InEpsilon(1*time.Millisecond, duration, float64(5*time.Microsecond))
		is.False(ok)
	})
}

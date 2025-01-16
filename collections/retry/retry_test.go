package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/samber/lo"

	"github.com/duke-git/lancet/v2/retry"
	"github.com/stretchr/testify/assert"
)

func TestRetryContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())

	number := 0
	increaseNumber := func() error {
		number++
		if number > 3 {
			// 成功则取消重试
			cancel()
		}
		// 失败才会重试
		return errors.New("error occurs")
	}

	retry.Retry(increaseNumber,
		// 重试间隔
		retry.RetryWithLinearBackoff(time.Microsecond*50),
		retry.Context(ctx),
	)

	assert.Equal(t, 4, number)
}

func TestRetryWithLinearBackoff(t *testing.T) {
	number := 0
	increaseNumber := func() error {
		number++
		if number == 3 {
			// 成功取消重试
			return nil
		}
		// 失败才会重试
		return errors.New("error occurs")
	}

	retry.Retry(increaseNumber,
		retry.RetryWithLinearBackoff(time.Microsecond*50),
	)

	assert.Equal(t, 3, number)
}

type ExampleCustomBackoffStrategy struct {
	interval time.Duration
}

func (c *ExampleCustomBackoffStrategy) CalculateInterval() time.Duration {
	return c.interval + 1
}

func TestRetryWithCustomBackoff(t *testing.T) {
	number := 0
	increaseNumber := func() error {
		number++
		if number == 3 {
			// 成功取消重试
			return nil
		}
		// 失败才会重试
		return errors.New("error occurs")
	}

	err := retry.Retry(increaseNumber, retry.RetryWithCustomBackoff(&ExampleCustomBackoffStrategy{interval: time.Microsecond * 50}))
	if err != nil {
		return
	}

	assert.Equal(t, 3, number)
}

func TestRetryWithExponentialWithJitterBackoff(t *testing.T) {
	number := 0
	increaseNumber := func() error {
		number++
		if number == 3 {
			return nil
		}
		return errors.New("error occurs")
	}

	err := retry.Retry(increaseNumber, retry.RetryWithExponentialWithJitterBackoff(time.Microsecond*50, 2, time.Microsecond*25))
	if err != nil {
		return
	}

	assert.Equal(t, 3, number)
}

func TestRetryTimes(t *testing.T) {
	number := 0
	increaseNumber := func() error {
		number++
		if number == 3 {
			// 成功取消重试
			return nil
		}
		// 失败才会重试
		return errors.New("error occurs")
	}

	// 最多重试两次
	err := retry.Retry(increaseNumber, retry.RetryTimes(2))
	if err != nil {
		assert.Equal(t, 2, number)
		assert.EqualError(t, err, "function retry.TestRetryTimes.func1 run failed after 2 times retry")
	}
}

// 调用一个函数N次，直到它返回有效输出。返回捕获到的错误或者nil（空值）。当第一个参数小于1时，该函数会一直运行直到返回成功的响应。
func TestAttempt(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	err := fmt.Errorf("failed")

	// 调用 1 次
	iter1, err1 := lo.Attempt(42, func(i int) error {
		return nil
	})
	is.Equal(iter1, 1)
	is.Equal(err1, nil)

	// 调用 6 次
	iter2, err2 := lo.Attempt(42, func(i int) error {
		if i == 5 {
			return nil
		}

		return err
	})
	is.Equal(iter2, 6)
	is.Equal(err2, nil)

	// 调用 2 次
	iter3, err3 := lo.Attempt(2, func(i int) error {
		if i == 5 {
			return nil
		}

		return err
	})
	is.Equal(iter3, 2)
	is.Equal(err3, err)

	// 调用 43 次（循环一直到返回成功）
	iter4, err4 := lo.Attempt(0, func(i int) error {
		if i < 42 {
			return err
		}

		return nil
	})
	is.Equal(iter4, 43)
	is.Equal(err4, nil)
}

// 设置重试的固定间隔
func TestAttemptWithDelay(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	err := fmt.Errorf("failed")

	iter1, dur1, err1 := lo.AttemptWithDelay(42, 10*time.Millisecond, func(i int, d time.Duration) error {
		return nil
	})
	iter2, dur2, err2 := lo.AttemptWithDelay(42, 10*time.Millisecond, func(i int, d time.Duration) error {
		if i == 5 {
			return nil
		}

		return err
	})
	iter3, dur3, err3 := lo.AttemptWithDelay(2, 10*time.Millisecond, func(i int, d time.Duration) error {
		if i == 5 {
			return nil
		}

		return err
	})
	iter4, dur4, err4 := lo.AttemptWithDelay(0, 10*time.Millisecond, func(i int, d time.Duration) error {
		if i < 10 {
			return err
		}

		return nil
	})

	is.Equal(iter1, 1)
	is.GreaterOrEqual(dur1, 0*time.Millisecond)
	is.Less(dur1, 1*time.Millisecond)
	is.Equal(err1, nil)
	is.Equal(iter2, 6)
	is.Greater(dur2, 50*time.Millisecond)
	is.Less(dur2, 60*time.Millisecond)
	is.Equal(err2, nil)
	is.Equal(iter3, 2)
	is.Greater(dur3, 10*time.Millisecond)
	is.Less(dur3, 20*time.Millisecond)
	is.Equal(err3, err)
	is.Equal(iter4, 11)
	is.Greater(dur4, 100*time.Millisecond)
	is.Less(dur4, 115*time.Millisecond)
	is.Equal(err4, nil)
}

// 返回任务返回的标识决定是否终止重试
func TestAttemptWhile(t *testing.T) {
	is := assert.New(t)

	err := fmt.Errorf("failed")

	// 执行一次
	iter1, err1 := lo.AttemptWhile(42, func(i int) (error, bool) {
		return nil, true
	})
	is.Equal(iter1, 1)
	is.Nil(err1)

	iter2, err2 := lo.AttemptWhile(42, func(i int) (error, bool) {
		if i == 5 {
			return nil, true
		}

		return err, true
	})
	is.Equal(iter2, 6)
	is.Nil(err2)

	iter3, err3 := lo.AttemptWhile(2, func(i int) (error, bool) {
		if i == 5 {
			return nil, true
		}

		return err, true
	})
	is.Equal(iter3, 2)
	is.Equal(err3, err)

	iter4, err4 := lo.AttemptWhile(0, func(i int) (error, bool) {
		if i < 42 {
			return err, true
		}

		return nil, true
	})
	is.Equal(iter4, 43)
	is.Nil(err4)

	iter5, err5 := lo.AttemptWhile(0, func(i int) (error, bool) {
		if i == 5 {
			return nil, false
		}

		return err, true
	})
	is.Equal(iter5, 6)
	is.Nil(err5)

	// 执行一次，因为任务返回 false，所以不会执行
	iter6, err6 := lo.AttemptWhile(0, func(i int) (error, bool) {
		return nil, false
	})
	is.Equal(iter6, 1)
	is.Nil(err6)

	// 执行第 43 次时终止执行
	iter7, err7 := lo.AttemptWhile(42, func(i int) (error, bool) {
		if i == 42 {
			// 执行第 43 次时终止执行
			return nil, false
		}
		if i < 41 {
			// 执行 42次
			return err, true
		}

		return nil, true
	})

	is.Equal(iter7, 42)
	is.Nil(err7)
}

func TestAttemptWhileWithDelay(t *testing.T) {
	is := assert.New(t)

	err := fmt.Errorf("failed")

	iter1, dur1, err1 := lo.AttemptWhileWithDelay(42, 10*time.Millisecond, func(i int, d time.Duration) (error, bool) {
		return nil, true
	})

	is.Equal(iter1, 1)
	is.GreaterOrEqual(dur1, 0*time.Millisecond)
	is.Less(dur1, 1*time.Millisecond)
	is.Nil(err1)

	iter2, dur2, err2 := lo.AttemptWhileWithDelay(42, 10*time.Millisecond, func(i int, d time.Duration) (error, bool) {
		if i == 5 {
			return nil, true
		}

		return err, true
	})

	is.Equal(iter2, 6)
	is.Greater(dur2, 50*time.Millisecond)
	is.Less(dur2, 60*time.Millisecond)
	is.Nil(err2)

	iter3, dur3, err3 := lo.AttemptWhileWithDelay(2, 10*time.Millisecond, func(i int, d time.Duration) (error, bool) {
		if i == 5 {
			return nil, true
		}

		return err, true
	})

	is.Equal(iter3, 2)
	is.Greater(dur3, 10*time.Millisecond)
	is.Less(dur3, 20*time.Millisecond)
	is.Equal(err3, err)

	iter4, dur4, err4 := lo.AttemptWhileWithDelay(0, 10*time.Millisecond, func(i int, d time.Duration) (error, bool) {
		if i < 10 {
			return err, true
		}

		return nil, true
	})

	is.Equal(iter4, 11)
	is.Greater(dur4, 100*time.Millisecond)
	is.Less(dur4, 115*time.Millisecond)
	is.Nil(err4)

	iter5, dur5, err5 := lo.AttemptWhileWithDelay(0, 10*time.Millisecond, func(i int, d time.Duration) (error, bool) {
		if i == 5 {
			return nil, false
		}

		return err, true
	})

	is.Equal(iter5, 6)
	is.Greater(dur5, 10*time.Millisecond)
	is.Less(dur5, 115*time.Millisecond)
	is.Nil(err5)

	iter6, dur6, err6 := lo.AttemptWhileWithDelay(0, 10*time.Millisecond, func(i int, d time.Duration) (error, bool) {
		return nil, false
	})

	is.Equal(iter6, 1)
	is.Less(dur6, 10*time.Millisecond)
	is.Less(dur6, 115*time.Millisecond)
	is.Nil(err6)

	iter7, dur7, err7 := lo.AttemptWhileWithDelay(42, 10*time.Millisecond, func(i int, d time.Duration) (error, bool) {
		if i == 42 {
			return nil, false
		}
		if i < 41 {
			return err, true
		}

		return nil, true
	})

	is.Equal(iter7, 42)
	is.Less(dur7, 500*time.Millisecond)
	is.Nil(err7)
}

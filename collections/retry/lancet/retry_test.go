package lancet

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/retry"
	"github.com/stretchr/testify/assert"
)

// 测试根据上下文取消重试
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

	// 执行任务
	retry.Retry(increaseNumber,
		// 重试间隔（50 毫秒）
		retry.RetryWithLinearBackoff(time.Microsecond*50),
		// 重试上下文
		retry.Context(ctx),
	)

	assert.Equal(t, 4, number)
}

// 测试线性重试
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

// 测试自定义重试间隔
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

	// 使用自定义重试间隔启动任务
	err := retry.Retry(increaseNumber, retry.RetryWithCustomBackoff(&ExampleCustomBackoffStrategy{
		interval: time.Microsecond * 50,
	}))
	if err != nil {
		return
	}

	assert.Equal(t, 3, number)
}

// 测试指数退避重试
func TestRetryWithExponentialWithJitterBackoff(t *testing.T) {
	number := 0
	increaseNumber := func() error {
		number++
		if number == 3 {
			return nil
		}
		return errors.New("error occurs")
	}

	err := retry.Retry(increaseNumber, retry.RetryWithExponentialWithJitterBackoff(
		time.Microsecond*50, // 初始重试间隔
		2,                   // 重试次数
		time.Microsecond*25, // 最大重试间隔
	))
	if err != nil {
		return
	}

	assert.Equal(t, 3, number)
}

// 测试重试次数
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
		assert.EqualError(
			t,
			err,
			"function retry.TestRetryTimes.func1 run failed after 2 times retry",
		)
	}
}

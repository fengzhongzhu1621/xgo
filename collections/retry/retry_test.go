package retry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/retry"
	"github.com/stretchr/testify/assert"
)

func TestRetryContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())

	number := 0
	increaseNumber := func() error {
		number++
		if number > 3 {
			// 取消重试
			cancel()
		}
		return errors.New("error occurs")
	}

	retry.Retry(increaseNumber,
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
			return nil
		}
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
			return nil
		}
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
			return nil
		}
		return errors.New("error occurs")
	}

	err := retry.Retry(increaseNumber, retry.RetryTimes(2))
	if err != nil {
		assert.Equal(t, 2, number)
		assert.EqualError(t, err, "function retry.TestRetryTimes.func1 run failed after 2 times retry")
	}

}

package backoff

import (
	"fmt"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// 模拟了一个可能失败的操作（unreliableOperation），并使用指数退避策略（exponential backoff）来自动重试这个操作，直到成功或达到最大重试时间限制。

var failureCount = 0

// 让它先失败3次，再成功
func unreliableOperation() error {
	failureCount++
	fmt.Println("failureCount", failureCount)
	if failureCount < 4 {
		return fmt.Errorf("transient error")
	}
	return nil
}

func retryOperationWithBackoff() error {
	// // 设置指数退避策略
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 30 * time.Second // 最大重试时间限制（超过30秒就停止重试）。
	bo.MaxInterval = 10 * time.Second    // 最大重试间隔（每次重试的等待时间不会超过10秒）。

	// 定义重试逻辑
	return backoff.Retry(func() error {
		err := unreliableOperation()
		if err != nil {
			return err // 如果操作失败，返回错误并重试
		}
		return nil // 操作成功
	}, bo)
}

func TestRetryOperationWithBackoff(t *testing.T) {
	err := retryOperationWithBackoff()
	if err != nil {
		t.Fatalf("Operation failed after retries: %v", err)
	}

	// failureCount 1
	// failureCount 2
	// failureCount 3
	// failureCount 4
}

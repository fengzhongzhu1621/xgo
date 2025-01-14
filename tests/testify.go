package tests

import (
	"testing"
	"time"
)

// TestWithTimeout https://github.com/stretchr/testify/issues/1101
// timeout 表示测试应该在这个时间段内完成，否则将被视为超时。
func TestWithTimeout(t *testing.T, timeout time.Duration) {
	t.Helper()

	// 用于在测试完成时发送信号
	testFinished := make(chan struct{})
	// 注册清理函数
	t.Cleanup(func() { close(testFinished) })

	go func() { //nolint:staticcheck
		select {
		case <-testFinished:
		case <-time.After(timeout):
			t.Errorf("test timed out after %s", timeout)
			// 立即停止当前测试，并标记为失败。
			t.FailNow() //nolint:govet,staticcheck
		}
	}()
}

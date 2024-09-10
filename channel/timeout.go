package channel

import (
	"errors"
	"time"
)

// TimeoutCaller 调用函数并设置超时时间, 在指定的超时时间内等待函数调用的结果，如果超时则返回一个错误。
func TimeoutCaller(functionCall func(chan error), timeout time.Duration) error {
	var (
		err error
	)
	// 执行任务
	echan := make(chan error)
	go functionCall(echan)

	// 等待任务执行完成
	select {
	case <-time.After(timeout):
		return errors.New("Timed out while initiating calling")
	case err = <-echan:
		return err
	}
}

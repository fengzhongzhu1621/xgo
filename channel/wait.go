package channel

// GoWait is a basic promise implementation: it wraps calls a function in a goroutine
// and returns a channel which will later return the function's return value.
func GoWait(f func() error) error {
	ch := make(chan error)
	go func() {
		ch <- f()
	}()

	// 等待子协程执行完成
	return <-ch
}

package channel

// IsClosed 判断是否收到关闭事件
// 当发送者 close(done)时，会发送一个零值，标记已关闭.
func IsClosed(done chan struct{}) bool {
	// select 语句用于在多个通道操作中进行选择
	select {
	case <-done:
		return true
	default:
		// 当没有其他分支可以执行时，将执行默认情况
		return false
	}
}

// IsChannelClosed returns true if provided `chan struct{}` is closed.
// IsChannelClosed panics if message is sent to this channel.
func IsChannelClosed(channel chan struct{}) bool {
	select {
	case _, ok := <-channel:
		if ok {
			// 如果管道从计划关闭的管道中收到一个消息
			panic("received unexpected message")
		}
		// 管道被关闭
		return true
	default:
		// 管道中没有消息
		return false
	}
}

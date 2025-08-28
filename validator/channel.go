package validator

// IsClosed 判断是否收到关闭事件(非阻塞)
// 当发送者 close(done)时，会发送一个零值，标记已关闭.
func IsClosed(done chan struct{}) bool {
	// select 语句用于在多个通道操作中进行选择
	select {
	case <-done: // 从管道接收到数据 或 接收到关闭信号，则返回 true
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

// ChannelIsNotFull 判断一个通道是否未满
// ch <-chan T 接受一个只读通道
// 未满 返回 true；已满返回 false
func ChannelIsNotFull[T any](ch <-chan T) bool {
	// cap(ch) == 0: 判断通道是否为无缓冲通道。无缓冲通道在发送和接收操作完成之前会阻塞，因此从某种意义上说，它们永远不会“满”。
	// len(ch) < cap(ch): 对于有缓冲通道，检查当前长度是否小于容量，以确定通道是否未满。
	return cap(ch) == 0 || len(ch) < cap(ch)
}

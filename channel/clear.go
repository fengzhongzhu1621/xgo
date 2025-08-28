package channel

func ClearChannel[T any](ch chan T) {
	for {
		select {
		case <-ch:
			// 读取一个值并丢弃
		default:
			// channel 已空，退出循环
			return
		}
	}
}

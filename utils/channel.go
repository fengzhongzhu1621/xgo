package utils

// IsClosed 判断是否收到关闭事件
// 当发送者 close(done)时，会发送一个零值，标记已关闭
func IsClosed(done chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

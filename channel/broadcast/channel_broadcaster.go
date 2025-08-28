package broadcast

type ChannelBroadcaster struct {
	signal chan struct{} // 定义广播信号
}

func NewChannelBroadcaster() *ChannelBroadcaster {
	return &ChannelBroadcaster{
		signal: make(chan struct{}),
	}
}

// Go 等待接受广播
func (b *ChannelBroadcaster) Go(fn func()) {
	go func() {
		<-b.signal
		fn()
	}()
}

// Broadcast 发送广播
func (b *ChannelBroadcaster) Broadcast() {
	close(b.signal)
}

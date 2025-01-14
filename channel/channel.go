package channel

func CreateChannels[T any](count int, channelBufferCap int) []chan T {
	children := make([]chan T, 0, count)

	for i := 0; i < count; i++ {
		children = append(children, make(chan T, channelBufferCap))
	}

	return children
}

func ChannelsToReadOnly[T any](children []chan T) []<-chan T {
	roChildren := make([]<-chan T, 0, len(children))

	for i := range children {
		roChildren = append(roChildren, children[i])
	}

	return roChildren
}

func CloseChannels[T any](children []chan T) {
	for i := 0; i < len(children); i++ {
		close(children[i])
	}
}

func ChannelIsNotFull[T any](ch <-chan T) bool {
	return cap(ch) == 0 || len(ch) < cap(ch)
}

package channel

func CloseChannels[T any](children []chan T) {
	for i := 0; i < len(children); i++ {
		close(children[i])
	}
}

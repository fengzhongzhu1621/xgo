package broadcast

import (
	"context"
)

type ContextBroadcaster struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewContextBroadcaster() *ContextBroadcaster {
	ctx, cancel := context.WithCancel(context.Background())
	return &ContextBroadcaster{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (b *ContextBroadcaster) Go(fn func()) {
	go func() {
		<-b.ctx.Done()
		fn()
	}()
}

func (b *ContextBroadcaster) Broadcast() {
	b.cancel()
}

package channel

import (
	"sync"
	"testing"
)

func TestNewChannelBroadcaster(t *testing.T) {
	b := NewChannelBroadcaster()

	var wg sync.WaitGroup
	wg.Add(2)

	b.Go(func() {
		t.Log("receiver 1 finished")
		wg.Done()
	})
	b.Go(func() {
		t.Log("receiver 2 finished")
		wg.Done()
	})

	b.Broadcast()

	wg.Wait()
}

package broadcast

import (
	"sync"
	"testing"
)

func TestNewContextBroadcaster(t *testing.T) {
	b := NewContextBroadcaster()

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

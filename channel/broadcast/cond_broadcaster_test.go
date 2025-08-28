package broadcast

import (
	"sync"
	"testing"
)

func TestNewCondBroadcaster(t *testing.T) {
	b := NewCondBroadcaster()

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

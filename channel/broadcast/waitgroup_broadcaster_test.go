package broadcast

import (
	"sync"
	"testing"
	"time"
)

func TestNewWaitGroupBroadcaster(t *testing.T) {
	b := NewWaitGroupBroadcaster()

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

	time.Sleep(100 * time.Millisecond)
	b.Broadcast()

	wg.Wait()
}

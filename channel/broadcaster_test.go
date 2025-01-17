package channel

import (
	"sync"
	"testing"
	"time"
)

func TestNewChannelBroadcaster(t *testing.T) {
	b := NewChannelBroadcaster()

	var wg sync.WaitGroup
	wg.Add(2)

	// 创建两个协程，等待接受广播信号
	b.Go(func() {
		t.Log("receiver 1 finished")
		wg.Done()
	})
	b.Go(func() {
		t.Log("receiver 2 finished")
		wg.Done()
	})

	// 给所有的子任务发送开始执行信号
	b.Broadcast()

	// 等待所有子任务执行完成
	wg.Wait()
}

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

func TestNewRWMutexBroadcaster(t *testing.T) {
	b := NewRWMutexBroadcaster()

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

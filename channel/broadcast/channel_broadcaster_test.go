package broadcast

import (
	"sync"
	"testing"
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

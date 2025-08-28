// Package queue implements infinite queue, supporting blocking data acquisition.

package listqueue

import (
	"container/list"
	"sync"
)

// Queue uses list and channel to achieve blocking acquisition and infinite queue.
// 基于链表实现的无限队列，支持阻塞获取数据
type Queue[T any] struct {
	list    *list.List      // 链表存储队列元素
	notify  chan struct{}   // 通知通道，用于唤醒阻塞的Get操作
	mu      sync.Mutex      // 互斥锁，保护队列操作
	waiting bool            // 标记是否有Get操作在等待数据
	done    <-chan struct{} // 外部取消通道，用于中断阻塞
}

// New initializes a queue, dones is used to notify Queue.Get() from blocking.
// 创建新的队列实例，done通道用于中断阻塞的Get操作
func New[T any](done <-chan struct{}) *Queue[T] {
	q := &Queue[T]{
		list:   list.New(),             // 初始化空链表
		notify: make(chan struct{}, 1), // 创建缓冲大小为1的通知通道
		done:   done,                   // 设置取消通道
	}
	return q
}

// Put puts an element into the queue.
// Put and Get can be concurrent, multiple Put can be concurrent.
// 向队列中添加元素，支持并发Put操作
func (q *Queue[T]) Put(v T) {
	var wakeUp bool

	q.mu.Lock()

	if q.waiting {
		wakeUp = true     // 标记需要唤醒等待的Get操作
		q.waiting = false // 重置等待标志
	}
	q.list.PushBack(v) // 将元素添加到链表尾部

	q.mu.Unlock()

	if wakeUp {
		select {
		case q.notify <- struct{}{}: // 发送通知唤醒阻塞的Get
		default: // 如果通知通道已满，跳过发送
		}
	}
}

// Get gets an element from the queue, blocking if there is no content.
// Put and Get can be concurrent, but not concurrent Get.
// If done channel notify it from blocking, it will return false.
// 从队列中获取元素，如果队列为空则阻塞等待
func (q *Queue[T]) Get() (T, bool) {
	for {
		q.mu.Lock()
		// 有数据，直接获取链表头部元素
		if e := q.list.Front(); e != nil {
			q.list.Remove(e) // 移除链表头部元素
			q.mu.Unlock()
			return e.Value.(T), true // 返回元素值和成功标志
		}

		// 无数据，设置等待标志，等待Put操作唤醒
		q.waiting = true // 设置等待标志，通知Put操作需要唤醒
		q.mu.Unlock()

		select {
		case <-q.notify: // 收到Put操作的通知
			continue // 继续循环尝试获取元素
		case <-q.done: // 收到取消信号
			var zero T
			return zero, false // 返回零值和失败标志
		}
	}
}

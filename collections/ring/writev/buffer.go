package writev

// 这是一个基于环形队列的批量写入缓冲区实现，主要功能：

// 核心特性：
// 使用环形队列缓存待发送的消息
// 支持批量写入（writev模式）提高IO效率
// 通过goroutine异步处理发送任务
// 提供优雅的退出机制

// 关键组件：
// Buffer - 主结构体，管理发送队列和goroutine
// 环形队列 - 使用 ring.Ring[[]byte] 缓存数据
// 唤醒通道 - 用于通知发送goroutine有新数据
// 错误处理 - 支持外部中断和错误传递

import (
	"errors"
	"io"
	"net"
	"runtime"

	"github.com/fengzhongzhu1621/xgo/collections/ring/ring"
)

const (
	// default buffer queue length.
	defaultBufferSize = 128
	// The maximum number of data packets that can be sent by writev (from Go source code definition).
	maxWritevBuffers = 1024
)

var (
	// ErrAskQuit sends a close request externally.
	// ErrAskQuit 外部发送关闭请求的错误
	ErrAskQuit = errors.New("writev goroutine is asked to quit")
	// ErrStopped Buffer stops receiving data.
	// ErrStopped 缓冲区停止接收数据的错误
	ErrStopped = errors.New("writev buffer stop to receive data")
)

// Buffer records the messages to be sent and sends them in batches using goroutines.
// Buffer 记录待发送的消息并使用goroutine批量发送
type Buffer struct {
	opts           *Options           // configuration items. 配置项
	w              io.Writer          // The underlying io.Writer that sends data. 底层数据写入器
	queue          *ring.Ring[[]byte] // queue for buffered messages. 缓冲消息的环形队列
	wakeupCh       chan struct{}      // used to wake up the sending goroutine. 唤醒发送goroutine的通道
	done           chan struct{}      // notify the sending goroutine to exit. 通知发送goroutine退出的通道
	err            error              // record error message. 错误信息记录
	errCh          chan error         // internal error notification. 内部错误通知通道
	isQueueStopped bool               // whether the cache queue stops receiving packets. 缓存队列是否停止接收数据包
}

var defaultQuitHandler = func(b *Buffer) {
	b.SetQueueStopped(true) // 默认退出处理函数：设置队列停止接收
}

// NewBuffer 创建环形缓存buffer
func NewBuffer(opt ...Option) *Buffer {
	// 自定义配置
	opts := &Options{
		bufferSize: defaultBufferSize,  // 默认缓冲区大小
		handler:    defaultQuitHandler, // 默认退出处理函数
	}
	for _, o := range opt {
		o(opts)
	}

	b := &Buffer{
		opts:           opts,                                      // 配置选项
		queue:          ring.New[[]byte](uint32(opts.bufferSize)), // 环形队列（需要uint32参数）
		wakeupCh:       make(chan struct{}, 1),                    // 唤醒通道（带缓冲）
		done:           make(chan struct{}),                       // 完成/退出信号通道
		errCh:          make(chan error, 1),                       // 错误通道
		isQueueStopped: false,                                     // 队列停止标志
	}
	return b
}

// SetQueue极pped sets whether the buffer queue stops receiving packets.
// 设置队列是否停止接收数据包
func (b *Buffer) SetQueueStopped(stopped bool) {
	b.isQueueStopped = stopped
	if b.err == nil {
		b.err = ErrStopped
	}
}

// Error returns the reason why the sending goroutine exited.
func (b *Buffer) Error() error {
	return b.err
}

// Done return to exit the Channel.
func (b *Buffer) Done() chan struct{} {
	return b.done
}

func (b *Buffer) wakeUp() {
	// Based on performance optimization considerations: due to the poor
	// efficiency of concurrent select write channel, check here first
	// whether wakeupCh already has a wakeup message, reducing the chance
	// of concurrently writing to the channel.
	if len(b.wakeupCh) > 0 {
		// 已经被唤醒
		return
	}
	// try to send wakeup signal, don't wait.
	select {
	case b.wakeupCh <- struct{}{}:
	default:
	}
}

// Write writes p to the buffer queue and returns the length of the data written.
// How to write a packet smaller than len(p), err returns the specific reason.
func (b *Buffer) Write(p []byte) (int, error) {
	if b.opts.dropFull {
		// 队列满时，写操作丢弃数据
		return b.writeNoWait(p)
	}
	return b.writeOrWait(p)
}

// WriteNoWait writes p to the buffer queue and returns the length of the data written.
// How to write a packet smaller than len(p), err returns the specific reason.
// 队列满时，写操作丢弃数据
func (b *Buffer) writeNoWait(p []byte) (int, error) {
	// The buffer queue stops receiving packets and returns directly.
	if b.isQueueStopped {
		return 0, b.err
	}
	// return directly when the queue is full.
	if err := b.queue.Put(p); err != nil {
		return 0, err
	}
	// Write the buffer queue successfully, wake up the sending goroutine.
	b.wakeUp()
	return len(p), nil
}

// WriteOrWait writes p to the buffer queue and returns the length of the data written.
func (b *Buffer) writeOrWait(p []byte) (int, error) {
	for {
		// The buffer queue stops receiving packets and returns directly.
		if b.isQueueStopped {
			return 0, b.err
		}
		// Write the buffer queue successfully, wake up the sending goroutine.
		if err := b.queue.Put(p); err == nil {
			b.wakeUp()
			return len(p), nil
		}

		// 处理队列满的场景，先从队列中获取1k的数据，写入到b.w中；在下一个循环中，再尝试写入p
		// The queue is full, send the package directly.
		if err := b.writeDirectly(); err != nil {
			return 0, err
		}
	}
}

// Writev writes p to the buffer queue and returns the length of the data written.
// writeDirectly 直接写入数据到底层writer
// 使用writev批量发送机制提高IO效率
func (b *Buffer) writeDirectly() error {
	if b.queue.IsEmpty() { // 队列为空，无需写入
		return nil
	}
	vals := make([][]byte, 0, maxWritevBuffers) // 准备批量数据缓冲区
	size, _ := b.queue.Gets(&vals)              // 从队列获取批量数据
	if size == 0 {                              // 没有获取到数据
		return nil
	}
	bufs := make(net.Buffers, 0, maxWritevBuffers) // 创建net.Buffers用于writev
	for _, v := range vals {                       // 将数据转换为Buffers格式
		bufs = append(bufs, v)
	}
	if _, err := bufs.WriteTo(b.w); err != nil {
		// Notify the sending goroutine setting error and exit.
		select {
		case b.errCh <- err:
		default:
		}
		return err
	}
	return nil
}

// Restart recreates a Buffer when restarting, reusing the buffer queue and configuration of the original Buffer.
func (b *Buffer) Restart(writer io.Writer, done chan struct{}) *Buffer {
	buffer := &Buffer{
		queue:    b.queue,
		opts:     b.opts,
		wakeupCh: make(chan struct{}, 1),
		errCh:    make(chan error, 1),
	}
	buffer.Start(writer, done)
	return buffer
}

// Start starts the sending goroutine, you need to set writer and done at startup.
// 启动发送goroutine，需要在启动时设置writer和done
func (b *Buffer) Start(writer io.Writer, done chan struct{}) {
	b.w = writer
	b.done = done
	go b.start()
}

func (b *Buffer) start() {
	initBufs := make(net.Buffers, 0, maxWritevBuffers)
	vals := make([][]byte, 0, maxWritevBuffers)
	bufs := initBufs

	// 设置goroutine退出时的清理函数
	defer b.opts.handler(b)

	for {
		// 从环形缓冲区读取数据，阻塞等待唤醒信号，如果环形缓存区有数据，则读取后立即返回
		if err := b.getOrWait(&vals); err != nil {
			b.err = err
			break
		}

		for _, v := range vals {
			bufs = append(bufs, v)
		}
		// 重置切片底层buf，防止append时重新分配内存
		vals = vals[:0]

		// 写入 io.Writer
		if _, err := bufs.WriteTo(b.w); err != nil {
			b.err = err
			break
		}

		// Reset bufs to the initial position to prevent `append` from generating new memory allocations.
		bufs = initBufs
	}
}

// getOrWait 获取数据或等待新数据到达
// 需要处理goroutine退出和错误通知
func (b *Buffer) getOrWait(values *[][]byte) error {
	for {
		// 检查是否收到退出或错误通知
		select {
		case <-b.done: // 收到退出信号
			return ErrAskQuit
		case err := <-b.errCh: // 收到错误通知
			return err
		default: // 没有通知，继续执行
		}
		// 从缓存队列批量接收数据包
		size, _ := b.queue.Gets(values) // 批量获取数据
		if size > 0 {                   // 成功获取到数据
			return nil
		}

		// 快速路径：由于使用select唤醒goroutine性能较差，
		// 这里优先使用Gosched()延迟检查队列，提高命中率和
		// 批量获取数据包的效率，从而降低使用select唤醒goroutine的概率
		runtime.Gosched() // 主动让出CPU时间片
		if !b.queue.IsEmpty() {
			// 队列仍有数据，则在下一个循环立即读取
			continue
		}

		// 队列为空，则阻塞等待唤醒信号
		// 慢速路径：延迟检查队列后仍然没有数据包，
		// 表明系统相对空闲。goroutine使用select机制
		// 等待唤醒。休眠的优势是在系统空闲时减少CPU空转损耗
		select {
		case <-b.done: // 收到退出信号
			return ErrAskQuit
		case err := <-b.errCh: // 收到错误通知
			return err
		case <-b.wakeupCh: // 收到唤醒信号
		}
	}
}

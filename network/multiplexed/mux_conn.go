package multiplexed

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	queue "github.com/fengzhongzhu1621/xgo/collections/queue/listqueue"
)

var _ IMuxConn = (*VirtualConnection)(nil)

// VirtualConnection 多路复用虚拟连接
// 在单个物理连接上复用多个虚拟连接
type VirtualConnection struct {
	id        uint32               // 虚拟连接ID
	conn      *Connection          // 所属的物理连接
	recvQueue *queue.Queue[[]byte] // 接收数据队列

	ctx        context.Context    // 连接上下文
	cancelFunc context.CancelFunc // 取消函数
	closed     uint32             // 连接关闭状态（原子操作）

	err error        // 连接错误
	mu  sync.RWMutex // 读写锁，保护错误状态
}

// RemoteAddr 获取连接的远程地址
// 返回:
//
//	net.Addr: 远程地址，对于数据报连接返回固定地址，对于流式连接返回实际远程地址
func (vc *VirtualConnection) RemoteAddr() net.Addr {
	if !vc.conn.isStream {
		return vc.conn.addr
	}
	if vc.conn == nil {
		return nil
	}
	conn := vc.conn.getRawConn()
	if conn == nil {
		return nil
	}
	return conn.RemoteAddr()
}

// LocalAddr 获取连接的本地地址
// 返回:
//
//	net.Addr: 本地地址，如果连接不存在则返回nil
func (vc *VirtualConnection) LocalAddr() net.Addr {
	if vc.conn == nil {
		return nil
	}
	conn := vc.conn.getRawConn()
	if conn == nil {
		return nil
	}
	return conn.LocalAddr()
}

// Write 写入请求数据包
// Write和Read可以并发执行，多个Write也可以并发执行
// 参数:
//
//	b: 要写入的数据字节切片
//
// 返回:
//
//	error: 写入过程中发生的错误
func (vc *VirtualConnection) Write(b []byte) error {
	if err := vc.loadErr(); err != nil {
		return err
	}
	select {
	case <-vc.ctx.Done():
		// 当上下文超时或取消时清理虚拟连接
		vc.Close()
		return vc.ctx.Err()
	default:
	}
	if err := vc.conn.send(b); err != nil {
		// 发送失败时清理虚拟连接
		vc.Close()
		return err
	}
	return nil
}

// Read 读取返回的数据包
// Write和Read可以并发执行，但Read不能并发执行
// 返回:
//
//	[]byte: 读取到的数据
//	error: 读取过程中发生的错误
func (vc *VirtualConnection) Read() ([]byte, error) {
	if err := vc.loadErr(); err != nil {
		return nil, err
	}
	rsp, ok := vc.recvQueue.Get()
	if !ok {
		vc.Close()
		if err := vc.loadErr(); err != nil {
			return nil, err
		}
		return nil, vc.ctx.Err()
	}
	return rsp, nil
}

// Close 关闭虚拟连接
// 使用原子操作确保只关闭一次，并从物理连接中移除该虚拟连接
func (vc *VirtualConnection) Close() {
	if atomic.CompareAndSwapUint32(&vc.closed, 0, 1) {
		vc.conn.remove(vc.id)
	}
}

// loadErr 加载连接错误状态
// 使用读锁保护错误状态的读取
// 返回:
//
//	error: 当前的连接错误
func (vc *VirtualConnection) loadErr() error {
	vc.mu.RLock()
	defer vc.mu.RUnlock()
	return vc.err
}

// storeErr 存储连接错误状态
// 如果已经有错误存在则不覆盖，使用写锁保护错误状态的写入
// 参数:
//
//	err: 要存储的错误
func (vc *VirtualConnection) storeErr(err error) {
	if vc.loadErr() != nil {
		return
	}
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.err = err
}

// cancel 取消虚拟连接
// 存储错误并调用取消函数
// 参数:
//
//	err: 取消的原因错误
func (vc *VirtualConnection) cancel(err error) {
	vc.storeErr(err)
	vc.cancelFunc()
}

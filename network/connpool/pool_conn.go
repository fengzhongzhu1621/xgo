package connpool

import (
	"errors"
	"net"
	"time"

	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/network/checker"
)

var globalBuffer []byte = make([]byte, 1)

var _ codec.IFramer = (*PoolConn)(nil)

// PoolConn is the connection in the connection pool.
// 连接池中的连接，实现 codec.IFramer 接口，用于读取来自网络的的二进制数据
type PoolConn struct {
	net.Conn
	fr         codec.IFramer // 用于读取来自网络的的二进制数据
	t          time.Time     // 成为空闲连接的时间
	created    time.Time     // 创建时间
	next, prev *PoolConn
	pool       *ConnectionPool // 连接属于哪个连接池
	closed     bool            // 标记连接已关闭
	forceClose bool            // 标记关闭连接时，放到连接池后是否强制关闭
	copyFrame  bool            // 帧读取器是否支持并发读取安全
	inPool     bool            // 标记是否被放到连接池中
}

// ReadFrame 从连接中读取帧数据
// 返回:
//
//	[]byte: 读取到的帧数据
//	error: 读取过程中发生的错误
func (pc *PoolConn) ReadFrame() ([]byte, error) {
	// 检查连接是否已关闭
	if pc.closed {
		return nil, ErrConnClosed
	}

	// 检查帧读取器是否已设置
	if pc.fr == nil {
		// 如果帧读取器未设置，将连接放回连接池并标记为强制关闭
		pc.pool.put(pc, true)
		return nil, errors.New("framer not set")
	}

	// 如果帧读取器已经设置，则使用帧读取器读取数据
	// 读取来自网络的的二进制数据
	frameData, err := pc.fr.ReadFrame()
	if err != nil {
		// ReadFrame 失败可能是由于 socket Read 接口超时失败
		// 或者解包失败，这两种情况都应该关闭连接
		// 将连接放回连接池并标记为强制关闭
		pc.pool.put(pc, true)
		return nil, err
	}

	// 如果帧读取器不支持并发读取安全，需要复制数据
	if pc.copyFrame {
		buf := make([]byte, len(frameData))
		copy(buf, frameData)
		return buf, err
	}

	return frameData, err
}

// isRemoteError 尝试接收一个字节来检测对端是否主动关闭了连接
// 如果对端返回 io.EOF 错误，表示对端已关闭
// 空闲连接不应该读取数据，如果读取到数据，说明上层的粘包处理未完成，连接也应该被丢弃
// 参数:
//
//	isFast: 是否使用快速检查模式
//
// 返回:
//
//	bool: 如果连接存在错误则返回 true
func (pc *PoolConn) isRemoteError(isFast bool) bool {
	var err error

	// 检查空闲连接是否存在连接错误
	if isFast {
		// 使用非阻塞方式快速检查连接状态
		err = checker.CheckConnErrUnblock(pc.Conn, globalBuffer)
	} else {
		// 使用阻塞方式检查连接状态
		err = checker.CheckConnErr(pc.Conn, globalBuffer)
	}

	if err != nil {
		return true
	}
	return false
}

// reset 重置连接状态，清除连接的超时设置
func (pc *PoolConn) reset() {
	if pc == nil {
		return
	}
	// 清除连接的超时设置，恢复为无限制等待
	pc.Conn.SetDeadline(time.Time{})
}

// Write 在连接上发送数据
// 参数:
//
//	b: 要发送的数据字节数组
//
// 返回:
//
//	int: 实际发送的字节数
//	error: 发送过程中发生的错误
func (pc *PoolConn) Write(b []byte) (int, error) {
	if pc.closed {
		return 0, ErrConnClosed
	}

	n, err := pc.Conn.Write(b)
	if err != nil {
		// 如果写入失败，将连接放回连接池并标记为强制关闭
		pc.pool.put(pc, true)
	}

	return n, err
}

// Read 从连接中读取数据
// 参数:
//
//	b: 用于存储读取数据的缓冲区
//
// 返回:
//
//	int: 实际读取的字节数
//	error: 读取过程中发生的错误
func (pc *PoolConn) Read(b []byte) (int, error) {
	if pc.closed {
		return 0, ErrConnClosed
	}

	n, err := pc.Conn.Read(b)
	if err != nil {
		// 如果读取失败，将连接放回连接池并标记为强制关闭
		pc.pool.put(pc, true)
	}

	return n, err
}

// Close 重写 net.Conn 的 Close 方法，将连接放回连接池
// 返回:
//
//	error: 关闭过程中发生的错误
func (pc *PoolConn) Close() error {
	if pc.closed {
		return ErrConnClosed
	}
	if pc.inPool {
		return ErrConnInPool
	}

	// 重置连接状态
	pc.reset()

	// 将连接放回连接池，根据 forceClose 标志决定是否强制关闭
	return pc.pool.put(pc, pc.forceClose)
}

// GetRawConn 获取 PoolConn 中的原始连接
// 返回:
//
//	net.Conn: 原始的网络连接
func (pc *PoolConn) GetRawConn() net.Conn {
	return pc.Conn
}

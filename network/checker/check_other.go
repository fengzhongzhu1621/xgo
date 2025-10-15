//go:build !aix && !darwin && !dragonfly && !freebsd && !netbsd && !openbsd && !solaris && !linux
// +build !aix,!darwin,!dragonfly,!freebsd,!netbsd,!openbsd,!solaris,!linux

package checker

import (
	"errors"
	"net"
	"time"
)

func CheckConnErr(conn net.Conn, buf []byte) error {
	// 设置极短超时：将连接的读超时设置为 1 毫秒 后。
	// 这是一个核心技巧，目的是让接下来的 Read 操作几乎立即超时，从而快速判断连接状态，而不是长时间阻塞等待数据
	conn.SetReadDeadline(time.Now().Add(time.Millisecond))

	// 尝试读取：尝试从连接中读取数据到缓冲区 buf。由于设置了极短超时，这个操作会很快返回
	n, err := conn.Read(buf)

	// 检查意外数据：如果读取没有错误 (err == nil) 或者实际读到了数据 (n > 0)，
	// 对于一个预期为“空闲”的连接来说，这属于异常情况。函数因此返回一个“unexpected read from socket”错误。
	// Idle connections should not read data, it is an unexpected read error.
	if err == nil || n > 0 {
		return errors.New("unexpected read from socket")
	}

	// 处理超时错误：这是关键逻辑。如果返回的错误是一个网络错误 (net.Error)，并且是超时类型 (Timeout() 方法返回 true)，
	// 这恰恰说明连接是正常的——它只是因为超时时间内没有数据可读而返回。
	// The idle connection is normal and returns timeout.
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		// 取消设置的读超时，避免影响该连接后续的正常使用
		conn.SetReadDeadline(time.Time{})
		return nil
	}

	// other connection errors, including connection closed.
	return err
}

func CheckConnErrUnblock(conn net.Conn, buf []byte) error {
	// Currently non-blocking mode is not supported.
	return nil
}

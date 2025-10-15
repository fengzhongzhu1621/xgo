//go:build aix || darwin || dragonfly || freebsd || netbsd || openbsd || solaris || linux
// +build aix darwin dragonfly freebsd netbsd openbsd solaris linux

package checker

import (
	"errors"
	"io"
	"net"
	"syscall"
)

// CheckConnErr 检查连接错误，使用阻塞方式检查连接状态
// 参数:
//
//	conn: 要检查的网络连接
//	buf: 用于读取操作的缓冲区
//
// 返回:
//
//	error: 连接错误，如果连接正常则返回 nil
func CheckConnErr(conn net.Conn, buf []byte) error {
	return CheckConnErrUnblock(conn, buf)
}

// CheckConnErrUnblock 使用非阻塞方式快速检查连接错误
// 这个方法通过系统调用直接检查连接状态，避免阻塞操作
// Go 1.7+ 的零字节读取：在 Go 1.7 及更高版本中，对未设置超时的连接进行零字节读取会立即返回，不会阻塞，
// 但其行为在设置超时后有所不同。此函数通过设置超时并要求至少尝试读取一个字节来规避了潜在问题。

// 参数:
//
//	conn: 要检查的网络连接
//	buf: 用于读取操作的缓冲区
//
// 返回:
//
//	error: 连接错误，如果连接正常则返回 nil
func CheckConnErrUnblock(conn net.Conn, buf []byte) error {
	// 检查连接是否支持系统调用接口
	// 如果 conn 不支持 syscall.Conn 接口（例如某些自定义包装的连接类型），则直接返回 nil，表示无法通过此方法检测，认为连接正常。
	sysConn, ok := conn.(syscall.Conn)
	if !ok {
		return nil // 认为连接正常
	}

	// 获取原始连接：通过 SyscallConn() 方法获取一个 syscall.RawConn 接口对象 rawConn。
	// 它允许在连接上执行底层的、不阻塞Go调度器的系统调用。
	rawConn, err := sysConn.SyscallConn()
	if err != nil {
		return err // 连接异常
	}

	var sysErr error
	var n int // 接收系统调用实际读取的字节数

	// 使用系统调用进行非阻塞读取操作，rawConn.Read 内部会协调Go运行时网络轮询器（netpoller），确保系统调用不会阻塞调度。
	err = rawConn.Read(func(fd uintptr) bool {
		// 进行底层的文件读取操作
		// Go 默认将 socket 设置为非阻塞模式，调用 syscall 可以直接返回
		// 参考 Go 源码: src/net/sock_cloexec.go 中的 sysSocket() 函数
		// n: 实际读取的字节数。可能为 0（表示文件末尾）或小于 buf 的长度（表示虽未出错但无更多数据）。
		n, sysErr = syscall.Read(int(fd), buf)
		// 返回 true，不会执行 net 库封装的阻塞等待，直接返回
		// 回调函数返回 true，告知 rawConn.Read 方法“我已处理完毕，无需重试”。这确保了操作的即时性，不会进入netpoller的阻塞等待逻辑
		return true
	})
	if err != nil {
		// 如果 rawConn.Read 方法本身执行出错（非回调函数内的 sysErr），返回该错误
		return err // 连接异常
	}

	// 如果系统调用读取的字节数 n 为0且系统错误 sysErr 为 nil，这通常表示对端已正常关闭连接（收到了FIN包），因此返回 io.EOF 错误
	if n == 0 && sysErr == nil {
		return io.EOF // 服务端连接已关闭
	}

	// 判断意外数据：如果读取到数据（n > 0），但对于一个预期为“空闲”的连接来说这是意外的，返回一个自定义错误。
	// 空闲连接不应该读取到数据
	if n > 0 {
		return errors.New("unexpected read from socket")
	}

	// 忽略指定的错误
	// 处理预期错误：如果系统错误是 EAGAIN 或 EWOULDBLOCK，这在意料之中，它表示在非阻塞模式下当前没有数据可读，但连接本身是正常的。因此返回 nil。
	if sysErr == syscall.EAGAIN || sysErr == syscall.EWOULDBLOCK {
		return nil
	}

	return sysErr
}

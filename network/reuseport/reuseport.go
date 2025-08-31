//go:build linux || darwin || dragonfly || freebsd || netbsd || openbsd
// +build linux darwin dragonfly freebsd netbsd openbsd

// Package reuseport provides a function that returns a net.Listener powered
// by a net.FileListener with a SO_REUSEPORT option set to the socket.
// 这是一个端口复用（SO_REUSEPORT）的实现，允许多个进程或线程绑定到相同的IP地址和端口。主要功能：

// 核心特性：
// 支持TCP、UDP、Unix域套接字的端口复用
// 提供统一的Listen和ListenPacket接口
// 自动检测系统是否支持SO_REUSEPORT
// 优雅降级到普通监听模式
//
// 使用场景：
// 负载均衡 - 多个进程共享同一端口接收连接
// 热升级 - 新旧版本同时监听，平滑切换
// 性能优化 - 减少锁竞争，提高并发处理能力
package reuseport

import (
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
)

const fileNameTemplate = "reuseport.%d.%s.%s"

var errUnsupportedProtocol = errors.New("only tcp, tcp4, tcp6, udp, udp4, udp6 are supported")

// getSockaddr parses protocol and address and returns implementor
// of syscall.Sockaddr: syscall.SockaddrInet4 or syscall.SockaddrInet6.
// 解析协议和地址，返回对应的syscall.Sockaddr实现
func getSockaddr(proto, addr string) (sa syscall.Sockaddr, soType int, err error) {
	switch proto {
	case "tcp", "tcp4", "tcp6":
		return getTCPSockaddr(proto, addr)
	case "udp", "udp4", "udp6":
		return getUDPSockaddr(proto, addr)
	default:
		return nil, -1, errUnsupportedProtocol
	}
}

func getSocketFileName(proto, addr string) string {
	return fmt.Sprintf(fileNameTemplate, os.Getpid(), proto, addr)  // 生成唯一的socket文件名
}

// Listen function is an alias for NewReusablePortListener.
func Listen(proto, addr string) (l net.Listener, err error) {
	return NewReusablePortListener(proto, addr)
}

// ListenPacket is an alias for NewReusablePortPacketConn.
func ListenPacket(proto, addr string) (l net.PacketConn, err error) {
	return NewReusablePortPacketConn(proto, addr)
}

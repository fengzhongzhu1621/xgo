//go:build linux || darwin || dragonfly || freebsd || netbsd || openbsd
// +build linux darwin dragonfly freebsd netbsd openbsd

package reuseport

import (
	"errors"
	"net"
	"os"
	"syscall"
)

var (
	// ListenerBacklogMaxSize setting backlog size
	ListenerBacklogMaxSize    = maxListenerBacklog()
	errUnsupportedTCPProtocol = errors.New("only tcp, tcp4, tcp6 are supported")
)

// getTCPSockaddr 解析TCP协议和地址，返回对应的系统调用sockaddr结构
// 参数：proto - 协议类型，addr - 地址字符串
// 返回值：sockaddr结构、socket类型、错误信息
func getTCPSockaddr(proto, addr string) (sa syscall.Sockaddr, soType int, err error) {
	tcp, tcpVersion, err := getTCPAddr(proto, addr) // 解析TCP地址
	if err != nil {
		return nil, -1, err
	}
	switch tcpVersion {
	case "tcp": // 通用TCP协议，默认使用IPv4
		return &syscall.SockaddrInet4{Port: tcp.Port}, syscall.AF_INET, nil
	case "tcp4": // IPv4 TCP协议
		return getTCP4Sockaddr(tcp)
	default: // 必须是 "tcp6" - IPv6 TCP协议
		return getTCP6Sockaddr(tcp)
	}
}

func getTCPAddr(proto, addr string) (*net.TCPAddr, string, error) {
	var tcp *net.TCPAddr

	// fix bugs https://github.com/kavu/go_reuseport/pull/33
	tcp, err := net.ResolveTCPAddr(proto, addr)
	if err != nil {
		return nil, "", err
	}

	tcpVersion, err := determineTCPProto(proto, tcp)
	if err != nil {
		return nil, "", err
	}
	return tcp, tcpVersion, nil
}

// getTCP4Sockaddr 将net.TCPAddr转换为IPv4的sockaddr结构
// 参数：tcp - TCP地址信息
// 返回值：IPv4 sockaddr结构、socket类型(AF_INET)、错误信息
func getTCP4Sockaddr(tcp *net.TCPAddr) (syscall.Sockaddr, int, error) {
	sa := &syscall.SockaddrInet4{Port: tcp.Port} // 创建IPv4 sockaddr结构，设置端口

	if tcp.IP != nil {
		if len(tcp.IP) == 16 { // IPv4映射的IPv6地址（16字节）
			// 复制IPv6地址的最后4字节（实际IPv4地址）到数组
			copy(sa.Addr[:], tcp.IP[12:16])
		} else { // 纯IPv4地址（4字节）
			// 复制所有字节到数组
			copy(sa.Addr[:], tcp.IP)
		}
	}

	return sa, syscall.AF_INET, nil // 返回IPv4 socket类型
}

// getTCP6Sockaddr 将net.TCPAddr转换为IPv6的sockaddr结构
// 参数：tcp - TCP地址信息
// 返回值：IPv6 sockaddr结构、socket类型(AF_INET6)、错误信息
func getTCP6Sockaddr(tcp *net.TCPAddr) (syscall.Sockaddr, int, error) {
	sa := &syscall.SockaddrInet6{Port: tcp.Port} // 创建IPv6 sockaddr结构，设置端口

	if tcp.IP != nil {
		copy(sa.Addr[:], tcp.IP) // 复制所有字节到数组
	}

	if tcp.Zone != "" { // 处理IPv6区域标识符（如eth0, wlan0等）
		iface, err := net.InterfaceByName(tcp.Zone)
		if err != nil {
			return nil, -1, err
		}

		sa.ZoneId = uint32(iface.Index) // 设置网络接口索引
	}

	return sa, syscall.AF_INET6, nil // 返回IPv6 socket类型
}

func determineTCPProto(proto string, ip *net.TCPAddr) (string, error) {
	// If the protocol is set to "tcp", we try to determine the actual protocol
	// version from the size of the resolved IP address. Otherwise, we simple use
	// the protocol given to us by the caller.

	if ip.IP.To4() != nil {
		return "tcp4", nil
	}

	if ip.IP.To16() != nil {
		return "tcp6", nil
	}

	switch proto {
	case "tcp", "tcp4", "tcp6":
		return proto, nil
	default:
		return "", errUnsupportedTCPProtocol
	}
}

// NewReusablePortListener returns net.FileListener that created from
// a file descriptor for a socket with SO_REUSEPORT option.
func NewReusablePortListener(proto, addr string) (l net.Listener, err error) {
	var (
		soType, fd int
		sockaddr   syscall.Sockaddr
	)
	// 解析协议和地址，返回对应的系统调用sockaddr结构
	if sockaddr, soType, err = getSockaddr(proto, addr); err != nil {
		return nil, err
	}

	// 创建socket文件描述符
	syscall.ForkLock.RLock()
	if fd, err = syscall.Socket(soType, syscall.SOCK_STREAM, syscall.IPPROTO_TCP); err != nil {
		syscall.ForkLock.RUnlock()
		return nil, err
	}
	syscall.ForkLock.RUnlock()

	// 创建可重用的文件描述符，设置socket选项并绑定地址
	if err = createReusableFd(fd, sockaddr); err != nil {
		return nil, err
	}

	// 从文件描述符创建可重用的网络监听器
	return createReusableListener(fd, proto, addr)
}

// createReusableFd 创建可重用的文件描述符，设置socket选项并绑定地址
// 参数：fd - 文件描述符，sockaddr - socket地址结构
// 返回值：错误信息
func createReusableFd(fd int, sockaddr syscall.Sockaddr) (err error) {
	defer func() {
		if err != nil { // 如果发生错误，确保关闭文件描述符
			syscall.Close(fd)
		}
	}()

	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return err
	}

	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, reusePort, 1); err != nil { // 设置SO_REUSEPORT选项，允许多个进程绑定相同端口
		return err
	}

	if err = syscall.Bind(fd, sockaddr); err != nil { // 绑定socket到指定地址
		return err
	}

	// Set backlog size to the maximum
	if err = syscall.Listen(fd, ListenerBacklogMaxSize); err != nil { // 开始监听，设置最大连接队列大小
		return err
	}

	return nil
}

// createReusableListener 从文件描述符创建可重用的网络监听器
// 参数：fd - 文件描述符，proto - 协议类型，addr - 地址字符串
// 返回值：网络监听器、错误信息
func createReusableListener(fd int, proto, addr string) (l net.Listener, err error) {
	file := os.NewFile(uintptr(fd), getSocketFileName(proto, addr)) // 从文件描述符创建os.File对象
	if l, err = net.FileListener(file); err != nil {                // 将文件转换为网络监听器
		file.Close() // 如果转换失败，关闭文件
		return nil, err
	}

	if err = file.Close(); err != nil { // 关闭文件，监听器会继续使用底层文件描述符
		return nil, err
	}
	return l, err
}

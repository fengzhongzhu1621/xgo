//go:build linux || darwin || dragonfly || freebsd || netbsd || openbsd
// +build linux darwin dragonfly freebsd netbsd openbsd

package reuseport

import (
	"errors"
	"net"
	"os"
	"syscall"
)

var errUnsupportedUDPProtocol = errors.New("only udp, udp4, udp6 are supported")

func getUDPSockaddr(proto, addr string) (sa syscall.Sockaddr, soType int, err error) {
	udp, udpVersion, err := getUDPAddr(proto, addr)
	if err != nil {
		return nil, -1, err
	}

	switch udpVersion {
	case "udp":
		return &syscall.SockaddrInet4{Port: udp.Port}, syscall.AF_INET, nil
	case "udp4":
		return getUDP4Sockaddr(udp)
	default:
		// must be "udp6"
		return getUDP6Sockaddr(udp)
	}
}

func getUDP4Sockaddr(udp *net.UDPAddr) (syscall.Sockaddr, int, error) {
	sa := &syscall.SockaddrInet4{Port: udp.Port}

	if udp.IP != nil {
		if len(udp.IP) == 16 {
			copy(sa.Addr[:], udp.IP[12:16]) // copy last 4 bytes of slice to array
		} else {
			copy(sa.Addr[:], udp.IP) // copy all bytes of slice to array
		}
	}

	return sa, syscall.AF_INET, nil
}

func getUDP6Sockaddr(udp *net.UDPAddr) (syscall.Sockaddr, int, error) {
	sa := &syscall.SockaddrInet6{Port: udp.Port}

	if udp.IP != nil {
		copy(sa.Addr[:], udp.IP) // copy all bytes of slice to array
	}

	if udp.Zone != "" {
		iface, err := net.InterfaceByName(udp.Zone)
		if err != nil {
			return nil, -1, err
		}

		sa.ZoneId = uint32(iface.Index)
	}

	return sa, syscall.AF_INET6, nil
}

func getUDPAddr(proto, addr string) (*net.UDPAddr, string, error) {

	var udp *net.UDPAddr

	udp, err := net.ResolveUDPAddr(proto, addr)
	if err != nil {
		return nil, "", err
	}

	udpVersion, err := determineUDPProto(proto, udp)
	if err != nil {
		return nil, "", err
	}

	return udp, udpVersion, nil
}
func determineUDPProto(proto string, ip *net.UDPAddr) (string, error) {
	// If the protocol is set to "udp", we try to determine the actual protocol
	// version from the size of the resolved IP address. Otherwise, we simple use
	// the protocol given to us by the caller.

	if ip.IP.To4() != nil {
		return "udp4", nil
	}

	if ip.IP.To16() != nil {
		return "udp6", nil
	}

	switch proto {
	case "udp", "udp4", "udp6":
		return proto, nil
	default:
		return "", errUnsupportedUDPProtocol
	}
}

// NewReusablePortPacketConn returns net.FilePacketConn that created from
// a file descriptor for a socket with SO_REUSEPORT option.
func NewReusablePortPacketConn(proto, addr string) (net.PacketConn, error) {
	sockaddr, soType, err := getSockaddr(proto, addr)
	if err != nil {
		return nil, err
	}

	syscall.ForkLock.RLock()
	fd, err := syscall.Socket(soType, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err == nil {
		syscall.CloseOnExec(fd)
	}
	syscall.ForkLock.RUnlock()
	if err != nil {
		syscall.Close(fd)
		return nil, err
	}

	// 创建可重用的数据包连接（UDP）
	return createPacketConn(fd, sockaddr, getSocketFileName(proto, addr))
}

// createPacketConn 创建可重用的数据包连接（UDP）
// 参数：fd - 文件描述符，sockaddr - socket地址结构，fdName - 文件描述符名称
// 返回值：数据包连接、错误信息
func createPacketConn(fd int, sockaddr syscall.Sockaddr, fdName string) (net.PacketConn, error) {
	// 设置socket选项
	if err := setPacketConnSockOpt(fd, sockaddr); err != nil {
		syscall.Close(fd) // 如果设置失败，关闭文件描述符
		return nil, err
	}

	file := os.NewFile(uintptr(fd), fdName) // 从文件描述符创建os.File对象
	l, err := net.FilePacketConn(file)      // 将文件转换为数据包连接
	if err != nil {
		syscall.Close(fd) // 如果转换失败，关闭文件描述符
		return nil, err
	}

	if err = file.Close(); err != nil {
		syscall.Close(fd)
		return nil, err
	}
	return l, err
}

// setPacketConnSockOpt 设置数据包连接的socket选项
// 参数：fd - 文件描述符，sockaddr - socket地址结构
// 返回值：错误信息
func setPacketConnSockOpt(fd int, sockaddr syscall.Sockaddr) error {
	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil { // 设置SO_REUSEADDR选项
		return err
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, reusePort, 1); err != nil { // 设置SO_REUSEPORT选项，允许多个进程绑定相同端口
		return err
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1); err != nil { // 设置SO_BROADCAST选项，允许广播
		return err
	}

	return syscall.Bind(fd, sockaddr) // 绑定socket到指定地址
}

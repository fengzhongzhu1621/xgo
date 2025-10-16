package ip

import (
	"errors"
	"net"
)

// GetFreePort 获取可用的空闲端口
// 参数:
//   - network: 网络协议类型 (tcp/tcp4/tcp6/udp/udp4/udp6)
//
// 返回值:
//   - int: 可用的端口号
//   - error: 错误信息
func GetFreePort(network string) (int, error) {
	// 处理TCP协议类型的端口获取
	if network == "tcp" || network == "tcp4" || network == "tcp6" {
		// 解析TCP地址，使用localhost:0让系统自动分配端口
		addr, err := net.ResolveTCPAddr(network, "localhost:0")
		if err != nil {
			return -1, err
		}

		// 创建TCP监听器
		l, err := net.ListenTCP(network, addr)
		if err != nil {
			return -1, err
		}
		defer l.Close() // 确保监听器被关闭

		// 获取实际分配的端口号
		return l.Addr().(*net.TCPAddr).Port, nil
	}

	// 处理UDP协议类型的端口获取
	if network == "udp" || network == "udp4" || network == "udp6" {
		// 解析UDP地址，使用localhost:0让系统自动分配端口
		addr, err := net.ResolveUDPAddr(network, "localhost:0")
		if err != nil {
			return -1, err
		}

		// 创建UDP监听器
		l, err := net.ListenUDP(network, addr)
		if err != nil {
			return -1, err
		}
		defer l.Close() // 确保监听器被关闭

		// 获取实际分配的端口号
		return l.LocalAddr().(*net.UDPAddr).Port, nil
	}

	// 不支持的协议类型
	return -1, errors.New("invalid network")
}

package ip

import (
	"fmt"
	"net"
	"strings"
)

// AddrToKey combines local and remote address into a string.
func AddrToKey(local, remote net.Addr) string {
	return strings.Join([]string{local.Network(), local.String(), remote.String()}, "_")
}

// GetFreeAddr 获取可用的空闲地址（格式为 ":端口号"）
// 参数:
//   - network: 网络协议类型 (tcp/tcp4/tcp6/udp/udp4/udp6)
//
// 返回值: 格式为 ":端口号" 的地址字符串
func GetFreeAddr(network string) string {
	// 调用 GetFreePort 函数获取可用的端口号
	p, err := GetFreePort(network)
	if err != nil {
		// 如果获取端口失败，直接panic终止程序
		// 这通常用于测试环境，表示端口获取是必需的
		panic(err)
	}

	// 返回格式化的地址字符串，如 ":8080"
	return fmt.Sprintf(":%d", p)
}

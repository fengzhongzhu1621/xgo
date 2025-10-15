package multiplexed

import (
	"net"
	"strings"
	"time"

	"github.com/fengzhongzhu1621/xgo/network/dial"
)

// isStream 判断网络协议是否为流式连接
// 参数:
//
//	network: 网络协议类型字符串
//
// 返回:
//
//	bool: true表示流式连接，false表示数据报连接
//	error: 不支持的协议类型错误
func isStream(network string) (bool, error) {
	switch network {
	case "tcp", "tcp4", "tcp6", "unix":
		return true, nil
	case "udp", "udp4", "udp6":
		return false, nil
	default:
		return false, ErrNetworkNotSupport
	}
}

// makeNodeKey 生成节点唯一标识键
// 参数:
//
//	network: 网络协议类型
//	address: 目标地址
//
// 返回:
//
//	string: 格式为"network_address"的唯一节点键
func makeNodeKey(network, address string) string {
	var key strings.Builder
	key.Grow(len(network) + len(address) + 1)
	key.WriteString(network)
	key.WriteString("_")
	key.WriteString(address)
	return key.String()
}

// filterOutConnection 从连接列表中过滤掉指定的连接
// 参数:
//
//	in: 原始连接列表
//	exclude: 需要排除的连接
//
// 返回:
//
//	[]*Connection: 过滤后的连接列表
func filterOutConnection(in []*Connection, exclude *Connection) []*Connection {
	// 创建新切片：基于原切片 in 创建一个长度为 0 但容量与 in 相同的新切片 out。
	// 这是一种高效的内存复用技巧，避免了立即分配新数组；int 和 out 指向切片内部的同一块内存空间
	out := in[:0]

	// 由于 out 的容量足够，这些追加操作不会引起新的内存分配，直接复用了 in 的底层数组
	for _, v := range in {
		if v != exclude {
			out = append(out, v)
		}
	}

	// 成功移除连接后，将切片末尾的值置空以避免内存泄漏
	// 过滤后，out 切片的长度小于等于原 in 切片的长度。循环从 out 的长度开始，直到原 in 切片的末尾，将底层数组中剩余位置的指针设置为 nil
	// 这样做是为了打破对已过滤掉连接对象的引用。如果不置为 nil，原切片 in 的底层数组中未被 out 使用的部分仍然会持有这些 Connection 对象的引用，
	// 即使程序逻辑上已经不再需要它们，垃圾回收器（GC）也无法回收这些对象占用的内存，从而导致内存泄漏
	for i := len(out); i < len(in); i++ {
		in[i] = nil
	}

	return out
}

// dialTCP 建立TCP连接
// 参数:
//
//	timeout: 连接超时时间
//	opts: 连接配置选项
//
// 返回:
//
//	net.Conn: 建立的TCP连接
//	*dial.DialOptions: 使用的拨号选项
//	error: 连接过程中发生的错误
func dialTCP(timeout time.Duration, opts *GetOptions) (net.Conn, *dial.DialOptions, error) {
	dialOpts := &dial.DialOptions{
		Network:       opts.network,
		Address:       opts.address,
		Timeout:       timeout,
		CACertFile:    opts.CACertFile,
		TLSCertFile:   opts.TLSCertFile,
		TLSKeyFile:    opts.TLSKeyFile,
		TLSServerName: opts.TLSServerName,
		LocalAddr:     opts.LocalAddr,
	}
	conn, err := dial.TryConnect(dialOpts)
	return conn, dialOpts, err
}

// dialUDP 建立UDP连接
// 参数:
//
//	opts: 连接配置选项
//
// 返回:
//
//	net.PacketConn: 建立的UDP数据包连接
//	*net.UDPAddr: 解析后的目标UDP地址
//	error: 连接过程中发生的错误
func dialUDP(opts *GetOptions) (net.PacketConn, *net.UDPAddr, error) {
	addr, err := net.ResolveUDPAddr(opts.network, opts.address)
	if err != nil {
		return nil, nil, err
	}
	const defaultLocalAddr = ":"
	localAddr := defaultLocalAddr
	if opts.LocalAddr != "" {
		localAddr = opts.LocalAddr
	}
	conn, err := net.ListenPacket(opts.network, localAddr)
	if err != nil {
		return nil, nil, err
	}
	return conn, addr, nil
}

package ip

import "errors"

var (
	// ErrNetworkNotSupport does not support network type.
	ErrNetworkNotSupport = errors.New("network not support")
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
func IsStream(network string) (bool, error) {
	switch network {
	case "tcp", "tcp4", "tcp6", "unix":
		return true, nil
	case "udp", "udp4", "udp6":
		return false, nil
	default:
		return false, ErrNetworkNotSupport
	}
}

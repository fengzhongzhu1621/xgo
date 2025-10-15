package dial

import "time"

// DialOptions 拨号请求参数配置
type DialOptions struct {
	Network       string        // 网络协议类型，如 "tcp", "udp" 等
	Address       string        // 目标地址，格式为 "host:port"
	LocalAddr     string        // 本地绑定地址，格式为 "host:port"
	Timeout       time.Duration // 连接超时时间
	CACertFile    string        // CA 证书文件路径，用于 TLS 认证
	TLSCertFile   string        // 客户端证书文件路径
	TLSKeyFile    string        // 客户端私钥文件路径
	TLSServerName string        // 客户端验证服务器的服务名称
	// 如果未填写，默认为 http 主机名
	IdleTimeout time.Duration // 连接空闲超时时间
}

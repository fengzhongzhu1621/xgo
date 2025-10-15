package dial

import (
	"crypto/tls"
	"net"

	"github.com/fengzhongzhu1621/xgo/network/ssl"
	"trpc.group/trpc-go/trpc-go/errs"
)

// DialFunc 使用选项中的信息连接到端点的函数类型
type DialFunc func(opts *DialOptions) (net.Conn, error)

// Dial 发起连接请求，根据选项配置建立网络连接
// 参数:
//
//	opts: 拨号选项配置
//
// 返回:
//
//	net.Conn: 建立的网络连接
//	error: 连接过程中发生的错误
func Dial(opts *DialOptions) (net.Conn, error) {
	var localAddr net.Addr
	// 解析本地绑定地址
	if opts.LocalAddr != "" {
		var err error
		localAddr, err = net.ResolveTCPAddr(opts.Network, opts.LocalAddr)
		if err != nil {
			return nil, err
		}
	}

	// 创建拨号器
	dialer := &net.Dialer{
		Timeout:   opts.Timeout, // 设置连接超时
		LocalAddr: localAddr,    // 设置本地绑定地址
	}

	// 如果不需要 TLS 认证，直接使用普通连接
	if opts.CACertFile == "" {
		return dialer.Dial(opts.Network, opts.Address)
	}

	// 设置 TLS 服务器名称，如果未指定则使用目标地址
	if opts.TLSServerName == "" {
		opts.TLSServerName = opts.Address
	}

	// 获取 TLS 客户端配置
	tlsConf, err := ssl.GetClientConfig(opts.TLSServerName, opts.CACertFile, opts.TLSCertFile, opts.TLSKeyFile)
	if err != nil {
		return nil, errs.NewFrameError(errs.RetClientDecodeFail, "client dial tls fail: "+err.Error())
	}

	// 使用 TLS 拨号器建立安全连接
	return tls.DialWithDialer(dialer, opts.Network, opts.Address, tlsConf)
}

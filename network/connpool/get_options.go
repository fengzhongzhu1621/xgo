package connpool

import (
	"context"
	"io"
	"time"

	"github.com/fengzhongzhu1621/xgo/buildin/buffer"
	"github.com/fengzhongzhu1621/xgo/codec"
)

// GetOptions 是获取连接的配置参数
// 用于配置 dial.DialFunc 的参数
type GetOptions struct {
	FramerBuilder codec.IFramerBuilder      // 帧构建器，用于创建帧读取器
	CustomReader  func(io.Reader) io.Reader // 自定义读取器，用于包装底层连接
	Ctx           context.Context           // 请求上下文，用于传递上下文参数

	CACertFile    string // CA 证书文件路径，用于 TLS 认证
	TLSCertFile   string // 客户端证书文件路径
	TLSKeyFile    string // 客户端私钥文件路径
	TLSServerName string // 客户端验证服务器的服务名称
	// 如果未填写，默认为 http 主机名

	LocalAddr   string        // 建立连接时的本地地址，默认随机选择多网卡
	DialTimeout time.Duration // 连接建立超时时间
	Protocol    string        // 协议类型
}

// NewGetOptions 创建并初始化 GetOptions
// 返回:
//
//	GetOptions: 初始化后的配置选项
func NewGetOptions() GetOptions {
	return GetOptions{
		CustomReader: buffer.NewReader, // 默认使用缓冲读取器
	}
}

// getDialCtx 获取拨号操作的上下文
// 参数:
//
//	dialTimeout: 默认拨号超时时间
//
// 返回:
//
//	context.Context: 拨号上下文
//	context.CancelFunc: 取消函数
func (o *GetOptions) getDialCtx(dialTimeout time.Duration) (context.Context, context.CancelFunc) {
	ctx := o.Ctx
	defer func() {
		// opts.Ctx 仅用于传递 ctx 参数，不建议数据结构持有 ctx
		o.Ctx = nil
	}()

	for {
		// 如果 RPC 请求没有设置 ctx，创建新的 ctx
		if ctx == nil {
			break
		}
		// 如果 RPC 请求没有设置 ctx 超时，创建新的 ctx
		deadline, ok := ctx.Deadline()
		if !ok {
			break
		}
		// 如果 RPC 请求超时大于设置的超时，创建新的 ctx
		d := time.Until(deadline)
		if o.DialTimeout > 0 && o.DialTimeout < d {
			break
		}

		// 需要 o.DialTimeout >= d（ctx 设置的超时）
		return ctx, nil
	}

	// 使用配置的超时时间或默认超时时间
	if o.DialTimeout > 0 {
		dialTimeout = o.DialTimeout
	}
	if dialTimeout == 0 {
		dialTimeout = defaultDialTimeout
	}
	// 创建带超时的上下文
	return context.WithTimeout(context.Background(), dialTimeout)
}

// WithFramerBuilder 设置帧构建器的选项
// 参数:
//
//	fb: 帧构建器接口
func (o *GetOptions) WithFramerBuilder(fb codec.IFramerBuilder) {
	o.FramerBuilder = fb
}

// WithDialTLS 设置客户端支持 TLS 的选项
// 参数:
//
//	certFile: 客户端证书文件路径
//	keyFile: 客户端私钥文件路径
//	caFile: CA 证书文件路径
//	serverName: 服务器名称
func (o *GetOptions) WithDialTLS(certFile, keyFile, caFile, serverName string) {
	o.TLSCertFile = certFile
	o.TLSKeyFile = keyFile
	o.CACertFile = caFile
	o.TLSServerName = serverName
}

// WithContext 设置请求上下文的选项
// 参数:
//
//	ctx: 上下文对象
func (o *GetOptions) WithContext(ctx context.Context) {
	o.Ctx = ctx
}

// WithLocalAddr 设置建立连接时本地地址的选项
// 当有多个网卡时，默认随机选择
// 参数:
//
//	addr: 本地地址
func (o *GetOptions) WithLocalAddr(addr string) {
	o.LocalAddr = addr
}

// WithDialTimeout 设置连接超时的选项
// 参数:
//
//	dur: 超时时间
func (o *GetOptions) WithDialTimeout(dur time.Duration) {
	o.DialTimeout = dur
}

// WithProtocol 设置后端服务协议名称的选项
// 参数:
//
//	s: 协议名称
func (o *GetOptions) WithProtocol(s string) {
	o.Protocol = s
}

// WithCustomReader 设置自定义读取器的选项
// 连接池将使用此自定义读取器创建包装底层连接的读取器，通常用于创建缓冲区
// 参数:
//
//	customReader: 自定义读取器函数
func (o *GetOptions) WithCustomReader(customReader func(io.Reader) io.Reader) {
	o.CustomReader = customReader
}

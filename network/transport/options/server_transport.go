package options

import (
	"runtime"
	"time"
)

const (
	defaultRecvMsgChannelSize      = 100
	defaultSendMsgChannelSize      = 100
	defaultRecvUDPPacketBufferSize = 65536
	defaultIdleTimeout             = time.Minute
)

// ServerTransportOptions 服务器传输配置选项
type ServerTransportOptions struct {
	RecvMsgChannelSize      int           // TCP接收消息通道大小
	SendMsgChannelSize      int           // TCP发送消息通道大小
	RecvUDPPacketBufferSize int           // UDP预分配缓冲区大小
	RecvUDPRawSocketBufSize int           // UDP原始套接字接收缓冲区大小
	IdleTimeout             time.Duration // 连接空闲超时时间
	KeepAlivePeriod         time.Duration // TCP保活周期
	ReusePort               bool          // 是否启用端口复用
}

// DefaultServerTransportOptions 返回默认的服务器传输选项
func DefaultServerTransportOptions() *ServerTransportOptions {
	return &ServerTransportOptions{
		RecvMsgChannelSize:      defaultRecvMsgChannelSize,
		SendMsgChannelSize:      defaultSendMsgChannelSize,
		RecvUDPPacketBufferSize: defaultRecvUDPPacketBufferSize,
	}
}

// ServerTransportOption 服务器传输选项修改函数类型
type ServerTransportOption func(*ServerTransportOptions)

// WithReusePort 返回启用端口复用的服务器传输选项
// 注意：Windows系统不支持端口复用
func WithReusePort(reuse bool) ServerTransportOption {
	return func(options *ServerTransportOptions) {
		options.ReusePort = reuse
		if runtime.GOOS == "windows" {
			options.ReusePort = false
		}
	}
}

// WithRecvMsgChannelSize 设置TCP服务器传输接收缓冲区大小
func WithRecvMsgChannelSize(size int) ServerTransportOption {
	return func(options *ServerTransportOptions) {
		options.RecvMsgChannelSize = size
	}
}

// WithSendMsgChannelSize 设置TCP服务器传输发送通道大小
func WithSendMsgChannelSize(size int) ServerTransportOption {
	return func(options *ServerTransportOptions) {
		options.SendMsgChannelSize = size
	}
}

// WithRecvUDPPacketBufferSize 设置UDP服务器传输预分配缓冲区大小
func WithRecvUDPPacketBufferSize(size int) ServerTransportOption {
	return func(options *ServerTransportOptions) {
		options.RecvUDPPacketBufferSize = size
	}
}

// WithRecvUDPRawSocketBufSize 设置UDP连接操作系统接收缓冲区大小
func WithRecvUDPRawSocketBufSize(size int) ServerTransportOption {
	return func(options *ServerTransportOptions) {
		options.RecvUDPRawSocketBufSize = size
	}
}

// WithIdleTimeout 设置服务器连接空闲超时时间
func WithIdleTimeout(timeout time.Duration) ServerTransportOption {
	return func(options *ServerTransportOptions) {
		options.IdleTimeout = timeout
	}
}

// WithKeepAlivePeriod 设置TCP连接保活周期
// 注意：TLS连接不支持此选项，因为TLS不使用net.TCPConn或net.Conn
func WithKeepAlivePeriod(d time.Duration) ServerTransportOption {
	return func(options *ServerTransportOptions) {
		options.KeepAlivePeriod = d
	}
}

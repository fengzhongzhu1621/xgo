package options

import (
	"net"
	"time"

	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/network/transport/handler"
)

// ListenServeOption 服务器监听服务选项修改函数类型
type ListenServeOption func(*ListenServeOptions)

// ListenServeOptions 服务器启动时的监听服务选项
type ListenServeOptions struct {
	ServiceName   string               // 服务名称
	Address       string               // 监听地址
	Network       string               // 网络类型 (tcp/udp)，多个逗号分割
	Handler       handler.IHandler     // 业务处理器
	FramerBuilder codec.IFramerBuilder // 帧构建器
	Listener      net.Listener         // 自定义监听器

	CACertFile  string        // CA证书文件路径
	TLSCertFile string        // 服务器证书文件路径
	TLSKeyFile  string        // 服务器密钥文件路径
	Routines    int           // 协程池大小
	ServerAsync bool          // 是否启用服务器异步模式
	Writev      bool          // 是否启用writev优化
	CopyFrame   bool          // 是否复制帧数据
	IdleTimeout time.Duration // 连接空闲超时时间

	// DisableKeepAlives 如果为true，禁用keep-alive，每个连接只处理单个请求
	// 用于RPC传输层（如HTTP），与TCP keep-alive无关
	DisableKeepAlives bool

	// StopListening 用于通知服务器传输停止监听
	StopListening <-chan struct{}
}

// WithServiceName returns a ListenServeOption which sets the service name.
func WithServiceName(name string) ListenServeOption {
	return func(opts *ListenServeOptions) {
		opts.ServiceName = name
	}
}

// WithServerFramerBuilder returns a ListenServeOption which sets server frame builder.
func WithServerFramerBuilder(fb codec.IFramerBuilder) ListenServeOption {
	return func(opts *ListenServeOptions) {
		opts.FramerBuilder = fb
	}
}

// WithListenAddress returns a ListenServerOption which sets listening address.
func WithListenAddress(address string) ListenServeOption {
	return func(opts *ListenServeOptions) {
		opts.Address = address
	}
}

// WithListenNetwork returns a ListenServeOption which sets listen network.
func WithListenNetwork(network string) ListenServeOption {
	return func(opts *ListenServeOptions) {
		opts.Network = network
	}
}

// WithListener returns a ListenServeOption which allows users to use their customized listener for
// specific accept/read/write logics.
func WithListener(lis net.Listener) ListenServeOption {
	return func(opts *ListenServeOptions) {
		opts.Listener = lis
	}
}

// WithHandler returns a ListenServeOption which sets business Handler.
func WithHandler(handler handler.IHandler) ListenServeOption {
	return func(opts *ListenServeOptions) {
		opts.Handler = handler
	}
}

// WithServeTLS returns a ListenServeOption which sets TLS relatives.
func WithServeTLS(certFile, keyFile, caFile string) ListenServeOption {
	return func(opts *ListenServeOptions) {
		opts.TLSCertFile = certFile
		opts.TLSKeyFile = keyFile
		opts.CACertFile = caFile
	}
}

// WithServerAsync returns a ListenServeOption which enables server async.
// When another frameworks call trpc, they may use long connections. tRPC server can not handle
// them concurrently, thus timeout.
// This option takes effect for each TCP connections.
func WithServerAsync(serverAsync bool) ListenServeOption {
	return func(opts *ListenServeOptions) {
		opts.ServerAsync = serverAsync
	}
}

// WithWritev returns a ListenServeOption which enables writev.
func WithWritev(writev bool) ListenServeOption {
	return func(opts *ListenServeOptions) {
		opts.Writev = writev
	}
}

// WithMaxRoutines returns a ListenServeOption which sets the max number of async goroutines.
// It's recommended to reserve twice of expected goroutines, but no less than MAXPROCS. The default
// value is (1<<31 - 1).
// This option takes effect only when async mod is enabled. It's ignored on sync mod.
func WithMaxRoutines(routines int) ListenServeOption {
	return func(opts *ListenServeOptions) {
		opts.Routines = routines
	}
}

// WithCopyFrame returns a ListenServeOption which sets whether copy frames.
// In stream RPC, even server use sync mod, stream is asynchronous, we need to copy frame to avoid
// over writing.
func WithCopyFrame(copyFrame bool) ListenServeOption {
	return func(opts *ListenServeOptions) {
		opts.CopyFrame = copyFrame
	}
}

// WithDisableKeepAlives returns a ListenServeOption which disables keep-alives.
func WithDisableKeepAlives(disable bool) ListenServeOption {
	return func(options *ListenServeOptions) {
		options.DisableKeepAlives = disable
	}
}

// WithServerIdleTimeout returns a ListenServeOption which sets the server idle timeout.
func WithServerIdleTimeout(timeout time.Duration) ListenServeOption {
	return func(options *ListenServeOptions) {
		options.IdleTimeout = timeout
	}
}

// WithStopListening returns a ListenServeOption which notifies the transport to stop listening.
func WithStopListening(ch <-chan struct{}) ListenServeOption {
	return func(options *ListenServeOptions) {
		options.StopListening = ch
	}
}

package options

import (
	"time"

	"github.com/fengzhongzhu1621/xgo/codec"
	"github.com/fengzhongzhu1621/xgo/network/connpool"
	"github.com/fengzhongzhu1621/xgo/network/multiplexed"
)

// RoundTripOptions 单次往返请求的配置选项
type RoundTripOptions struct {
	Address               string               // IP:端口地址（已从命名服务解析）
	Password              string               // 连接密码
	Network               string               // 网络类型 (tcp/udp)
	LocalAddr             string               // 接受连接时随机选择的本地地址
	DialTimeout           time.Duration        // 拨号超时时间
	Pool                  connpool.IPool       // 客户端连接池
	ReqType               RequestType          // 请求类型：SendAndRecv, SendOnly
	FramerBuilder         codec.IFramerBuilder // 帧构建器
	ConnectionMode        ConnectionMode       // UDP连接模式
	DisableConnectionPool bool                 // 禁用连接池
	EnableMultiplexed     bool                 // 启用多路复用
	Multiplexed           multiplexed.IPool    // 多路复用连接池
	Msg                   codec.IMsg           // 消息对象
	Protocol              string               // 协议类型

	CACertFile    string // CA证书文件路径
	TLSCertFile   string // 客户端证书文件路径
	TLSKeyFile    string // 客户端密钥文件路径
	TLSServerName string // 客户端验证服务器时的名称，默认为HTTP主机名
}

// ConnectionMode is the connection mode, either Connected or NotConnected.
type ConnectionMode bool

// ConnectionMode of UDP.
const (
	Connected    = false // UDP which isolates packets from non-same path
	NotConnected = true  // UDP which allows returning packets from non-same path
)

// RequestType is the client request type, such as SendAndRecv or SendOnly.
type RequestType = codec.RequestType

// Request types.
const (
	SendAndRecv RequestType = codec.SendAndRecv // send and receive
	SendOnly    RequestType = codec.SendOnly    // send only
)

// RoundTripOption modifies the RoundTripOptions.
type RoundTripOption func(*RoundTripOptions)

// WithDialAddress returns a RoundTripOption which sets dial address.
func WithDialAddress(address string) RoundTripOption {
	return func(opts *RoundTripOptions) {
		opts.Address = address
	}
}

// WithDialPassword returns a RoundTripOption which sets dial password.
func WithDialPassword(password string) RoundTripOption {
	return func(opts *RoundTripOptions) {
		opts.Password = password
	}
}

// WithDialNetwork returns a RoundTripOption which sets dial network.
func WithDialNetwork(network string) RoundTripOption {
	return func(opts *RoundTripOptions) {
		opts.Network = network
	}
}

// WithDialPool returns a RoundTripOption which sets dial pool.
func WithDialPool(pool connpool.IPool) RoundTripOption {
	return func(opts *RoundTripOptions) {
		opts.Pool = pool
	}
}

// WithClientFramerBuilder returns a RoundTripOption which sets FramerBuilder.
func WithClientFramerBuilder(builder codec.IFramerBuilder) RoundTripOption {
	return func(opts *RoundTripOptions) {
		opts.FramerBuilder = builder
	}
}

// WithReqType returns a RoundTripOption which sets request type.
func WithReqType(reqType RequestType) RoundTripOption {
	return func(opts *RoundTripOptions) {
		opts.ReqType = reqType
	}
}

// WithConnectionMode returns a RoundTripOption which sets UDP connection mode.
func WithConnectionMode(connMode ConnectionMode) RoundTripOption {
	return func(opts *RoundTripOptions) {
		opts.ConnectionMode = connMode
	}
}

// WithDialTLS returns a RoundTripOption which sets UDP TLS relatives.
func WithDialTLS(certFile, keyFile, caFile, serverName string) RoundTripOption {
	return func(opts *RoundTripOptions) {
		opts.TLSCertFile = certFile
		opts.TLSKeyFile = keyFile
		opts.CACertFile = caFile
		opts.TLSServerName = serverName
	}
}

// WithDisableConnectionPool returns a RoundTripOption which disables connection pool.
func WithDisableConnectionPool() RoundTripOption {
	return func(opts *RoundTripOptions) {
		opts.DisableConnectionPool = true
	}
}

// WithMultiplexed returns a RoundTripOption which enables multiplexed.
func WithMultiplexed(enable bool) RoundTripOption {
	return func(opts *RoundTripOptions) {
		opts.EnableMultiplexed = enable
	}
}

// WithMultiplexedPool returns a RoundTripOption which sets multiplexed pool.
// This function also enables multiplexed.
func WithMultiplexedPool(p multiplexed.IPool) RoundTripOption {
	return func(opts *RoundTripOptions) {
		opts.EnableMultiplexed = true
		opts.Multiplexed = p
	}
}

// WithMsg returns a RoundTripOption which sets msg.
func WithMsg(msg codec.IMsg) RoundTripOption {
	return func(opts *RoundTripOptions) {
		opts.Msg = msg
	}
}

// WithLocalAddr returns a RoundTripOption which sets local address.
// Random selection by default when there are multiple NICs.
func WithLocalAddr(addr string) RoundTripOption {
	return func(o *RoundTripOptions) {
		o.LocalAddr = addr
	}
}

// WithDialTimeout returns a RoundTripOption which sets dial timeout.
func WithDialTimeout(dur time.Duration) RoundTripOption {
	return func(o *RoundTripOptions) {
		o.DialTimeout = dur
	}
}

// WithProtocol returns a RoundTripOption which sets protocol name, such as trpc.
func WithProtocol(s string) RoundTripOption {
	return func(o *RoundTripOptions) {
		o.Protocol = s
	}
}

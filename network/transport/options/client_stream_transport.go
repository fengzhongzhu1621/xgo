package options

// CstOptions 客户端流传输配置选项
type CstOptions struct {
	MaxConcurrentStreams int // 每个TCP连接的最大并发流数量
	MaxIdleConnsPerHost  int // 每个主机的最大空闲连接数
}

// ClientStreamTransportOption 客户端流传输选项修改函数类型
type ClientStreamTransportOption func(*CstOptions)

// WithMaxConcurrentStreams 设置每个TCP连接的最大并发流数量
func WithMaxConcurrentStreams(n int) ClientStreamTransportOption {
	return func(opts *CstOptions) {
		opts.MaxConcurrentStreams = n
	}
}

// WithMaxIdleConnsPerHost 设置每个主机的最大空闲连接数
func WithMaxIdleConnsPerHost(n int) ClientStreamTransportOption {
	return func(opts *CstOptions) {
		opts.MaxIdleConnsPerHost = n
	}
}

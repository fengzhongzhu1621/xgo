package options

// ClientTransportOptions 客户端传输配置选项
type ClientTransportOptions struct {
	DisableHTTPEncodeTransInfoBase64 bool // 禁用HTTP传输信息值的base64编码
}

// ClientTransportOption 客户端传输选项修改函数类型
type ClientTransportOption func(*ClientTransportOptions)

// WithDisableEncodeTransInfoBase64 返回禁用HTTP传输信息base64编码的客户端传输选项
func WithDisableEncodeTransInfoBase64() ClientTransportOption {
	return func(opts *ClientTransportOptions) {
		opts.DisableHTTPEncodeTransInfoBase64 = true
	}
}

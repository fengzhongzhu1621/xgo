package referer

// Option 设置参数选项
type Option func(*options)

// WithRefererDomain 设置每个接口对应允许Referer域名
func WithRefererDomain(rpcName string, domains ...string) Option {
	return func(opts *options) {
		if opts.AllowReferer == nil {
			opts.AllowReferer = make(map[string][]string)
		}
		domain, ok := opts.AllowReferer[rpcName]
		if !ok {
			domain = make([]string, 0)
		}
		domain = append(domain, domains...)
		opts.AllowReferer[rpcName] = domain
	}
}

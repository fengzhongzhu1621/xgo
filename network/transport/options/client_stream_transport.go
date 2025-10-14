package options

// cstOptions is the client stream transport options.
type CstOptions struct {
	MaxConcurrentStreams int
	MaxIdleConnsPerHost  int
}

// ClientStreamTransportOption sets properties of ClientStreamTransport.
type ClientStreamTransportOption func(*CstOptions)

// WithMaxConcurrentStreams sets the maximum concurrent streams in each TCP connection.
func WithMaxConcurrentStreams(n int) ClientStreamTransportOption {
	return func(opts *CstOptions) {
		opts.MaxConcurrentStreams = n
	}
}

// WithMaxIdleConnsPerHost sets the maximum idle connections per host.
func WithMaxIdleConnsPerHost(n int) ClientStreamTransportOption {
	return func(opts *CstOptions) {
		opts.MaxIdleConnsPerHost = n
	}
}

package recovery

// Option sets Recovery option.
type Option func(*options)

// WithRecoveryHandler sets Recovery handle function.
func WithRecoveryHandler(rh Handler) Option {
	return func(opts *options) {
		opts.rh = rh
	}
}

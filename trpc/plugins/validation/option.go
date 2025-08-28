package validation

// Option sets an option for the parameter.
type Option func(*options)

// WithErrorLog sets whether to log errors.
func WithErrorLog(allow bool) Option {
	return func(opts *options) {
		opts.EnableErrorLog = allow
	}
}

// WithServerValidateErrCode sets the error code for server-side request validation failure.
func WithServerValidateErrCode(code int) Option {
	return func(opts *options) {
		opts.ServerValidateErrCode = code
	}
}

// WithClientValidateErrCode sets the error code for client-side response validation failure.
func WithClientValidateErrCode(code int) Option {
	return func(opts *options) {
		opts.ClientValidateErrCode = code
	}
}

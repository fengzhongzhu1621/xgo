package debuglog

// Option sets the optiopns.
type Option func(*options)

// WithLogFunc sets the print body method.
func WithLogFunc(f LogFunc) Option {
	return func(opts *options) {
		opts.logFunc = f
	}
}

// WithErrLogLevelFunc sets the log level print method.
func WithErrLogLevelFunc(f LogLevelFunc) Option {
	return func(opts *options) {
		opts.errLogLevelFunc = f
	}
}

// WithNilLogLevelFunc sets the non-error log level print method.
func WithNilLogLevelFunc(f LogLevelFunc) Option {
	return func(opts *options) {
		opts.nilLogLevelFunc = f
	}
}

// WithInclude sets the include options.
func WithInclude(in *RuleItem) Option {
	return func(opts *options) {
		opts.include = append(opts.include, in)
	}
}

// WithExclude sets the exclude options.
func WithExclude(ex *RuleItem) Option {
	return func(opts *options) {
		opts.exclude = append(opts.exclude, ex)
	}
}

// WithEnableColor enable multiple color log output.
func WithEnableColor(enable bool) Option {
	return func(opts *options) {
		opts.enableColor = enable
	}
}

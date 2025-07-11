package validation

import (
	"trpc.group/trpc-go/trpc-go/errs"
)

// defaultOptions is the default options of parameter.
var defaultOptions = options{
	LogFile:               nil,
	EnableErrorLog:        false,
	ServerValidateErrCode: int(errs.RetServerValidateFail),
	ClientValidateErrCode: int(errs.RetClientValidateFail),
}

// options is the options for parameter validation.
type options struct {
	LogFile               []bool `yaml:"logfile"`
	EnableErrorLog        bool   `yaml:"enable_error_log"`
	ServerValidateErrCode int    `yaml:"server_validate_err_code"`
	ClientValidateErrCode int    `yaml:"client_validate_err_code"`
}

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

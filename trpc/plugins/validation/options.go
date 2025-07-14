package validation

import "trpc.group/trpc-go/trpc-go/errs"

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

package xerror

import (
	"errors"
	"fmt"
	"io"
)

// Err frame error value.
var (
	// ErrOK means success.
	ErrOK error
)

const (
	// Success is the success prompt string.
	Success = "success"
)

const (
	RetOK                    = 0   // 成功
	RetClientTimeout         = 101 // 客户端调用超时错误码
	RetClientFullLinkTimeout = 102 // 客户端全链路超时错误码
	RetServerTimeout         = 21  // 服务端超时错误码
	RetServerFullLinkTimeout = 24  // 服务端全链路超时错误码
	RetUnknown               = 999 // 未明确的错误码
)

// ErrorType is the error code type, including framework error code and business error code.
const (
	ErrorTypeFramework       = 1 // 框架错误码
	ErrorTypeBusiness        = 2 // 业务错误码
	ErrorTypeCalleeFramework = 3 // 下游框架错误码（客户端调用返回的错误码，代表下游框架错误码）
)

var (
	// ErrUnknown is an unknown error.
	ErrUnknown = NewFrameError(RetUnknown, "unknown error")
)

// typeDesc returns the error type description.
func typeDesc(t int) string {
	switch t {
	case ErrorTypeFramework:
		return "framework"
	case ErrorTypeCalleeFramework:
		return "callee framework"
	default:
		return "business"
	}
}

// Error 是错误码结构体，包含错误码类型和错误信息
type Error struct {
	Type int    // 错误码类型：1 框架错误码，2 业务错误码，3 下游框架错误码
	Code int32  // 错误码
	Msg  string // 错误信息描述
	Desc string // 错误额外描述，主要用于监控前缀

	cause error      // 内部错误，形成错误链
	stack stackTrace // 调用堆栈，如果错误链已经有堆栈，则不会设置
}

// Error 实现 error 接口，返回错误描述
func (e *Error) Error() string {
	if e == nil {
		return Success
	}

	if e.cause != nil {
		return fmt.Sprintf("type:%s, code:%d, msg:%s, caused by %s",
			typeDesc(e.Type), e.Code, e.Msg, e.cause.Error())
	}
	return fmt.Sprintf("type:%s, code:%d, msg:%s", typeDesc(e.Type), e.Code, e.Msg)
}

// Format 实现 fmt.Formatter 接口，支持不同的格式化输出
func (e *Error) Format(s fmt.State, verb rune) {
	var stackTrace stackTrace
	defer func() {
		if stackTrace != nil {
			stackTrace.Format(s, verb) // 格式化输出堆栈跟踪
		}
	}()

	switch verb {
	case 'v':
		if s.Flag('+') { // 详细模式
			_, _ = fmt.Fprintf(s, "type:%s, code:%d, msg:%s", typeDesc(e.Type), e.Code, e.Msg)
			if e.stack != nil {
				stackTrace = e.stack
			}
			if e.Unwrap() != nil {
				_, _ = fmt.Fprintf(s, "\nCause by %+v", e.Unwrap())
			}
			return
		}
		fallthrough // 执行当前 case 后，无条件继续执行下一个 case 的代码块
	case 's': // 字符串模式
		_, _ = io.WriteString(s, e.Error())
	case 'q': // 引号模式
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	default:
		_, _ = fmt.Fprintf(s, "%%!%c(errs.Error=%s)", verb, e.Error())
	}
}

// Unwrap 支持 Go 1.13+ 错误链，返回内部错误
func (e *Error) Unwrap() error { return e.cause }

// IsTimeout 检查此错误是否为指定类型的超时错误
func (e *Error) IsTimeout(typ int) bool {
	return e.Type == typ &&
		(e.Code == RetClientTimeout ||
			e.Code == RetClientFullLinkTimeout ||
			e.Code == RetServerTimeout ||
			e.Code == RetServerFullLinkTimeout)
}

// ErrCode 允许使用在 https://go.dev/ref/spec#Numeric_types  中定义的任何整数类型
type ErrCode interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~int | ~uintptr
}

// New 创建一个错误，默认为业务错误类型，提高业务开发效率
func New[T ErrCode](code T, msg string) error {
	err := &Error{
		Type: ErrorTypeBusiness,
		Code: int32(code),
		Msg:  msg,
	}
	if traceable {
		err.stack = callers() // 添加调用堆栈
	}

	return err
}

// Newf 创建一个错误，默认为业务错误类型，msg 支持格式化字符串
func Newf[T ErrCode](code T, format string, params ...interface{}) error {
	msg := fmt.Sprintf(format, params...)
	err := &Error{
		Type: ErrorTypeBusiness,
		Code: int32(code),
		Msg:  msg,
	}
	if traceable {
		err.stack = callers()
	}
	return err
}

// Wrap 创建一个包含输入错误的新错误
// 只有当 traceable 为 true 且输入类型不是 Error 时才添加堆栈，确保错误链中没有多个堆栈
func Wrap[T ErrCode](err error, code T, msg string) error {
	if err == nil {
		return nil
	}
	wrapErr := &Error{
		Type:  ErrorTypeBusiness,
		Code:  int32(code),
		Msg:   msg,
		cause: err, // 保存原始错误
	}
	var e *Error
	// 如果错误链不包含 Error 类型的项，添加堆栈
	if traceable && !errors.As(err, &e) {
		wrapErr.stack = callers()
	}
	return wrapErr
}

// Wrapf 与 Wrap 相同，msg 支持格式化字符串
func Wrapf[T ErrCode](err error, code T, format string, params ...interface{}) error {
	if err == nil {
		return nil
	}
	msg := fmt.Sprintf(format, params...)
	wrapErr := &Error{
		Type:  ErrorTypeBusiness,
		Code:  int32(code),
		Msg:   msg,
		cause: err,
	}
	var e *Error
	// the error chain does not contain item which type is Error, add stack.
	if traceable && !errors.As(err, &e) {
		wrapErr.stack = callers()
	}
	return wrapErr
}

// NewFrameError 创建一个框架错误
func NewFrameError[T ErrCode](code T, msg string) error {
	err := &Error{
		Type: ErrorTypeFramework,
		Code: int32(code),
		Msg:  msg,
		Desc: "xgo",
	}
	if traceable {
		err.stack = callers()
	}
	return err
}

// WrapFrameError 与 Wrap 相同，只是类型为 ErrorTypeFramework
func WrapFrameError[T ErrCode](err error, code T, msg string) error {
	if err == nil {
		return nil
	}
	wrapErr := &Error{
		Type:  ErrorTypeFramework,
		Code:  int32(code),
		Msg:   msg,
		Desc:  "xgo",
		cause: err,
	}
	var e *Error
	// the error chain does not contain item which type is Error, add stack.
	if traceable && !errors.As(err, &e) {
		wrapErr.stack = callers()
	}
	return wrapErr
}

// Code 通过错误获取错误码
func Code(e error) int32 {
	if e == nil {
		return RetOK
	}

	// 先进行类型断言比直接使用 errors.As 有轻微性能提升
	err, ok := e.(*Error)
	if !ok && !errors.As(e, &err) { // 如果不是 Error 类型
		return RetUnknown
	}
	if err == nil {
		return RetOK
	}
	return err.Code
}

// Msg 通过错误获取错误消息
func Msg(e error) string {
	if e == nil {
		return Success // 成功返回 "success"
	}
	err, ok := e.(*Error)
	if !ok && !errors.As(e, &err) { // 如果不是 Error 类型
		return e.Error() // 返回原始错误的字符串表示
	}
	if err == (*Error)(nil) {
		return Success
	}

	// 对于错误链的情况，err.Error() 会打印整个链，包括当前错误和嵌套的错误消息
	if err.Unwrap() != nil {
		return err.Error()
	}
	return err.Msg
}

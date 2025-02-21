package xgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ErrType 错误类型 / 错误编码 ErrorType represents the type of error.
type ErrType uint

const (
	// ErrUnknownDefault indicates a generic error.
	ErrUnknownDefault ErrType = iota

	// ErrExpectedArgument indicates that an argument was expected.
	ErrExpectedArgument

	// ErrUnknownFlag indicates an unknown flag.
	ErrUnknownFlag

	// ErrUnknownGroup indicates an unknown group.
	ErrUnknownGroup

	// ErrMarshal indicates a marshalling error while converting values.
	ErrMarshal

	// ErrHelp indicates that the built-in help was shown (the error
	// contains the help message).
	ErrHelp

	// ErrNoArgumentForBool indicates that an argument was given for a
	// boolean flag (which don't not take any arguments).
	ErrNoArgumentForBool

	// ErrRequired indicates that a required flag was not provided.
	ErrRequired

	// ErrShortNameTooLong indicates that a short flag name was specified,
	// longer than one character.
	ErrShortNameTooLong

	// ErrDuplicatedFlag indicates that a short or long flag has been
	// defined more than once.
	ErrDuplicatedFlag

	// ErrTag indicates an error while parsing flag tags.
	ErrTag

	// ErrCommandRequired indicates that a command was required but not
	// specified.
	ErrCommandRequired

	// ErrUnknownCommand indicates that an unknown command was specified.
	ErrUnknownCommand

	// ErrInvalidChoice indicates an invalid option value which only allows
	// a certain number of choices.
	ErrInvalidChoice

	// ErrInvalidTag indicates an invalid tag or invalid use of an existing tag.
	ErrInvalidTag
)

// String 将错误码转换为错误信息，设计缺陷会导致每个错误码忘记添加错误信息
func (e ErrType) String() string {
	switch e {
	case ErrUnknown:
		return "unknown"
	case ErrExpectedArgument:
		return "expected argument"
	case ErrUnknownFlag:
		return "unknown flag"
	case ErrUnknownGroup:
		return "unknown group"
	case ErrMarshal:
		return "marshal"
	case ErrHelp:
		return "help"
	case ErrNoArgumentForBool:
		return "no argument for bool"
	case ErrRequired:
		return "required"
	case ErrShortNameTooLong:
		return "short name too long"
	case ErrDuplicatedFlag:
		return "duplicated flag"
	case ErrTag:
		return "tag"
	case ErrCommandRequired:
		return "command required"
	case ErrUnknownCommand:
		return "unknown command"
	case ErrInvalidChoice:
		return "invalid choice"
	case ErrInvalidTag:
		return "invalid tag"
	}

	return "unrecognized error type"
}

// ErrType 也是一种错误
func (e ErrType) Error() string {
	return e.String()
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Error represents a parser error. The error returned from Parse is of this
// type. The error contains both a Type and Message.
type Error struct {
	// The type of error
	Type ErrType

	// The error message
	Message string
}

// Error returns the error's message.
func (e *Error) Error() string {
	return e.Message
}

func NewError(tp ErrType, message string) *Error {
	return &Error{
		Type:    tp,
		Message: message,
	}
}

func NewErrorf(tp ErrType, format string, args ...interface{}) *Error {
	return NewError(tp, fmt.Sprintf(format, args...))
}

// WrapError 将其他错误转换为Error对象.
func WrapError(err error) *Error {
	ret, ok := err.(*Error)

	if !ok {
		return NewError(ErrUnknownDefault, err.Error())
	}

	return ret
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
const (
	NoError                = 0
	ParamError             = 1901002
	BadRequestError        = 1901400
	UnauthorizedError      = 1901401
	ForbiddenError         = 1901403
	NotFoundError          = 1901404
	ConflictError          = 1901409
	SystemError            = 1901500
	TooManyRequests        = 1901429
	StaffUnauthorizedError = 2901403
	IAMNoPermission        = 9900403
)

type Code struct {
	ErrCode int    `json:"Code" form:"Code"`
	Message string `json:"Message" form:"Message"`
}

// Error 将 code 转换为字符串
func (code *Code) Error() string {
	errs, _ := json.Marshal(code)
	return string(errs)
}

func (code *Code) ToError() error {
	errs, _ := json.Marshal(code)
	return fmt.Errorf("%s", string(errs))
}

func NewErrCode(code int, message string) *Code {
	return &Code{
		ErrCode: code,
		Message: message,
	}
}

// ParseErrCode 将字符串转换为 Code 对象
func ParseErrCode(s string) (*Code, error) {
	var (
		code Code
	)
	err := json.Unmarshal([]byte(s), &code)
	if err != nil {
		return nil, err
	}

	return &code, nil
}

var (
	JwtTokenNoneErr = &Code{
		ErrCode: 1006,
		Message: "jwt token not found",
	}

	JwtTokenInvalidErr = &Code{
		ErrCode: 1007,
		Message: "jwt token is invalid",
	}
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 如下是更好的设计方式.
var (
	ErrUnauthorized   = fmt.Errorf("%s", "unauthorized")
	ErrInvalidArg     = fmt.Errorf("%s", "invalid argument")
	ErrInvalidAddress = fmt.Errorf("%s", "invalid address")
	ErrUnknown        = fmt.Errorf("%s", "unknown")
)

var (
	ErrInvalidCopyDestination        = errors.New("copy destination must be non-nil and addressable")
	ErrInvalidCopyFrom               = errors.New("copy from must be non-nil and addressable")
	ErrMapKeyNotMatch                = errors.New("map's key type doesn't match")
	ErrNotSupported                  = errors.New("not supported")
	ErrFieldNameTagStartNotUpperCase = errors.New("copier field name tag must be start upper case")
)

var (
	// ErrTxNotWritable is returned when performing a write operation on a
	// read-only transaction.
	ErrTxNotWritable = errors.New("tx not writable")

	// ErrTxClosed is returned when committing or rolling back a transaction
	// that has already been committed or rolled back.
	ErrTxClosed = errors.New("tx closed")

	// ErrNotFound is returned when an item or index is not in the database.
	ErrNotFound = errors.New("not found")

	// ErrInvalid is returned when the database file is an invalid format.
	ErrInvalid = errors.New("invalid database")

	// ErrDatabaseClosed is returned when the database is closed.
	ErrDatabaseClosed = errors.New("database closed")

	// ErrIndexExists is returned when an index already exists in the database.
	ErrIndexExists = errors.New("index exists")

	// ErrInvalidOperation is returned when an operation cannot be completed.
	ErrInvalidOperation = errors.New("invalid operation")

	// ErrInvalidSyncPolicy is returned for an invalid SyncPolicy value.
	ErrInvalidSyncPolicy = errors.New("invalid sync policy")

	// ErrShrinkInProcess is returned when a shrink operation is in-process.
	ErrShrinkInProcess = errors.New("shrink is in-process")

	// ErrPersistenceActive is returned when post-loading data from an database
	// not opened with Open(":memory:").
	ErrPersistenceActive = errors.New("persistence active")

	// ErrTxIterating is returned when Set or Delete are called while iterating.
	ErrTxIterating = errors.New("tx is iterating")
)

var errTable = map[int32]error{}

// Code2Error 将错误码转换为error.
func Code2Error(code int32) error {
	if err, ok := errTable[code]; ok {
		return err
	}
	return ErrUnknown
}

// MultipleErrors 合并多个错误信息，用空格分割
// Errors is a type of []error
// This is used to pass multiple errors when using parallel or concurrent methods
// and yet subscribe to the error interface.
type MultipleErrors []error

// Prints all errors from asynchronous tasks separated by space.
func (e MultipleErrors) Error() string {
	b := bytes.NewBufferString(EmptyStr)

	for _, err := range e {
		b.WriteString(err.Error())
		b.WriteString(" ")
	}

	return strings.TrimSpace(b.String())
}

// EncodingError 基于字符串的错误.
type EncodingError string

func (e EncodingError) Error() string {
	return string(e)
}

// NoProtoMessageError is returned when the given value does not implement proto.Message.
type NoProtoMessageError struct {
	V interface{}
}

func (e NoProtoMessageError) Error() string {
	rv := reflect.ValueOf(e.V)
	if rv.Kind() != reflect.Ptr {
		return "v is not proto.Message, you must pass pointer value to implement proto.Message"
	}

	return "v is not proto.Message"
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Errorx is a struct for wrap raw err with message
type Errorx struct {
	message string
	err     error
}

// Error return the error message
func (e Errorx) Error() string {
	return e.message
}

// Is reports whether any error in err's chain matches target.
func (e Errorx) Is(target error) bool {
	if target == nil || e.err == nil {
		return e.err == target
	}

	return errors.Is(e.err, target)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func (e *Errorx) Unwrap() error {
	u, ok := e.err.(interface {
		Unwrap() error
	})
	if !ok {
		return e.err
	}

	return u.Unwrap()
}

// makeMessage make the message for error wrap
func makeMessage(err error, layer, function, msg string) string {
	var message string
	var e Errorx
	// func As(err error, target interface{}) bool
	// 将一个错误值转换为特定的错误类型
	// err：要检查的错误值。
	// target：一个指向你想要转换的目标错误类型的指针。
	if errors.As(err, &e) {
		message = fmt.Sprintf("[%s:%s] %s => %s", layer, function, msg, err.Error())
	} else {
		message = fmt.Sprintf("[%s:%s] %s => [Raw:Error] %v", layer, function, msg, err.Error())
	}

	return message
}

// Wrap the error with message
func Wrap(err error, layer string, function string, message string) error {
	if err == nil {
		return nil
	}

	return Errorx{
		message: makeMessage(err, layer, function, message),
		err:     err,
	}
}

// Wrapf the error with formatted message, shortcut for
func Wrapf(err error, layer string, function string, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, args...)

	return Errorx{
		message: makeMessage(err, layer, function, msg),
		err:     err,
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// WrapFuncWithLayerFunction define the func of wrapError for partial specific layer name and function name
type WrapFuncWithLayerFunction func(err error, message string) error

// WrapfFuncWithLayerFunction define the func of wrapfError for partial specific layer name and function name
type WrapfFuncWithLayerFunction func(err error, format string, args ...interface{}) error

// NewLayerFunctionErrorWrap 偏函数  create the wrapError func with specific layer and func
func NewLayerFunctionErrorWrap(layer string, function string) WrapFuncWithLayerFunction {
	return func(err error, message string) error {
		return Wrap(err, layer, function, message)
	}
}

// NewLayerFunctionErrorWrapf 偏函数 create the wrapfError func with specific layer and func
func NewLayerFunctionErrorWrapf(layer string, function string) WrapfFuncWithLayerFunction {
	return func(err error, format string, args ...interface{}) error {
		return Wrapf(err, layer, function, format, args...)
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Validate is a helper that creates an error when a condition is not met.
// 根据条件创建错误
// Play: https://go.dev/play/p/vPyh51XpCBt
func Validate(ok bool, format string, args ...any) error {
	if !ok {
		return fmt.Errorf(format, args...)
	}
	return nil
}

// ErrorsAs is a shortcut for errors.As(err, &&T).
// Play: https://go.dev/play/p/8wk5rH8UfrE
// 错误类型转换
func ErrorsAs[T error](err error) (T, bool) {
	var t T
	ok := errors.As(err, &t)
	return t, ok
}

// messageFromMsgAndArgs 将对象转换为错误字符串
func messageFromMsgAndArgs(msgAndArgs ...any) string {
	if len(msgAndArgs) == 1 {
		if msgAsStr, ok := msgAndArgs[0].(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msgAndArgs[0])
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}

// must panics if err is error or false.
func must(err any, messageArgs ...any) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case bool:
		// false 则 panic 错误信息
		if !e {
			message := messageFromMsgAndArgs(messageArgs...)
			if message == "" {
				message = "not ok"
			}

			panic(message)
		}
		// true 则不处理

	case error:
		message := messageFromMsgAndArgs(messageArgs...)
		if message != "" {
			panic(message + ": " + e.Error())
		} else {
			panic(e.Error())
		}

	default:
		panic("must: invalid err type '" + reflect.TypeOf(err).Name() + "', should either be a bool or an error")
	}
}

// Must is a helper that wraps a call to a function returning a value and an error
// and panics if err is error or false.
// Play: https://go.dev/play/p/TMoWrRp3DyC
func Must[T any](val T, err any, messageArgs ...any) T {
	must(err, messageArgs...)
	return val
}

// Must0 has the same behavior as Must, but callback returns no variable.
// Play: https://go.dev/play/p/TMoWrRp3DyC
func Must0(err any, messageArgs ...any) {
	must(err, messageArgs...)
}

// Must1 is an alias to Must
// Play: https://go.dev/play/p/TMoWrRp3DyC
func Must1[T any](val T, err any, messageArgs ...any) T {
	return Must(val, err, messageArgs...)
}

// Must2 has the same behavior as Must, but callback returns 2 variables.
// Play: https://go.dev/play/p/TMoWrRp3DyC
func Must2[T1, T2 any](val1 T1, val2 T2, err any, messageArgs ...any) (T1, T2) {
	must(err, messageArgs...)
	return val1, val2
}

// Must3 has the same behavior as Must, but callback returns 3 variables.
// Play: https://go.dev/play/p/TMoWrRp3DyC
func Must3[T1, T2, T3 any](val1 T1, val2 T2, val3 T3, err any, messageArgs ...any) (T1, T2, T3) {
	must(err, messageArgs...)
	return val1, val2, val3
}

// Must4 has the same behavior as Must, but callback returns 4 variables.
// Play: https://go.dev/play/p/TMoWrRp3DyC
func Must4[T1, T2, T3, T4 any](val1 T1, val2 T2, val3 T3, val4 T4, err any, messageArgs ...any) (T1, T2, T3, T4) {
	must(err, messageArgs...)
	return val1, val2, val3, val4
}

// Must5 has the same behavior as Must, but callback returns 5 variables.
// Play: https://go.dev/play/p/TMoWrRp3DyC
func Must5[T1, T2, T3, T4, T5 any](val1 T1, val2 T2, val3 T3, val4 T4, val5 T5, err any, messageArgs ...any) (T1, T2, T3, T4, T5) {
	must(err, messageArgs...)
	return val1, val2, val3, val4, val5
}

// Must6 has the same behavior as Must, but callback returns 6 variables.
// Play: https://go.dev/play/p/TMoWrRp3DyC
func Must6[T1, T2, T3, T4, T5, T6 any](val1 T1, val2 T2, val3 T3, val4 T4, val5 T5, val6 T6, err any, messageArgs ...any) (T1, T2, T3, T4, T5, T6) {
	must(err, messageArgs...)
	return val1, val2, val3, val4, val5, val6
}

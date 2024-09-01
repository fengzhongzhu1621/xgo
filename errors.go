package xgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ErrType ErrorType represents the type of error.
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

// 将错误码转换为错误信息.
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

func (e ErrType) Error() string {
	return e.String()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Code struct {
	ErrCode int    `json:"Code" form:"Code"`
	Message string `json:"Message" form:"Message"`
}

func (code *Code) Error() string {
	errs, _ := json.Marshal(code)
	return string(errs)
}

func NewErrCode(code int, message string) *Code {
	return &Code{
		ErrCode: code,
		Message: message,
	}
}

func (code *Code) ToError() error {
	errs, _ := json.Marshal(code)
	return fmt.Errorf("%s", string(errs))
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

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

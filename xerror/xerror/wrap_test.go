package xerror

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/duke-git/lancet/v2/xerror"
	"github.com/pkg/errors"
)

// TestWrapAndUnWrap 根据错误对象创建一个新的XError指针实例，并添加消息。
// func Wrap(cause error, message ...any) *XError
// func Unwrap(err error) *XError
// func (e *XError) Wrap(cause error) *XError
// func (e *XError) Unwrap() error
// func (e *XError) With(key string, value any) *XError
// func (e *XError) Values() map[string]any
func TestWrapAndUnWrap(t *testing.T) {
	err1 := xerror.New("error").With("level", "high")
	fmt.Println(err1.Error()) // error

	// oops: error
	wrapErr := errors.Wrap(err1, "oops")
	fmt.Println(wrapErr.Error())
	// error
	err := xerror.Unwrap(wrapErr)
	fmt.Println(err.Error())
	values := err.Values()
	fmt.Println(values["level"]) // high

	// error: oops
	err2 := err1.Wrap(errors.New("oops"))
	fmt.Println(err2.Error())

	// oops
	err3 := err2.Unwrap()
	fmt.Println(err3.Error())
}

// TestTryUnwrap 如果 err 为 nil，则 TryUnwrap 返回一个有效值。如果 err 不为 nil，Unwrap 将因 err 而引发 panic。
// func TryUnwrap[T any](val T, err error) T
func TestTryUnwrap(t *testing.T) {
	// 返回正确结果
	result1 := xerror.TryUnwrap(strconv.Atoi("42"))
	fmt.Println(result1)

	_, err := strconv.Atoi("4o2")
	defer func() {
		v := recover()
		result2 := reflect.DeepEqual(err.Error(), v.(*strconv.NumError).Error())
		fmt.Println(result2)
	}()

	xerror.TryUnwrap(strconv.Atoi("4o2"))
}

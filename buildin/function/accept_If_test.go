package function

// `AcceptIf`返回一个函数，这个函数与`apply`函数具有相同的签名，并且还包含一个布尔值用以指示成功或失败。
// 一个谓词函数，它接受一个类型为`T`的参数并返回一个布尔值；
// 一个`apply`函数，它也接受一个类型为`T`的参数并返回一个相同类型的修改后的值。
// func AcceptIf[T any](predicate func(T) bool, apply func(T) T) func(T) (T, bool)
import (
	"errors"
	"fmt"
	"testing"

	"github.com/gookit/goutil/basefn"
	"github.com/stretchr/testify/assert"

	"github.com/duke-git/lancet/v2/function"
)

// TestAcceptIf 根据条件判断确定是否执行函数
func TestAcceptIf(t *testing.T) {

	adder := function.AcceptIf(
		function.And(
			func(x int) bool {
				return x > 10
			},
			func(x int) bool {
				return x%2 == 0
			}),
		// 条件为 true 时才执行次函数
		func(x int) int {
			return x + 1
		},
	)

	result, ok := adder(20)
	fmt.Println(result) // 21
	fmt.Println(ok)     // true

	result, ok = adder(21)
	fmt.Println(result) // 0
	fmt.Println(ok)     // false
}

func TestPanicIf(t *testing.T) {
	basefn.PanicIf(false, "")
	assert.Panics(t, func() {
		basefn.PanicIf(true, "a error msg")
	})

	assert.Panics(t, func() {
		basefn.PanicIf(true)
	})
	assert.Panics(t, func() {
		basefn.PanicIf(true, 234)
	})
	assert.Panics(t, func() {
		basefn.PanicIf(true, 234, "abc")
	})

	assert.Panics(t, func() {
		basefn.PanicIf(true, "a error %s", "msg")
	})
}

func TestPanicErr(t *testing.T) {
	basefn.MustOK(nil)
	basefn.PanicErr(nil)
	assert.Panics(t, func() {
		basefn.PanicErr(errors.New("a error"))
	})

	// must ignore
	assert.NotPanics(t, func() {
		basefn.MustIgnore(nil, nil)
	})
	assert.Panics(t, func() {
		basefn.MustIgnore(nil, errors.New("a error"))
	})
}

func TestPanicf(t *testing.T) {
	basefn.MustOK(nil)
	assert.Panics(t, func() {
		basefn.Panicf("hi %s", "inhere")
	})

	assert.Equal(t, "hi", basefn.Must("hi", nil))
	assert.Panics(t, func() {
		basefn.Must("hi", errors.New("a error"))
	})
	assert.Panics(t, func() {
		basefn.MustOK(errors.New("a error"))
	})
}

func TestErrOnFail(t *testing.T) {
	err := errors.New("a error")
	assert.Error(t, basefn.ErrOnFail(false, err))
	assert.NoError(t, basefn.ErrOnFail(true, err))
}

func TestOrValue(t *testing.T) {
	assert.Equal(t, "ab", basefn.OrValue(true, "ab", "dc"))
	assert.Equal(t, "dc", basefn.OrValue(false, "ab", "dc"))
	assert.Equal(t, 1, basefn.FirstOr([]int{1, 2}, 3))
	assert.Equal(t, 3, basefn.FirstOr(nil, 3))
}

func TestOrReturn(t *testing.T) {
	assert.Equal(t, "ab", basefn.OrReturn(true, func() string {
		return "ab"
	}, func() string {
		return "dc"
	}))
	assert.Equal(t, "dc", basefn.OrReturn(false, func() string {
		return "ab"
	}, func() string {
		return "dc"
	}))
}

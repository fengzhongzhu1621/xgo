package function

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/duke-git/lancet/v2/function"
	"github.com/gookit/goutil/basefn"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
	"github.com/stretchr/testify/assert"
)

// TestAfter 创建一个函数，限制最少执行次数，当被调用达到n次或n次以上时调用给定的func（函数）。
func TestAfter(t *testing.T) {
	fn := function.After(2, func() {
		fmt.Println("hello")
	})

	// hello
	fn() // 一次不调用
	fn()
	// hello
	fn()
	// hello
	fn()
}

// TestBefore 创建一个函数，限制最多执行次数，当被调用的次数小于n次时调用func（函数）。
// func Before(n int, fn any) func(args ...any) []reflect.Value
func TestBefore(t *testing.T) {
	fn := function.Before(2, func() {
		fmt.Println("hello")
	})

	// hello
	fn()
	// hello
	fn()

	// 调用次数超过 2 次，不再调用
	// not call
	fn()
}

// 调用（执行）迭代函数n次，返回每次调用结果的数组。
func TestTimes(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := lo.Times(3, func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	})

	is.Equal(len(result1), 3)
	is.Equal(result1, []string{"0", "1", "2"})
}

func TestParallelTimes(t *testing.T) {
	is := assert.New(t)

	result1 := parallel.Times(3, func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	})

	is.Equal(len(result1), 3)
	is.Equal(result1, []string{"0", "1", "2"})
}

func TestCallOn(t *testing.T) {
	assert.NoError(t, basefn.CallOn(false, func() error {
		return errors.New("a error")
	}))
	assert.Error(t, basefn.CallOn(true, func() error {
		return errors.New("a error")
	}))

	assert.Error(t, basefn.CallOrElse(true, func() error {
		return errors.New("a error 001")
	}, func() error {
		return errors.New("a error 002")
	}), "a error 001")
	assert.Error(t, basefn.CallOrElse(false, func() error {
		return errors.New("a error 001")
	}, func() error {
		return errors.New("a error 002")
	}), "a error 002")
}

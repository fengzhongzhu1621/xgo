package channel

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSafeGo(t *testing.T) {
	wg := new(sync.WaitGroup)
	type args struct {
		do_func func()
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "正常函数",
			args: args{do_func: func() {
				defer wg.Done()
				time.Sleep(500 * time.Millisecond)
			}},
		},
		{
			name: "异常函数",
			args: args{do_func: func() {
				defer wg.Done()
				time.Sleep(500 * time.Millisecond)
				panic("error")
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				assert.Nil(t, recover())
			}()
			wg.Add(1)
			SafeGo(tt.args.do_func)
		})
	}
	wg.Wait()
}

func TestSafeGoWait(t *testing.T) {
	type args struct {
		fs []func()
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "正常函数+异常函数",
			args: args{fs: []func(){
				func() {
					time.Sleep(500 * time.Millisecond)
				},
				func() {
					panic("error")
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				assert.Nil(t, recover())
			}()
			SafeGoWait(tt.args.fs...)
		})
	}
}

func TestTry(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	is.False(Try(func() error {
		panic("error")
	}))
	is.True(Try(func() error {
		return nil
	}))
	is.False(Try(func() error {
		return fmt.Errorf("fail")
	}))
}

func TestTryX(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	is.True(Try1(func() error {
		return nil
	}))

	is.True(Try2(func() (string, error) {
		return "", nil
	}))

	is.True(Try3(func() (string, string, error) {
		return "", "", nil
	}))

	is.True(Try4(func() (string, string, string, error) {
		return "", "", "", nil
	}))

	is.True(Try5(func() (string, string, string, string, error) {
		return "", "", "", "", nil
	}))

	is.True(Try6(func() (string, string, string, string, string, error) {
		return "", "", "", "", "", nil
	}))

	is.False(Try1(func() error {
		panic("error")
	}))

	is.False(Try2(func() (string, error) {
		panic("error")
	}))

	is.False(Try3(func() (string, string, error) {
		panic("error")
	}))

	is.False(Try4(func() (string, string, string, error) {
		panic("error")
	}))

	is.False(Try5(func() (string, string, string, string, error) {
		panic("error")
	}))

	is.False(Try6(func() (string, string, string, string, string, error) {
		panic("error")
	}))

	is.False(Try1(func() error {
		return errors.New("foo")
	}))

	is.False(Try2(func() (string, error) {
		return "", errors.New("foo")
	}))

	is.False(Try3(func() (string, string, error) {
		return "", "", errors.New("foo")
	}))

	is.False(Try4(func() (string, string, string, error) {
		return "", "", "", errors.New("foo")
	}))

	is.False(Try5(func() (string, string, string, string, error) {
		return "", "", "", "", errors.New("foo")
	}))

	is.False(Try6(func() (string, string, string, string, string, error) {
		return "", "", "", "", "", errors.New("foo")
	}))
}

func TestTryOr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	a1, ok1 := TryOr(func() (int, error) { panic("error") }, 42)
	a2, ok2 := TryOr(func() (int, error) { return 21, assert.AnError }, 42)
	a3, ok3 := TryOr(func() (int, error) { return 21, nil }, 42)

	is.Equal(42, a1)
	is.False(ok1)

	is.Equal(42, a2)
	is.False(ok2)

	is.Equal(21, a3)
	is.True(ok3)
}

func TestTryOrX(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	{
		a1, ok1 := TryOr1(func() (int, error) { panic("error") }, 42)
		a2, ok2 := TryOr1(func() (int, error) { return 21, assert.AnError }, 42)
		a3, ok3 := TryOr1(func() (int, error) { return 21, nil }, 42)

		is.Equal(42, a1)
		is.False(ok1)

		is.Equal(42, a2)
		is.False(ok2)

		is.Equal(21, a3)
		is.True(ok3)
	}

	{
		a1, b1, ok1 := TryOr2(func() (int, string, error) { panic("error") }, 42, "hello")
		a2, b2, ok2 := TryOr2(func() (int, string, error) { return 21, "world", assert.AnError }, 42, "hello")
		a3, b3, ok3 := TryOr2(func() (int, string, error) { return 21, "world", nil }, 42, "hello")

		is.Equal(42, a1)
		is.Equal("hello", b1)
		is.False(ok1)

		is.Equal(42, a2)
		is.Equal("hello", b2)
		is.False(ok2)

		is.Equal(21, a3)
		is.Equal("world", b3)
		is.True(ok3)
	}

	{
		a1, b1, c1, ok1 := TryOr3(func() (int, string, bool, error) { panic("error") }, 42, "hello", false)
		a2, b2, c2, ok2 := TryOr3(func() (int, string, bool, error) { return 21, "world", true, assert.AnError }, 42, "hello", false)
		a3, b3, c3, ok3 := TryOr3(func() (int, string, bool, error) { return 21, "world", true, nil }, 42, "hello", false)

		is.Equal(42, a1)
		is.Equal("hello", b1)
		is.Equal(false, c1)
		is.False(ok1)

		is.Equal(42, a2)
		is.Equal("hello", b2)
		is.Equal(false, c2)
		is.False(ok2)

		is.Equal(21, a3)
		is.Equal("world", b3)
		is.Equal(true, c3)
		is.True(ok3)
	}

	{
		a1, b1, c1, d1, ok1 := TryOr4(func() (int, string, bool, int, error) { panic("error") }, 42, "hello", false, 42)
		a2, b2, c2, d2, ok2 := TryOr4(func() (int, string, bool, int, error) { return 21, "world", true, 21, assert.AnError }, 42, "hello", false, 42)
		a3, b3, c3, d3, ok3 := TryOr4(func() (int, string, bool, int, error) { return 21, "world", true, 21, nil }, 42, "hello", false, 42)

		is.Equal(42, a1)
		is.Equal("hello", b1)
		is.Equal(false, c1)
		is.Equal(42, d1)
		is.False(ok1)

		is.Equal(42, a2)
		is.Equal("hello", b2)
		is.Equal(false, c2)
		is.Equal(42, d2)
		is.False(ok2)

		is.Equal(21, a3)
		is.Equal("world", b3)
		is.Equal(true, c3)
		is.Equal(21, d3)
		is.True(ok3)
	}

	{
		a1, b1, c1, d1, e1, ok1 := TryOr5(func() (int, string, bool, int, int, error) { panic("error") }, 42, "hello", false, 42, 42)
		a2, b2, c2, d2, e2, ok2 := TryOr5(func() (int, string, bool, int, int, error) { return 21, "world", true, 21, 21, assert.AnError }, 42, "hello", false, 42, 42)
		a3, b3, c3, d3, e3, ok3 := TryOr5(func() (int, string, bool, int, int, error) { return 21, "world", true, 21, 21, nil }, 42, "hello", false, 42, 42)

		is.Equal(42, a1)
		is.Equal("hello", b1)
		is.Equal(false, c1)
		is.Equal(42, d1)
		is.Equal(42, e1)
		is.False(ok1)

		is.Equal(42, a2)
		is.Equal("hello", b2)
		is.Equal(false, c2)
		is.Equal(42, d2)
		is.Equal(42, e2)
		is.False(ok2)

		is.Equal(21, a3)
		is.Equal("world", b3)
		is.Equal(true, c3)
		is.Equal(21, d3)
		is.Equal(21, e3)
		is.True(ok3)
	}

	{
		a1, b1, c1, d1, e1, f1, ok1 := TryOr6(func() (int, string, bool, int, int, int, error) { panic("error") }, 42, "hello", false, 42, 42, 42)
		a2, b2, c2, d2, e2, f2, ok2 := TryOr6(func() (int, string, bool, int, int, int, error) { return 21, "world", true, 21, 21, 21, assert.AnError }, 42, "hello", false, 42, 42, 42)
		a3, b3, c3, d3, e3, f3, ok3 := TryOr6(func() (int, string, bool, int, int, int, error) { return 21, "world", true, 21, 21, 21, nil }, 42, "hello", false, 42, 42, 42)

		is.Equal(42, a1)
		is.Equal("hello", b1)
		is.Equal(false, c1)
		is.Equal(42, d1)
		is.Equal(42, e1)
		is.Equal(42, f1)
		is.False(ok1)

		is.Equal(42, a2)
		is.Equal("hello", b2)
		is.Equal(false, c2)
		is.Equal(42, d2)
		is.Equal(42, e2)
		is.Equal(42, f2)
		is.False(ok2)

		is.Equal(21, a3)
		is.Equal("world", b3)
		is.Equal(true, c3)
		is.Equal(21, d3)
		is.Equal(21, e3)
		is.Equal(21, f3)
		is.True(ok3)
	}
}

func TestTryWithErrorValue(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	err, ok := TryWithErrorValue(func() error {
		// getting error in case of panic, using recover function
		panic("error")
	})
	is.False(ok)
	is.Equal("error", err)

	err, ok = TryWithErrorValue(func() error {
		return errors.New("foo")
	})
	is.False(ok)
	is.EqualError(err.(error), "foo")

	err, ok = TryWithErrorValue(func() error {
		return nil
	})
	is.True(ok)
	is.Equal(nil, err)
}

func TestTryCatch(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	caught := false
	TryCatch(func() error {
		panic("error")
	}, func() {
		// error was caught
		caught = true
	})
	is.True(caught)

	caught = false
	TryCatch(func() error {
		return nil
	}, func() {
		// no error to be caught
		caught = true
	})
	is.False(caught)
}

func TestTryCatchWithErrorValue(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	caught := false
	TryCatchWithErrorValue(func() error {
		panic("error")
	}, func(val any) {
		// error was caught
		caught = val == "error"
	})
	is.True(caught)

	caught = false
	TryCatchWithErrorValue(func() error {
		return nil
	}, func(val any) {
		// no error to be caught
		caught = true
	})
	is.False(caught)
}

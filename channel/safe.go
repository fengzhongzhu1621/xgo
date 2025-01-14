package channel

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var StackTraceBufferSize = 4 << 10

// ///////////////////////////////////////////////////////////////////////////////////////

// 调用函数 func() error，忽略 panic，返回执行是否成功
// Try calls the function and return false in case of error.
func Try(callback func() error) (ok bool) {
	ok = true

	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	err := callback()
	if err != nil {
		ok = false
	}

	return
}

// Try0 has the same behavior as Try, but callback returns no variable.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try0(callback func()) bool {
	return Try(func() error {
		callback()
		return nil
	})
}

// Try1 is an alias to Try.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try1(callback func() error) bool {
	return Try(callback)
}

// Try2 has the same behavior as Try, but callback returns 2 variables.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try2[T any](callback func() (T, error)) bool {
	return Try(func() error {
		_, err := callback()
		return err
	})
}

// Try3 has the same behavior as Try, but callback returns 3 variables.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try3[T, R any](callback func() (T, R, error)) bool {
	return Try(func() error {
		_, _, err := callback()
		return err
	})
}

// Try4 has the same behavior as Try, but callback returns 4 variables.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try4[T, R, S any](callback func() (T, R, S, error)) bool {
	return Try(func() error {
		_, _, _, err := callback()
		return err
	})
}

// Try5 has the same behavior as Try, but callback returns 5 variables.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try5[T, R, S, Q any](callback func() (T, R, S, Q, error)) bool {
	return Try(func() error {
		_, _, _, _, err := callback()
		return err
	})
}

// Try6 has the same behavior as Try, but callback returns 6 variables.
// Play: https://go.dev/play/p/mTyyWUvn9u4
func Try6[T, R, S, Q, U any](callback func() (T, R, S, Q, U, error)) bool {
	return Try(func() error {
		_, _, _, _, _, err := callback()
		return err
	})
}

// TryOr has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr[A any](callback func() (A, error), fallbackA A) (A, bool) {
	return TryOr1(callback, fallbackA)
}

// TryOr1 has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr1[A any](callback func() (A, error), fallbackA A) (A, bool) {
	ok := false

	Try0(func() {
		a, err := callback()
		if err == nil {
			fallbackA = a
			ok = true
		}
	})

	// 执行成功：返回回调执行结果 和 成功状态
	// 执行失败：返回默认值 和 失败状态
	return fallbackA, ok
}

// TryOr2 has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr2[A, B any](callback func() (A, B, error), fallbackA A, fallbackB B) (A, B, bool) {
	ok := false

	Try0(func() {
		a, b, err := callback()
		if err == nil {
			fallbackA = a
			fallbackB = b
			ok = true
		}
	})

	return fallbackA, fallbackB, ok
}

// TryOr3 has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr3[A, B, C any](callback func() (A, B, C, error), fallbackA A, fallbackB B, fallbackC C) (A, B, C, bool) {
	ok := false

	Try0(func() {
		a, b, c, err := callback()
		if err == nil {
			fallbackA = a
			fallbackB = b
			fallbackC = c
			ok = true
		}
	})

	return fallbackA, fallbackB, fallbackC, ok
}

// TryOr4 has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr4[A, B, C, D any](callback func() (A, B, C, D, error), fallbackA A, fallbackB B, fallbackC C, fallbackD D) (A, B, C, D, bool) {
	ok := false

	Try0(func() {
		a, b, c, d, err := callback()
		if err == nil {
			fallbackA = a
			fallbackB = b
			fallbackC = c
			fallbackD = d
			ok = true
		}
	})

	return fallbackA, fallbackB, fallbackC, fallbackD, ok
}

// TryOr5 has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr5[A, B, C, D, E any](callback func() (A, B, C, D, E, error), fallbackA A, fallbackB B, fallbackC C, fallbackD D, fallbackE E) (A, B, C, D, E, bool) {
	ok := false

	Try0(func() {
		a, b, c, d, e, err := callback()
		if err == nil {
			fallbackA = a
			fallbackB = b
			fallbackC = c
			fallbackD = d
			fallbackE = e
			ok = true
		}
	})

	return fallbackA, fallbackB, fallbackC, fallbackD, fallbackE, ok
}

// TryOr6 has the same behavior as Must, but returns a default value in case of error.
// Play: https://go.dev/play/p/B4F7Wg2Zh9X
func TryOr6[A, B, C, D, E, F any](callback func() (A, B, C, D, E, F, error), fallbackA A, fallbackB B, fallbackC C, fallbackD D, fallbackE E, fallbackF F) (A, B, C, D, E, F, bool) {
	ok := false

	Try0(func() {
		a, b, c, d, e, f, err := callback()
		if err == nil {
			fallbackA = a
			fallbackB = b
			fallbackC = c
			fallbackD = d
			fallbackE = e
			fallbackF = f
			ok = true
		}
	})

	return fallbackA, fallbackB, fallbackC, fallbackD, fallbackE, fallbackF, ok
}

// TryWithErrorValue has the same behavior as Try, but also returns value passed to panic.
// 返回 recover 错误和执行状态
// Play: https://go.dev/play/p/Kc7afQIT2Fs
func TryWithErrorValue(callback func() error) (errorValue any, ok bool) {
	ok = true

	defer func() {
		if r := recover(); r != nil {
			ok = false
			errorValue = r
		}
	}()

	err := callback()
	if err != nil {
		ok = false
		errorValue = err
	}

	return
}

// TryCatch has the same behavior as Try, but calls the catch function in case of error.
// Play: https://go.dev/play/p/PnOON-EqBiU
func TryCatch(callback func() error, catch func()) {
	if !Try(callback) {
		catch()
	}
}

// TryCatchWithErrorValue has the same behavior as TryWithErrorValue, but calls the catch function in case of error.
// Play: https://go.dev/play/p/8Pc9gwX_GZO
func TryCatchWithErrorValue(callback func() error, catch func(any)) {
	if err, ok := TryWithErrorValue(callback); !ok {
		catch(err)
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////

// ErrFunc type
type ErrFunc func() error

// GoWait is a basic promise implementation: it wraps calls a function in a goroutine
// and returns a channel which will later return the function's return value.
func GoWait(f func() error) error {
	ch := make(chan error)
	go func() {
		ch <- f()
	}()

	// 等待子协程执行完成
	return <-ch
}

// CallOn call func on condition is true
func CallOn(cond bool, fn ErrFunc) error {
	if cond {
		return fn()
	}
	return nil
}

// CallOrElse call okFunc() on condition is true, else call elseFn()
func CallOrElse(cond bool, okFn, elseFn ErrFunc) error {
	if cond {
		return okFn()
	}
	return elseFn()
}

// ///////////////////////////////////////////////////////////////////////////////////////
// SafeGo 安全的异步执行任务

func defaultRecoverHandle() {
	if err := recover(); err != nil {
		buf := make([]byte, StackTraceBufferSize)
		ss := runtime.Stack(buf, false)
		if ss > StackTraceBufferSize {
			ss = StackTraceBufferSize
		}
		buf = buf[:ss]
		log.Error(fmt.Sprintf("SafeGoWait.recover error %s[%v]", string(buf), err))
	}
}

func SafeGo(f func()) {
	safeFunc := SafeFunc(f)
	go safeFunc()
}

// SafeGoTicker 安全的异步执行周期性任务
func SafeGoTicker(f func(), tick *time.Ticker) {
	go func() {
		for {
			// 在每次循环中，通过 <-tick.C 语句等待 tick 的通道 C 接收到一个值。
			// 当 tick 的 Tick 方法被调用时，它会向通道 C 发送一个值，从而唤醒等待的 goroutine。
			<-tick.C
			SafeGo(f)
		}
	}()
}

// SafeGoWait 安全并发执行多个任务
func SafeGoWait(funcs ...func()) {
	if len(funcs) == 0 {
		return
	}
	var wg sync.WaitGroup
	for _, f := range funcs {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				defaultRecoverHandle()
			}()
			f()
		}()
	}
	// 等待所有异步任务执行完成
	wg.Wait()
}

func SafeFunc(f func()) func() {
	return func() {
		defer defaultRecoverHandle()
		f()
	}
}

func SafeLoop(ctx context.Context, f func()) {
	safeFunc := SafeFunc(f)
	for {
		safeFunc()

		select {
		// 子协程收到上下文取消通知，不会停止而是重新执行
		case <-ctx.Done():
		default:
		}
	}
}

func SafeLoopGo(ctx context.Context, f func()) {
	safeFunc := SafeFunc(f)
	go func() {
		for {
			safeFunc()
			select {
			// 子协程收到上下文取消通知，不会停止而是重新执行
			case <-ctx.Done():
			default:
			}
		}
	}()
}

func SafeLoopex(ctx context.Context, f func(), onStart, onStop func()) {
	onStart()
	defer onStop()
	safeFunc := SafeFunc(f)
	for {
		safeFunc()
		select {
		case <-ctx.Done():
		default:
		}
	}
}

func SafeLoopGoex(ctx context.Context, f func(), onStart, onStop func()) {
	onStart()
	defer onStop()
	safeFunc := SafeFunc(f)
	go func() {
		for {
			safeFunc()
			select {
			case <-ctx.Done():
			default:
			}
		}
	}()
}

// ///////////////////////////////////////////////////////////////////////////////////////
// SafeGoV2 async run a func.
// If the func panics, the panic value will be handle by errHandler.
func SafeGoV2(fn func(), errHandler func(error)) {
	go func() {
		if err := SafeRun(fn); err != nil {
			errHandler(err)
		}
	}()
}

// SafeGoWithError async run a func with error.
// If the func panics, the panic value will be handle by errHandler.
func SafeGoWithError(fn func() error, errHandler func(error)) {
	go func() {
		if err := SafeRunWithError(fn); err != nil {
			errHandler(err)
		}
	}()
}

// SafeRun sync run a func. If the func panics, the panic value is returned as an error.
func SafeRun(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	fn()
	return nil
}

// SafeRunWithError sync run a func with error.
// If the func panics, the panic value is returned as an error.
func SafeRunWithError(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	return fn()
}

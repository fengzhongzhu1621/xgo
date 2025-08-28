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

// ErrFunc type
type ErrFunc func() error

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

func SafeLoopEx(ctx context.Context, f func(), onStart, onStop func()) {
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

func SafeLoopGoEx(ctx context.Context, f func(), onStart, onStop func()) {
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

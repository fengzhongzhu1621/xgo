package channel

var ReallyCrash = true

func logPanic(r interface{}) {
}

// 全局默认的Panic处理
var PanicHandlers = []func(any){logPanic}

// 允许外部传入额外的异常处理
func HandleCrash(additionalHandlers ...func(any)) {
	if r := recover(); r != nil {
		for _, fn := range PanicHandlers {
			fn(r)
		}
		for _, fn := range additionalHandlers {
			fn(r)
		}
		if ReallyCrash {
			panic(r)
		}
	}
}

// Go 使用自定义的 go 函数，避免了自己忘记增加 panic 的处理。
// 如果没有对相应的Goroutine 进行异常处理，会导致主线程 panic
func Go(fn func()) {
	go func() {
		defer HandleCrash()
		fn()
	}()
}

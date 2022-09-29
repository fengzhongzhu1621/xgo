package async

import (
	"reflect"
	"runtime"

	"github.com/fengzhongzhu1621/xgo"
)

// Internal usage to gather results from tasks.
type execResult struct {
	err     error
	results []reflect.Value
	key     string
}

// funcs is the struct used to control the stack
// of functions to be executed.
type funcs struct {
	Stack interface{} // 函数任务列表
}

// 任务顺序执行，前一个任务的执行结果作为下一个任务的输入
// ExecInSeries executes recursively each task of the stack until it reachs
// the bottom of the stack or it is interrupted by an error.
func (f *funcs) ExecInSeries(args ...reflect.Value) ([]interface{}, error) {
	var (
		// 堆栈是一个切片，表示多个任务
		fns = f.Stack.([]reflect.Value)
		// 任务的数量
		fnl = len(fns)
		// 任务输出参数最后一个参数是否为 error
		returnsError bool // true if function has the last return value is of type `error`
		// 执行的任务
		fn reflect.Value // Get function to be executed
		// 执行的任务类型
		fnt reflect.Type // Get type of the function to be executed
		// 任务执行后的返回结果
		outArgs []reflect.Value // Parameters to be sent to the next function
	)

	// end of stack, no need to proceed
	if fnl == 0 {
		result := xgo.EmptyResult
		if l := len(args); l > 0 {
			for i := 0; i < l; i++ {
				result = append(result, args[i].Interface())
			}
		}
		return result, nil
	}

	// Get function to be executed
	// 从堆栈获取一个任务
	fn = fns[0]
	// Get type of the function to be executed
	// 获取任务类型
	fnt = fn.Type()

	// If function expect any argument
	// 判断任务输出参数的最后一个参数是否是 error
	if l := fnt.NumOut(); l > 0 {
		// Get last argument of the function
		// 获得最后一个输出参数类型
		lastArg := fnt.Out(l - 1)

		// Check if the last argument is a error
		// 判断参数的类型是否是 error
		returnsError = reflect.Zero(lastArg).Interface() == xgo.EmptyError
	}

	// Remove current function from the stack
	// 任务出栈
	f.Stack = fns[1:fnl]
	// 执行任务
	outArgs = fn.Call(args)

	// 函数输出值的数量
	lr := len(outArgs)

	// If function is expecting an `error`
	// 去掉 error返回
	if lr > 0 && returnsError {
		// check if the error occurred, if so returns the error and break the execution
		if e, ok := outArgs[lr-1].Interface().(error); ok {
			return xgo.EmptyResult, e
		}
		lr--
	}

	// 将返回结果作为下一个任务的输入继续执行
	return f.ExecInSeries(outArgs[:lr]...)
}

// 并发执行多个任务
// parallel: true表示并发执行，限制并发的数量为CPU核数
// ExecInParallel executes all functions in the stack in Parallel.
func (f *funcs) ExecConcurrent(parallel bool) (Results, error) {
	var (
		results Results
		errs    xgo.MultipleErrors
	)

	if funcs, ok := f.Stack.([]reflect.Value); ok {
		// 任务是一个切片
		results, errs = execSlice(funcs, parallel)
	} else if mapFuncs, mapOk := f.Stack.(map[string]reflect.Value); mapOk {
		// 任务是一个map
		results, errs = execMap(mapFuncs, parallel)
	} else {
		// Incorret t.Stack type
		panic("Stack type must be of type []reflect.Value or map[string]reflect.Value.")
	}

	if len(errs) == 0 {
		return results, nil
	}

	return results, errs
}

// 并发执行多个任务，多个任务是切片类型.
func execSlice(funcs []reflect.Value, parallel bool) (SliceResults, xgo.MultipleErrors) {
	var (
		errs    xgo.MultipleErrors          // 任务执行失败结果
		results = SliceResults{}            // 任务执行成功给结果
		ls      = len(funcs)                // Length of the functions to execute 任务数量
		cr      = make(chan execResult, ls) // Creates buffered channel for errors 任务执行结果管道
	)
	// If parallel, tries to distribute the go routines among the cores, creating
	// at most `runtime.GOMAXPROCS` go routine.
	if parallel {
		// 设置并发数量，默认是CPU核数
		// Creates bufferd channel for controlling CPU usage and guarantee Paralellism
		sem := make(chan int, runtime.GOMAXPROCS(0))
		for i := 0; i < ls; i++ {
			// Fill the buffered channel, if it gets full, go will block the execution
			// until any routine frees the channel
			sem <- 1 // the value doesn't matter
			go execRoutineParallel(funcs[i], cr, sem, xgo.EmptyStr)
		}
	} else {
		for i := 0; i < ls; i++ {
			go execRoutine(funcs[i], cr, xgo.EmptyStr)
		}
	}

	// Consumes the results from the channel
	for i := 0; i < ls; i++ {
		// 获得任务执行结果
		r := <-cr

		if r.err != nil {
			// 记录任务的执行错误
			errs = append(errs, r.err)
		} else if lcr := len(r.results); lcr > 0 {
			// 记录任务的执行结果，每个任务的执行结果是一个切片
			res := make([]interface{}, lcr)
			for j, v := range r.results {
				res[j] = v.Interface()
			}
			results = append(results, res)
		}
	}
	// 返回任务执行成功结果和任务执行失败结果
	return results, errs
}

// 并发执行多个任务，多个任务是字典格式.
func execMap(funcs map[string]reflect.Value, parallel bool) (MapResults, xgo.MultipleErrors) {
	var (
		errs    xgo.MultipleErrors
		results = MapResults{}
		ls      = len(funcs)                // Length of the functions to execute
		cr      = make(chan execResult, ls) // Creates buffered channel for errors
	)

	// If parallel, tries to distribute the go routines among the cores, creating
	// at most `runtime.GOMAXPROCS` go routine.
	if parallel {
		// Creates bufferd channel for controlling CPU usage and guarantee Paralellism
		sem := make(chan int, runtime.GOMAXPROCS(0))
		for k, f := range funcs {
			// Fill the buffered channel, if it gets full, go will block the execution
			// until any routine frees the channel
			sem <- 1 // the value doesn't matter
			go execRoutineParallel(f, cr, sem, k)
		}
	} else {
		for k, f := range funcs {
			go execRoutine(f, cr, k)
		}
	}

	for i := 0; i < ls; i++ {
		r := <-cr

		if r.err != nil {
			errs = append(errs, r.err)
		} else if lcr := len(r.results); lcr > 0 {
			res := make([]interface{}, lcr)
			for j, v := range r.results {
				res[j] = v.Interface()
			}
			results[r.key] = res
		}
	}

	return results, errs
}

// Executes the task and consumes the message of `sem` channel.
func execRoutineParallel(f reflect.Value, c chan execResult, sem chan int, k string) {
	// execute routine
	execRoutine(f, c, k)

	// Once the task has done its job, consumes message from channel `sem`
	<-sem
}

// Executes the task and sends error to the `c` channel.
func execRoutine(f reflect.Value, c chan execResult, key string) {
	var (
		exr = execResult{} // Result
		// 任务无输入参数
		res = f.Call(xgo.EmptyArgs) // Calls the function
	)

	// Get type of the function to be executed
	// 获得函数类型对象
	fnt := f.Type()

	// Check if function returns any value
	if l := fnt.NumOut(); l > 0 {
		// Gets last return value type of the function
		lastArg := fnt.Out(l - 1)

		// Check if the last return value is error
		if reflect.Zero(lastArg).Interface() == xgo.EmptyError {
			// If so and an error occurred, set the execResult.error to the occurred error
			if e, ok := res[l-1].Interface().(error); ok {
				// 记录任务执行的错误
				exr.err = e
			}
			// Decrements l so the results returned doesn't contain the error
			l--
		}

		// If no error occurred, fills the exr.results
		if exr.err == nil && l > 0 {
			// 记录任务执行成功的结果
			exr.results = res[:l]
			// If result has a key
			// 给执行结果打个标签
			if key != "" {
				exr.key = key
			}
		}
	}
	// Sends message to the error channel
	c <- exr
}

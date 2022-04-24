package async

import (
	"fmt"
	"reflect"
)

// Interface fot async to handle user user functions
type taskier interface {
	GetFuncs() (interface{}, error)
}

// Type used as a slice of tasks
// 任务列表
type Tasks []interface{}

// 获得任务列表
func (t Tasks) GetFuncs() (interface{}, error) {
	var (
		l   = len(t) // 任务的数量
		fns = make([]reflect.Value, l)
	)

	for i := 0; i < l; i++ {
		// 获得任务的运行时表示
		f := reflect.Indirect(reflect.ValueOf(t[i]))
		// 判断任务是否是函数
		if f.Kind() != reflect.Func {
			return fns, fmt.Errorf("%T must be a Function ", f)
		}
		// 将函数转换为 reflect.Value 格式
		fns[i] = f
	}

	return fns, nil
}

// Type used as a map of tasks
// 任务字典
type MapTasks map[string]interface{}

func (mt MapTasks) GetFuncs() (interface{}, error) {
	var fns = map[string]reflect.Value{}

	for k, v := range mt {
		f := reflect.Indirect(reflect.ValueOf(v))

		if f.Kind() != reflect.Func {
			return fns, fmt.Errorf("%T must be a Function ", f)
		}

		fns[k] = f
	}

	return fns, nil
}

// Waterfall executes every task sequencially.
// The execution flow may be interrupted by returning an error.
// `firstArgs` is a slice of parameters to be passed to the first task of the stack.
func Waterfall(stack Tasks, firstArgs ...interface{}) ([]interface{}, error) {
	var (
		err  error
		args []reflect.Value
		f    = &funcs{}
	)
	// Checks if the Tasks passed are valid functions.
	// 获得任务列表
	f.Stack, err = stack.GetFuncs()
	if err != nil {
		panic(err.Error())
	}

	// 构造任务输入参数
	// transform interface{} to reflect.Value for execution
	for i := 0; i < len(firstArgs); i++ {
		args = append(args, reflect.ValueOf(firstArgs[i]))
	}

	// 执行任务列表
	return f.ExecInSeries(args...)
}

// 限制并发数目
func Parallel(stack taskier) (Results, error) {
	return execConcurrentStack(stack, true)
}

// 不限制并发数目
func Concurrent(stack taskier) (Results, error) {
	return execConcurrentStack(stack, false)
}

func execConcurrentStack(stack taskier, parallel bool) (Results, error) {
	var (
		err error
		f   = &funcs{}
	)

	// Checks if the Tasks passed are valid functions.
	f.Stack, err = stack.GetFuncs()

	if err != nil {
		panic(err)
	}
	return f.ExecConcurrent(parallel)
}

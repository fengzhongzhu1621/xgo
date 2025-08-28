# 可捕获

当执行panic时，程序会立即停止当前函数执行，逐层向上回溯调用栈，同时递归执行沿途的defer语句，直到被recover捕获或者导致整个程序崩溃。

```go
func f() {
  	defer func() {
		if r := recover(); r != nil {
  	  	  	fmt.Println("recovered", r)
		}
	}()
	panic("panic in f")
}
```

# 捕获不了的panic
显然仅依靠原生的处理机制还不能解决所有问题，无法”捕获“的场景主要有两个：

1. 系统终端
例如 map 并发写

每个协程都有自己独立的堆栈，panic和recover的作用范围仅限于当前协程的调用栈。所以当某个协程发生panic时，如果没有在该协程内部处理，整个程序会崩溃。go语言中存在调用exit方法退出程序的错误处理，如fatalthrow()、fatalpanic()，这类错误属于运行时的硬中断，会绕过recover的保护机制，直接导致服务崩溃。

2. 跨协程
go原生的recover并不支持跨协程捕获；每个协程都有自己独立的堆栈，panic和recover的作用范围仅限于当前协程的调用栈。所以当某个协程发生panic时，如果没有在该协程内部处理，整个程序会崩溃。

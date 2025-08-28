# 简介
用于在 API 边界和进程之间传递截止时间、取消信号以及请求范围的值。

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

* Deadline()：返回该Context被取消的截止时间。
* Done()：返回一个Channel，当Context被取消时，这个Channel会被关闭。
* Err()：返回Context被取消的原因。
* Value()：获取Context中存储的键值对。


# 特性

* ***可取消性***：Context允许我们在不同的goroutine之间传播取消信号，优雅地终止不再需要的操作。

* ***层次结构***：Context可以派生出子Context，形成一个树状的层次结构。当父Context被取消时，其所有的子Context也会被取消。

* ***值传递***：Context可以携带请求范围的值，这些值可以在整个调用链中传递。

* ***协程安全***：Context被设计为在多个goroutine之间安全使用，无需额外的同步机制。

# 创建 Context

```go
func Background() Context // 返回一个空的Context，通常用作整个请求的顶层Context。 ctx := context.Background()
func TODO() Context // 当不确定应该使用哪种Context时，可以使用TODO()。
```

# 创建派生Context

```go
// 创建可取消的 Context
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)

// 创建带有截止时间的 Context
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc)

// 创建带有超时的 Context
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)

// 创建携带键值对的 Context
func WithValue(parent Context, key, val interface{}) Context
```

# 最佳实践

* 将Context作为函数的第一个参数传递，并命名为ctx
* 不要存储Context， Context应该在函数调用链中传递，而不是存储在结构体中。
* 使用WithCancel的defer模式：在创建可取消的Context时，立即使用defer调用cancel函数。
* 使用自定义类型作为Context值的键，以避免冲突。
* 不要传递nil Context：如果不确定使用哪个Context，可以使用context.TODO()
* 谨慎使用Value：Context的Value方法在查找键时需要遍历整个Context链，对于频繁访问的值，考虑使用其他方式传递。要避免传递可选参数。
* 注意goroutine泄漏：确保在父Context取消时，所有子goroutine都能正确退出，避免goroutine泄漏。
* Context应该是并发安全的：不要修改通过Context传递的数据。

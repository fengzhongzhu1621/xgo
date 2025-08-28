# 简介
CyclicBarrier 允许一组 goroutine 彼此等待，到达一个共同的执行点。
同时，因为它可以被重复使用，所以叫循环栅栏。具体的机制是，大家都在栅栏前等待，等全部都到齐了，就抬起栅栏放行。

WaitGroup 更适合用在“一个 goroutine 等待一组 goroutine 到达同一个执行点”的场景中，或者是不需要重用的场景中。

```go
go get github.com/marusama/cyclicbarrier
```

```go
func New(parties int) CyclicBarrier
func NewWithAction(parties int, barrierAction func() error) CyclicBarrier
```

```go
type CyclicBarrier interface {
    // 等待所有的参与者到达，如果被ctx.Done()中断，会返回ErrBrokenBarrier
    Await(ctx context.Context) error
    // 重置循环栅栏到初始化状态。如果当前有等待者，那么它们会返回ErrBrokenBarrier
    Reset()
    // 返回当前等待者的数量
    GetNumberWaiting() int
    // 参与者的数量
    GetParties() int
    // 循环栅栏是否处于中断状态
    IsBroken() bool
}
```

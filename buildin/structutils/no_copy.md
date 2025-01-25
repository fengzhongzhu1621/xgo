# 作用

通常与 sync.WaitGroup 等同步原语一起出现，可以有效地防止开发者错误地复制包含锁或其他重要数据的结构体，从而避免潜在的错误。

noCopy 结构体本身是一个空结构体，它实现了 Lock 和 Unlock 方法，这两个方法都是空操作。它没有实际的功能属性。

Go Vet 的 copylocks 检查器会检测到这种错误，并提示我们 "Locks Erroneously Passed by Value"。
noCopy 结构体通过与 Go Vet 的 copylocks 检查器配合，来阻止开发者错误地通过值传递包含锁的结构体。
当一个结构体包含 noCopy 类型时，Go Vet 会检查该结构体是否被通过值传递。如果被通过值传递，Go Vet 会发出警告，提醒开发者该结构体不应该被复制。


```go
type WaitGroup struct {
    noCopy noCopy

    state atomic.Uint64
    sema uint32
}
```

# 示例
```go
func func1(wg sync.WaitGroup) {
    wg.Add(1)
    // ...
    wg.Done()
}

func func2() {
    var wg sync.WaitGroup
    func1(wg)
    wg.Wait() // 这里会造成错误
}
```

在 func2 中，我们通过值传递了 wg，导致 func1 中的 wg 成了一个新的副本。
当 func1 中的 wg 执行 Done() 操作时，func2 中的 wg 仍然是 Add(1) 的状态，最终导致 wg.Wait() 无法正常工作。

# 场景
可以有效地防止开发者错误地复制包含锁或其他重要数据的结构体，从而避免潜在的错误。

# 引用
> https://mp.weixin.qq.com/s/cm9kh7KaYbMfwGNIc_vUyg

# struct{}{}
## zerobase
空结构体是没有内存大小的结构体。这句话是没有错的，但是更准确的来说，其实是有一个特殊起点的，那就是 zerobase 变量，这是一个 uintptr 全局变量，占用 8 个字节。

当在任何地方定义无数个 struct {} 类型的变量，编译器都只是把这个 zerobase 变量的地址给出去。涉及到所有内存 size 为 0 的内存分配，那么就是用的同一个地址 &zerobase 。

## 内存对齐
一般情况下，struct 中包含 empty struct ，这个字段是不占用内存空间的，但是有一种情况是特殊的，那就是 empty struct 位于最后一位，它会触发内存对齐 。
```go
//https://go.dev/play/p/HcxlywljovS
// x 是 int 类型，通常占用8字节。
// y 是 string 类型，通常占用16字节（在64位系统上，string 由两个指针组成，每个指针8字节）。
// z 是一个空结构体 struct{}，理论上不占用空间，但在Go中，为了内存对齐，可能会有填充字节。
// | x (8) | y (16) | z (0) | 填充 (8) |
// |-------|--------|-------|----------|
// |  8    |  16    |   0   |    8     |  = 32 字节
type A struct {
    x int
    y string
    z struct{}
}

// x 是 int 类型，占用8字节。
// z 是空结构体 struct{}，不占用空间。
// y 是 string 类型，占用16字节。
// | x (8) | z (0) | y (16) |
// |-------|-------|--------|
// |  8    |   0   |   16   |  = 24 字节
type B struct {
    x int
    z struct{}
    y string
}

func main() {
    // unsafe.Alignof 返回一个类型的对齐值，即该类型变量在内存中的起始地址必须是该对齐值的倍数。
    // 在64位系统上，大多数基本数据类型（如int和string）的对齐值通常是8。
    // 结构体的对齐值由其最大字段的对齐值决定。在结构体A和B中，最大的字段是int和string，它们的对齐值都是8，因此整个结构体的对齐值也是8。
    println(unsafe.Alignof(A{})) // 8
    println(unsafe.Alignof(B{})) // 8
    println(unsafe.Sizeof(A{})) // 32
    println(unsafe.Sizeof(B{})) // 24
}
```


# noCopy

## 简介
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

## 示例
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

## 场景
可以有效地防止开发者错误地复制包含锁或其他重要数据的结构体，从而避免潜在的错误。

## 引用
> https://mp.weixin.qq.com/s/cm9kh7KaYbMfwGNIc_vUyg

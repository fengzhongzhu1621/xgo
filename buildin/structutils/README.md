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

```
// Note that it must not be embedded, due to the Lock and Unlock methods.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()    {}
func (*noCopy) Unlock()  {}
```

Go Vet 的 copylocks 检查器会检测到这种错误，并提示我们 "Locks Erroneously Passed by Value"。
noCopy 结构体通过与 Go Vet 的 copylocks 检查器配合，来阻止开发者错误地通过值传递包含锁的结构体。
当一个结构体包含 noCopy 类型时，Go Vet 会检查该结构体是否被通过值传递。如果被通过值传递，Go Vet 会发出警告，提醒开发者该结构体不应该被复制。

go vet 的 noCopy 机制是一种防止结构体被拷贝的方法，尤其是那些包含同步原语（如 sync.Mutex 和 sync.WaitGroup）的结构，目的是防止意外的锁拷贝，但这种防止并不是强制性的，是否拷贝需要由开发者检测。

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


# go1.23.0 builder.go abi.NoEscape

```go
// A Builder is used to efficiently build a string using [Builder.Write] methods.
// It minimizes memory copying. The zero value is ready to use.
// Do not copy a non-zero Builder.
type Builder struct {
    // addr stores the address of the Builder to detect copies by value.
    // It is initialized to nil and set to the Builder's address on first use.
    addr *Builder // of receiver, to detect copies by value

    // buf is the underlying byte slice that stores the string being built.
    // External users should never get direct access to this buffer,
    // since the slice at some point will be converted to a string using unsafe,
    // also data between len(buf) and cap(buf) might be uninitialized.
    buf []byte
}

// copyCheck checks if the Builder has been copied by value.
// If it has, it panics to prevent incorrect usage.
func (b *Builder) copyCheck() {
    if b.addr == nil {
        // This hack works around a failing of Go's escape analysis
        // that was causing b to escape and be heap allocated.
        // See issue 23382.
        // TODO: once issue 7921 is fixed, this should be reverted to
        // just "b.addr = b".
        b.addr = (*Builder)(abi.NoEscape(unsafe.Pointer(b)))
    } else if b.addr != b {
        panic("strings: illegal use of non-zero Builder copied by value")
    }
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b *Builder) Write(p []byte) (int, error) {
    // Check if the Builder has been copied by value.
    b.copyCheck()

    // Append the data to the buffer.
    b.buf = append(b.buf, p...)

    // Return the number of bytes written and nil error.
    return len(p), nil
}

```


1. 如果 addr == nil：
这是第一次调用 copyCheck()（即 Builder 刚被创建）。
通过 unsafe.Pointer 获取当前 Builder 的地址，并存储到 addr 中。
abi.NoEscape 的作用：绕过 Go 的逃逸分析（避免 b 被错误地分配到堆上）。
2. 如果 addr != b：
说明 Builder 被复制了（addr 仍然指向原来的 Builder，而不是当前的 Builder）。
直接 panic，防止错误使用。

* Go 的逃逸分析可能会错误地认为 b 需要分配到堆上（即使它实际上可以留在栈上）。
* abi.NoEscape 是一个内部函数（Go 运行时提供），用于告诉编译器“这个指针不会逃逸到堆上”，从而优化内存分配。
* 注意：abi.NoEscape 是 Go 内部 API，普通代码不能直接使用（strings.Builder 是标准库的一部分，可以访问）。

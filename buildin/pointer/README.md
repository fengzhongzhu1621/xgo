# 指针限制

* 不能通过递增或递减指针来遍历内存地址，无法使用 p++ 或 p-- 来遍历数组。
* 不能执行数学运算。
* 不同类型的指针不能相互转换。
* 不同类型的指针不能进行比较。
* 不同类型的指针不能互相赋值。


# new and make 的区别

## new
new 是一个分配内存的内置函数。new(T)为 T 类型分配内存，用该类型的零值初始化内存，并返回一个 *T 类型的指针。

## make
make 也是一个内置函数，但它只创建和初始化切片、映射和通道。

```go
slice := make([]int, 5)
m := make(map[string]int)
ch := make(chan int, 2)
```

# unsafe.Pointer

unsafe 包供了一些绕过 Go 类型安全的操作，允许进行低级内存操作。
是一种通用指针类型，用于转换不同类型的指针。它不能参与指针运算。

unsafe.Pointer 可转换为 uintptr 或从 uintptr 转换为 unsafe.Pointer。

## 方法

```go
type ArbitraryType int

type Pointer *ArbitraryType

func Sizeof(x ArbitraryType) uintptr：返回 x 的大小（字节），不包括 x 指向的内容的大小。
func Offsetof(x ArbitraryType) uintptr：返回作为参数传递的结构体成员从结构体起始位置开始的偏移量。
func Alignof(x ArbitraryType) uintptr：返回 x 类型的对齐要求。
```

# uintptr
uintptr 是一种整数类型，其大小足以容纳任何指针的位模式。不过，它没有指针语义，会被垃圾回收。

不安全指针类型不能直接执行算术运算，但可以转换为 uintptr，在 uintptr 上执行算术运算，然后再转换回指针。


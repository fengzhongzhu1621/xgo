# new
用于为变量分配内存，并返回指向该内存的指针。且该内存区域会被初始化为零值（例如整数为 0，布尔为 false，指针为 nil）。
可以用于所有类型，包括基本类型（如 int、float 等）和复合类型（如 struct、数组等）

new(T)，其中 T 是要分配的类型。

```go
ptr := new(int)

// 使用 new 分配一个 struct 的指针
type Person struct {
    Name string
    Age  int
}
p := new(Person)
```

# make
用于为特定类型的数据结构分配和初始化内存，这些类型包括切片（slice）、映射（map）和信道（channel）,返回初始化后的值，而不是指针。
只能用于切片、映射和信道，不能用于其他类型。

make(T, size)，其中 T 是要分配的类型，size 是指定的大小（对于映射和信道是容量）。

```go
slice := make([]int, 5) // 创建一个长度为 5 的切片，初始值为 0
myMap := make(map[string]int)
ch := make(chan int, 2) // 创建一个容量为 2 的缓冲信道
```

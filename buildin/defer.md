# defer历史

* 1.12 及其以前是使用先注册、后执行的方式，defer结构体分配在堆上；
* 1.13 加入了部分局部变量，将defer信息保存到栈上；1.13 版本中并不是所有defer都能够在栈上分配。循环中的defer，无论是显示的for循环，还是goto形成的隐式循环，都只能使用堆上分配，即使循环一次也是只能使用堆上分配。
* 1.14 开始使用开发编码（open coded） defer，不再使用链表；该机制会defer调用直接插入函数返回之前，省去了运行时的 deferproc 或 deferprocStack 操作。

```go
func main() {
    // 每次调用都会生成一个_defer对象，此对象会加入到defer链表的表头，即指向g中的_defer对象g._defer。
    // 函数结束后，会从链表头开始遍历执行所有的defer函数。因为是从链表头部遍历执行，所以先注册的defer函数后被调用。
	defer deferFunc()
	fmt.Println("a")
}

//go:noinline
func deferFunc() {
	fmt.Println("b")
}
```

defer对象是一个链表，在执行时，每个g中都会包含一个defer对象，指向了这个链表的头。
在1.12版本中，_defer对象会被创建在堆上。
```go
// A _defer holds an entry on the list of deferred calls.
// If you add a field here, add code to clear it in freedefer.
type _defer struct {
	siz     int32    // 函数入参和返回值的大小
	started bool     // 是否已运行
	sp      uintptr  // sp at time of defer（调用者的栈指针）
	pc      uintptr  // 调用者的程序计数器
	fn      *funcval // 对应的funcval
	_panic  *_panic  // panic that is running defer
	link    *_defer  // next defer
}

type g struct {
	......

	_panic    *_panic // innermost panic - offset known to liblink
	_defer    *_defer // innermost defer

	......
}
```

# deferproc


# deferreturn

# 1. 赋值

```go
// any is an alias for interface{} and is equivalent to interface{} in all ways.
type any = interface{}
```

```go
var a interface{} = 1       // 存储 int
var b interface{} = "hello world!"  // 存储 string
var c interface{} = []nti{1, 2, 3}  // 存储切片
```


# 2. 比较规则
两个 interface{} 能否比较，取决于它们存储的具体类型和值。

两个 interface{} 时，Go 会依次检查以下内容：比较 type 是否相同。
* 如果 type 相同，再比较 value 是否相等。
* 如果 type 是不可比较的类型（如切片、map、函数等），Go 会在运行时直接 panic。

## 2.1 动态类型和动态值都相同
```go
var a interface{} = 1
var b interface{} = 1
fmt.Println(a == b) // true
```

## 2.2 动态类型不同
即使它们的动态值“看起来”相等，比较结果也是 false。
```go
type MyInt int
var a interface{} = 1
var b interface{} = MyInt(1)
fmt.Println(a == b) // false（动态类型不同）
```

## 2.3 动态类型不可比较
动态类型是不可比较的类型（如切片、map、函数等），直接比较会导致 运行时 panic。
```go
var a interface{} = []int{1, 2, 3}
var b interface{} = []int{1, 2, 3}
fmt.Println(a == b) // panic: runtime error: comparing uncomparable type []int
```

## 2.4 nil 的特殊情况
* 如果两个 interface{} 都是 nil，则它们相等。
* 但如果一个 interface{} 存储了具体类型的 nil（如 *int(nil)），另一个是纯 nil，则它们不相等。
```go
var a interface{} = nil
var b interface{} = (*int)(nil)
fmt.Println(a == b) // false
```

# 3. 安全比较
```go
import "reflect"

func isEqual(a, b interface{}) bool {
    return reflect.DeepEqual(a, b)
}
```

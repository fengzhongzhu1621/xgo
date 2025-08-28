# reflecdt

## reflect.TypeOf
返回一个reflect.Type类型的值，表示给定值的类型。


## reflect.ValueOf 返回给定值的值


## reflect.Value 表示运行时数据的结构类型


## reflect.Type 表示运行时类型信息的结构类型 int, string, float, struct
表示一个具体的、完整的 Go 类型，包括它的结构、方法、字段等所有信息。
* 类型的名称（如 Person、int）
* 类型的种类（通过 Kind() 方法获取）
* 类型的字段（如果是结构体）
* 类型的方法（如果是带有方法的类型）
* 类型的大小、对齐方式等

```go
t := reflect.TypeOf(value)  // value 可以是任意类型的值
```

## reflect.Kind 类型的枚举 reflect.Int, reflect.Struct
是一个枚举类型，表示一个类型的底层基础类别（如 int、struct、ptr、slice 等），是一种更抽象的分类。
Go 中所有的类型，无论是基本类型（如 int、string）、复合类型（如 struct、slice）、还是引用类型（如 map、chan），最终都可以归类为某一种 Kind。

```go
k := reflect.TypeOf(value).Kind()  // 先获取 reflect.Type，再调用 Kind() 方法
```

| Kind 值         | 对应的 Go 类型示例                  |
|-----------------|-------------------------------------|
| `reflect.Int`   | `int`, `int32`, `int64`（注意：不是所有整数类型都是 `reflect.Int`！） |
| `reflect.String`| `string`                           |
| `reflect.Struct`| `struct` 类型                       |
| `reflect.Ptr`   | 指针类型，如 `*int`、`*Person`      |
| `reflect.Slice` | 切片类型，如 `[]int`、`[]string`    |
| `reflect.Map`   | 映射类型，如 `map[string]int`       |
| `reflect.Func`  | 函数类型，如 `func()`               |
| `reflect.Interface` | 接口类型，如 `interface{}`         |


注意：reflect.Kind 只表示“底层类别”，而不是具体的类型。例如，int、int32、int64 在 Kind 上可能都表现为不同的值（取决于具体平台），但通常我们关注的是它们是否属于整数类别。

## reflect.Value.Method 返回与给定名称的方法对应的函数值


# reflect.DeepEqual
比较逻辑
* 基本类型：直接比较值是否相等。
* 数组：比较每个元素是否相等。
* 切片：比较切片长度和每个元素是否相等。
* 映射：比较键值对的数量和每个键对应的值是否相等。
* 结构体：比较每个字段是否相等。
* 指针：比较指针指向的值是否相等。
* 接口：比较接口的动态类型和值是否相等。

局限性
* 未导出字段：对于结构体的未导出字段，reflect.DeepEqual无法访问，因此无法比较。
* 函数比较：如果结构体或容器中包含函数，reflect.DeepEqual不会比较函数逻辑是否相同，只会比较函数指针是否相等。

# 通用格式化动词
| 动词 | 含义 | 示例 | 输出 |
|------|------|------|------|
| `%v` | 默认格式的值 | `fmt.Printf("%v", 42)` | `42` |
| `%+v` | 结构体时打印字段名 | `fmt.Printf("%+v", Person{"Alice", 25})` | `{Name:Alice Age:25}` |
| `%#v` | Go 语法表示（带类型） | `fmt.Printf("%#v", Person{"Alice", 25})` | `main.Person{Name:"Alice", Age:25}` |
| `%T` | 变量的类型 | `fmt.Printf("%T", 42)` | `int` |
| `%%` | 打印 % 字符 | `fmt.Printf("%%")` | `%` |

# 布尔值（bool）
| 动词 | 含义 | 示例 | 输出 |
|------|------|------|------|
| `%t` | `true` 或 `false` | `fmt.Printf("%t", true)` | `true` |

# 整数（int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64）
| 动词 | 含义 | 示例 | 输出 |
|------|------|------|------|
| `%d` | 十进制整数 | `fmt.Printf("%d", 42)` | `42` |
| `%b` | 二进制 | `fmt.Printf("%b", 42)` | `101010` |
| `%o` | 八进制 | `fmt.Printf("%o", 42)` | `52` |
| `%x` | 十六进制（小写） | `fmt.Printf("%x", 42)` | `2a` |
| `%X` | 十六进制（大写） | `fmt.Printf("%X", 42)` | `2A` |
| `%U` | Unicode 码点 | `fmt.Printf("%U", 'A')` | `U+0041` |

# 浮点数（float32, float64）
| 动词 | 含义 | 示例 | 输出 |
|------|------|------|------|
| `%f` | 默认精度（6位小数） | `fmt.Printf("%f", 3.14)` | `3.140000` |
| `%e` | 科学计数法（小写 e） | `fmt.Printf("%e", 3.14)` | `3.140000e+00` |
| `%E` | 科学计数法（大写 E） | `fmt.Printf("%E", 3.14)` | `3.140000E+00` |
| `%g` | 自动选择 `%f` 或 `%e`（更紧凑） | `fmt.Printf("%g", 3.14)` | `3.14` |
| `%G` | 自动选择 `%f` 或 `%E`（更紧凑） | `fmt.Printf("%G", 3.14)` | `3.14` |

* %f 可以指定小数位数，如 %.2f 表示保留 2 位小数：

```go
fmt.Printf("%.2f", 3.14159) // 输出: 3.14
```

# 字符串（string）和字节（[]byte）
| 动词 | 含义 | 示例 | 输出 |
|------|------|------|------|
| `%s` | 字符串 | `fmt.Printf("%s", "hello")` | `hello` |
| `%q` | 带引号的字符串（Go 语法） | `fmt.Printf("%q", "hello")` | `"hello"` |
| `%x` | 十六进制（小写） | `fmt.Printf("%x", []byte("hello"))` | `68656c6c6f` |
| `%X` | 十六进制（大写） | `fmt.Printf("%X", []byte("hello"))` | `68656C6C6F` |


# 指针（pointer）
| 动词 | 含义 | 示例 | 输出 |
|------|------|------|------|
| `%p` | 十六进制内存地址 | `fmt.Printf("%p", &x)` | `0xc00001a0a8`（具体值可能不同） |

* %p 用于打印指针变量的内存地址，以十六进制格式输出
* 示例中的 x 应该是一个已声明的变量
* 实际输出的内存地址每次运行可能会不同，因为 Go 运行时会动态分配内存
* 地址格式通常以 0x 开头，后跟十六进制数字

# 复数（complex64, complex128）


# 宽度和精度控制
可以在动词前指定 宽度 和 精度：

* %5d：至少 5 位宽度，不足补空格。
* %.2f：保留 2 位小数。
* %10.2f：至少 10 位宽度，保留 2 位小数。

```go
fmt.Printf("|%5d|\n", 42)      // 输出: |   42|
fmt.Printf("|%.2f|\n", 3.14159) // 输出: |3.14|
fmt.Printf("|%10.2f|\n", 3.14159) // 输出: |      3.14|
```

左对齐
* 使用 - 标志，如 %-5d 表示左对齐

```go
fmt.Printf("|%-5d|\n", 42) // 输出: |42   |
```

# 打印结构体

```go
type Person struct {
    Name string
    Age  int
}
p := Person{"Alice", 25}
fmt.Printf("%v\n", p)   // 输出: {Alice 25}
fmt.Printf("%+v\n", p)  // 输出: {Name:Alice Age:25}
fmt.Printf("%#v\n", p)  // 输出: main.Person{Name:"Alice", Age:25}
```

# 打印错误
```go
err := errors.New("file not found")
fmt.Printf("%v\n", err) // 输出: file not found
fmt.Printf("%+v\n", err) // 输出: file not found（错误类型可能不同）
```

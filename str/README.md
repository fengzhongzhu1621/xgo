# 字符串

## 不能修改
Go 的字符串本质上是只读的 UTF-8 字节序列，也就是 string 类型底层其实就是：
```go
type stringStruct struct {
    str unsafe.Pointer
    len int
}
```
也就是说，它本质是一个不可变的只读 byte slice

字符串虽然不能直接改，但你可以先转成可变的 []byte，再进行修改，然后再转回来。
```go
s := "hello"
b := []byte(s)
b[0] = 'H'
s2 := string(b) // 内存拷贝
fmt.Println(s2) // 输出：Hallo
```

| 类型   | 用途                     | 底层类型  | 说明                          |
|--------|--------------------------|-----------|-------------------------------|
| `[]byte` | 操作 UTF-8 编码的字节     | `[]uint8` | 适用于二进制数据或 ASCII/UTF-8 文本 |
| `[]rune` | 操作 Unicode 字符        | `[]int32` | 每个 `rune` 表示一个 Unicode 码点   |

* 如果你处理的是英文、数字，[]byte 更高效；
* 如果你处理的是中文、emoji 等字符，建议用 []rune，不会乱码。

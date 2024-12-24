# json.RawMessage
是 Golang 中的一个类型，用于处理 JSON 数据。它是 []byte 的别名，但具有特殊的 JSON 解码行为。当你解码 JSON 数据到一个包含 `json.RawMessage` 的结构体字段时，该字段将保留原始的 JSON 数据，而不会对其进行解析。

这在以下情况下非常有用：

* 当你不知道 JSON 数据的结构时。
* 当你想在稍后的时间点解析 JSON 数据时。
* 当你想将 JSON 数据传递给另一个函数进行处理时。


# niljson
## 可空类型
提供了一系列可空类型，例如 NilString、NilInt、NilFloat、NilBool 等，可以方便地集成到现有的Go结构体中。
为每一种基础类型都定义了一个对应的可空类型，这些可空类型都包含一个指针类型的字段，用于存储实际的值。

### NilString
```go
type NilString struct {
  *string
}
```

* ***IsNil() bool***: 判断值是否为 nil。
* ***NotNil() bool***: 判断值是否不为 nil。
* ***Value() string***: 获取实际的字符串值，如果值为 nil，则返回空字符串。


## JSON 序列化和反序列化支持
可以自动处理JSON字段的序列化和反序列化，将 null JSON值转换为Go语言中的 nil 或零值，反之亦然。




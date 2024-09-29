# json.RawMessage 
是 Golang 中的一个类型，用于处理 JSON 数据。它是 []byte 的别名，但具有特殊的 JSON 解码行为。当你解码 JSON 数据到一个包含 `json.RawMessage` 的结构体字段时，该字段将保留原始的 JSON 数据，而不会对其进行解析。

这在以下情况下非常有用：

* 当你不知道 JSON 数据的结构时。
* 当你想在稍后的时间点解析 JSON 数据时。
* 当你想将 JSON 数据传递给另一个函数进行处理时。

# 简介
mapstructure用于将通用的map[string]interface{}解码到对应的 Go 结构体中，或者执行相反的操作。
先使用标准的encoding/json库将数据解码为map[string]interface{}类型，然后根据标识字段利用mapstructure库转为相应的 Go 结构体以便使用。

# 字段标签

mapstructure解码时会在map[string]interface{}中查找键名name。注意，这里的name是大小写不敏感的！

```go
type Person struct {
  Name string `mapstructure:"username"`
}
```

# 内嵌结构
mapstructure:",squash"将该结构体的字段提到父结构中

```go
type Person struct {
  Name string
}

type Friend1 struct {
  Person
}

type Friend2 struct {
  Person `mapstructure:",squash"`
}
```

# 未映射的值
如果源数据中有未映射的值（即结构体中无对应的字段），mapstructure默认会忽略它。

可以在结构体中定义一个字段，为其设置
mapstructure:",remain"
标签。这样未映射的值就会添加到这个字段中。注意，这个字段的类型只能为map[string]interface{}或map[interface{}]interface{}


# 忽略字段映射
mapstructure:",omitempty"

# Metadata
```go
// mapstructure.go
type Metadata struct {
  Keys   []string // 解码成功的键名
  Unused []string // 在源数据中存在，但是目标结构中不存在的键名
}
```

# 错误处理
mapstructure执行转换的过程中不可避免地会产生错误，例如 JSON 中某个键的类型与对应 Go 结构体中的字段类型不一致。Decode/DecodeMetadata会返回这些错误

# 弱类型输入
如果不想对结构体字段类型和map[string]interface{}的对应键值做强类型一致的校验。这时可以使用WeakDecode/WeakDecodeMetadata方法，它们会尝试做类型转换。

# 解码器
```go
// mapstructure.go
type DecoderConfig struct {
	ErrorUnused       bool
	ZeroFields        bool
	WeaklyTypedInput  bool
	Metadata          *Metadata
	Result            interface{}
	TagName           string
}
```

* ErrorUnused：为true时，如果输入中的键值没有与之对应的字段就返回错误；
* ZeroFields：为true时，在Decode前清空目标map。为false时，则执行的是map的合并。用在struct到map的转换中；
* WeaklyTypedInput：实现WeakDecode/WeakDecodeMetadata的功能；
* Metadata：不为nil时，收集Metadata数据；
* Result：为结果对象，在map到struct的转换中，Result为struct类型。在struct到map的转换中，Result为map类型；
* TagName：默认使用mapstructure作为结构体的标签名，可以通过该字段设置。

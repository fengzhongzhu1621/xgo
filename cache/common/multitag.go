package common

// MultiTag 用于管理多标签的缓存，提供了一个基本的缓存管理机制，可以用于存储和检索与标签相关的多个值
type MultiTag struct {
	// 嵌入了一个名为 Cache 的接口或类型（代码中没有给出 Cache 的定义，可能是其他地方定义的）
	ICache
	// 一个字符串字段，存储了 MultiTag 的值。
	value string
	// 一个映射，键是字符串，值是字符串切片，用于存储缓存数据
	cache map[string][]string
}

var _ ICache = (*MultiTag)(nil)

// NewMultiTag 创建并返回一个新的 MultiTag 实例，初始化 value 字段
func NewMultiTag(value string) MultiTag {
	return MultiTag{
		value: value,
	}
}

// Parse 初始化缓存
func (x *MultiTag) Parse() error {
	// 获取初始缓存数据，并将其存储在 cache 字段中
	vals, err := x.Scan()
	x.cache = vals

	return err
}

// cached 获得缓存对象，没有则初始化缓存对象
func (x *MultiTag) cached() map[string][]string {
	if x.cache == nil {
		// 获取初始缓存数据，并将其存储在 cache 字段中
		cache, _ := x.Scan()

		if cache == nil {
			cache = make(map[string][]string)
		}

		x.cache = cache
	}

	// 返回当前的缓存映射
	return x.cache
}

// Get 从缓存中获取指定键的值
func (x *MultiTag) Get(key string) string {
	// 获得缓存对象
	c := x.cached()

	// 如果键存在，它返回该键对应的字符串切片中的最后一个元素；否则，返回空字符串。
	if v, ok := c[key]; ok {
		return v[len(v)-1]
	}

	return ""
}

// GetMany 返回缓存中指定键的所有值
func (x *MultiTag) GetMany(key string) []string {
	c := x.cached()
	return c[key]
}

// Set 将一个值设置到缓存中指定键的位置
func (x *MultiTag) Set(key string, value string) {
	c := x.cached()
	c[key] = []string{value}
}

// SetMany 将多个值设置到缓存中指定键的位置
func (x *MultiTag) SetMany(key string, value []string) {
	c := x.cached()
	c[key] = value
}

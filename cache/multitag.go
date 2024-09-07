package cache

type MultiTag struct {
	Cache
	value string
	cache map[string][]string
}

// NewMultiTag
func NewMultiTag(v string) MultiTag {
	return MultiTag{
		value: v,
	}
}

// Parse 初始化缓存
func (x *MultiTag) Parse() error {
	vals, err := x.Scan()
	x.cache = vals

	return err
}

// cached 获得缓存对象，没有则初始化缓存对象
func (x *MultiTag) cached() map[string][]string {
	if x.cache == nil {
		cache, _ := x.Scan()

		if cache == nil {
			cache = make(map[string][]string)
		}

		x.cache = cache
	}

	return x.cache
}

// Get 从缓存中获取值
func (x *MultiTag) Get(key string) string {
	// 获得缓存对象
	c := x.cached()

	// 从缓存中获取值，取最新的一条记录
	if v, ok := c[key]; ok {
		return v[len(v)-1]
	}

	return ""
}

func (x *MultiTag) GetMany(key string) []string {
	c := x.cached()
	return c[key]
}

func (x *MultiTag) Set(key string, value string) {
	c := x.cached()
	c[key] = []string{value}
}

func (x *MultiTag) SetMany(key string, value []string) {
	c := x.cached()
	c[key] = value
}

package cache

type MultiTag struct {
	Cache
	value string
	cache map[string][]string
}

func NewMultiTag(v string) MultiTag {
	return MultiTag{
		value: v,
	}
}

func (x *MultiTag) Parse() error {
	vals, err := x.Scan()
	x.cache = vals

	return err
}

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

func (x *MultiTag) Get(key string) string {
	c := x.cached()

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

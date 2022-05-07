package utils

// 定义集合.
type StringSet map[string]struct{}

// 判断key是否在集合中.
func (s StringSet) Exist(key string) bool {
	_, ok := s[key]

	return ok
}

// 将key添加到集合中.
func (s StringSet) Append(key string) {
	s[key] = struct{}{}
}

// 从集合中删除key.
func (s StringSet) Remove(key string) {
	delete(s, key)
}

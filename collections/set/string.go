package collections

import "strings"

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

// StringSet is a set of string
type StringSet2 struct {
	Data map[string]struct{}
}

// Has return true if set contains the key
func (s *StringSet2) Has(key string) bool {
	_, ok := s.Data[key]
	return ok
}

// Add a key into set
func (s *StringSet2) Add(key string) {
	s.Data[key] = struct{}{}
}

// Append append keys into set
func (s *StringSet2) Append(keys ...string) {
	for _, key := range keys {
		s.Data[key] = struct{}{}
	}
}

// Size return the size of set
func (s *StringSet2) Size() int {
	return len(s.Data)
}

// ToSlice return key slice
func (s *StringSet2) ToSlice() []string {
	l := make([]string, 0, len(s.Data))
	for k := range s.Data {
		l = append(l, k)
	}
	return l
}

// ToString join the string with sep
func (s *StringSet2) ToString(sep string) string {
	l := s.ToSlice()
	return strings.Join(l, sep)
}

// Diff will return the difference of two set
func (s *StringSet2) Diff(b *StringSet2) *StringSet2 {
	diffSet := NewStringSet()

	for k := range s.Data {
		if !b.Has(k) {
			diffSet.Add(k)
		}
	}
	return diffSet
}

// NewStringSet make a string set
func NewStringSet() *StringSet2 {
	return &StringSet2{
		Data: map[string]struct{}{},
	}
}

// NewStringSetWithValues make a string set with values
func NewStringSetWithValues(keys []string) *StringSet2 {
	set := &StringSet2{
		Data: map[string]struct{}{},
	}
	for _, key := range keys {
		set.Add(key)
	}
	return set
}

// NewFixedLengthStringSet make a string set with fixed length
func NewFixedLengthStringSet(length int) *StringSet2 {
	return &StringSet2{
		Data: make(map[string]struct{}, length),
	}
}

// SplitStringToSet make a string set by split a string into parts
func SplitStringToSet(s string, sep string) *StringSet2 {
	if s == "" {
		return &StringSet2{Data: map[string]struct{}{}}
	}

	data := map[string]struct{}{}
	keys := strings.Split(s, sep)
	for _, key := range keys {
		data[key] = struct{}{}
	}
	return &StringSet2{Data: data}
}

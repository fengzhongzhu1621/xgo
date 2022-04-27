package async

// 返回结果分为两类
// 1. 结果是二维切片
// 2. 结果是字段：key为字符串，value是切片

// Results is an interface used to return sliceResult or mapResults
// from asynchronous tasks. It has methods that should be used to
// get data from the results.
type Results interface {
	Index(int) []interface{}  // Get value by index
	Key(string) []interface{} // Get value by key
	Len() int                 // Get the length of the result
	Keys() []string           // Get the keys of the result
}

// sliceResults is a slice of slice of interface used to return
// results from asynchronous tasks that were passed as slice.
type SliceResults [][]interface{}

// 根据数组索引获取指定结果
// Returns the values returned from ith task.
func (s SliceResults) Index(i int) []interface{} {
	return s[i]
}

// 获得结果的数量
// Returns the length of the results.
func (s SliceResults) Len() int {
	return len(s)
}

// Not supported by sliceResults.
func (s SliceResults) Keys() []string {
	panic("Cannot get map keys from Slice")
}

// Not supported by sliceResults.
func (s SliceResults) Key(k string) []interface{} {
	panic("Cannot get map key from Slice")
}

// sliceResults is a map of string of slice of interface used to return
// results from asynchronous tasks that were passed as map of string.
// key是字符串，value是切片.
type MapResults map[string][]interface{}

// Not supported by mapResults.
func (m MapResults) Index(i int) []interface{} {
	panic("Cannot get index from Map")
}

// Returns the length of the results.
func (m MapResults) Len() int {
	return len(m)
}

// Returns the keys of the result map
// 获得切片数组的key值，放到字符组切片.
func (m MapResults) Keys() []string {
	// 创建字符串切片
	var keys = make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

// Returns the result value by key.
func (m MapResults) Key(k string) []interface{} {
	return m[k]
}

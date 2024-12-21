package compare

import (
	"reflect"
)

// Deprecated: 切片比较.
func CompareStringSliceReflect(a, b []string) bool {
	// 不同类型的值不会深度相等
	// 当两个相同结构体的所有字段对应深度相等的时候，两个结构体深度相等
	// 当两个函数都为nil时，两个函数深度相等，其他情况不相等（相同函数也不相等）
	// 当两个interface的真实值深度相等时，两个interface深度相等
	//
	// 在比较 Map ，slice等是会考虑键值对的顺序的
	// 如果结构体中有切片，切片中的数据顺序如果不通，比较结果会判定为值不一致
	return reflect.DeepEqual(a, b)
}

// CompareStringSlice 切片比较.
func CompareStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	// Golang提供BCE特性，即Bounds-checking elimination
	// 通过b = b[:len(a)]处的bounds check能够明确保证v != b[i]中的b[i]不会出现越界错误，从而避免了b[i]中的越界检查从而提高效率
	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

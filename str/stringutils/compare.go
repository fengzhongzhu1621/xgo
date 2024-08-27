package stringutils

import (
	"reflect"
	"strconv"

	"github.com/tidwall/gjson"
)

// IndexString is a helper function that return true if 'a' is less than 'b'.
// This is a case-insensitive comparison. Use the IndexBinary() for comparing
// case-sensitive strings.
// 区分大小写
func IndexString(a, b string) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] >= 'A' && a[i] <= 'Z' {
			if b[i] >= 'A' && b[i] <= 'Z' {
				// both are uppercase, do nothing
				if a[i] < b[i] {
					return true
				} else if a[i] > b[i] {
					return false
				}
			} else {
				// a[i]是大写字母，b[i]是小写字母
				// a is uppercase, convert a to lowercase
				if a[i]+32 < b[i] {
					return true
				} else if a[i]+32 > b[i] {
					return false
				}
			}
		} else if b[i] >= 'A' && b[i] <= 'Z' {
			// a[i]是小写字母，b[i]是大写字母
			// b is uppercase, convert b to lowercase
			if a[i] < b[i]+32 {
				return true
			} else if a[i] > b[i]+32 {
				return false
			}
		} else {
			// a，b都是小写
			// neither are uppercase
			if a[i] < b[i] {
				return true
			} else if a[i] > b[i] {
				return false
			}
		}
	}
	// 一个字符串是另一个字符串的子串
	return len(a) < len(b)
}

// IndexBinary is a helper function that returns true if 'a' is less than 'b'.
// This compares the raw binary of the string.
func IndexBinary(a, b string) bool {
	return a < b
}

// IndexInt is a helper function that returns true if 'a' is less than 'b'.
func IndexInt(a, b string) bool {
	ia, _ := strconv.ParseInt(a, 10, 64)
	ib, _ := strconv.ParseInt(b, 10, 64)
	return ia < ib
}

// IndexUint is a helper function that returns true if 'a' is less than 'b'.
// This compares uint64s that are added to the database using the
// Uint() conversion function.
func IndexUint(a, b string) bool {
	ia, _ := strconv.ParseUint(a, 10, 64)
	ib, _ := strconv.ParseUint(b, 10, 64)
	return ia < ib
}

// IndexFloat is a helper function that returns true if 'a' is less than 'b'.
// This compares float64s that are added to the database using the
// Float() conversion function.
func IndexFloat(a, b string) bool {
	ia, _ := strconv.ParseFloat(a, 64)
	ib, _ := strconv.ParseFloat(b, 64)
	return ia < ib
}

// IndexJSON provides for the ability to create an index on any JSON field.
// When the field is a string, the comparison will be case-insensitive.
// It returns a helper function used by CreateIndex.
func IndexJSON(path string) func(a, b string) bool {
	return func(a, b string) bool {
		return gjson.Get(a, path).Less(gjson.Get(b, path), false)
	}
}

// IndexJSONCaseSensitive provides for the ability to create an index on
// any JSON field.
// When the field is a string, the comparison will be case-sensitive.
// It returns a helper function used by CreateIndex.
func IndexJSONCaseSensitive(path string) func(a, b string) bool {
	return func(a, b string) bool {
		return gjson.Get(a, path).Less(gjson.Get(b, path), true)
	}
}

// Desc is a helper function that changes the order of an index.
func Desc(less func(a, b string) bool) func(a, b string) bool {
	return func(a, b string) bool { return less(b, a) }
}

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

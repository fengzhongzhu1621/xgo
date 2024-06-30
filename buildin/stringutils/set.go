package stringutils

import (
	"sort"
)

// Intersect 两个字符串数组的交集
func Intersect(a, b []string) []string {
	// 字符串切片排序
	sort.Strings(a)
	sort.Strings(b)

	res := make([]string, 0, func() int {
		if len(a) < len(b) {
			return len(a)
		}
		return len(b)
	}())

	for _, v := range a {
		// 判断字符串是否在一个字符串切片中
		// 用二分法在一个有序字符串数组中寻找特定字符串的索引
		// 如果找到了，那么返回目标字符串在排序后的列表中第一次出现的索引。如果没有找到，那么返回数组中最后一个元素的索引。
		// 所以只要 index 小于最后一个元素的索引，那么目标字符串肯定存在；如果等于最后一个元素的索引，但是值不等于最后一个元素，那么目标字符串就不存在于字符串数组中。
		idx := sort.SearchStrings(b, v)
		if idx < len(b) && b[idx] == v {
			res = append(res, v)
		}
	}
	return res
}

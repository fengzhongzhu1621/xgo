package stringutils

// SliceRemoveDuplicates 从一个字符串切片（[]string）中移除重复的元素，并返回一个新的不包含重复元素的切片。
func SliceRemoveDuplicates(s []string) []string {
	m := make(map[string]bool)
	result := []string{}
	for _, v := range s {
		// 判断 key 是否存在
		if !m[v] {
			m[v] = true
			result = append(result, v)
		}
	}
	return result
}

// RemoveDuplicateElement 切片去重
// 空struct不占内存空间，使用它来实现我们的函数空间复杂度是最低的。
func RemoveDuplicateElement(items []string) []string {
	result := make([]string, 0, len(items))
	// 定义集合
	set := map[string]struct{}{}
	for _, item := range items {
		if _, ok := set[item]; !ok {
			// 如果集合中不存在此元素，则加入集合
			set[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

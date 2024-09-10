package stringutils

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
